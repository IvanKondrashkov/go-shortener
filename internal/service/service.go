package service

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/middleware/compress"
	"github.com/IvanKondrashkov/go-shortener/internal/middleware/logger"
	"github.com/go-chi/chi/v5"
)

type service interface {
	ShortenURL(res http.ResponseWriter, req *http.Request)
	ShortenAPI(res http.ResponseWriter, req *http.Request)
	GetURLByID(res http.ResponseWriter, req *http.Request)
}

type handlers struct {
	service service
}

func NewHandlers(service service) *handlers {
	return &handlers{
		service: service,
	}
}

func NewRouter(h *handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.RequestLogger, compress.Gzip)
	r.Route(`/`, func(r chi.Router) {
		r.Post(`/`, h.service.ShortenURL)
		r.Get(`/{id}`, h.service.GetURLByID)
	})
	r.Route(`/api`, func(r chi.Router) {
		r.Post(`/shorten`, h.service.ShortenAPI)
	})

	return r
}
