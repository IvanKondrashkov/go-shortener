package grpc

import (
	"context"
	"errors"
	"net"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	pb "github.com/IvanKondrashkov/go-shortener/internal/proto"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/service/grpc/interceptors/auth"
	log "github.com/IvanKondrashkov/go-shortener/internal/service/grpc/interceptors/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service/grpc/interceptors/subnet"
	customError "github.com/IvanKondrashkov/go-shortener/internal/storage"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Save сохраняет URL в хранилище, gRPC обертка
// Принимает:
// - ctx: контекст с информацией о пользователе
// - in: оригинальный URL SaveRequest
// Возвращает:
// - SaveResponse или ошибку, если URL уже существует (ErrConflict) или возникли проблемы при сохранении
func (s *Server) Save(ctx context.Context, in *pb.SaveRequest) (*pb.SaveResponse, error) {
	u, err := url.Parse(in.Url)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Url is invalidate!")
	}

	id, err := s.service.Save(ctx, uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String())), u)
	if err != nil && errors.Is(err, customError.ErrConflict) {
		return nil, status.Errorf(codes.AlreadyExists, "Url is conflict!")
	}

	return &pb.SaveResponse{Result: s.URL + id.String()}, err
}

// SaveBatch сохраняет несколько URL в хранилище, gRPC обертка
// Принимает:
// - ctx: контекст с информацией о пользователе
// - in: массив URL для сохранения SaveBatchRequest
// Возвращает:
// - SaveBatchResponse или ошибку, если batch пуст или возникли проблемы при сохранении
func (s *Server) SaveBatch(ctx context.Context, in *pb.SaveBatchRequest) (*pb.SaveBatchResponse, error) {
	batchReq, err := models.RequestBatchGrpcToRequestShortenAPIBatch(in.GetBatch())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Entity mapping is incorrect!")
	}

	err = s.service.SaveBatch(ctx, batchReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Save batch error!")
	}

	batchResp, err := models.RequestBatchGrpcToResponseBatchGrpc(in.GetBatch())
	return &pb.SaveBatchResponse{Batch: batchResp}, err
}

// GetByID получает оригинальный URL по его сокращенному идентификатору, gRPC обертка
// Принимает:
// - ctx: контекст с информацией о пользователе
// - in: UUID сокращенного URL GetByIDRequest
// Возвращает:
// - GetByIDResponse или ошибку, если URL не найден или был удален
func (s *Server) GetByID(ctx context.Context, in *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	u, err := s.service.GetByID(ctx, uuid.MustParse(in.Id))
	if err != nil && errors.Is(err, customError.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "Url by id not found!")
	}

	if err != nil && errors.Is(err, customError.ErrDeleteAccepted) {
		return nil, status.Errorf(codes.FailedPrecondition, "Delete url accepted!")
	}
	return &pb.GetByIDResponse{Url: u.String()}, err
}

// GetAllByUserID получает все URL, принадлежащие текущему пользователю, gRPC обертка
// Принимает:
// - ctx: контекст с информацией о пользователе
// Возвращает:
// - GetAllByUserIDResponse или ошибку, если пользователь не авторизован или возникли проблемы при получении данных
func (s *Server) GetAllByUserID(ctx context.Context, in *emptypb.Empty) (*pb.GetAllByUserIDResponse, error) {
	urls, err := s.service.GetAllByUserID(ctx)
	if err != nil && errors.Is(err, service.ErrUserUnauthorized) {
		return nil, status.Errorf(codes.Unauthenticated, "User unauthorized!")
	}

	if len(urls) == 0 {
		return nil, status.Errorf(codes.NotFound, "Urls by user id not found!")
	}

	batchResp, err := models.ResponseShortenAPIUserToURLGrpc(urls)
	return &pb.GetAllByUserIDResponse{Urls: batchResp}, err
}

// DeleteBatchByUserID удаляет несколько URL текущего пользователя, gRPC обертка
// Принимает:
// - ctx: контекст с информацией о пользователе
// - in: массив UUID URL для удаления DeleteBatchRequest
// Возвращает:
// - ошибку, если пользователь не авторизован или возникли проблемы при удалении
func (s *Server) DeleteBatchByUserID(ctx context.Context, in *pb.DeleteBatchRequest) (*emptypb.Empty, error) {
	ids := make([]uuid.UUID, len(in.Ids))
	for _, i := range in.Ids {
		ids = append(ids, uuid.MustParse(i))
	}

	err := s.service.DeleteBatchByUserID(ctx, ids)
	if err != nil && errors.Is(err, service.ErrUserUnauthorized) {
		return nil, status.Errorf(codes.Unauthenticated, "User unauthorized!")
	}

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Delete batch error!")
	}
	return nil, err
}

// Ping проверяет доступность хранилища, gRPC обертка
// Принимает:
// - ctx: контекст
// Возвращает:
// - ошибку, если хранилище недоступно
func (s *Server) Ping(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.service.Ping(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Database is not active!")
	}
	return nil, err
}

// GetStats получить статистику сервиса, gRPC обертка
// Принимает:
// - ctx: контекст
// Возвращает:
// - StatsResponse или ошибку, если запрос не удался
func (s *Server) GetStats(ctx context.Context, in *emptypb.Empty) (*pb.StatsResponse, error) {
	stats, err := s.service.GetStats(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error!")
	}
	return &pb.StatsResponse{Urls: int64(stats.URLs), Users: int64(stats.Users)}, err
}

// Start запуск gRPC сервера
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", config.ServerAddressGrpc)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			log.RequestLogger,
			auth.Authentication,
			methodSpecificInterceptor("/shortener.ShortenerService/GetStats", subnet.TrustedSubnet),
		),
	}

	if config.EnableHTTPS {
		creds, err := credentials.NewServerTLSFromFile("cert/server.crt", "cert/server.key")
		if err != nil {
			return err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s.server = grpc.NewServer(opts...)
	pb.RegisterShortenerServiceServer(s.server, s)

	return s.server.Serve(lis)
}

// Stop остановка gRPC сервера
func (s *Server) Stop() {
	s.server.GracefulStop()
}
