package storage

import (
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/models"
	"github.com/google/uuid"
)

type MemRepositoryImpl struct {
	memRepository map[uuid.UUID]*url.URL
}

func NewMemRepositoryImpl() *MemRepositoryImpl {
	return &MemRepositoryImpl{
		memRepository: make(map[uuid.UUID]*url.URL),
	}
}

func (m *MemRepositoryImpl) Save(id uuid.UUID, u *url.URL) (res uuid.UUID, err error) {
	_, ok := m.memRepository[id]
	if ok {
		return id, models.ErrConflict
	}

	m.memRepository[id] = u
	return id, err
}

func (m *MemRepositoryImpl) GetByID(id uuid.UUID) (res *url.URL, err error) {
	_, ok := m.memRepository[id]
	if !ok {
		return nil, models.ErrNotFound
	}

	return m.memRepository[id], err
}
