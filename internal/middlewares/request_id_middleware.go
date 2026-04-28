package middlewares

import (
	"github.com/WebCraftersGH/User-service/internal/requestctx"
	"github.com/google/uuid"
	"net/http"
)

func GenerateRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		ctx := requestctx.WithRequestID(r.Context(), reqID)

		r = r.WithContext(ctx)

		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r)
	})
}
