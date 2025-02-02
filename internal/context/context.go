package context

import (
	"context"

	"github.com/google/uuid"
)

type userKeyID int

const (
	keyPrincipalID userKeyID = iota
)

func SetContextUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, keyPrincipalID, userID)
}

func GetContextUserID(ctx context.Context) *uuid.UUID {
	v := ctx.Value(keyPrincipalID)
	if v == nil {
		return nil
	}
	userID := v.(uuid.UUID)
	return &userID
}
