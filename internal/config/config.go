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
}

var (
	BaseServerAddress      = "localhost:8080"
	BaseURL                = "http://localhost:8080/"
	BaseTerminationTimeout = time.Second * 10
	BaseLogLevel           = "INFO"
)

func ParseConfig() error {
	flag.StringVar(&BaseServerAddress, "a", BaseServerAddress, "Base host host:port")
	flag.StringVar(&BaseURL, "b", BaseURL, "Base url protocol://host:port/")
	flag.StringVar(&BaseLogLevel, "l", BaseLogLevel, "Base log level info")
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

	if !strings.HasSuffix(BaseURL, "/") {
		BaseURL += "/"
	}
	return nil
}
