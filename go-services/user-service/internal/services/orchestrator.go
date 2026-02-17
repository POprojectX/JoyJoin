package services

import (
	"context"
	"log"
	"sync"
	"time"
	"user-service/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

type OrchestratorService interface {
	// ==================== AUTH & USER ====================
	RegisterAndCreateProfile(ctx context.Context, email, password, firstName, lastName string) (*domain.UserWithEventsDTO, string, error) // возвращает user + JWT
	LoginAndGetProfile(ctx context.Context, email, password string) (*domain.UserWithEventsDTO, string, error)
	GetUserFullProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	DeleteUserAndCleanup(ctx context.Context, userID uuid.UUID) error // удаляет пользователя + все участия

	// ==================== EVENT MANAGEMENT ====================
	CreateEventWithOwner(ctx context.Context, title, description, location string, slots int, dateFrom, dateTo time.Time, ownerID uuid.UUID, access domain.Access) (*domain.Event, error)
	PublishEventAtomic(ctx context.Context, eventID, requesterID uuid.UUID) (*domain.Event, error) // проверяет права + публикует
	CancelEventWithCleanup(ctx context.Context, eventID, requesterID uuid.UUID) error              // отмена + уведомление участников (async)
	DeleteEventWithPermissions(ctx context.Context, eventID, requesterID uuid.UUID) error          // проверка прав + удаление

	// ==================== PARTICIPANT & SLOTS ====================
	joinEventAsGuest(ctx context.Context, requesterID, userID, eventID uuid.UUID) error                                              // занять слот + добавить участника (ATOMIC)
	JoinEventAsGuestPublic(ctx context.Context, userID, eventID uuid.UUID) error                                            // использование JoinEventAsGuest без приглашения
	JoinEventAsGuestPrivate(ctx context.Context, requesterID, userID, eventID uuid.UUID) error                                            // использование JoinEventAsGuest по приглашению/оплате
	LeaveEventAndFreeSlot(ctx context.Context, userID, eventID uuid.UUID) error                                         // удалить участника + освободить слот
	AssignStaffWithOutSlotCheck(ctx context.Context, requesterID, targetUserID, eventID uuid.UUID, profRoleID *uint) error // добавить staff без слота
	//BulkAssignGuests(ctx context.Context, requesterID, eventID uuid.UUID, userIDs []uuid.UUID) []BulkAssignmentResult   // пакетное добавление с семафором

	// ==================== COMPLEX QUERIES ====================
	GetEventFullDetails(ctx context.Context, eventID, requesterID uuid.UUID) (*domain.SystemRole, []domain.EventParticipant, error) // ивент + участники + права запрашивающего
	TransferOwnership(ctx context.Context, currentOwnerID, newOwnerID, eventID uuid.UUID) error            // смена владельца с проверками
}

type BulkAssignmentResult struct {
	UserID  uuid.UUID
	Success bool
	Error   error
}

type EventFullDetailsDTO struct {
	Event        domain.Event
	Participants []domain.EventParticipant
	IsOwner      bool
	CanEdit      bool
	CanManage    bool
}

type orchestratorService struct {
    userService        UserService
    eventService       EventService
    participantService ParticipantService
    authService        *AuthService // для генерации токенов
    
    semaphores map[string]chan struct{} // ограничиваем запросы, которые не мешают друг другу, до определенного числа к конкретным event.
    semMu      sync.RWMutex // лочит доступ к семафору и посику ключей в нем. Разлочиваем когда ивента ключ был найден/создан и даем другим горутинам работать с семафором

    distributedLocks map[string]*sync.RWMutex // ограничивает именно критически важные запросы, логика выполнения которых может повредить друг другу. Например смена статуса, набор слотов и тд.
    lockMu           sync.RWMutex //логика такая же как с semMu только для distributedLocks
    
    orchestratorLimiter *rate.Limiter // ограничиваем общее количество запросов к сервису.
    
    //breaker *gobreaker.CircuitBreaker // некий надзиратель метода. Его задача такая: если метод не выдает ошибки при запросах - мы пропускаем все следуйщие запросы. 
	// Если метод выдал ошибку при n-ым количестве запросов ПОДРЯД, тогда breaker видет что метод "болен" и просто всем последуйщим запросам выдает поментальную ошибку без лишней траны на обработку запросов больным методом. 
    
    taskQueue chan func()  // канал для воркеров (одновременно выполняющихся задач) для WorkerPool
    rateMap map[string]time.Time // устанавливаем для emal-ов кд по запросам. Например один пользователь может делать 1 запрос в секунду. 
	// При больших нагрузках в паре с orchestratorLimiter дает баланс между стабильной нагрузкой на сервер и балансом запросов среди пользователей
	rateMu sync.RWMutex // мутекс для пользователей из rateMap
	rateLimit time.Duration // устанавливаем время кд _
}

func NewEventOrchestrator(
    us UserService,
    es EventService,
    ps ParticipantService,
    as *AuthService,
) OrchestratorService {
    o := &orchestratorService{
        userService:        us,
        eventService:       es,
        participantService: ps,
        authService:        as,
        semaphores:         make(map[string]chan struct{}), //буфер на 5 запросов одновременно к одноум ивенту
        distributedLocks:   make(map[string]*sync.RWMutex),
        orchestratorLimiter: rate.NewLimiter(rate.Every(time.Second), 10000), // 10к req/sec на сервак йоу
        // breaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
        //     Name:        "orchestrator",
        //     MaxRequests: 6, // 6 прав на ошибку у ивента
        //     Interval:    10 * time.Second,
        //     Timeout:     30 * time.Second,
        //     ReadyToTrip: func(counts gobreaker.Counts) bool {
        //         failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
        //         return counts.Requests >= 3 && failureRatio >= 0.6
        //     },
        // }),
        taskQueue: make(chan func(), 10000), // 10к воркеров
		rateMap: make(map[string]time.Time),
		rateLimit: time.Second, // 1 секунда кд для запроса у пользователя
    }
    
    o.starterWorkerPool(20) // воркеры для фоновых задач
    return o
}

func (o *orchestratorService) RegisterAndCreateProfile(ctx context.Context, email, password, firstName, lastName string) (*domain.UserWithEventsDTO, string, error) {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return nil, "", ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return nil, "", ErrContextCancelled
	}

	user, err := o.userService.Register(ctx, email, password, firstName, lastName)
	if err != nil {
		return nil, "", err
	}

	var (
		jwt string
		dto *domain.UserWithEventsDTO
	)
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		jwt, err = o.authService.GenerateToken(user.ID)
		return err
	})
	g.Go(func() error {
		var err error
		dto, err = o.userService.GetUserWithEvents(ctx, user.ID)
		return err
	})

	if err := g.Wait(); err != nil {
		if jwt == "" {
			return dto, "", err
		}
		if dto == nil {
			localDto := &domain.UserWithEventsDTO{
				ID : user.ID,
				Email: user.Email,
				FirstName: user.FirstName,
				LastName: user.LastName,
				Events: nil,
			}
			return localDto, "", err
		}
		return nil, "", err
	}

	return dto, jwt, nil
}

