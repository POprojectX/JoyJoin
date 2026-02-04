package services

import (
	"context"
	"errors"
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
var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists = errors.New("email already exist")
	ErrInvalidPassword = errors.New("invalid password")
	ErrValidationFailed = errors.New("validation failed")
	ErrShortPassword = errors.New("password is too short, must be more then 5 symbols!")
	ErrTooManyRequests = errors.New("too many requests! (1 req/sec)")
	ErrContextCancelled = errors.New("context was cancelled")
	ErrNoEvents = errors.New("user has no events")
)

type UserService interface {
	// Auth
	Register(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	
	// CRUD
	Update(ctx context.Context, userID uuid.UUID, firstName, lastName string) (*domain.User, error)
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

//тут мы запускаем наш воркер пул
func (s *userService) starterWorkerPool(workers int) {
	for i:= 0; i < workers; i++ {
		go func(id int) {
			for task := range s.taskQueue {
				task()
			}
		}(i)
	}
}

//тут будет лимитер который будет следить за лемитом на каждый email, что бы один еблан не положил сервак.
//Если будем писать тесты: один пользователь (то есть с одного мыла) может делать ОДИН запрос В СЕКУНДУ!
func (s *userService) checkRateLimitPerEmail(key string) bool {
	s.rateLimiter.Lock()

	now := time.Now()
	if lastTime, exists := s.rateMap[key]; exists {
		if now.Sub(lastTime) < time.Second {
			return false
		}
	}
	return true
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
	if err != nil && exist.ID != uuid.Nil {
		return nil, ErrEmailExists
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

func (s *userService) Update(ctx context.Context, userID uuid.UUID, firstName, lastName string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.UpdatedAt = time.Now()

	s.userRepo.Update(ctx, user)

	return user, nil
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
			Date:       p.Event.Date,
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

//participant-service
