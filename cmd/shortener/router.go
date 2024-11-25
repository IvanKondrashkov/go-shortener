package main

import (
	"github.com/IvanKondrashkov/go-shortener/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route(`/`, func(r chi.Router) {
		r.Post(`/`, app.ShortenURL)
		r.Get(`/{id}`, app.GetURLByID)
	})

	return r
}
