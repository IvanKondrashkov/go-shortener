package storage

import (
	"context"
	"net/url"
	"sync"

	"github.com/IvanKondrashkov/go-shortener/internal/errors"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/google/uuid"
)

type MemRepositoryImpl struct {
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

func (m *MemRepositoryImpl) Save(ctx context.Context, id uuid.UUID, u *url.URL) (res uuid.UUID, err error) {
	if ctx.Err() != nil {
		m.Logger.Log.Warn("Context is canceled!")
		return
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	_, ok := m.memRepository[id]
	if ok {
		m.memRepository[id] = u
		return id, err
	}

	m.memRepository[id] = u
	return id, err
}

func (m *MemRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (res *url.URL, err error) {
	if ctx.Err() != nil {
		m.Logger.Log.Warn("Context is canceled!")
		return
	}

	m.mux.Lock()
	defer m.mux.Unlock()

	u, ok := m.memRepository[id]
	if !ok {
		return nil, errors.ErrNotFound
	}

	return u, err
}
