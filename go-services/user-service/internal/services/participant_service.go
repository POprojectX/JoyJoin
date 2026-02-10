package services

import (
	"context"
	"sync"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"

	"github.com/google/uuid"
)

type ParticipantService interface {
	Create(ctx context.Context, participant *domain.EventParticipant) error
	GetByID(ctx context.Context, id uint) (*domain.EventParticipant, error)
	Update(ctx context.Context, participant *domain.EventParticipant) error
	Delete(ctx context.Context, id uint) error
	
	AssignRole(ctx context.Context, requesterID, targetUserID, eventID uuid.UUID, role domain.SystemRole, profRoleID *uint) error
	ChangeRole(ctx context.Context, requesterID uuid.UUID, participantID uint, newRole domain.SystemRole) error
	RemoveFromEvent(ctx context.Context, requesterID, targetUserID, eventID uuid.UUID) error
	SelfRemove(ctx context.Context, userID, eventID uuid.UUID) error
	
	GetUserRoleInEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.SystemRole, error)
	IsUserOwner(ctx context.Context, userID, eventID uuid.UUID) (bool, error)
	HasAnyRole(ctx context.Context, userID, eventID uuid.UUID, roles ...domain.SystemRole) (bool, error)
	
	// Batch операции
	AssignRolesBulk(ctx context.Context, requesterID, eventID uuid.UUID, assignments []RoleAssignment) []BulkResult
}

type RoleAssignment struct {
	UserID         uuid.UUID
	Role           domain.SystemRole
	ProfessionalID *uint
}

type BulkResult struct {
	UserID uuid.UUID
	Error  error
}

type participantService struct {
	participantRepo repo.ParticipantRepository
	eventRepo       repo.EventRepository
	userRepo        repo.UserRepository
	
	rateLimiter sync.RWMutex
	rateMap     map[string]time.Time
	taskQueue   chan func()
}

func NewParticipantService(
	pr repo.ParticipantRepository,
	er repo.EventRepository,
	ur repo.UserRepository,
) ParticipantService {
	s := &participantService{
		participantRepo: pr,
		eventRepo:       er,
		userRepo:        ur,
		rateMap:         make(map[string]time.Time),
		taskQueue:       make(chan func(), 1000),
	}
	
	s.starterWorkerPool(10)
	
	return s
}


func (s *participantService) Create(ctx context.Context, participant *domain.EventParticipant) error {
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	
	key := participant.UserID.String() + ":" + participant.EventID.String() + ":create"
	if !s.checkRateLimitPerEmail(key) {
		return ErrTooManyRequests
	}
	
	return s.participantRepo.Create(ctx, participant)
}

func (s *participantService) GetByID(ctx context.Context, id uint) (*domain.EventParticipant, error) {
	if err := ctx.Err(); err != nil {
		return nil, ErrContextCancelled
	}
	return s.participantRepo.GetByID(ctx, id)
}

func (s *participantService) Update(ctx context.Context, participant *domain.EventParticipant) error {
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	return s.participantRepo.Update(ctx, participant)
}

func (s *participantService) Delete(ctx context.Context, id uint) error {
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}
	return s.participantRepo.Delete(ctx, id)
}

func (s *participantService) GetUserRoleInEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.SystemRole, error) {
	return s.participantRepo.GetUserRoleInEvent(ctx, userID, eventID)
}

func (s *participantService) IsUserOwner(ctx context.Context, userID, eventID uuid.UUID) (bool, error) {
	return s.participantRepo.IsUserOwner(ctx, userID, eventID)
}

func (s *participantService) HasAnyRole(ctx context.Context, userID, eventID uuid.UUID, roles ...domain.SystemRole) (bool, error) {
	return s.participantRepo.HasAnyRole(ctx, userID, eventID, roles...)
}

