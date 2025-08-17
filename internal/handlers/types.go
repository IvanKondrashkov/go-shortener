package handlers

import (
	"bufio"
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

var (
	readerPool = sync.Pool{
		New: func() interface{} {
			return bufio.NewReaderSize(nil, 128*1024)
		},
	}
	writerPool = sync.Pool{
		New: func() interface{} {
			return bufio.NewWriterSize(nil, 128*1024)
		},
	}
)

type Service interface {
	ShortenURL(res http.ResponseWriter, req *http.Request)
	ShortenAPI(res http.ResponseWriter, req *http.Request)
	ShortenAPIBatch(res http.ResponseWriter, req *http.Request)
	GetURLByID(res http.ResponseWriter, req *http.Request)
	GetAllURLByUserID(res http.ResponseWriter, req *http.Request)
	DeleteBatchByUserID(res http.ResponseWriter, req *http.Request)
	Ping(res http.ResponseWriter, req *http.Request)
}

type App struct {
	URL     string
	service *api.Service
	worker  *worker.Worker
}

type Handler struct {
	Logger  *logger.ZapLogger
	service Service
}

type Suite struct {
	*testing.T
	app *App
}

func NewApp(s *api.Service, w *worker.Worker) *App {
	return &App{
		URL:     config.URL,
		service: s,
		worker:  w,
	}
}

func NewHandler(zl *logger.ZapLogger, s Service) *Handler {
	return &Handler{
		Logger:  zl,
		service: s,
	}
}

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

func NewServer(r *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         config.ServerAddress,
		Handler:      r,
		ReadTimeout:  config.TerminationTimeout,
		WriteTimeout: config.TerminationTimeout,
	}
}

func NewSuite(t *testing.T) *Suite {
	t.Helper()
	t.Parallel()

	zl, _ := logger.NewZapLogger(config.LogLevel)
	newRepository := mem.NewRepository(zl)
	newRunner := newRepository
	newService := api.NewService(zl, newRunner, newRepository)
	newWorker := worker.NewWorker(config.WorkerCount, zl, newService)
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
