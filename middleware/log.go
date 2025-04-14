//nolint:all
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
)

type (
	// Aud интерфейс сервиса аудит лога
	Aud interface {
		StoreStatus(h entity.UpdateStatus)
		StoreHTTP(h entity.HTTPInfo)
	}
	// responseWriterWrapper обёртка
	responseWriterWrapper struct {
		http.ResponseWriter
		statusCode int
		body       []byte
	}
	// ProcessOrderRequest запрос ProcessClient
	ProcessOrderRequest struct {
		UserID    int64   `json:"user_id"`
		Operation string  `json:"operation"`
		OrderIDs  []int64 `json:"order_ids"`
	}

	// ProcessOrderResponce ответ на ProcessClient
	ProcessOrderResponce struct {
		Orders []int64 `json:"orders"`
	}
)

// WriteHeader перехватывает вызов WriteHeader
func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(formJSON []byte) (int, error) {
	w.body = formJSON
	n, err := w.ResponseWriter.Write(formJSON)

	return n, err
}

// LogHandlerMiddleware логирование
func LogHandlerMiddleware(a Aud) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			defer r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			handler.ServeHTTP(wrappedWriter, r)

			info := entity.HTTPInfo{
				Method:       r.Method,
				Request:      r.URL.Path,
				ResponseCode: wrappedWriter.statusCode,
			}
			a.StoreHTTP(info)
		})
	}
}

func LogStatusUpdateMiddleware(a Aud) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK, body: nil}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)

				return
			}
			defer r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			handler.ServeHTTP(wrappedWriter, r)

			var req ProcessOrderRequest
			if err := json.NewDecoder(bytes.NewReader(body)).Decode(&req); err != nil {
				return
			}

			var res ProcessOrderResponce
			if err := json.NewDecoder(bytes.NewReader(wrappedWriter.body)).Decode(&res); err != nil {
				return
			}
			if req.Operation == consts.IssueAction {
				for _, val := range res.Orders {
					a.StoreStatus(entity.UpdateStatus{OrderID: val, NewStatus: consts.AcceptedSt})
				}
			} else if req.Operation == consts.ReturnAction {
				for _, val := range res.Orders {
					a.StoreStatus(entity.UpdateStatus{OrderID: val, NewStatus: consts.ReturnedSt})
				}
			}
		})
	}
}
