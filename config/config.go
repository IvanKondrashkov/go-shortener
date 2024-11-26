package config

import (
	"flag"
	"strings"
)

var (
	BaseHost = "localhost:8080"
	BaseURL  = "http://localhost:8080/"
)

func ParseFlags() {
	flag.StringVar(&BaseHost, "a", BaseHost, "Base host host:port")
	flag.StringVar(&BaseURL, "b", BaseURL, "Base url protocol://host:port/")
	flag.Parse()

	if !strings.HasSuffix(BaseURL, "/") {
		BaseURL += "/"
	}
}
