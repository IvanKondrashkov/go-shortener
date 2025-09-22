package logger

import (
	"context"
	"time"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequestLogger логирует gRPC запросы
func RequestLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	zl, err := logger.NewZapLogger(LogLevel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Logger not work!")
	}
	defer zl.Sync()

	resp, err := handler(ctx, req)
	duration := time.Since(start)
	if err != nil {
		zl.Log.Debug("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Duration("duration", duration),
			zap.String("status", err.Error()),
		)
		return nil, err
	}

	zl.Log.Debug("gRPC request",
		zap.String("method", info.FullMethod),
		zap.Duration("duration", duration),
	)

	return resp, err
}
