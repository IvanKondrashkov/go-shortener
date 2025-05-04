package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrUserUnauthorized = errors.New("user unauthorized")
)

type Runner interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type UserRepository interface {
	SaveUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, id uuid.UUID, url *url.URL) (uuid.UUID, error)
	SaveBatchUser(ctx context.Context, userID uuid.UUID, batch []*models.RequestShortenAPIBatch) error
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ResponseShortenAPIUser, error)
	DeleteBatchByUserID(ctx context.Context, userID uuid.UUID, batch []uuid.UUID) error
}

type Repository interface {
	Runner
	UserRepository
	Save(ctx context.Context, tx pgx.Tx, id uuid.UUID, url *url.URL) (uuid.UUID, error)
	SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error
	GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error)
	Load(ctx context.Context) error
	Ping(ctx context.Context) error
	Close()
}

type Service struct {
	Runner
	Logger     *logger.ZapLogger
	Repository Repository
}

func NewService(zl *logger.ZapLogger, ru Runner, r Repository) *Service {
	return &Service{
		Logger:     zl,
		Runner:     ru,
		Repository: r,
	}
}
