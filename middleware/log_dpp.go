package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// LogMiddleware логирование http.MethodDelete, http.MethodPost, http.MethodPut
func LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete, http.MethodPost, http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			defer r.Body.Close()
			log.Printf("%s, %s, %s", r.Method, r.URL, string(body))
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		handler.ServeHTTP(w, r)
	})
}
