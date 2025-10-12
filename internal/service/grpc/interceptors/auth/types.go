package auth

import (
	"fmt"
	"time"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	authHeader      string = "Authorization"
	tokenExpiration        = 24 * time.Hour
)

// Claims представляют собой утверждения токена JWT для аутентификации пользователя
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

// validateToken валидация JWT токена
func validateToken(tokenString string) (*uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("validate token error: %w", service.ErrUserUnauthorized)
		}
		return config.AuthKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate token error: %w", service.ErrUserUnauthorized)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &claims.UserID, nil
	}

	return nil, fmt.Errorf("validate token error: %w", service.ErrUserUnauthorized)
}

// generateToken генерирует JWT токен
func generateToken(userID uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(tokenExpiration)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.AuthKey)
}
