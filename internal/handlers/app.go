package handlers

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	api "github.com/IvanKondrashkov/go-shortener/internal/service"

	"github.com/go-chi/chi/v5"
)

type App struct {
	URL     string
	service *api.Service
}

func NewApp(service *api.Service) *App {
	return &App{
		URL:     config.URL,
		service: service,
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
