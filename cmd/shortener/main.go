package main

import (
	"log"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/handlers"
	"github.com/IvanKondrashkov/go-shortener/internal/middleware/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/storage"
	"go.uber.org/zap"
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
	if err := logger.Initialize(config.BaseLogLevel); err != nil {
		return err
	}

	memRepositoryImpl := storage.NewMemRepositoryImpl()
	app := handlers.NewApp(memRepositoryImpl)
	h := service.NewHandlers(app)
	r := service.NewRouter(h)
	s := handlers.NewServer(r)

	logger.Log.Info("Running server", zap.String("address", config.BaseServerAddress))
	return s.ListenAndServe()
}
