package httph

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/KrllF/pvz_service/internal/entity"
)

// AcceptOrder принять заказ
func (h *Handler) AcceptOrder(w http.ResponseWriter, r *http.Request) {
	var orderEntry entity.OrderEntry
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&orderEntry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	ok := ValidateReq(orderEntry)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := h.acceptServ.AcceptOrder(r.Context(), orderEntry); err != nil {
		w.WriteHeader(GetStatusCode(err))

		return
	}

	w.WriteHeader(http.StatusCreated)
	response := "Заказ успешно принят"
	if _, err := w.Write([]byte(response)); err != nil {
		log.Printf("ошибка при отправке ответа: %v\n", err)
	}
}

// ValidateReq валидация заказа
func ValidateReq(ord entity.OrderEntry) bool {
	if ord.PackType == "" {
		return false
	}
	if ord.Price < 0 {
		return false
	}
	if ord.Weight < 0 {
		return false
	}
	if ord.UserID < 0 {
		return false
	}
	if ord.OrderID < 0 {
		return false
	}

	return true
}
