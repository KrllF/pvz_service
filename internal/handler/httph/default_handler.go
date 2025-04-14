package httph

import (
	"log"
	"net/http"
)

// DefaultHandler хэндлер по умолчанию
func (h *Handler) DefaultHandler(w http.ResponseWriter, _ *http.Request) {
	response := "Здравствуйте!"
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(response)); err != nil {
		log.Printf("ошибка при отправке ответа: %v\n", err)
	}
}
