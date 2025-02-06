package main

import (
	"context"
	"log"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	api "github.com/IvanKondrashkov/go-shortener/internal/controller"
	"github.com/IvanKondrashkov/go-shortener/internal/handlers"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/router"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/storage"
	"github.com/IvanKondrashkov/go-shortener/internal/worker"
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
	zl, err := logger.NewZapLogger(config.LogLevel)
	if err != nil {
		return err
	}
	defer zl.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), config.TerminationTimeout)
	defer cancel()

	var newRepository service.Repository
	newRepository = storage.NewMemRepositoryImpl(zl)
	if config.FileStoragePath != "" {
		newRepository, err = storage.NewFileRepositoryImpl(zl, newRepository, config.FileStoragePath)
		if err != nil {
			return err
		}

		err = newRepository.Load(ctx)
		if err != nil {
			return err
		}
	}

	if config.DatabaseDSN != "" {
		newRepository, err = storage.NewPgRepositoryImpl(ctx, zl, config.DatabaseDSN)
		if err != nil {
			return err
		}
		defer newRepository.Close()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		newService := service.NewService(zl, newRepository)
		newWorker := worker.NewWorker(config.WorkerCount, zl, newService)
		newApp := handlers.NewApp(newService, newWorker)
		newController := api.NewController(zl, newApp)
		newRouter := router.NewRouter(newController)
		newServer := handlers.NewServer(newRouter)

		defer newWorker.Close()

		zl.Log.Info("Running server", zap.String("address", config.ServerAddress))
		return newServer.ListenAndServe()
	}
}
