package storage

import (
	"net/url"
	"sync"

	"github.com/IvanKondrashkov/go-shortener/internal/errors"
	"github.com/google/uuid"
)

type MemRepositoryImpl struct {
	mux           sync.Mutex
	memRepository map[uuid.UUID]*url.URL
}

func NewMemRepositoryImpl() *MemRepositoryImpl {
	return &MemRepositoryImpl{
		mux:           sync.Mutex{},
		memRepository: make(map[uuid.UUID]*url.URL),
	}
}

func (m *MemRepositoryImpl) Save(id uuid.UUID, u *url.URL) (res uuid.UUID, err error) {
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

func (m *MemRepositoryImpl) GetByID(id uuid.UUID) (res *url.URL, err error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	u, ok := m.memRepository[id]
	if !ok {
		return nil, errors.ErrNotFound
	}

	return u, err
}
