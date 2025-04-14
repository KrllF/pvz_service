package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

type responseWriterStatusWrapper struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader перехватывает вызов WriteHeader
func (w *responseWriterStatusWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// SpanMiddleware span
func SpanMiddleware() func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span, ctx := opentracing.StartSpanFromContext(r.Context(), r.URL.Path)
			defer span.Finish()

			wrappedWriter := &responseWriterStatusWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			handler.ServeHTTP(wrappedWriter, r.WithContext(ctx))
			log.Println(wrappedWriter.statusCode)
			if wrappedWriter.statusCode >= http.StatusBadRequest &&
				wrappedWriter.statusCode < http.StatusInternalServerError {
				span.SetTag("bad request", true)
				span.LogKV("bad request", wrappedWriter.statusCode)
			} else if wrappedWriter.statusCode >= http.StatusInternalServerError {
				span.SetTag("error", true)
				span.LogKV("error", errors.New("wrappedWriter.statusCode"))
			}
		})
	}
}
