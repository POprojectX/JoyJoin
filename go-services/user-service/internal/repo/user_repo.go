package repo

import (
	"user-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository - интерфейс для работы с пользователями
type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uuid.UUID) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]domain.User, error)
	// Получить все события, где пользователь участвует
	GetUserEvents(userID uuid.UUID) ([]domain.EventParticipant, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// Create - создание пользователя
func (r *userRepo) Create(user *domain.User) error {
	// BeforeCreate хук в GORM автоматически генерирует UUID
	return r.db.Create(user).Error
}

// GetByID - получение по UUID
func (r *userRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err // вернёт gorm.ErrRecordNotFound если не найден наш пользователь
	}
	return &user, nil
}

// GetByEmail - для логина/проверки уникальности
func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

// Update - обновление (только переданные поля)
func (r *userRepo) Update(user *domain.User) error {
	// Updates обновляет только ненулевые поля
	return r.db.Model(user).Updates(user).Error
}

// Delete - soft delete или hard delete (зависит от GORM config)
func (r *userRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

// List - пагинация списка пользователей
func (r *userRepo) List(limit, offset int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// GetUserEvents - все события пользователя с ролями (ключевой метод!)
func (r *userRepo) GetUserEvents(userID uuid.UUID) ([]domain.EventParticipant, error) {
	var participations []domain.EventParticipant
	err := r.db.
		Preload("Event").        // подгружаем данные события
		Preload("ProfessionalRole"). // подгружаем специализацию
		Where("user_id = ?", userID).
		Find(&participations).Error
	return participations, err
}