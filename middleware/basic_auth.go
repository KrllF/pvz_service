package middleware

import (
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
)

// BasicAuthMiddleware проверка авторизации
func BasicAuthMiddleware(config config.ConnectConfig) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok || username != config.Login || password != config.Password {
				w.WriteHeader(http.StatusUnauthorized)

				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
