package main

import (
	"context"
	"log"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	memRepositoryImpl := storage.NewMemRepositoryImpl(zl)
	fileRepositoryImpl, err := storage.NewFileRepositoryImpl(zl, memRepositoryImpl, config.FileStoragePath)
	if err != nil {
		return err
	}

	err = fileRepositoryImpl.Load(ctx)
	if err != nil {
		return err
	}

	pgRepositoryImpl, err := storage.NewPgRepositoryImpl(zl, config.DatabaseDsn)
	if err != nil {
		return err
	}
	defer pgRepositoryImpl.Close()

	app := handlers.NewApp(memRepositoryImpl, fileRepositoryImpl, pgRepositoryImpl)
	c := api.NewController(zl, app)
	r := router.NewRouter(c)
	s := handlers.NewServer(r)

	zl.Log.Info("Running server", zap.String("address", config.ServerAddress))
	return s.ListenAndServe()
}
