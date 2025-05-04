package logger

import "net/http"

const (
	LogLevel = "DEBUG"
)

type responseData struct {
	http.ResponseWriter
	status int
	size   int
}
