package config

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	URL             string `env:"URL"`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	AuthKey         string `env:"AUTH_KEY"`

	TerminationTimeout int `env:"TERMINATION_TIMEOUT"`
	WorkerCount        int `env:"WORKER_COUNT"`
}

var (
	ServerAddress   = "localhost:8080"
	URL             = "http://localhost:8080/"
	LogLevel        = "INFO"
	FileStoragePath = "internal/storage/urls.json"
	DatabaseDSN     = ""
	AuthKey         = []byte("6368616e676520746869732070617373776f726420746f206120736563726574")

	TerminationTimeout = time.Second * 30
	WorkerCount        = 10
)

func ParseConfig() error {
	flag.StringVar(&ServerAddress, "a", ServerAddress, "Base host host:port")
	flag.StringVar(&URL, "b", URL, "Base url protocol://host:port/")
	flag.StringVar(&LogLevel, "l", LogLevel, "Base log level info")
	flag.StringVar(&FileStoragePath, "f", FileStoragePath, "Base file storage path")
	flag.StringVar(&DatabaseDSN, "d", DatabaseDSN, "Base url db connection")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("config parse error: %w", err)
	}

	if envServerAddress := cfg.ServerAddress; envServerAddress != "" {
		ServerAddress = envServerAddress
	}

	if envBaseURL := cfg.URL; envBaseURL != "" {
		URL = envBaseURL
	}

	if envLogLevel := cfg.LogLevel; envLogLevel != "" {
		LogLevel = envLogLevel
	}

	if envFileStoragePath := cfg.FileStoragePath; envFileStoragePath != "" {
		FileStoragePath = envFileStoragePath
	}

	if envDatabaseDsn := cfg.DatabaseDSN; envDatabaseDsn != "" {
		DatabaseDSN = envDatabaseDsn
	}

	if envAuthKey := cfg.AuthKey; envAuthKey != "" {
		AuthKey = []byte(envAuthKey)
	}

	if envTerminationTimeout := cfg.TerminationTimeout; envTerminationTimeout != 0 {
		TerminationTimeout = time.Duration(envTerminationTimeout)
	}

	if envWorkerCount := cfg.WorkerCount; envWorkerCount != 0 {
		WorkerCount = envWorkerCount
	}

	if !strings.HasSuffix(URL, "/") {
		URL += "/"
	}
	return nil
}
