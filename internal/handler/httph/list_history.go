package httph

import (
	"encoding/json"
	"log"
	"net/http"
)

// ListHistory получить историю заказов
func (h *Handler) ListHistory(w http.ResponseWriter, r *http.Request) {
	history, err := h.showServ.ListHistory(r.Context())
	if err != nil {
		w.WriteHeader(GetStatusCode(err))

		return
	}

	b, err := json.MarshalIndent(history, "", "\t")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Printf("ошибка при отправке ответа: %v\n", err)
	}
}
