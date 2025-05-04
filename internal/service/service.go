package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/models"
	customContext "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	customError "github.com/IvanKondrashkov/go-shortener/internal/storage"

	"github.com/google/uuid"
)

func (s *Service) Save(ctx context.Context, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	ok, _ := s.Repository.GetByID(ctx, id)
	if ok != nil {
		return id, fmt.Errorf("save error: %w", customError.ErrConflict)
	}

	tx, err := s.BeginTx(ctx)
	if err != nil {
		return id, fmt.Errorf("open transactional error: %w", err)
	}

	if tx == nil {
		userID := customContext.GetContextUserID(ctx)
		if userID != nil {
			return s.Repository.SaveUser(ctx, nil, *userID, id, u)
		}
		return s.Repository.Save(ctx, nil, id, u)
	}

	userID := customContext.GetContextUserID(ctx)
	if userID != nil {
		id, err = s.Repository.SaveUser(ctx, tx, *userID, id, u)
		if err != nil {
			_ = tx.Rollback(ctx)
			return id, err
		}
		err = tx.Commit(ctx)
		if err != nil {
			return id, err
		}
	} else {
		id, err = s.Repository.Save(ctx, tx, id, u)
		if err != nil {
			_ = tx.Rollback(ctx)
			return id, err
		}
		err = tx.Commit(ctx)
		if err != nil {
			return id, err
		}
	}
	return id, err
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
	return nil, fmt.Errorf("get all url by user id error: %w", ErrUserUnauthorized)
}

func (s *Service) DeleteBatchByUserID(ctx context.Context, batch []uuid.UUID) error {
	userID := customContext.GetContextUserID(ctx)
	if userID != nil {
		err := s.Repository.DeleteBatchByUserID(ctx, *userID, batch)
		if err != nil {
			return fmt.Errorf("user delete batch error: %w", err)
		}
		return nil
	}
	return fmt.Errorf("delete batch by user id error: %w", ErrUserUnauthorized)
}

func (s *Service) Ping(ctx context.Context) error {
	err := s.Repository.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database ping error: %w", err)
	}
	return nil
}
