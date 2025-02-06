package handlers

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	api "github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/worker"

	"github.com/go-chi/chi/v5"
)

type App struct {
	URL     string
	service *api.Service
	worker  *worker.Worker
}

func NewApp(service *api.Service, deleteWorker *worker.Worker) *App {
	return &App{
		URL:     config.URL,
		service: service,
		worker:  deleteWorker,
	}
}

func NewServer(router *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         config.ServerAddress,
		Handler:      router,
		ReadTimeout:  config.TerminationTimeout,
		WriteTimeout: config.TerminationTimeout,
	}
}
