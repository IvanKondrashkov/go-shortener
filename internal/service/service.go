package service

import (
	"context"
	"fmt"
	"net/url"

	customErr "github.com/IvanKondrashkov/go-shortener/internal/errors"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, id uuid.UUID, url *url.URL) (uuid.UUID, error)
	SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error
	GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error)
	Load(ctx context.Context) error
	Ping(ctx context.Context) error
	Close(ctx context.Context)
}

type Service struct {
	Logger     *logger.ZapLogger
	Repository Repository
}

func NewService(zl *logger.ZapLogger, repository Repository) *Service {
	return &Service{
		Logger:     zl,
		Repository: repository,
	}
}

func (s *Service) Save(ctx context.Context, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	ok, _ := s.Repository.GetByID(ctx, id)

	if ok != nil {
		return id, fmt.Errorf("save error: %w", customErr.ErrConflict)
	}

	id, err := s.Repository.Save(ctx, id, u)
	if err != nil {
		return id, fmt.Errorf("save error: %w", err)
	}
	return id, nil
}

func (s *Service) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	err := s.Repository.SaveBatch(ctx, batch)
	if err != nil {
		return fmt.Errorf("save batch error: %w", err)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	u, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get url by id error: %w", err)
	}
	return u, nil
}

func (s *Service) Ping(ctx context.Context) error {
	err := s.Repository.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database ping error: %w", err)
	}
	return nil
}
