package middlewares

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"

	"github.com/WebCraftersGH/User-service/internal/authclient"
	"github.com/WebCraftersGH/User-service/internal/requestctx"
)

type AuthChecker interface {
	Check(ctx context.Context, token string) (uuid.UUID, error)
}

func AuthFromToken(authChecker AuthChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
			if token == "" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			userID, err := authChecker.Check(r.Context(), token)
			if err != nil {
				switch {
				case errors.Is(err, authclient.ErrUnauthorized):
					http.Error(w, "", http.StatusUnauthorized)
					return
				default:
					http.Error(w, "", http.StatusServiceUnavailable)
					return
				}
			}

			ctx := requestctx.WithUserID(r.Context(), userID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
