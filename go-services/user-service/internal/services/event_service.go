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
	Create(event *domain.Event) error
	GetByID(id uuid.UUID) (*domain.Event, error)
	Update(event *domain.Event) error
	Delete(id uuid.UUID) error
	ListByDateRange(start, end time.Time) ([]domain.Event, error)
	// Получить всех участников события с их ролями
	GetEventParticipants(eventID uuid.UUID) ([]domain.EventParticipant, error)
	// Получить события по системной роли пользователя
	GetEventsByUserRole(userID uuid.UUID, role domain.SystemRole) ([]domain.Event, error)
}

type eventService struct {
	eventRepo repo.EventRepository
	rateLimiter sync.RWMutex
	rateMap map[string]time.Time
	taskQueue chan func()
}

func NewEventService(r repo.EventRepository) EventService {
	s := &eventService{
		eve:  r,
		rateMap:   make(map[string]time.Time),
		taskQueue: make(chan func(), 1000), 
	}
	
	// Запускаем пул воркеров для фоновых задач
	s.starterWorkerPool(10)
	
	return s
}

// evet-service
// Create, GetByID, Update, Delete, ListByDateRange, GetEventParticipants, GetEventsByUserRole
func (s *userService) Create(ctx context.Context, )