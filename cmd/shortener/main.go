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
	"github.com/IvanKondrashkov/go-shortener/internal/service/grpc"
	"github.com/IvanKondrashkov/go-shortener/internal/service/worker"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/db"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/file"
	"github.com/IvanKondrashkov/go-shortener/internal/storage/mem"

	"go.uber.org/zap"
)

// Глобальные переменные сервера со значениями по умолчанию
var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

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
		newHTTPServer := handlers.NewServer(newRouter)
		newGrpcServer := grpc.NewServer(newService)

		defer newWorker.Close()

		return runServer(zl, newHTTPServer, newGrpcServer)
	}
}

func runServer(zl *logger.ZapLogger, httpServer *http.Server, grpcServer *grpc.Server) error {
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error, 2)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		zl.Log.Info("HTTP server starting", zap.String("address", config.ServerAddress), zap.Bool("tls", config.EnableHTTPS))
		if config.EnableHTTPS {
			errChan <- httpServer.ListenAndServeTLS("cert/server.crt", "cert/server.key")
		} else {
			errChan <- httpServer.ListenAndServe()
		}
	}()

	go func() {
		zl.Log.Info("gRPC server starting", zap.String("address", config.ServerAddressGrpc), zap.Bool("tls", config.EnableHTTPS))
		errChan <- grpcServer.Start()
	}()

	select {
	case sig := <-sigChan:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.TerminationTimeout)
		defer cancel()

		zl.Log.Info("Received signal, shutting down gracefully", zap.String("signal", sig.String()))
		zl.Log.Info("Stopping gRPC server...")
		grpcServer.Stop()

		zl.Log.Info("Stopping HTTP server...")
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
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

func printBuildInfo() {
	fmt.Printf("Build version: %s\n", BuildVersion)
	fmt.Printf("Build date: %s\n", BuildDate)
	fmt.Printf("Build commit: %s\n", BuildCommit)
}
