package main

import (
	"log"
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/handlers"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/storage"
)

func main() {
	err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	memRepositoryImpl := storage.NewMemRepositoryImpl()
	memService := handlers.NewApp(config.BaseURL, memRepositoryImpl)
	h := service.NewHandlers(memService)
	r := service.NewRouter(h)
	return http.ListenAndServe(config.BaseServerAddress, r)
}
