package storage

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	customErr "github.com/IvanKondrashkov/go-shortener/internal/errors"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/google/uuid"
)

type MemRepositoryImpl struct {
	service.Repository
	mux            sync.Mutex
	Logger         *logger.ZapLogger
	memRepository  map[uuid.UUID]*url.URL
	userRepository map[uuid.UUID]map[uuid.UUID]*url.URL
}

func NewMemRepositoryImpl(zl *logger.ZapLogger) *MemRepositoryImpl {
	return &MemRepositoryImpl{
		mux:            sync.Mutex{},
		Logger:         zl,
		memRepository:  make(map[uuid.UUID]*url.URL),
		userRepository: make(map[uuid.UUID]map[uuid.UUID]*url.URL),
	}
}

func (m *MemRepositoryImpl) Save(ctx context.Context, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	_, ok := m.memRepository[id]
	if ok {
		m.memRepository[id] = u
		return id, fmt.Errorf("save in mem storage error: %w", customErr.ErrConflict)
	}

	m.memRepository[id] = u
	return id, nil
}

func (m *MemRepositoryImpl) SaveUser(ctx context.Context, userID, id uuid.UUID, u *url.URL) (uuid.UUID, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	_, ok := m.userRepository[userID]
	if !ok {
		m.userRepository[userID] = make(map[uuid.UUID]*url.URL)
	}

	_, ok = m.userRepository[userID][id]
	if ok {
		m.userRepository[userID][id] = u
		return id, fmt.Errorf("save in mem storage error: %w", customErr.ErrConflict)
	}

	m.userRepository[userID][id] = u
	for k, v := range m.userRepository[userID] {
		m.memRepository[k] = v
	}
	return id, nil
}

func (m *MemRepositoryImpl) SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	if len(batch) == 0 {
		return fmt.Errorf("save batch in mem storage error: %w", customErr.ErrBatchIsEmpty)
	}

	for _, it := range batch {
		u, err := url.Parse(it.OriginalURL)
		if err != nil {
			return fmt.Errorf("save batch in mem storage error: %w", customErr.ErrURLNotValid)
		}
		m.memRepository[uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String()))] = u
	}
	return nil
}

func (m *MemRepositoryImpl) SaveBatchUser(ctx context.Context, userID uuid.UUID, batch []*models.RequestShortenAPIBatch) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	if len(batch) == 0 {
		return fmt.Errorf("save batch in mem storage error: %w", customErr.ErrBatchIsEmpty)
	}

	_, ok := m.userRepository[userID]
	if !ok {
		m.userRepository[userID] = make(map[uuid.UUID]*url.URL)
	}

	for _, it := range batch {
		u, err := url.Parse(it.OriginalURL)
		if err != nil {
			return fmt.Errorf("save batch in mem storage error: %w", customErr.ErrURLNotValid)
		}
		m.userRepository[userID][uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String()))] = u
	}

	for k, v := range m.userRepository[userID] {
		m.memRepository[k] = v
	}
	return nil
}

func (m *MemRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	u, ok := m.memRepository[id]
	if !ok {
		return nil, fmt.Errorf("get in mem storage error: %w", customErr.ErrNotFound)
	}
	return u, nil
}

func (m *MemRepositoryImpl) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ResponseShortenAPIUser, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, cancel := context.WithTimeout(ctx, config.TerminationTimeout)
	defer cancel()

	urls, ok := m.userRepository[userID]
	if !ok {
		return nil, fmt.Errorf("get all in mem storage error: %w", customErr.ErrNotFound)
	}

	res := make([]*models.ResponseShortenAPIUser, 0, len(urls))
	for k, v := range urls {
		u := models.ResponseShortenAPIUser{
			ShortURL:    config.URL + k.String(),
			OriginalURL: v.String(),
		}
		res = append(res, &u)
	}
	return res, nil
}
