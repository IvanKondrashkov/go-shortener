package config

import (
	"flag"
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
	DatabaseDsn        string `env:"DATABASE_DSN"`
}

var (
	ServerAddress      = "localhost:8080"
	URL                = "http://localhost:8080/"
	TerminationTimeout = time.Second * 10
	LogLevel           = "INFO"
	FileStoragePath    = "internal/storage/urls.json"
	DatabaseDsn        = ""
)

func ParseConfig() error {
	flag.StringVar(&ServerAddress, "a", ServerAddress, "Base host host:port")
	flag.StringVar(&URL, "b", URL, "Base url protocol://host:port/")
	flag.StringVar(&LogLevel, "l", LogLevel, "Base log level info")
	flag.StringVar(&FileStoragePath, "f", FileStoragePath, "Base file storage path")
	flag.StringVar(&DatabaseDsn, "d", DatabaseDsn, "Base url db connection")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return err
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

	if envDatabaseDsn := cfg.DatabaseDsn; envDatabaseDsn != "" {
		DatabaseDsn = envDatabaseDsn
	}

	if !strings.HasSuffix(URL, "/") {
		URL += "/"
	}
	return nil
}
