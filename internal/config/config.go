package config

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress      string `env:"SERVER_ADDRESS"`
	URL                string `env:"URL"`
	TerminationTimeout int    `env:"TERMINATION_TIMEOUT"`
	LogLevel           string `env:"LOG_LEVEL"`
	FileStoragePath    string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN        string `env:"DATABASE_DSN"`
	AuthKey            string `env:"AUTH_KEY"`
}

var (
	ServerAddress      = "localhost:8080"
	URL                = "http://localhost:8080/"
	TerminationTimeout = time.Second * 10
	LogLevel           = "INFO"
	FileStoragePath    = "internal/storage/urls.json"
	DatabaseDSN        = ""
	AuthKey            = []byte("6368616e676520746869732070617373776f726420746f206120736563726574")
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

	if envTerminationTimeout := cfg.TerminationTimeout; envTerminationTimeout != 0 {
		TerminationTimeout = time.Duration(envTerminationTimeout)
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

	if !strings.HasSuffix(URL, "/") {
		URL += "/"
	}
	return nil
}
