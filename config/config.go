package config

import (
	"flag"
)

var (
	BaseHost = ":8080"
	BaseURL  = "http://localhost:8080"
)

func ParseFlags() {
	flag.StringVar(&BaseHost, "a", BaseHost, "Base host host:port")
	flag.StringVar(&BaseURL, "b", BaseURL, "Base url protocol://host:port/")
	flag.Parse()
}
