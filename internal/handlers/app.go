package handlers

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/go-chi/chi/v5"
)

type App struct {
	URL            string
	repository     repository
	fileRepository fileRepository
	pgRepository   pgRepository
}

func NewApp(repository repository, fileRepository fileRepository, pgRepository pgRepository) *App {
	return &App{
		URL:            config.URL,
		repository:     repository,
		fileRepository: fileRepository,
		pgRepository:   pgRepository,
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
