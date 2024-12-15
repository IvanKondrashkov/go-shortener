package handlers

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/go-chi/chi/v5"
)

type App struct {
	BaseURL    string
	repository repository
}

func NewApp(repository repository) *App {
	return &App{
		BaseURL:    config.BaseURL,
		repository: repository,
	}
}

func NewServer(router *chi.Mux) *http.Server {
	return &http.Server{
		Addr:         config.BaseServerAddress,
		Handler:      router,
		ReadTimeout:  config.BaseTerminationTimeout,
		WriteTimeout: config.BaseTerminationTimeout,
	}
}
