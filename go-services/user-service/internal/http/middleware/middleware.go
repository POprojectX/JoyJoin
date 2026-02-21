package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
	customErrors "user-service/internal/errors"

	"github.com/golang-jwt/jwt"
)

func ParseJWT(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing auth header", http.StatusUnauthorized)
				return 
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid auth header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			
			user, err := ValidateToken(secretKey, tokenString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ValidateToken(secretKey, tokenString  string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, customErrors.ErrorUnexpectedSigningMethodJWT
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", customErrors.ErrInvalidToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", customErrors.ErrCantParseToken
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", customErrors.ErrTokenNoExpiration
	}
	if int64(exp) < time.Now().Unix() {
		return "", customErrors.ErrTokenExpired
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", customErrors.ErrNoUserInToken
	}
	return userID, nil
}