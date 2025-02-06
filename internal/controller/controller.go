package controller

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
)

type service interface {
	ShortenURL(res http.ResponseWriter, req *http.Request)
	ShortenAPI(res http.ResponseWriter, req *http.Request)
	ShortenAPIBatch(res http.ResponseWriter, req *http.Request)
	GetURLByID(res http.ResponseWriter, req *http.Request)
	GetAllURLByUserID(res http.ResponseWriter, req *http.Request)
	DeleteBatchByUserID(res http.ResponseWriter, req *http.Request)
	Ping(res http.ResponseWriter, req *http.Request)
}

type Controller struct {
	Logger  *logger.ZapLogger
	Service service
}

func NewController(zl *logger.ZapLogger, service service) *Controller {
	return &Controller{
		Logger:  zl,
		Service: service,
	}
}
