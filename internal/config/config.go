package config

import (
	"flag"
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	BaseServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL           string `env:"BASE_URL"`
}

var (
	BaseServerAddress = "localhost:8080"
	BaseURL           = "http://localhost:8080/"
)

func ParseConfig() error {
	flag.StringVar(&BaseServerAddress, "a", BaseServerAddress, "Base host host:port")
	flag.StringVar(&BaseURL, "b", BaseURL, "Base url protocol://host:port/")
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

	if !strings.HasSuffix(BaseURL, "/") {
		BaseURL += "/"
	}
	return nil
}
