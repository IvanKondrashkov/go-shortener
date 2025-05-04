package main

import (
	"context"
	"log"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	"github.com/IvanKondrashkov/go-shortener/internal/handlers"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/service"
	"github.com/IvanKondrashkov/go-shortener/internal/service/worker"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/db"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/file"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/mem"

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
	var newRunner service.Runner

	newRepository = mem.NewRepository(zl)
	newRunner = newRepository
	if config.FileStoragePath != "" {
		newRepository, err = file.NewRepository(zl, newRepository, config.FileStoragePath)
		newRunner = newRepository
		if err != nil {
			return err
		}

		err = newRepository.Load(ctx)
		if err != nil {
			return err
		}
	}

	if config.DatabaseDSN != "" {
		newRepository, err = db.NewRepository(ctx, zl, config.DatabaseDSN)
		newRunner = newRepository
		if err != nil {
			return err
		}
		defer newRepository.Close()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		newService := service.NewService(zl, newRunner, newRepository)
		newWorker := worker.NewWorker(config.WorkerCount, zl, newService)
		newApp := handlers.NewApp(newService, newWorker)
		newHandler := handlers.NewHandler(zl, newApp)
		newRouter := handlers.NewRouter(newHandler)
		newServer := handlers.NewServer(newRouter)

		defer newWorker.Close()

		zl.Log.Info("Running server", zap.String("address", config.ServerAddress))
		return newServer.ListenAndServe()
	}
}
