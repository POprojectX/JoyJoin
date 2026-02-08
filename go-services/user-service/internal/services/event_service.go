package services

import (
	"context"
	"sync"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"

	"github.com/google/uuid"
)

type EventService interface {
	Create(ctx context.Context, title, description, location string, dateFrom, dateTo time.Time, ownerID uuid.UUID) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, input domain.UpdateEventInput, id uuid.UUID) (*domain.Event, error)
	Delete(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	ListByDateRange(ctx context.Context, start, end time.Time) ([]domain.Event, error)
	// Получить всех участников события с их ролями
	GetEventParticipants(ctx context.Context, eventID uuid.UUID) ([]domain.EventParticipant, error)
	// Получить события по системной роли пользователя
	GetEventsByUserRole(ctx context.Context, userID uuid.UUID, role domain.SystemRole) ([]domain.Event, error)
}

type eventService struct {
	eventRepo repo.EventRepository
	rateLimiter sync.RWMutex
	rateMap map[string]time.Time
	taskQueue chan func()
}

func NewEventService(r repo.EventRepository) EventService {
	s := &eventService{
		eventRepo:  r,
		rateMap:   make(map[string]time.Time),
		taskQueue: make(chan func(), 1000), 
	}
	
	// Запускаем пул воркеров для фоновых задач
	s.starterWorkerPool(10)
	
	return s
}

func (s *eventService) Create(ctx context.Context, title, description, location string, dateFrom, dateTo time.Time, ownerID uuid.UUID) (*domain.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, ErrContextCancelled
	}
	if title == "" {
		return nil, ErrNoEventTitle
	}
	
	event := &domain.Event{
		Status: domain.StatusDraft,
		Title: title,
		Description: description,
		DateFrom: dateFrom,
		DateTo: dateTo,
		Location: location,
		OwnerID: ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.eventRepo.Create(ctx, event)
	return event, nil
}

func (s *eventService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrEventNotFound
	}

	return event, nil
}

func (s *eventService) Update(ctx context.Context, input domain.UpdateEventInput, id uuid.UUID) (*domain.Event, error) {
	currentEvent, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if currentEvent == nil {
		return nil, ErrEventNotFound
	}
	if currentEvent.Status == domain.StatusOngoing || currentEvent.Status == domain.StatusCompleted || currentEvent.Status == domain.StatusCancelled {
		return nil, ErrCantUpdateEvent
	}
	updates := make(map[string]interface{})
	
	// Чисто, читаемо, без дублирования
	CollectUpdates(updates, input.Status != nil, "status", input.Status)
	CollectUpdates(updates, input.Title != nil, "title", input.Title)
	CollectUpdates(updates, input.Description != nil, "description", input.Description)
	CollectUpdates(updates, input.Location != nil, "location", input.Location)
	CollectUpdates(updates, input.DateFrom != nil, "date_from", input.DateFrom)
	CollectUpdates(updates, input.DateTo != nil, "date_to", input.DateTo)

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
	if event.Status != domain.StatusDraft {
		return nil, ErrCantPublishEvent
	}
	event.Status = domain.StatusAnnounced
	input := domain.UpdateEventInput{
		Status: &event.Status,
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
	if event.Status != domain.StatusAnnounced {
		return nil, ErrCantDraftEvent
	}
	event.Status = domain.StatusDraft
	input := domain.UpdateEventInput{
		Status: &event.Status,
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
	if event.Status == domain.StatusDraft || event.Status == domain.StatusAnnounced {
		event.Status = domain.StatusCancelled
	} else {
		return nil, ErrCantCancelEvent
	}
	
	input := domain.UpdateEventInput{
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
		return nil, ErrEventNotFound
	}
	if event.Status == domain.StatusAnnounced || event.Status == domain.StatusDraft {
		s.eventRepo.Delete(ctx, id)
	}else {
		return nil, ErrDeleteLiveEvent
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
	if participants == nil && len(participants) == 0{
		return nil, ErrEventNoParticipants
	}
	return participants, nil
}

func (s *eventService) GetEventsByUserRole(ctx context.Context, userID uuid.UUID, role domain.SystemRole) ([]domain.Event, error) {
	events, err := s.eventRepo.GetEventsByUserRole(ctx, userID, role)
	if err != nil {
		return nil, err
	}
	if events == nil && len(events) == 0 {
		return nil, ErrNoEventsWithRole
	}

	return events, nil
}

//используем утилиты
func (s *eventService) starterWorkerPool(workers int) {
	StarterWorkerPool(workers, s.taskQueue)
}

func (s *eventService) checkRateLimitPerEmail(email string) bool {
	return CheckRateLimitPerEmail(email, &s.rateLimiter, s.rateMap)
}