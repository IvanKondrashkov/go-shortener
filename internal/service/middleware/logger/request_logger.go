package logger

import (
	"net/http"
	"time"

	"github.com/IvanKondrashkov/go-shortener/internal/logger"
	"go.uber.org/zap"
)

func (r *responseData) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *responseData) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}

func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := responseData{
			ResponseWriter: w,
			status:         0,
			size:           0,
		}

		h.ServeHTTP(&responseData, r)
		duration := time.Since(start)

		zl, err := logger.NewZapLogger(LogLevel)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Logger not work!"))
			return
		}
		defer zl.Sync()

		zl.Log.Debug("HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
}
