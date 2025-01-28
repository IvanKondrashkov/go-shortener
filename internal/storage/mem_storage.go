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
	mux           sync.Mutex
	Logger        *logger.ZapLogger
	memRepository map[uuid.UUID]*url.URL
}

func NewMemRepositoryImpl(zl *logger.ZapLogger) *MemRepositoryImpl {
	return &MemRepositoryImpl{
		mux:           sync.Mutex{},
		Logger:        zl,
		memRepository: make(map[uuid.UUID]*url.URL),
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
