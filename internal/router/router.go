package router

import (
	"github.com/IvanKondrashkov/go-shortener/internal/controller"
	"github.com/IvanKondrashkov/go-shortener/internal/middleware/auth"
	"github.com/IvanKondrashkov/go-shortener/internal/middleware/compress"
	"github.com/IvanKondrashkov/go-shortener/internal/middleware/logger"

	"github.com/go-chi/chi/v5"
)

func NewRouter(c *controller.Controller) *chi.Mux {
	r := chi.NewRouter()

	r.Use(logger.RequestLogger, compress.Gzip, auth.Authentication)
	r.Route(`/`, func(r chi.Router) {
		r.Post(`/`, c.Service.ShortenURL)
		r.Get(`/{id}`, c.Service.GetURLByID)
		r.Get(`/ping`, c.Service.Ping)
	})
	r.Route(`/api`, func(r chi.Router) {
		r.Post(`/shorten`, c.Service.ShortenAPI)
		r.Post(`/shorten/batch`, c.Service.ShortenAPIBatch)
		r.Get(`/user/urls`, c.Service.GetAllURLByUserID)
		r.Delete(`/user/urls`, c.Service.DeleteBatchByUserID)
	})

	return r
}
