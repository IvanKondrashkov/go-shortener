package auth

import (
	"context"
	"strings"

	"github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authentication обрабатывает аутентификацию для gRPC
func Authentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	authHeaders := md.Get(authHeader)
	if !ok || len(authHeaders) == 0 {
		userID := uuid.New()
		token, err := generateToken(userID)
		if err != nil {
			return nil, err
		}

		err = sendHeader(ctx, token)
		if err != nil {
			return nil, err
		}
		ctx = auth.SetContextUserID(ctx, userID)
		return handler(ctx, req)
	}

	token := strings.Split(authHeaders[0], " ")[1]
	userID, err := validateToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Token is invalidate!")
	}

	ctx = auth.SetContextUserID(ctx, *userID)
	return handler(ctx, req)
}

// sendHeader отправляет заголовок на сервер gRPC
func sendHeader(ctx context.Context, token string) error {
	md := metadata.Pairs(
		authHeader, "Bearer "+token,
	)

	err := grpc.SendHeader(ctx, md)
	if err != nil {
		return status.Errorf(codes.Internal, "Send header authorization  error!")
	}
	return err
}
