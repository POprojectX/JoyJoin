package services

import (
	"context"
	"sync"
	"time"
	"user-service/internal/domain"
	"user-service/internal/repo"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repo.UserRepository
	secretKey string
}

func NewAuthService(r repo.UserRepository, secret string) *AuthService {
	return &AuthService{
		userRepo: r,
		secretKey: secret,
	}
}

func (s *AuthService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp" : time.Now().Add(time.Hour * 72).Unix(), //токен выдаем на 3 дня
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}


//user-service
type UserService interface {
	// Auth
	Register(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	
	// CRUD
	Update(ctx context.Context, input domain.UpdateUserInput, id uuid.UUID) (*domain.User, error)
	Delete(ctx context.Context, userID uuid.UUID) (*domain.User, error)

	// Business logic
	GetUserWithEvents(ctx context.Context, userID uuid.UUID) (*domain.UserWithEventsDTO, error)
}

type RegisterRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type userService struct {
	userRepo repo.UserRepository
	rateLimiter sync.RWMutex
	rateMap map[string]time.Time
	taskQueue chan func()
}

func NewUserService(r repo.UserRepository) UserService {
	s := &userService{
		userRepo:  r,
		rateMap:   make(map[string]time.Time),
		taskQueue: make(chan func(), 1000), 
	}
	
	// Запускаем пул воркеров для фоновых задач
	s.starterWorkerPool(10)
	
	return s
}


func (s *userService) Register(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error) {
	//тут куча if-ов так как я хз как по другому проверить на фулл ошибки, слой гавнокодика :)
	if !s.checkRateLimitPerEmail(email) {
		return nil, ErrTooManyRequests
	}
	if err := ctx.Err(); err != nil {
		return nil, ErrContextCancelled
	}
	if email == "" || password == "" {
		return nil, ErrValidationFailed
	}
	if len(password) < 6 {
		return nil, ErrShortPassword
	}
	exist, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, ErrEmailAlreadyExists
	}


	passwordChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func() {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			errChan <- err
			return 
		}
		passwordChan <- hash
	}()

	select {
	case <- ctx.Done():
		return nil, ErrContextCancelled
	case err := <- errChan: 
		return nil, err
	case hashedPassword := <- passwordChan:
		user := &domain.User{
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			Password:  hashedPassword,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}
}


func (s *userService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	if !s.checkRateLimitPerEmail(email) {
		return nil, ErrTooManyRequests
	}
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return  nil, ErrUserNotFound
	}
	type result struct {
		valid bool
		err error
	}
	resultChan := make(chan result, 1)

	go func () {
		err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
		resultChan <- result{valid: err == nil, err: err}
	}()

	select{
	case <- ctx.Done():
		return nil, ErrContextCancelled
	case res := <- resultChan:
		if res.err != nil {
			return nil, ErrInvalidPassword
		}
		if !res.valid {
			return nil, ErrInvalidPassword
		}

		user.Password = nil
		return user, nil
	}
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.Password = nil //очищаем пароль что бы чел через запрос id челика не узнал его парольчик, просто пароль сносим к хуям
	return user, nil
}

func (s *userService) Update(ctx context.Context, input domain.UpdateUserInput, id uuid.UUID) (*domain.User, error) {
	updates := make(map[string]interface{})

	CollectUpdates(updates, input.FirstName != nil, "firstName", input.FirstName)
	CollectUpdates(updates, input.LastName != nil, "firstName", input.LastName)

	if input.Password != nil {
		if len(*input.Password) < 6 {
			return nil, ErrShortPassword
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updates["password"] = hashed
	}

	if len(updates) == 0 {
		return s.userRepo.GetByID(ctx, id)
	}

	updates["updated_at"] = time.Now()
	
	if err := s.userRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}

	return s.userRepo.GetByID(ctx, id)	
}

func (s *userService) Delete(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	s.userRepo.Delete(ctx, user)
	return nil, nil
} 

func (s *userService) GetUserWithEvents(ctx context.Context, userID uuid.UUID) (*domain.UserWithEventsDTO, error) {
	var user *domain.User
	var participations []domain.EventParticipant
	var errUser, errEvents error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		user, errUser = s.userRepo.GetByID(ctx, userID)
	}()
	go func() {
		defer wg.Done()
		participations, errEvents = s.userRepo.GetUserEvents(ctx, userID)
	}()
	wg.Wait()

	if errUser != nil {
		return nil, errUser
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if errEvents != nil {
		return nil, errEvents
	}
	if participations == nil {
		return nil, ErrNoEvents
	}

	eventsDTO := make([]domain.EventSummaryDTO, len(participations))
	for i, p := range participations {
		eventsDTO[i] = domain.EventSummaryDTO{
			EventID:    p.EventID,
			Title:      p.Event.Title,
			DateFrom:       p.Event.DateFrom,
			DateTo:       p.Event.DateTo,
			SystemRole: p.SystemRole,
			JoinedAt:   p.JoinedAt,
		}
	}

	return &domain.UserWithEventsDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Events:    eventsDTO,
	}, nil
}

//используем утилиты
func (s *userService) starterWorkerPool(workers int) {
	StarterWorkerPool(workers, s.taskQueue)
}

func (s *userService) checkRateLimitPerEmail(email string) bool {
	return CheckRateLimitPerEmail(email, &s.rateLimiter, s.rateMap)
}