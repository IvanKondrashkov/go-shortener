package main

import (
	"net/http"

	"github.com/IvanKondrashkov/go-shortener/config"
	"github.com/IvanKondrashkov/go-shortener/internal/app"
	"github.com/IvanKondrashkov/go-shortener/storage"
)

func main() {
	config.ParseConfig()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	memRepositoryImpl := storage.NewMemRepositoryImpl()
	router := NewRouter(app.NewApp(
		config.BaseURL,
		memRepositoryImpl))
	return http.ListenAndServe(config.ServerAddress, router)
}
