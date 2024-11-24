package main

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/app"
)

func NewRouter(app *app.App) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, app.ShortenURL)
	mux.HandleFunc(`/{id}`, app.GetURLByID)
	return mux
}
