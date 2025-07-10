package logging

import (
	"forum/pkg/logger"
	"net/http"
	"time"
)

func LoggingMiddleware(log logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Вызов следующего обработчика
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Info(
			"HTTP Request",
			logger.F("remote", r.RemoteAddr),
			logger.F("method", r.Method),
			logger.F("path", r.URL.Path),
			logger.F("proto", r.Proto),
			logger.F("duration_ms", duration.Milliseconds()),
		)
	})
}
