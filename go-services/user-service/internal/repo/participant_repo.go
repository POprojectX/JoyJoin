package repo

import (
	"user-service/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type ParticipantRepository interface {
	// CRUD
	Create(ctx context.Context, participant *domain.EventParticipant) error
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