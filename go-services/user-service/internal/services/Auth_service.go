package services

import (
	"time"
	"user-service/internal/repo"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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