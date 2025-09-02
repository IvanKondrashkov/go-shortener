package service

import (
	"context"
	"errors"
	"net/url"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Пакет service содержит определения ошибок сервисного слоя
var (
	// ErrUserUnauthorized возвращается когда операция требует авторизации пользователя
	ErrUserUnauthorized = errors.New("user unauthorized")
)

// Runner интерфейс для работы с транзакциями
type Runner interface {
	// BeginTx начинает новую транзакцию
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

// UserRepository интерфейс для пользовательских операций с URL
type UserRepository interface {
	// SaveUser сохраняет URL для конкретного пользователя
	SaveUser(ctx context.Context, tx pgx.Tx, userID uuid.UUID, id uuid.UUID, url *url.URL) (uuid.UUID, error)
	// SaveBatchUser сохраняет несколько URL для конкретного пользователя
	SaveBatchUser(ctx context.Context, userID uuid.UUID, batch []*models.RequestShortenAPIBatch) error
	// GetAllByUserID получает все URL пользователя
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]*models.ResponseShortenAPIUser, error)
	// DeleteBatchByUserID удаляет несколько URL пользователя
	DeleteBatchByUserID(ctx context.Context, userID uuid.UUID, batch []uuid.UUID) error
}

// Repository объединяет интерфейсы для работы с хранилищем URL
type Repository interface {
	Runner
	UserRepository
	// Save сохраняет URL
	Save(ctx context.Context, tx pgx.Tx, id uuid.UUID, url *url.URL) (uuid.UUID, error)
	// SaveBatch сохраняет несколько URL
	SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error
	// GetByID получает URL по его идентификатору
	GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error)
	// Load загружает данные в хранилище
	Load(ctx context.Context) error
	// Ping проверяет доступность хранилища
	Ping(ctx context.Context) error
	// Close освобождает ресурсы хранилища
	Close()
}

// Service реализует бизнес-логику сервиса сокращения URL
type Service struct {
	Runner                       // Для работы с транзакциями
	Logger     *logger.ZapLogger // Логгер для записи событий
	Repository Repository        // Репозиторий для работы с данными
}

// NewService создает новый экземпляр сервиса
// Принимает:
// - zl: логгер
// - ru: реализация интерфейса Runner
// - r: реализация интерфейса Repository
// Возвращает инициализированный Service
func NewService(zl *logger.ZapLogger, ru Runner, r Repository) *Service {
	return &Service{
		Logger:     zl,
		Runner:     ru,
		Repository: r,
	}
}
