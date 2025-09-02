package auth

import (
	"context"

	"github.com/google/uuid"
)

type userKeyID int

const (
	keyPrincipalID userKeyID = iota
)

// SetContextUserID добавляет ID пользователя в контекст
// Принимает:
// ctx - исходный контекст
// userID - идентификатор пользователя
// Возвращает новый контекст с ID пользователя
func SetContextUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, keyPrincipalID, userID)
}

// GetContextUserID получает ID пользователя из контекста
// Принимает:
// ctx - контекст запроса
// Возвращает ID пользователя или nil если не установлен
func GetContextUserID(ctx context.Context) *uuid.UUID {
	v := ctx.Value(keyPrincipalID)
	if v == nil {
		return nil
	}
	userID := v.(uuid.UUID)
	return &userID
}
