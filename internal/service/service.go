package service

import (
	"context"
	"fmt"
	"net/url"

	customContext "github.com/IvanKondrashkov/go-shortener/internal/context"
	customErr "github.com/IvanKondrashkov/go-shortener/internal/errors"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, id uuid.UUID, url *url.URL) (uuid.UUID, error)
	SaveUser(ctx context.Context, userID uuid.UUID, id uuid.UUID, url *url.URL) (uuid.UUID, error)
	SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error
	SaveBatchUser(ctx context.Context, userID uuid.UUID, batch []*models.RequestShortenAPIBatch) error
	GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ResponseShortenAPIUser, error)
	DeleteBatchByUserID(ctx context.Context, userID uuid.UUID, batch []*uuid.UUID) error
	Load(ctx context.Context) error
	Ping(ctx context.Context) error
	Close()
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

	userID := customContext.GetContextUserID(ctx)
	if userID != nil {
		id, err := s.Repository.SaveUser(ctx, *userID, id, u)
		if err != nil {
			return id, fmt.Errorf("user save error: %w", err)
		}
		return id, nil
	}

	id, err := s.Repository.Save(ctx, id, u)
	if err != nil {
		return id, fmt.Errorf("save error: %w", err)
	}
	return id, nil
}

func (s *Service) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	userID := customContext.GetContextUserID(ctx)
	if userID != nil {
		err := s.Repository.SaveBatchUser(ctx, *userID, batch)
		if err != nil {
			return fmt.Errorf("user save batch error: %w", err)
		}
		return nil
	}

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

func (s *Service) GetAllByUserID(ctx context.Context) ([]*models.ResponseShortenAPIUser, error) {
	userID := customContext.GetContextUserID(ctx)
	if userID != nil {
		urls, err := s.Repository.GetAllByUserID(ctx, *userID)
		if err != nil {
			return nil, fmt.Errorf("user get all urls error: %w", err)
		}
		return urls, nil
	}
	return nil, fmt.Errorf("get all url by user id error: %w", customErr.ErrUserUnauthorized)
}

func (s *Service) DeleteBatchByUserID(ctx context.Context, batch []*uuid.UUID) error {
	userID := customContext.GetContextUserID(ctx)
	if userID != nil {
		err := s.Repository.DeleteBatchByUserID(ctx, *userID, batch)
		if err != nil {
			return fmt.Errorf("user delete batch error: %w", err)
		}
		return nil
	}
	return fmt.Errorf("delete batch by user id error: %w", customErr.ErrUserUnauthorized)
}

func (s *Service) Ping(ctx context.Context) error {
	err := s.Repository.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database ping error: %w", err)
	}
	return nil
}
