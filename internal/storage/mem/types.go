package mem

import (
	"net/url"
	"sync"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service"

	"github.com/google/uuid"
)

// Repository реализует in-memory хранилище для сервиса сокращения URL.
// Обеспечивает потокобезопасные операции с использованием мьютексов.
type Repository struct {
	service.Runner
	service.Repository
	Logger         *logger.ZapLogger                    // Логгер для записи событий
	mux            sync.Mutex                           // Мьютекс для потокобезопасного доступа
	memRepository  map[uuid.UUID]*url.URL               // Основное хранилище URL
	userRepository map[uuid.UUID]map[uuid.UUID]*url.URL // Хранилище URL по пользователям
}

// NewRepository создает новый экземпляр in-memory хранилища.
// Принимает логгер и возвращает инициализированный Repository.
func NewRepository(zl *logger.ZapLogger) *Repository {
	return &Repository{
		Logger:         zl,
		mux:            sync.Mutex{},
		memRepository:  make(map[uuid.UUID]*url.URL),
		userRepository: make(map[uuid.UUID]map[uuid.UUID]*url.URL),
	}
}
