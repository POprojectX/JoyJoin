package services

import (
	"context"
	"sync"
	"time"
	"user-service/internal/domain"
	customErrors "user-service/internal/errors"
	"user-service/internal/repo"

	"github.com/google/uuid"
)

type EventService interface {
	Create(ctx context.Context, title, description, location string, slots int, dateFrom, dateTo time.Time, ownerID uuid.UUID, access domain.Access) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, input domain.UpdateEventInput, id uuid.UUID) (*domain.Event, error)
	Delete(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	ListByDateRange(ctx context.Context, start, end time.Time) ([]domain.Event, error)
	PublishEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	DraftEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	CancelledEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	// Получить всех участников события с их ролями
	GetEventParticipants(ctx context.Context, eventID uuid.UUID) ([]domain.EventParticipant, error)
	// Получить события по системной роли пользователя
	GetEventsByUserRole(ctx context.Context, userID uuid.UUID, role domain.SystemRole) ([]domain.Event, error)

	SetSlots(ctx context.Context, eventID uuid.UUID, slots int) (*domain.Event, error)
    OccupySlot(ctx context.Context, eventID uuid.UUID) error
    ReleaseSlot(ctx context.Context, eventID uuid.UUID) error
    HasAvailableSlots(ctx context.Context, eventID uuid.UUID) (bool, error)
}

type eventService struct {
	eventRepo repo.EventRepository
	rateLimiter sync.RWMutex
	rateMap map[string]time.Time
	taskQueue chan func()
	workerPool chan struct{}
}

type slotsOperation struct {
    eventID uuid.UUID
    delta   int
    done    chan error
}

func NewEventService(r repo.EventRepository) EventService {
	s := &eventService{
		eventRepo:  r,
		rateMap:   make(map[string]time.Time),
		taskQueue: make(chan func(), 1000),
		workerPool: make(chan struct{}, 1),
	}
	
	// Запускаем пул воркеров для фоновых задач
	s.starterWorkerPool(10)
	
	return s
}

func (s *eventService) Create(ctx context.Context, title, description, location string, slots int, dateFrom, dateTo time.Time, ownerID uuid.UUID, access domain.Access) (*domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, customErrors.ErrContextCancelled
	}
	if title == "" {
		return nil, customErrors.ErrNoEventTitle
	}
	
	event := &domain.Event{
		Status: domain.StatusDraft,
		Access: access,
		Title: title,
		Description: description,
		Slots: slots,
		AvailableSlots: slots,
		DateFrom: dateFrom,
		DateTo: dateTo,
		Location: location,
		OwnerID: ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := s.eventRepo.Create(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, customErrors.ErrEventNotFound
	}

	return event, nil
}

func (s *eventService) Update(ctx context.Context, input domain.UpdateEventInput, id uuid.UUID) (*domain.Event, error) {
	currentEvent, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if currentEvent == nil {
		return nil, customErrors.ErrEventNotFound
	}
	if currentEvent.Status == domain.StatusOngoing || currentEvent.Status == domain.StatusCompleted || currentEvent.Status == domain.StatusCancelled {
		return nil, customErrors.ErrCantUpdateEvent
	}
	updates := make(map[string]interface{})
	
	// Чисто, читаемо, без дублирования
	CollectUpdates(updates, input.Access != nil, "access", input.Access)
	CollectUpdates(updates, input.Status != nil, "status", input.Status)
	CollectUpdates(updates, input.Title != nil, "title", input.Title)
	CollectUpdates(updates, input.Description != nil, "description", input.Description)
	CollectUpdates(updates, input.Location != nil, "location", input.Location)
	CollectUpdates(updates, input.DateFrom != nil, "date_from", input.DateFrom)
	CollectUpdates(updates, input.DateTo != nil, "date_to", input.DateTo)
	CollectUpdates(updates, input.Slots != nil, "slots", input.Slots)

	if len(updates) == 0 {
		return s.eventRepo.GetByID(ctx, id)
	}

	updates["updated_at"] = time.Now()

	if err := s.eventRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}

	return s.eventRepo.GetByID(ctx, id)
}

func (s *eventService) PublishEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, customErrors.ErrEventNotFound
	}
	if event.Status != domain.StatusDraft {
		return nil, customErrors.ErrCantPublishEvent
	}
	event.Status = domain.StatusAnnounced
	input := domain.UpdateEventInput{
		Status: &event.Status,
		Access: &event.Access,
		Title: &event.Title,
		Description: &event.Description,
		Location: &event.Location,
		DateFrom: &event.DateFrom,
		DateTo: &event.DateTo,
	}
	return s.Update(ctx, input, id)
}

