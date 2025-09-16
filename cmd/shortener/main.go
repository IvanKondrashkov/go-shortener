package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

const defaultBuildInfo = "N/A"

// @title Go Shortener API
// @version 1.0
// @description API сервиса сокращения URL

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	printBuildInfo()

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
		newWorker := worker.NewWorker(ctx, config.WorkerCount, zl, newService)
		newApp := handlers.NewApp(newService, newWorker)
		newHandler := handlers.NewHandler(zl, newApp)
		newRouter := handlers.NewRouter(newHandler)
		newServer := handlers.NewServer(newRouter)

		defer newWorker.Close()

		sigChan := make(chan os.Signal, 1)
		errChan := make(chan error, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		go func() {
			zl.Log.Info("Running server", zap.String("address", config.ServerAddress))
			if config.EnableHTTPS {
				errChan <- newServer.ListenAndServeTLS("cert/server.crt", "cert/server.key")
			} else {
				errChan <- newServer.ListenAndServe()
			}
		}()

		select {
		case sig := <-sigChan:
			zl.Log.Info("Received signal, shutting down gracefully", zap.String("signal", sig.String()))
			if err := newServer.Shutdown(ctx); err != nil {
				zl.Log.Error("Server shutdown failed", zap.Error(err))
				return err
			}

			zl.Log.Info("Server stopped gracefully")
			return nil

		case err := <-errChan:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				zl.Log.Error("Server error", zap.Error(err))
				return err
			}
			return nil
		}
	}
}

func printBuildInfo() {
	version := buildVersion
	if version == "" {
		version = defaultBuildInfo
	}

	date := buildDate
	if date == "" {
		date = defaultBuildInfo
	}

	commit := buildCommit
	if commit == "" {
		commit = defaultBuildInfo
	}

	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", date)
	fmt.Printf("Build commit: %s\n", commit)
}