func (o *orchestratorService) LoginAndGetProfile(ctx context.Context, email, password string) (*domain.UserWithEventsDTO, string, error) {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return nil, "", ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return nil, "", ErrContextCancelled
	}

	user, err := o.userService.Login(ctx, email, password)
	if err != nil {
		return nil, "", err
	}
	g, _ := errgroup.WithContext(ctx)
	var (
		jwt string
		dto *domain.UserWithEventsDTO
	)
	g.Go(func() error {
		var err error
		jwt, err = o.authService.GenerateToken(user.ID)
		return err
	})
	g.Go(func() error {
		var err error
		dto, err = o.userService.GetUserWithEvents(ctx, user.ID)
		return err
	})
                                     
	if err := g.Wait(); err != nil {
		if jwt == "" {
			return dto, "", err
		}
		if dto == nil {
			localDto := &domain.UserWithEventsDTO{
				ID : user.ID,
				Email: user.Email,
				FirstName: user.FirstName,
				LastName: user.LastName,
				Events: nil,
			}
			return localDto, "", err
		}
		return nil, "", err
	}

	return dto, jwt, nil
}

func (o *orchestratorService) GetUserFullProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return nil, ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return nil, ErrContextCancelled
	}

	return o.userService.GetByID(ctx, userID)
}

