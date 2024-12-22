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
	fileRepositoryImpl, err := storage.NewFileRepositoryImpl(memRepositoryImpl, config.BaseFileStoragePath)
	if err != nil {
		return err
	}

	err = fileRepositoryImpl.Load()
	if err != nil {
		return err
	}

	app := handlers.NewApp(memRepositoryImpl, fileRepositoryImpl)
	h := service.NewHandlers(app)
	r := service.NewRouter(h)
	s := handlers.NewServer(r)

	logger.Log.Info("Running server", zap.String("address", config.BaseServerAddress))
	return s.ListenAndServe()
}
