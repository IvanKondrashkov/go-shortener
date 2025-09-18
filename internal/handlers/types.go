package handlers

import (
	"bufio"
	"context"
	"net/http"
	"sync"
	"testing"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	api "github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/service/middleware/auth"
	"github.com/IvanKondrashkov/go-shortener/internal/service/middleware/compress"
	customLogger "github.com/IvanKondrashkov/go-shortener/internal/service/middleware/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service/worker"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/mem"

	"github.com/go-chi/chi/v5"
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

// Service определяет интерфейс для работы с URL
type Service interface {
	// Сокращение URL (text)
	ShortenURL(res http.ResponseWriter, req *http.Request)
	// Сокращение URL (json)
	ShortenAPI(res http.ResponseWriter, req *http.Request)
	// Пакетное сокращение URL
	ShortenAPIBatch(res http.ResponseWriter, req *http.Request)
	// Получение оригинального URL по ID
	GetURLByID(res http.ResponseWriter, req *http.Request)
	// Получение всех URL пользователя
	GetAllURLByUserID(res http.ResponseWriter, req *http.Request)
	// Пакетное удаление URL пользователя
	DeleteBatchByUserID(res http.ResponseWriter, req *http.Request)
	// Пакетное удаление URL пользователя
	Ping(res http.ResponseWriter, req *http.Request)
}

// App представляет основное приложение с сервисом и воркером
type App struct {
	URL     string         // Базовый URL сервиса
	service *api.Service   // Сервис для работы с URL
	worker  *worker.Worker // Воркер для фоновых задач
}

// Handler обрабатывает HTTP-запросы
type Handler struct {
	Logger  *logger.ZapLogger // Логгер
	service Service           // Сервис для работы с URL
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
func NewHandler(zl *logger.ZapLogger, s Service) *Handler {
	return &Handler{
		Logger:  zl,
		service: s,
	}
}

// NewRouter создает маршрутизатор с middleware и обработчиками
func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(customLogger.RequestLogger, compress.Gzip, auth.Authentication)
	r.Route(`/`, func(r chi.Router) {
		r.Post(`/`, h.service.ShortenURL)
		r.Get(`/{id}`, h.service.GetURLByID)
		r.Get(`/ping`, h.service.Ping)
	})
	r.Route(`/api`, func(r chi.Router) {
		r.Post(`/shorten`, h.service.ShortenAPI)
		r.Post(`/shorten/batch`, h.service.ShortenAPIBatch)
		r.Get(`/user/urls`, h.service.GetAllURLByUserID)
		r.Delete(`/user/urls`, h.service.DeleteBatchByUserID)
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
