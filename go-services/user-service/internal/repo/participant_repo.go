package repo

import (
	"user-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantRepository interface {
	// CRUD
	Create(participant *domain.EventParticipant) error
	GetByID(id uint) (*domain.EventParticipant, error)
	Update(participant *domain.EventParticipant) error
	Delete(id uint) error
	
	// Бизнес логика ролей
	AssignRole(userID, eventID uuid.UUID, role domain.SystemRole, profRoleID *uint) error
	ChangeRole(participantID uint, newRole domain.SystemRole) error
	RemoveFromEvent(userID, eventID uuid.UUID) error
	
	// Проверки прав
	GetUserRoleInEvent(userID, eventID uuid.UUID) (domain.SystemRole, error)
	IsUserOwner(userID, eventID uuid.UUID) (bool, error)
	HasAnyRole(userID, eventID uuid.UUID, roles ...domain.SystemRole) (bool, error)
}

type participantRepo struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) ParticipantRepository {
	return &participantRepo{db: db}
}

func (r *participantRepo) Create(participant *domain.EventParticipant) error {
	return r.db.Create(participant).Error
}

func (r *participantRepo) GetByID(id uint) (*domain.EventParticipant, error) {
	var p domain.EventParticipant
	err := r.db.Preload("User").Preload("Event").Preload("ProfessionalRole").First(&p, id).Error
	return &p, err
}

func (r *participantRepo) Update(participant *domain.EventParticipant) error {
	return r.db.Model(participant).Updates(participant).Error
}

func (r *participantRepo) Delete(id uint) error {
	return r.db.Delete(&domain.EventParticipant{}, id).Error
}

// AssignRole - назначить роль пользователю в событии
func (r *participantRepo) AssignRole(userID, eventID uuid.UUID, role domain.SystemRole, profRoleID *uint) error {
	participant := &domain.EventParticipant{
		UserID:             userID,
		EventID:            eventID,
		SystemRole:         role,
		ProfessionalRoleID: profRoleID,
	}
	return r.db.Create(participant).Error
}

// ChangeRole - смена роли (например, Guest → Organizer)
func (r *participantRepo) ChangeRole(participantID uint, newRole domain.SystemRole) error {
	return r.db.Model(&domain.EventParticipant{}).
		Where("id = ?", participantID).
		Update("system_role", newRole).Error
}

// RemoveFromEvent - удалить участника из события
func (r *participantRepo) RemoveFromEvent(userID, eventID uuid.UUID) error {
	return r.db.
		Where("user_id = ? AND event_id = ?", userID, eventID).
		Delete(&domain.EventParticipant{}).Error
}

// GetUserRoleInEvent - какая роль у пользователя в конкретном событии
func (r *participantRepo) GetUserRoleInEvent(userID, eventID uuid.UUID) (domain.SystemRole, error) {
	var participant domain.EventParticipant
	err := r.db.
		Select("system_role").
		Where("user_id = ? AND event_id = ?", userID, eventID).
		First(&participant).Error
	
	if err != nil {
		return "", err
	}
	return participant.SystemRole, nil
}

// IsUserOwner - проверка для защищённых операций
func (r *participantRepo) IsUserOwner(userID, eventID uuid.UUID) (bool, error) {
	role, err := r.GetUserRoleInEvent(userID, eventID)
	if err != nil {
		return false, err
	}
	return role == domain.RoleOwner, nil
}

// HasAnyRole - проверка нескольких ролей (например, Owner или Organizer)
func (r *participantRepo) HasAnyRole(userID, eventID uuid.UUID, roles ...domain.SystemRole) (bool, error) {
	var count int64
	err := r.db.
		Model(&domain.EventParticipant{}).
		Where("user_id = ? AND event_id = ? AND system_role IN ?", userID, eventID, roles).
		Count(&count).Error
	return count > 0, err
}