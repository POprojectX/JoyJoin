package repo

import (
	"user-service/internal/domain"
	customErrors "user-service/internal/errors"

	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ParticipantRepository interface {
	// CRUD
	Create(ctx context.Context, participant *domain.EventParticipant) error
	CreateWithSlotAtomic(ctx context.Context, participant *domain.EventParticipant, eventID uuid.UUID) error
	GetByID(ctx context.Context, id uint) (*domain.EventParticipant, error)
	Update(ctx context.Context, participant *domain.EventParticipant) error
	Delete(ctx context.Context, id uint) error
	
	// Бизнес логика ролей
	AssignRole(ctx context.Context, userID, eventID uuid.UUID, role domain.SystemRole, profRoleID *uint) error
	ChangeRole(ctx context.Context, participantID uint, newRole domain.SystemRole) error
	RemoveFromEvent(ctx context.Context, userID, eventID uuid.UUID) error
	
	// Проверки прав
	GetUserRoleInEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.SystemRole, error)
	IsUserOwner(ctx context.Context, userID, eventID uuid.UUID) (bool, error)
	HasAnyRole(ctx context.Context, userID, eventID uuid.UUID, roles ...domain.SystemRole) (bool, error)

	GetByEventID(ctx context.Context, eventID uuid.UUID) ([]domain.EventParticipant, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.EventParticipant, error)
	CountOwners(ctx context.Context, eventID uuid.UUID) (int64, error)
}

type participantRepo struct {
	db *gorm.DB
	workerPool chan struct{}
}

func NewParticipantRepository(db *gorm.DB) ParticipantRepository {
	return &participantRepo{db: db}
}

func (r *participantRepo) Create(ctx context.Context, participant *domain.EventParticipant) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

func (r *participantRepo) CreateWithSlotAtomic(ctx context.Context, participant *domain.EventParticipant, eventID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Блокируем строку события (SELECT FOR UPDATE)
		var event domain.Event
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&event, "id = ?", eventID).Error; err != nil {
			return err
		}
		
		// 2. Проверяем что можно назначить (статус, слоты)
		if event.Status != domain.StatusAnnounced && event.Status != domain.StatusOngoing {
			return customErrors.ErrCantBeAssignToEvent
		}
		
		// 3. Только Guest занимает слот
		if participant.SystemRole == domain.RoleGuest {
			if event.AvailableSlots <= 0 {
				return customErrors.ErrNoSlotsAvailable
			}
			result := tx.Model(&domain.Event{}).Where("id = ? AND available_slots > 0 AND status IN ('Announced', 'Ongoing')", eventID).
			UpdateColumn("available_slots", gorm.Expr("available_slots - 1"))

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return customErrors.ErrNoSlotsAvailable
			}
		}
		
		// 4. Создаём участника
		if err := tx.Create(participant).Error; err != nil {
            return err
        }
		return nil
	})
}

func (r *participantRepo) GetByID(ctx context.Context, id uint) (*domain.EventParticipant, error) {
	var p domain.EventParticipant
	err := r.db.WithContext(ctx).Preload("User").Preload("Event").Preload("ProfessionalRole").First(&p, id).Error
	return &p, err
}

func (r *participantRepo) Update(ctx context.Context, participant *domain.EventParticipant) error {
	return r.db.WithContext(ctx).Model(participant).Updates(participant).Error
}

func (r *participantRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.EventParticipant{}, id).Error
}

// AssignRole - назначить роль пользователю в событии
func (r *participantRepo) AssignRole(ctx context.Context, userID, eventID uuid.UUID, role domain.SystemRole, profRoleID *uint) error {
	participant := &domain.EventParticipant{
		UserID:             userID,
		EventID:            eventID,
		SystemRole:         role,
		ProfessionalRoleID: profRoleID,
	}
	return r.db.WithContext(ctx).Create(participant).Error
}

// ChangeRole - смена роли (например, Guest → Organizer)
func (r *participantRepo) ChangeRole(ctx context.Context, participantID uint, newRole domain.SystemRole) error {
	return r.db.WithContext(ctx).Model(&domain.EventParticipant{}).
		Where("id = ?", participantID).
		Update("system_role", newRole).Error
}

// RemoveFromEvent - удалить участника из события
func (r *participantRepo) RemoveFromEvent(ctx context.Context, userID, eventID uuid.UUID) error {
	return r.db.
		Where("user_id = ? AND event_id = ?", userID, eventID).
		Delete(&domain.EventParticipant{}).Error
}

// GetUserRoleInEvent - какая роль у пользователя в конкретном событии
func (r *participantRepo) GetUserRoleInEvent(ctx context.Context, userID, eventID uuid.UUID) (domain.SystemRole, error) {
	var participant domain.EventParticipant
	err := r.db.
		WithContext(ctx).
		Select("system_role").
		Where("user_id = ? AND event_id = ?", userID, eventID).
		First(&participant).Error
	
	if err != nil {
		return "", err
	}
	return participant.SystemRole, nil
}

// IsUserOwner - проверка для защищённых операций
func (r *participantRepo) IsUserOwner(ctx context.Context, userID, eventID uuid.UUID) (bool, error) {
	role, err := r.GetUserRoleInEvent(ctx, userID, eventID)
	if err != nil {
		return false, err
	}
	return role == domain.RoleOwner, nil
}

// HasAnyRole - проверка нескольких ролей (например, Owner или Organizer)
func (r *participantRepo) HasAnyRole(ctx context.Context, userID, eventID uuid.UUID, roles ...domain.SystemRole) (bool, error) {
	var count int64
	err := r.db.
		WithContext(ctx).
		Model(&domain.EventParticipant{}).
		Where("user_id = ? AND event_id = ? AND system_role IN ?", userID, eventID, roles).
		Count(&count).Error
	return count > 0, err
}
func (r *participantRepo) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]domain.EventParticipant, error) {
	var participants []domain.EventParticipant
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("ProfessionalRole").
		Where("event_id = ?", eventID).
		Find(&participants).Error
	return participants, err
}

func (r *participantRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.EventParticipant, error) {
	var participants []domain.EventParticipant
	err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("ProfessionalRole").
		Where("user_id = ?", userID).
		Find(&participants).Error
	return participants, err
}

func (r *participantRepo) CountOwners(ctx context.Context, eventID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.EventParticipant{}).
		Where("event_id = ? AND system_role = ?", eventID, domain.RoleOwner).
		Count(&count).Error
	return count, err
}