//метод который приписывает кого то к ивенту + выдаем мы ему роль
func (s *participantService) AssignRole(
	ctx context.Context,
	requesterID,
	targetUserID,
	eventID uuid.UUID,
	role domain.SystemRole,
	profRoleID *uint,
) error {
	key := requesterID.String() + ":" + eventID.String() + ":assign"
	if !s.checkRateLimitPerEmail(key) {
		return ErrTooManyRequests
	}

	if err := ctx.Err(); err != nil {
        return ErrContextCancelled
    }

	// Проверяем валидность роли для слотов (только Guest занимает слот)
    needsSlot := role == domain.RoleGuest

    // Параллельные проверки: права + существование пользователя + событие
    type checkResult struct {
        hasRight bool
        user     *domain.User
        event    *domain.Event
        err      error
    }

    resultChan := make(chan checkResult, 1)

    go func() {
        var res checkResult
        
        // Проверяем права
        hasRight, err := s.participantRepo.HasAnyRole(ctx, requesterID, eventID, domain.RoleOwner, domain.RoleOrganizer)
        if err != nil || !hasRight {
            res.err = err
            if !hasRight && err == nil {
                res.err = ErrNotOwner
            }
            resultChan <- res
            return
        }

        // Проверяем пользователя
        user, err := s.userRepo.GetByID(ctx, targetUserID)
        if err != nil {
            res.err = err
            resultChan <- res
            return
        }
        if user == nil {
            res.err = ErrUserNotFound
            resultChan <- res
            return
        }
        res.user = user

        // Проверяем событие и статус
        event, err := s.eventRepo.GetByID(ctx, eventID)
        if err != nil {
            res.err = err
            resultChan <- res
            return
        }
        if event == nil {
            res.err = ErrEventNotFound
            resultChan <- res
            return
        }
        if event.Status == domain.StatusCompleted || event.Status == domain.StatusCancelled {
            res.err = ErrCantBeAssignToEvent
            resultChan <- res
            return
        }
        res.event = event

        // Проверяем, не участник ли уже
        existingRole, err := s.participantRepo.GetUserRoleInEvent(ctx, targetUserID, eventID)
        if err == nil && existingRole != "" {
            res.err = ErrAlreadyParticipant
            resultChan <- res
            return
        }

        res.hasRight = true
        resultChan <- res
    }()

    check := <-resultChan
    if check.err != nil {
        return check.err
    }

    // Валидация профессиональной роли
    if profRoleID != nil && role != domain.RoleStaff {
        return ErrInvalidRole
    }

    // === КЛЮЧЕВОЙ МОМЕНТ: атомарное занятие слота ===
    // Если нужен слот — занимаем ДО создания участника
    if needsSlot {
        ok, err := s.eventRepo.TakeSlot(ctx, eventID)
        if err != nil {
            return err
        }
        if !ok {
            return ErrNoSlotsAvailable
        }
    }

    // Создаём участника
    participant := &domain.EventParticipant{
        UserID:             targetUserID,
        EventID:            eventID,
        SystemRole:         role,
        ProfessionalRoleID: profRoleID,
        Notes:              "Added by " + requesterID.String(),
    }

    if err := s.participantRepo.Create(ctx, participant); err != nil {
        // ОТКАТ: если создать не удалось, но слот заняли — освобождаем
        if needsSlot {
            // Асинхронно освобождаем, не блокируем ответ
            go s.eventRepo.FreeUpSlot(context.Background(), eventID)
        }
        if isDuplicateError(err) {
            return ErrAlreadyParticipant
        }
        return err
    }

    return nil
}

func (s *participantService) ChangeRole(
	ctx context.Context,
	requesterID uuid.UUID,
	participantID uint,
	newRole domain.SystemRole,
) error {
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}

	participant, err := s.participantRepo.GetByID(ctx, participantID)
	if err != nil {
		return err
	}

	if participant.UserID == requesterID {
		return ErrSelfRoleChange
	}

	requesterRole, err := s.participantRepo.GetUserRoleInEvent(ctx, requesterID, participant.EventID)
	if err != nil {
		return err
	}

	if !s.canChangeRole(requesterRole, participant.SystemRole, newRole) {
		return ErrNotOwner
	}

	if participant.SystemRole == domain.RoleOwner && newRole != domain.RoleOwner {
		isLast, err := s.isLastOwner(ctx, participant.EventID, participant.UserID)
		if err != nil {
			return err
		}
		if isLast {
			return ErrLastOwnerCannotLeave
		}
	}

	return s.participantRepo.ChangeRole(ctx, participantID, newRole)
}

func (s *participantService) RemoveFromEvent(
	ctx context.Context,
	requesterID,
	targetUserID,
	eventID uuid.UUID,
) error {
	if err := ctx.Err(); err != nil {
		return ErrContextCancelled
	}

	if requesterID == targetUserID {
		return s.SelfRemove(ctx, requesterID, eventID)
	}

	return s.removeOther(ctx, requesterID, targetUserID, eventID)
}

