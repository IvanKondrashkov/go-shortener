package service

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Service interface {
	ShortenURL(res http.ResponseWriter, req *http.Request)
	GetURLByID(res http.ResponseWriter, req *http.Request)
}

type Handlers struct {
	service Service
}

func NewHandlers(service Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

func NewRouter(h *Handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Route(`/`, func(r chi.Router) {
		r.Post(`/`, h.service.ShortenURL)
		r.Get(`/{id}`, h.service.GetURLByID)
	})

	return r
}
