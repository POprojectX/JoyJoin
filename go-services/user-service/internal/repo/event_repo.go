package repo

import (
	"time"
	"user-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
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

type eventRepo struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) Create(event *domain.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepo) GetByID(id uuid.UUID) (*domain.Event, error) {
	var event domain.Event
	// Preload загружает связанные данные участников ивента
	err := r.db.Preload("Participants.User").Preload("Participants.ProfessionalRole").First(&event, "id = ?", id).Error
	return &event, err
}

func (r *eventRepo) Update(event *domain.Event) error {
	return r.db.Model(event).Updates(event).Error
}

func (r *eventRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Event{}, "id = ?", id).Error
}

func (r *eventRepo) ListByDateRange(start, end time.Time) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.Where("date BETWEEN ? AND ?", start, end).Find(&events).Error
	return events, err
}

// GetEventParticipants - мейн метод для управления мероприятием
func (r *eventRepo) GetEventParticipants(eventID uuid.UUID) ([]domain.EventParticipant, error) {
	var participants []domain.EventParticipant
	err := r.db.
		Preload("User").              // данные пользователя
		Preload("ProfessionalRole").  // специализация
		Where("event_id = ?", eventID).
		Find(&participants).Error
	return participants, err
}

// GetEventsByUserRole - например, все события, где я Owner
func (r *eventRepo) GetEventsByUserRole(userID uuid.UUID, role domain.SystemRole) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.
		Joins("JOIN event_participants ON event_participants.event_id = events.id").
		Where("event_participants.user_id = ? AND event_participants.system_role = ?", userID, role).
		Find(&events).Error
	return events, err
}