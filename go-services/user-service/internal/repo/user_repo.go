package repo

import (
	"context"
	"errors"
	"sync"
	"user-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository - интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, user *domain.User) error
	ListByEmail(ctx context.Context, emails []string) (map[string]*domain.User, error)
	// Получить все события, где пользователь участвует
	GetUserEvents(ctx context.Context, userID uuid.UUID) ([]domain.EventParticipant, error)
}

type userRepo struct {
	db *gorm.DB
	workerPool chan struct{}
}

func NewUserRepository(db *gorm.DB, maxWorkers int) UserRepository {
	return &userRepo{
		db: db,
		workerPool: make(chan struct{}, maxWorkers),
	}
}

// Create - создание пользователя
func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	// BeforeCreate хук в GORM автоматически генерирует UUID
	return r.db.Create(user).Error
}

// GetByID - получение по UUID
func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, nil
}

// GetByEmail - для логина/проверки уникальности
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update - обновление (только переданные поля)
func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	// Updates обновляет только ненулевые поля
	return r.db.WithContext(ctx).Where("id = ?", user.ID).Updates(user).Error
}

// Delete - soft delete или hard delete (зависит от GORM config)
func (r *userRepo) Delete(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Where("id = ?", user.ID).Delete(user).Error
}

// List - пагинация списка пользователей
func (r *userRepo) ListByEmail(ctx context.Context, emails []string) (map[string]*domain.User, error) {
	result := make(map[string]*domain.User)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(emails))

	maxGoroutines := make(chan struct{}, 10)

	for _, email := range emails {
		wg.Add(1)
		maxGoroutines <- struct{}{}

		go func(e string) {
			defer wg.Done()
			defer func() {<-maxGoroutines}()

			user, err := r.GetByEmail(ctx, e)
			if err != nil {
				errChan <- err
			}
			if user != nil {
				mu.Lock()
				result[e] = user
				mu.Unlock()
			}
		}(email)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// GetUserEvents - все события пользователя с ролями (ключевой метод!)
func (r *userRepo) GetUserEvents(ctx context.Context, userID uuid.UUID) ([]domain.EventParticipant, error) {
	var participations []domain.EventParticipant
	err := r.db.
		Preload("Event").        // подгружаем данные события
		Preload("ProfessionalRole"). // подгружаем специализацию
		Where("user_id = ?", userID).
		Find(&participations).Error
	return participations, err
}