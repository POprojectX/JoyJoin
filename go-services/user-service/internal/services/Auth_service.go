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
)

type UserService interface {
	// Auth
	Register(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	
	// CRUD
	UpdateProfile(ctx context.Context, userID uuid.UUID, firstName, lastName string) error
	DeleteAccount(ctx context.Context, userID uuid.UUID) error

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
	return &userService{
		userRepo: r,
	}
}

func (s *userService) Register(email, password, firstName, lastName string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, ErrValidationFailed
	}
	if len(password) < 6 {
		return nil, ErrShortPassword
	}

	exist, err := s.userRepo.GetByEmail(email)
	if err != nil && exist.ID != uuid.Nil {
		return nil, ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Password: hashedPassword,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
//evet-service
//participant-service
