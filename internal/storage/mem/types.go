package mem

import (
	"net/url"
	"sync"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service"

	"github.com/google/uuid"
)

type Repository struct {
	service.Runner
	service.Repository
	Logger         *logger.ZapLogger
	mux            sync.Mutex
	memRepository  map[uuid.UUID]*url.URL
	userRepository map[uuid.UUID]map[uuid.UUID]*url.URL
}

func NewRepository(zl *logger.ZapLogger) *Repository {
	return &Repository{
		Logger:         zl,
		mux:            sync.Mutex{},
		memRepository:  make(map[uuid.UUID]*url.URL),
		userRepository: make(map[uuid.UUID]map[uuid.UUID]*url.URL),
	}
}
