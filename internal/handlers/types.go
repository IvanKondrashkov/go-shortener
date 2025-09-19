package handlers

import (
	"bufio"
	"context"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/models"
	api "github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/service/middleware/admin"
	"github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	"github.com/IvanKondrashkov/go-shortener/internal/service/middleware/compress"
	customLogger "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service/worker"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/mem"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Размеры буферов для чтения и записи
const (
	bufReaderSize = 128 * 1024 // 128KB размер буфера для чтения
	bufWriterSize = 128 * 1024 // 128KB размер буфера для записи
)

// Пул буферизированных ридеров и райтеров для повторного использования
var (
	readerPool = sync.Pool{
		New: func() interface{} {
			return bufio.NewReaderSize(nil, bufReaderSize)
		},
	}
	writerPool = sync.Pool{
		New: func() interface{} {
			return bufio.NewWriterSize(nil, bufWriterSize)
		},
	}
)

// AdminService определяет интерфейс для работы администратора
type AdminService interface {
	// Получить статистику сервиса
	GetStats(ctx context.Context) (*models.Stats, error)
}

// Service определяет интерфейс для работы с URL
type Service interface {
	AdminService
	// Сокращение URL (text/json)
	Save(ctx context.Context, id uuid.UUID, u *url.URL) (uuid.UUID, error)
	// Пакетное сокращение URL
	SaveBatch(ctx context.Context, batch []*models.RequestShortenAPIBatch) error
	// Получение оригинального URL по ID
	GetByID(ctx context.Context, id uuid.UUID) (*url.URL, error)
	// Получение всех URL пользователя
	GetAllByUserID(ctx context.Context) ([]*models.ResponseShortenAPIUser, error)
	// Пакетное удаление URL пользователя
	DeleteBatchByUserID(ctx context.Context, batch []uuid.UUID) error
	// Проверяет доступность хранилища
	Ping(ctx context.Context) error
}

// App представляет основное приложение с сервисом и воркером
type App struct {
	URL     string         // Базовый URL сервиса
	service *api.Service   // Сервис для работы с URL
	worker  *worker.Worker // Воркер для фоновых задач
}

// Handler обрабатывает HTTP-запросы
type Handler struct {
	Logger *logger.ZapLogger // Логгер
	app    *App              // Экземпляр приложения
}

// Suite представляет тестовый набор
type Suite struct {
	*testing.T
	app *App
}

// NewApp создает новый экземпляр App
func NewApp(s *api.Service, w *worker.Worker) *App {
	return &App{
		URL:     config.URL,
		service: s,
		worker:  w,
	}
}

// NewHandler создает новый обработчик HTTP-запросов
func NewHandler(zl *logger.ZapLogger, app *App) *Handler {
	return &Handler{
		Logger: zl,
		app:    app,
	}
}

// NewRouter создает маршрутизатор с middleware и обработчиками
func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(customLogger.RequestLogger, compress.Gzip, auth.Authentication)
	r.Route(`/`, func(r chi.Router) {
		r.Post(`/`, h.app.ShortenURL)
		r.Get(`/{id}`, h.app.GetURLByID)
		r.Get(`/ping`, h.app.Ping)
	})
	r.Route(`/api`, func(r chi.Router) {
		r.Post(`/shorten`, h.app.ShortenAPI)
		r.Post(`/shorten/batch`, h.app.ShortenAPIBatch)
		r.Get(`/user/urls`, h.app.GetAllURLByUserID)
		r.Delete(`/user/urls`, h.app.DeleteBatchByUserID)

		r.Route(`/internal`, func(r chi.Router) {
			r.Use(admin.TrustedSubnet)
			r.Get(`/stats`, h.app.GetStats)
		})
	})
	return r
}

// NewServer создает HTTP-сервер с настройками
func NewServer(r *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         config.ServerAddress,
		Handler:      r,
		ReadTimeout:  config.TerminationTimeout,
		WriteTimeout: config.TerminationTimeout,
	}
}

// NewSuite создает тестовый набор
func NewSuite(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	zl, _ := logger.NewZapLogger(config.LogLevel)
	newRepository := mem.NewRepository(zl)
	newRunner := newRepository
	newService := api.NewService(zl, newRunner, newRepository)
	newWorker := worker.NewWorker(context.Background(), config.WorkerCount, zl, newService)
	app := &App{
		URL:     config.URL,
		service: newService,
		worker:  newWorker,
	}

	return &Suite{
		T:   t,
		app: app,
	}
}