func (s *eventService) DraftEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, customErrors.ErrEventNotFound
	}
	if event.Status != domain.StatusAnnounced {
		return nil, customErrors.ErrCantDraftEvent
	}
	event.Status = domain.StatusDraft
	input := domain.UpdateEventInput{
		Status: &event.Status,
		Access: &event.Access,
		Title: &event.Title,
		Description: &event.Description,
		Location: &event.Location,
		DateFrom: &event.DateFrom,
		DateTo: &event.DateTo,
	}
	return s.Update(ctx, input, id)
}

func (s *eventService) CancelledEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, customErrors.ErrEventNotFound
	}
	if event.Status == domain.StatusDraft || event.Status == domain.StatusAnnounced {
		event.Status = domain.StatusCancelled
	} else {
		return nil, customErrors.ErrCantCancelEvent
	}
	
	input := domain.UpdateEventInput{
		Access: &event.Access,
		Status: &event.Status,
		Title: &event.Title,
		Description: &event.Description,
		Location: &event.Location,
		DateFrom: &event.DateFrom,
		DateTo: &event.DateTo,
	}
	return s.Update(ctx, input, id)
}

func(s *eventService) Delete(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, customErrors.ErrEventNotFound
	}
	if event.Status == domain.StatusAnnounced || event.Status == domain.StatusDraft {
		s.eventRepo.Delete(ctx, id)
	}else {
		return nil, customErrors.ErrDeleteLiveEvent
	}

	return nil, nil
}

func(s *eventService) ListByDateRange(ctx context.Context, start, end time.Time) ([]domain.Event, error) {
	events, err := s.eventRepo.ListByDateRange(ctx, start, end)
	if err != nil {
		return nil, err
	}
	return events, nil
}


func (s *eventService) GetEventParticipants(ctx context.Context, eventID uuid.UUID) ([]domain.EventParticipant, error) {
	participants, err := s.eventRepo.GetEventParticipants(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if participants == nil || len(participants) == 0{
		return nil, customErrors.ErrEventNoParticipants
	}
	return participants, nil
}

func (s *eventService) GetEventsByUserRole(ctx context.Context, userID uuid.UUID, role domain.SystemRole) ([]domain.Event, error) {
	events, err := s.eventRepo.GetEventsByUserRole(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	if events == nil && len(events) == 0 {
		return nil, customErrors.ErrNoEventsWithRole
	}

	return events, nil
}

func (s *eventService) SetSlots(ctx context.Context, eventID uuid.UUID, slots int) (*domain.Event, error) {
    if err := ctx.Err(); err != nil {
        return nil, customErrors.ErrContextCancelled
    }
    if slots <= 0 {
        return nil, customErrors.ErrInvalidSlots
    }

    event, err := s.eventRepo.GetByID(ctx, eventID)
    if err != nil {
        return nil, err
    }
    if event == nil {
        return nil, customErrors.ErrEventNotFound
    }

    if event.Status != domain.StatusDraft && event.Status != domain.StatusAnnounced {
        return nil, customErrors.ErrCantUpdateEvent
    }

    participants, err := s.eventRepo.GetEventParticipants(ctx, eventID)
    if err != nil {
        return nil, err
    }
    
    guestCount := 0
	for _, p := range participants {
		if p.SystemRole == domain.RoleGuest {
			guestCount++
		}
	}

    updates := map[string]interface{}{
        "slots":      slots,
		"available_slots": slots - guestCount,
        "updated_at": time.Now(),
    }

    if err := s.eventRepo.Update(ctx, eventID, updates); err != nil {
        return nil, err
    }

    return s.eventRepo.GetByID(ctx, eventID)
}

// OccupySlot - атомарно занимаем слот при регистрации
func (s *eventService) OccupySlot(ctx context.Context, eventID uuid.UUID) error {
	ok, err := s.eventRepo.TakeSlot(ctx, eventID)
	if err != nil {
		return err
	}

	if !ok {
		return customErrors.ErrNoSlotsAvailable
	}

	return nil
}

// ReleaseSlot - освобождаем слот
func (s *eventService) ReleaseSlot(ctx context.Context, eventID uuid.UUID) error {
    result, err :=s.eventRepo.FreeUpSlot(ctx, eventID)
	if err != nil {
		return err
	}
	if !result {
		return customErrors.ErrNoSlotsToRelease 
	}
    return nil
}

func (s *eventService) HasAvailableSlots(ctx context.Context, eventID uuid.UUID) (bool, error) {
    result, err := s.eventRepo.HasAvailableSlots(ctx, eventID)
	if err != nil {
		return false, err
	}
	if !result {
		return false, customErrors.ErrNoSlotsAvailable
	}
	return true, nil
}

//используем утилиты
func (s *eventService) starterWorkerPool(workers int) {
	StarterWorkerPool(workers, s.taskQueue)
}

func (s *eventService) checkRateLimitPerEmail(email string) bool {
	return CheckRateLimitPerEmail(email, &s.rateLimiter, s.rateMap)
}