package main

import (
	"context"
	"log"

	"github.com/IvanKondrashkov/go-shortener/internal/config"
	api "github.com/IvanKondrashkov/go-shortener/internal/controller"
	"github.com/IvanKondrashkov/go-shortener/internal/handlers"
	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"github.com/IvanKondrashkov/go-shortener/internal/router"
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
	zl, err := logger.NewZapLogger(config.LogLevel)
	if err != nil {
		return err
	}
	defer zl.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), config.TerminationTimeout)
	defer cancel()

	memRepositoryImpl := storage.NewMemRepositoryImpl(zl)
	var fileRepositoryImpl *storage.FileRepositoryImpl
	if config.FileStoragePath != "" {
		fileRepositoryImpl, err = storage.NewFileRepositoryImpl(zl, memRepositoryImpl, config.FileStoragePath)
		if err != nil {
			return err
		}

		err = fileRepositoryImpl.Load()
		if err != nil {
			return err
		}
	}

	var pgRepositoryImpl *storage.PgRepositoryImpl
	if config.DatabaseDsn != "" {
		pgRepositoryImpl, err = storage.NewPgRepositoryImpl(zl, config.DatabaseDsn)
		if err != nil {
			return err
		}
		defer pgRepositoryImpl.Close()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		app := handlers.NewApp(memRepositoryImpl, fileRepositoryImpl, pgRepositoryImpl)
		c := api.NewController(zl, app)
		r := router.NewRouter(c)
		s := handlers.NewServer(r)

		zl.Log.Info("Running server", zap.String("address", config.ServerAddress))
		return s.ListenAndServe()
	}
}
