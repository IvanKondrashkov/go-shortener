package logger

import "net/http"

// LogLevel определяет уровень логирования (DEBUG по умолчанию).
const (
	LogLevel = "DEBUG"
)

// responseData расширяет http.ResponseWriter для отслеживания статуса и размера ответа.
type responseData struct {
	http.ResponseWriter
	status int
	size   int
}
