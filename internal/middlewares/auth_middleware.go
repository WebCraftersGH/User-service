package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/WebCraftersGH/User-service/internal/authclient"
	"github.com/WebCraftersGH/User-service/internal/requestctx"
	"github.com/WebCraftersGH/User-service/pkg/logging"
)

type AuthChecker interface {
	Check(ctx context.Context, token string) (uuid.UUID, error)
}

func AuthFromToken(authChecker AuthChecker, logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "empty auth http header", http.StatusUnauthorized)
				return
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				http.Error(w, "empty auth prefix (Bearer)", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
			if token == "" {
				http.Error(w, "empty token", http.StatusUnauthorized)
				return
			}

			userID, err := authChecker.Check(r.Context(), token)
			if err != nil {
				switch {
				case errors.Is(err, authclient.ErrUnauthorized):
					http.Error(w, "user is unauthorized", http.StatusUnauthorized)
					return
				default:
					logger.WithError(err).Error("auth check error")
					http.Error(w, "auth check unknown error", http.StatusServiceUnavailable)
					return
				}
			}

			ctx := requestctx.WithUserID(r.Context(), userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
