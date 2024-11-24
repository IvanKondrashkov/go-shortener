package storage

import (
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/models"
)

type MemRepositoryImpl struct {
	memRepository map[string]*url.URL
}

func NewMemRepositoryImpl() *MemRepositoryImpl {
	return &MemRepositoryImpl{
		memRepository: make(map[string]*url.URL),
	}
}

func (m *MemRepositoryImpl) Save(encodedURL string, u *url.URL) (id string, err error) {
	_, ok := m.memRepository[encodedURL]
	if ok {
		return encodedURL, models.ErrConflict
	}

	m.memRepository[encodedURL] = u
	return encodedURL, err
}

func (m *MemRepositoryImpl) GetByID(id string) (u *url.URL, err error) {
	_, ok := m.memRepository[id]
	if !ok {
		return nil, models.ErrNotFound
	}

	return m.memRepository[id], err
}
