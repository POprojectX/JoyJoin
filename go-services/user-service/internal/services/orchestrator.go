package services

import (
	"context"
	"sync"
	"time"
	"user-service/internal/domain"

	"github.com/google/uuid"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type OrchestratorService interface {
	// ==================== AUTH & USER ====================
	RegisterAndCreateProfile(ctx context.Context, email, password, firstName, lastName string) (*domain.UserWithEventsDTO, string, error) // возвращает user + JWT
	LoginAndGetProfile(ctx context.Context, email, password string) (*domain.UserWithEventsDTO, string, error)
	GetUserFullProfile(ctx context.Context, userID uuid.UUID) (*domain.UserWithEventsDTO, error)
	DeleteUserAndCleanup(ctx context.Context, userID uuid.UUID) error // удаляет пользователя + все участия

	// ==================== EVENT MANAGEMENT ====================
	CreateEventWithOwner(ctx context.Context, title, description, location string, slots int, dateFrom, dateTo time.Time, ownerID uuid.UUID) (*domain.Event, error)
	PublishEventAtomic(ctx context.Context, eventID, requesterID uuid.UUID) (*domain.Event, error) // проверяет права + публикует
	CancelEventWithCleanup(ctx context.Context, eventID, requesterID uuid.UUID) error              // отмена + уведомление участников (async)
	DeleteEventWithPermissions(ctx context.Context, eventID, requesterID uuid.UUID) error          // проверка прав + удаление

	// ==================== PARTICIPANT & SLOTS ====================
	JoinEventAsGuest(ctx context.Context, userID, eventID uuid.UUID) error                                              // занять слот + добавить участника (ATOMIC)
	LeaveEventAndFreeSlot(ctx context.Context, userID, eventID uuid.UUID) error                                         // удалить участника + освободить слот
	AssignStaffWithSlotCheck(ctx context.Context, requesterID, targetUserID, eventID uuid.UUID, profRoleID *uint) error // добавить staff без слота
	BulkAssignGuests(ctx context.Context, requesterID, eventID uuid.UUID, userIDs []uuid.UUID) []BulkAssignmentResult   // пакетное добавление с семафором

	// ==================== COMPLEX QUERIES ====================
	GetEventFullDetails(ctx context.Context, eventID, requesterID uuid.UUID) (*EventFullDetailsDTO, error) // ивент + участники + права запрашивающего
	TransferOwnership(ctx context.Context, currentOwnerID, newOwnerID, eventID uuid.UUID) error            // смена владельца с проверками

	// ==================== ADMIN & BATCH ====================
	CleanupCancelledEvents(ctx context.Context, olderThan time.Duration) (int, error) // фоновая задача
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
    
    breaker *gobreaker.CircuitBreaker // некий надзиратель метода. Его задача такая: если метод не выдает ошибки при запросах - мы пропускаем все следуйщие запросы. 
	// Если метод выдал ошибку при n-ым количестве запросов ПОДРЯД, тогда breaker видет что метод "болен" и просто всем последуйщим запросам выдает поментальную ошибку без лишней траны на обработку запросов больным методом. 
    
    taskQueue chan func()  // канал для воркеров (одновременно выполняющихся задач) для WorkerPool
    rateMap map[string]time.Time // устанавливаем для emal-ов кд по запросам. Например один пользователь может делать 1 запрос в секунду. 
	// При больших нагрузках в паре с orchestratorLimiter дает баланс между стабильной нагрузкой на сервер и балансом запросов среди пользователей
	rateMu sync.RWMutex // мутекс для пользователей из rateMap
	rateLimit time.Duration // устанавливаем время кд
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
        orchestratorLimiter: rate.NewLimiter(rate.Every(time.Second), 10000), // 10к req/sec на сервак
        breaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
            Name:        "orchestrator",
            MaxRequests: 6, // 6 прав на ошибку у ивента
            Interval:    10 * time.Second,
            Timeout:     30 * time.Second,
            ReadyToTrip: func(counts gobreaker.Counts) bool {
                failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
                return counts.Requests >= 3 && failureRatio >= 0.6
            },
        }),
        taskQueue: make(chan func(), 10000), // 10к воркеров
		rateMap: make(map[string]time.Time),
		rateLimit: time.Second, // 1 секунда кд для запроса у пользователя
    }
    
    o.starterWorkerPool(20) // воркеры для фоновых задач
    return o
}

func (o *orchestratorService) RegisterAndCreateProfile(ctx context.Context, email, password, firstName, lastName string) (*domain.UserWithEventsDTO, string, error) {
	// orchestratorLimiter, breaker, taskQueue, rateMap, rateMu
	
	// проверка на запросы сервака
	if err := o.orchestratorLimiter.Wait(ctx); err != nil {
		return nil, "", ErrTooManyRequests
	}
	// проверка на ошибки в context
	if err := ctx.Err(); err != nil {
		return nil, "", ErrContextCancelled
	}
	// проверка на кд по мылу
	if !o.checkRateLimitPerEmail(email) {
		return nil, "", ErrTooManyRequests
	}
	// проверка 

	user, err := o.userService.Register(ctx, email, password, firstName, lastName)
	if err != nil {
		return nil, "", err
	}
	jwt, err := o.authService.GenerateToken(user.ID) 
	if err != nil {
		return nil, "", err
	}
	userWithEvents, err := o.userService.GetUserWithEvents(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	return userWithEvents, jwt, nil
}



//используем утилиты
func (o *orchestratorService) starterWorkerPool(workers int) {
	StarterWorkerPool(workers, o.taskQueue)
}

func (o *orchestratorService) checkRateLimitPerEmail(email string) bool {
	return CheckRateLimitPerEmail(email, &o.rateMu, o.rateMap)
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