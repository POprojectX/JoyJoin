package repo

import (
	"context"
	"errors"
	"time"
	"user-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByDateRange(ctx context.Context, start, end time.Time) ([]domain.Event, error)
	// Получить всех участников события с их ролями
	GetEventParticipants(ctx context.Context, eventID uuid.UUID) ([]domain.EventParticipant, error)
	// Получить события по системной роли пользователя
	GetEventsByUserRole(ctx context.Context, userID uuid.UUID, role domain.SystemRole) ([]domain.Event, error)
	GetEventsByLocation(ctx context.Context, location string) ([]domain.Event, error)
	HasAvailableSlots(ctx context.Context, eventID uuid.UUID) (bool, error)
	TakeSlot(ctx context.Context, eventID uuid.UUID) (bool, error)
	FreeUpSlot(ctx context.Context, eventID uuid.UUID) (bool, error)
}

type eventRepo struct {
	db *gorm.DB
	workerPool chan struct{}
}

func NewEventRepository(db *gorm.DB, maxWorkers int) EventRepository {
	return &eventRepo{
		db: db,
		workerPool: make(chan struct{}, maxWorkers),
	}
}

func (r *eventRepo) Create(ctx context.Context, event *domain.Event) error {
	select {
	case <- ctx.Done():
		return ctx.Err()
	default:
	}

	return r.db.Create(event).Error
}

func (r *eventRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	var event domain.Event
	// Preload загружает связанные данные участников ивента
	err := r.db.WithContext(ctx).Preload("Participants.User").Preload("Participants.ProfessionalRole").First(&event, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &event, err
}

func (r *eventRepo) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.Event{}).Where("id = ?", id).Updates(updates).Error
}

func (r *eventRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Event{}, "id = ?", id).Error
}

func (r *eventRepo) ListByDateRange(ctx context.Context, start, end time.Time) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.WithContext(ctx).Where("date_from BETWEEN ? AND ?", start, end).Find(&events).Error
	return events, err
}

// GetEventParticipants - мейн метод для управления мероприятием
func (r *eventRepo) GetEventParticipants(ctx context.Context, eventID uuid.UUID) ([]domain.EventParticipant, error) {
	var participants []domain.EventParticipant
	err := r.db.
		Preload("User").              // данные пользователя
		Preload("ProfessionalRole").  // специализация
		Where("event_id = ?", eventID).
		Find(&participants).Error
	return participants, err
}

// GetEventsByUserRole - например, все события, где я Owner
func (r *eventRepo) GetEventsByUserRole(ctx context.Context, userID uuid.UUID, role domain.SystemRole) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.
		WithContext(ctx).
		Joins("JOIN event_participants ON event_participants.event_id = events.id").
		Where("event_participants.user_id = ? AND event_participants.system_role = ?", userID, role).
		Find(&events).Error
	return events, err
}

func (r *eventRepo) GetEventsByLocation(ctx context.Context, location string) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.WithContext(ctx).
		Where("location = ?", location).
		Find(&events).Error

	return events, err
}

func (r *eventRepo) HasAvailableSlots(ctx context.Context, eventID uuid.UUID) (bool, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).First(&event, "id = ?", eventID).Error
	if err != nil {
		return false, err
	}
	return event.AvailableSlots > 0, nil
}

func (r *eventRepo) TakeSlot(ctx context.Context, eventID uuid.UUID) (bool, error) {
	result := r.db.WithContext(ctx).
        Model(&domain.Event{}).
        Where("id = ? AND available_slots > 0 AND status IN ('Announced', 'Ongoing')", eventID).
        UpdateColumn("available_slots", gorm.Expr("available_slots - 1"))

    if result.Error != nil {
        return false, result.Error
    }

    return result.RowsAffected > 0, nil
}

func (r *eventRepo) FreeUpSlot(ctx context.Context, eventID uuid.UUID) (bool, error) {
	result := r.db.WithContext(ctx).
        Model(&domain.Event{}).
        Where("id = ? AND available_slots < slots AND status IN ('Announced', 'Ongoing')", eventID).
        UpdateColumn("available_slots", gorm.Expr("available_slots + 1"))

    if result.Error != nil {
        return false, result.Error
    }

    return result.RowsAffected > 0, nil
}