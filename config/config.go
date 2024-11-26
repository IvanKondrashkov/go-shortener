package config

import (
	"flag"
	"log"
	"strings"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

var (
	ServerAddress = "localhost:8080"
	BaseURL       = "http://localhost:8080/"
)

func ParseConfig() {
	flag.StringVar(&ServerAddress, "a", ServerAddress, "Base host host:port")
	flag.StringVar(&BaseURL, "b", BaseURL, "Base url protocol://host:port/")
	flag.Parse()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if envServerAddress := cfg.ServerAddress; envServerAddress != "" {
		ServerAddress = envServerAddress
	}

	if envBaseURL := cfg.BaseURL; envBaseURL != "" {
		BaseURL = envBaseURL
	}

	if !strings.HasSuffix(BaseURL, "/") {
		BaseURL += "/"
	}
}
