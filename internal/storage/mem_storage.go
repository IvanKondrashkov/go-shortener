package storage

import (
	"net/url"
	"sync"

	"github.com/IvanKondrashkov/go-shortener/internal/errors"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
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

func (m *MemRepositoryImpl) SaveBatch(batch []*models.RequestShortenAPIBatch) (err error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if len(batch) == 0 {
		return err
	}

	for _, it := range batch {
		u, err := url.Parse(it.OriginalURL)
		if err != nil {
			return err
		}
		m.memRepository[uuid.NewSHA1(uuid.NameSpaceURL, []byte(u.String()))] = u
	}
	return err
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
