package grpc

import (
	"context"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	pb "github.com/IvanKondrashkov/go-shortener/internal/proto"
	"github.com/IvanKondrashkov/go-shortener/internal/service"

	"google.golang.org/grpc"
)

// Server представляет gRPC сервер
type Server struct {
	pb.UnimplementedShortenerServiceServer
	URL     string           // Базовый URL сервиса
	service *service.Service // Сервис для работы с URL
	server  *grpc.Server     // gRPC сервер
}

// NewServer создает новый экземпляр gRPC сервера
// Принимает:
// - s: сервис для работы с URL
// Возвращает инициализированный Server
func NewServer(s *service.Service) *Server {
	return &Server{
		URL:     config.URL,
		service: s,
	}
}

// methodSpecificInterceptor вспомогательная функция для методов
func methodSpecificInterceptor(method string, interceptor grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == method {
			return interceptor(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}
}
