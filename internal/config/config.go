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
}

var (
	ServerAddress      = "localhost:8080"
	URL                = "http://localhost:8080/"
	TerminationTimeout = time.Second * 10
	LogLevel           = "INFO"
	FileStoragePath    = "internal/storage/urls.json"
)

func ParseConfig() error {
	flag.StringVar(&ServerAddress, "a", ServerAddress, "Base host host:port")
	flag.StringVar(&URL, "b", URL, "Base url protocol://host:port/")
	flag.StringVar(&LogLevel, "l", LogLevel, "Base log level info")
	flag.StringVar(&FileStoragePath, "f", FileStoragePath, "Base file storage path")
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

	if envBaseTerminationTimeout := cfg.TerminationTimeout; envBaseTerminationTimeout != 0 {
		TerminationTimeout = time.Duration(envBaseTerminationTimeout)
	}

	if envBaseLogLevel := cfg.LogLevel; envBaseLogLevel != "" {
		LogLevel = envBaseLogLevel
	}

	if envBaseFileStoragePath := cfg.FileStoragePath; envBaseFileStoragePath != "" {
		FileStoragePath = envBaseFileStoragePath
	}

	if !strings.HasSuffix(URL, "/") {
		URL += "/"
	}
	return nil
}