func (o *orchestratorService) DeleteUserAndCleanup(ctx context.Context, userID uuid.UUID) error {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}

	UserWithEventsDTO, err := o.userService.GetUserWithEvents(ctx, userID)
	if err != nil {
		return err
	}

	listOfEvents := UserWithEventsDTO.Events
	g, _ := errgroup.WithContext(ctx)

	for _, event := range listOfEvents {
		eventID := event.EventID
		g.Go(func() error {
			return o.participantService.SelfRemove(ctx, userID, eventID)
		})
	}
	if err := g.Wait(); err != nil {
		return  err
	}
	deletedUser, err := o.userService.Delete(ctx, userID)
	if err != nil {
		return err
	}
	if deletedUser != nil {
		return ErrDeleteUserFailed
	}
	return nil
}

func (o *orchestratorService) CreateEventWithOwner(ctx context.Context, title, description, location string, slots int, dateFrom, dateTo time.Time, ownerID uuid.UUID, access domain.Access) (*domain.Event, error) {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return nil, ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return nil, ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(ownerID.String()) {
        return nil, ErrTooManyRequests
    }

	event, err := o.eventService.Create(ctx, title, description, location, slots, dateFrom, dateTo, ownerID, access)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			o.eventService.Delete(ctx, event.ID)
		}
	}()

	participant := &domain.EventParticipant{
		UserID:   ownerID,
		EventID:  event.ID,
		SystemRole: domain.RoleOwner,
		JoinedAt: time.Now(),
		Notes: "Creator of the event.",
	}
	if err := o.participantService.Create(ctx, participant); err != nil {
		return nil, err
	}
	return event, nil
}

