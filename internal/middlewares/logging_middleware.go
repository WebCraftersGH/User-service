package middlewares

import (
	"github.com/WebCraftersGH/User-service/internal/requestctx"
	"github.com/WebCraftersGH/User-service/pkg/logging"
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoggingMiddleware(logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID, _ := requestctx.RequestID(r.Context())

			next.ServeHTTP(w, r)

			logger.WithFields(logrus.Fields{
				"request_id": reqID,
				"method":     r.Method,
				"path":       r.URL.Path,
				"remote":     r.RemoteAddr,
			}).Info("HTTP request")
		})
	}
}
