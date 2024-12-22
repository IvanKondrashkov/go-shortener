package config

import (
	"flag"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	BaseServerAddress      string `env:"SERVER_ADDRESS"`
	BaseURL                string `env:"BASE_URL"`
	BaseTerminationTimeout int    `env:"TERMINATION_TIMEOUT"`
	BaseLogLevel           string `env:"LOG_LEVEL"`
	BaseFileStoragePath    string `env:"FILE_STORAGE_PATH"`
}

var (
	BaseServerAddress      = "localhost:8080"
	BaseURL                = "http://localhost:8080/"
	BaseTerminationTimeout = time.Second * 10
	BaseLogLevel           = "INFO"
	BaseFileStoragePath    = "internal/storage/urls.json"
)

func ParseConfig() error {
	flag.StringVar(&BaseServerAddress, "a", BaseServerAddress, "Base host host:port")
	flag.StringVar(&BaseURL, "b", BaseURL, "Base url protocol://host:port/")
	flag.StringVar(&BaseLogLevel, "l", BaseLogLevel, "Base log level info")
	flag.StringVar(&BaseFileStoragePath, "f", BaseFileStoragePath, "Base file storage path")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return err
	}

	if envServerAddress := cfg.BaseServerAddress; envServerAddress != "" {
		BaseServerAddress = envServerAddress
	}

	if envBaseURL := cfg.BaseURL; envBaseURL != "" {
		BaseURL = envBaseURL
	}

	if envBaseTerminationTimeout := cfg.BaseTerminationTimeout; envBaseTerminationTimeout != 0 {
		BaseTerminationTimeout = time.Duration(envBaseTerminationTimeout)
	}

	if envBaseLogLevel := cfg.BaseLogLevel; envBaseLogLevel != "" {
		BaseLogLevel = envBaseLogLevel
	}

	if envBaseFileStoragePath := cfg.BaseFileStoragePath; envBaseFileStoragePath != "" {
		BaseFileStoragePath = envBaseFileStoragePath
	}

	if !strings.HasSuffix(BaseURL, "/") {
		BaseURL += "/"
	}
	return nil
}