func (o *orchestratorService) PublishEventAtomic(ctx context.Context, eventID, requesterID uuid.UUID) (*domain.Event, error) {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return nil, ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return nil, ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(requesterID.String()) {
		return nil, ErrTooManyRequests
	}
	// локамем наш ивент что бы если вдруг другой овнер решил его опопубликовать в то же время, то второй запрос будет ждать пока первый не закончится
	lock := o.getDistributedLock(eventID.String())
	lock.Lock()
	defer lock.Unlock()

	isOwner, err := o.participantService.HasAnyRole(ctx, requesterID, eventID, domain.RoleOwner)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, ErrCantPublishEvent
	}

	event, err := o.eventService.PublishEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (o *orchestratorService) CancelEventWithCleanup(ctx context.Context, eventID, requesterID uuid.UUID) error {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(requesterID.String()) {
		return ErrTooManyRequests
	}

	lock := o.getDistributedLock(eventID.String())
	lock.Lock()
	defer lock.Unlock()

	g, _ := errgroup.WithContext(ctx)
	var eventData *domain.Event
	
	g.Go(func() error {
		var err error
		eventData, err = o.eventService.GetByID(ctx, eventID)
		if err != nil {
			return err
		}
		if eventData == nil {
			return ErrEventNotFound
		}
		return nil
	})
	g.Go(func() error {
		isOwner, err := o.participantService.HasAnyRole(ctx, requesterID, eventID, domain.RoleOwner)
		if err != nil {
			return err
		}
		if !isOwner {
			return ErrCantCancelEvent
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	bgCtx := context.Background()
	for _, users := range eventData.Participants {
		if users.SystemRole == domain.RoleGuest || users.SystemRole == domain.RoleStaff {
			continue
		}
		u := users
		wg.Add(1)
		go func() {
			err := o.participantService.SelfRemove(bgCtx, u.UserID, eventID)
			defer wg.Done()
			if err != nil {
				select {
					case errChan <- err:
					default:
				}
			}
		}()
	}
	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		return ErrToRemoveUserFromEvent
	}

	return nil
}

func (o *orchestratorService) DeleteEventWithPermissions(ctx context.Context, eventID, requesterID uuid.UUID) error {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(requesterID.String()) {
		return ErrTooManyRequests
	}

	lock := o.getDistributedLock(eventID.String())
	lock.Lock()
	defer lock.Unlock()

	isOwner, err := o.participantService.HasAnyRole(ctx, requesterID, eventID, domain.RoleOwner)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrCantDeleteEvent
	}
	_ , err = o.eventService.Delete(ctx, eventID)
	if err != nil {
		return err
	}
	return nil
}

func (o *orchestratorService) GetEventFullDetails(ctx context.Context, eventID, requesterID uuid.UUID) (*domain.SystemRole, []domain.EventParticipant, error) {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return nil, nil, ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return nil, nil, ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(requesterID.String()) {
		return nil, nil,ErrTooManyRequests
	}
	sem := o.getEventSemaphore(eventID.String())
	sem <- struct{}{}
	defer func() {<-sem}()

	_, err := o.eventService.GetByID(ctx, eventID)
	if err != nil {
		return nil, nil, err
	}
	
	var (
		role domain.SystemRole
		litParticipants []domain.EventParticipant
	)
	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		role, err = o.participantService.GetUserRoleInEvent(ctx, requesterID, eventID)
		if err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		var err error
		litParticipants, err = o.eventService.GetEventParticipants(ctx, eventID)
		if err != nil {
			return err
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	return &role, litParticipants, nil
}

func (o *orchestratorService) TransferOwnership(ctx context.Context, currentOwnerID, newOwnerID, eventID uuid.UUID) error {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(currentOwnerID.String()) {
		return ErrTooManyRequests
	}
	sem := o.getEventSemaphore(eventID.String())
	sem <- struct{}{}
	defer func() {<-sem}()
	
	g, _ := errgroup.WithContext(ctx)

	var (participant *domain.EventParticipant)
	g.Go(func() error {
		isRequesterOwner, err := o.participantService.HasAnyRole(ctx, currentOwnerID, eventID, domain.RoleOwner)
		if err != nil {
			return err
		}
		if !isRequesterOwner {
			return ErrNotOwner
		}
		return nil
	})
	g.Go(func() error {
		var err error
		participant, err = o.participantService.GetParticipantByUserID(ctx, newOwnerID, eventID)
		if err != nil {
			return err
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	g2, ctx2 := errgroup.WithContext(ctx)
	g2.Go(func() error {
		err := o.participantService.ChangeRole(ctx2, currentOwnerID, participant.ID, domain.RoleOwner)
		if err != nil {
			return err
		}
		return nil
	})
	if err := g2.Wait(); err != nil {
		return err
	}

	return nil
}

func (o *orchestratorService) joinEventAsGuest(ctx context.Context, requesterID, userID, eventID uuid.UUID) error {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(userID.String()) {
		return ErrTooManyRequests
	}
	sem := o.getEventSemaphore(eventID.String())
	sem <- struct{}{}
	defer func() {<-sem}()
	var isEvent *domain.Event

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		isEvent, err = o.eventService.GetByID(ctx, eventID)
		if err != nil {
			return err
		}
		if isEvent == nil {
			return ErrEventNotFound
		}
		return nil
	})
	g.Go(func() error {
		isUser, err := o.userService.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		if isUser == nil {
			return ErrUserNotFound
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return o.participantService.AssignRole(ctx, requesterID, userID, eventID, domain.RoleGuest, nil)
}
// доп методы для использования под разные доступы к ивентам
func (o *orchestratorService) JoinEventAsGuestPublic(ctx context.Context, userID, eventID uuid.UUID) error {
	return o.joinEventAsGuest(ctx, uuid.Nil, userID, eventID)
}

func (o *orchestratorService) JoinEventAsGuestPrivate(ctx context.Context, requesterID, userID, eventID uuid.UUID) error {
	return o.joinEventAsGuest(ctx, requesterID, userID, eventID)
}

func (o *orchestratorService) LeaveEventAndFreeSlot(ctx context.Context, userID, eventID uuid.UUID) error {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(userID.String()) {
		return ErrTooManyRequests
	}

	sem := o.getEventSemaphore(eventID.String())
	sem <- struct{}{}
	defer func() {<-sem}()

	var isEvent *domain.Event

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		isEvent, err = o.eventService.GetByID(ctx, eventID)
		if err != nil {
			return err
		}
		if isEvent == nil {
			return ErrEventNotFound
		}
		return nil
	})
	g.Go(func() error {
		isUser, err := o.userService.GetByID(ctx, userID)
		if err != nil {
			return err
		}
		if isUser == nil {
			return ErrUserNotFound
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return o.participantService.SelfRemove(ctx, userID, eventID)
}

func (o *orchestratorService) AssignStaffWithOutSlotCheck(ctx context.Context, requesterID, targetUserID, eventID uuid.UUID, profRoleID *uint) error {
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	// проверка на кд по email
	if !o.checkRateLimitPerEmail(targetUserID.String()) {
		return ErrTooManyRequests
	}
	sem := o.getEventSemaphore(eventID.String())
	sem <- struct{}{}
	defer func() {<-sem}()
	var isEvent *domain.Event

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		isEvent, err = o.eventService.GetByID(ctx, eventID)
		if err != nil {
			return err
		}
		if isEvent == nil {
			return ErrEventNotFound
		}
		return nil
	})
	g.Go(func() error {
		isUser, err := o.userService.GetByID(ctx, targetUserID)
		if err != nil {
			return err
		}
		if isUser == nil {
			return ErrUserNotFound
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return o.participantService.AssignRole(ctx, requesterID, targetUserID, eventID, domain.RoleStaff, profRoleID)
}
//используем утилиты
func (o *orchestratorService) starterWorkerPool(workers int) {
	StarterWorkerPool(workers, o.taskQueue)
}

func (o *orchestratorService) checkRateLimitPerEmail(email string) bool {
	result := CheckRateLimitPerEmail(email, &o.rateMu, o.rateMap)
    if !result {
        log.Printf("RATE LIMIT HIT for %s", email)
    }
    return result
}

// доп функции для баланса запросоввы
func (o *orchestratorService) getDistributedLock(eventID string) *sync.RWMutex {
	o.lockMu.RLock()
	if lock, exists := o.distributedLocks[eventID]; exists {
		o.lockMu.RUnlock()
		return lock
	}
	o.lockMu.RUnlock()

	o.lockMu.Lock()
	defer o.lockMu.Unlock()
	
	if lock, ok := o.distributedLocks[eventID]; ok {
		return lock
	}

	lock := &sync.RWMutex{}
	o.distributedLocks[eventID] = lock
	return lock
}

func (o *orchestratorService) getEventSemaphore(eventID string) chan struct{} {
	o.semMu.RLock()
	if sem, exists := o.semaphores[eventID]; exists {
		o.semMu.RUnlock()
		return sem
	}
	o.semMu.RUnlock()

	o.semMu.Lock()
	defer o.semMu.Unlock()

	if sem, ok := o.semaphores[eventID]; ok {
		return sem
	}

	sem := make(chan struct{}, 10)
	o.semaphores[eventID] = sem

	return sem
} 