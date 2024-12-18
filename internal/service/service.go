package service

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type service interface {
	ShortenURL(res http.ResponseWriter, req *http.Request)
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

	r.Use(middleware.Logger)
	r.Route(`/`, func(r chi.Router) {
		r.Post(`/`, h.service.ShortenURL)
		r.Get(`/{id}`, h.service.GetURLByID)
	})

	return r
}