func (s *participantService) SelfRemove(ctx context.Context, userID, eventID uuid.UUID) error {
	role, err := s.participantRepo.GetUserRoleInEvent(ctx, userID, eventID)
	if err != nil {
		return err
	}

	if err := s.participantRepo.RemoveFromEvent(ctx, userID, eventID); err != nil {
        return err
    }
	
	errChan := make(chan error, 1)
    // Освобождаем слот только если был Guest
    if role == domain.RoleGuest {
        // Не блокируем ответ, логируем ошибку если что
        go func() {
            if _, err := s.eventRepo.FreeUpSlot(context.Background(), eventID); err != nil {
				errChan <- ErrToReleaseSlot
            }else {
				errChan <- nil
			}
        }()

		if err := <-errChan; err != nil {
			return err
		}
    }

	return nil
}

func (s *participantService) removeOther(ctx context.Context, requesterID, targetUserID, eventID uuid.UUID) error {
	requesterRole, err := s.participantRepo.GetUserRoleInEvent(ctx, requesterID, eventID)
	if err != nil {
		return err
	}

	targetRole, err := s.participantRepo.GetUserRoleInEvent(ctx, targetUserID, eventID)
	if err != nil {
		return err
	}

	canRemove := false
	switch requesterRole {
	case domain.RoleOwner:
		canRemove = true
	case domain.RoleOrganizer:
		canRemove = targetRole == domain.RoleGuest || targetRole == domain.RoleStaff
	case domain.RoleManager:
		canRemove = targetRole == domain.RoleStaff
	}

	if !canRemove {
		return ErrNotOwner
	}

	if err := s.participantRepo.RemoveFromEvent(ctx, targetUserID, eventID); err != nil {
        return err
    }
	
	errChan := make(chan error, 1)
	if targetRole == domain.RoleGuest {
		go func() {
			if _, err := s.eventRepo.FreeUpSlot(context.Background(), eventID); err != nil {
				errChan <- ErrToReleaseSlot
			}else {
				errChan <- nil
			}
		}()

		if err := <-errChan; err != nil {
			return err
		}
	}

	return nil
}

// ==================== BATCH ОПЕРАЦИИ ====================
//тут мы можем приписать пачкой кого то к ивенту, например если у нас идут Staff группой (группа официантов или же кого то еще)
func (s *participantService) AssignRolesBulk(
	ctx context.Context,
	requesterID,
	eventID uuid.UUID,
	assignments []RoleAssignment,
) []BulkResult {
	results := make([]BulkResult, len(assignments))
	var wg sync.WaitGroup
	
	semaphore := make(chan struct{}, 5)

	for i, assign := range assignments {
		wg.Add(1)
		go func(index int, a RoleAssignment) {
			defer wg.Done()
			
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := s.AssignRole(ctx, requesterID, a.UserID, eventID, a.Role, a.ProfessionalID)
			results[index] = BulkResult{
				UserID: a.UserID,
				Error:  err,
			}
		}(i, assign)
	}

	wg.Wait()
	return results
}

func (s *participantService) canChangeRole(requesterRole, currentRole, newRole domain.SystemRole) bool {
	if requesterRole == domain.RoleOwner {
		return true
	}
	
	if requesterRole == domain.RoleOrganizer {
		return (currentRole == domain.RoleGuest || currentRole == domain.RoleStaff) &&
			   (newRole == domain.RoleGuest || newRole == domain.RoleStaff || newRole == domain.RoleManager)
	}
	
	return false
}

func (s *participantService) isLastOwner(ctx context.Context, eventID, excludeUserID uuid.UUID) (bool, error) {
	participants, err := s.eventRepo.GetEventParticipants(ctx, eventID)
	if err != nil {
		return false, err
	}

	ownerCount := 0
	for _, p := range participants {
		if p.SystemRole == domain.RoleOwner {
			ownerCount++
		}
	}

	isExcludeOwner := false
	for _, p := range participants {
		if p.UserID == excludeUserID && p.SystemRole == domain.RoleOwner {
			isExcludeOwner = true
			break
		}
	}

	if isExcludeOwner {
		return ownerCount <= 1, nil
	}
	return false, nil
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "duplicate") || contains(errStr, "unique constraint") || contains(errStr, "23505")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && len(substr) > 0 && findSubstr(s, substr))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (s *participantService) starterWorkerPool(workers int) {
	StarterWorkerPool(workers, s.taskQueue)
}

func (s *participantService) checkRateLimitPerEmail(email string) bool {
	return CheckRateLimitPerEmail(email, &s.rateLimiter, s.rateMap)
}