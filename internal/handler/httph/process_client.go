package httph

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// ProcessOrderRequest ответ на ProcessClient
type ProcessOrderRequest struct {
	UserID    int64   `json:"user_id"`
	Operation string  `json:"operation"`
	OrderIDs  []int64 `json:"order_ids"`
}

// ProcessOrderResponce ответ на ProcessClient
type ProcessOrderResponce struct {
	Orders []int64 `json:"orders"`
}

// ProcessClient принять или вернуть заказы клиента
func (h *Handler) ProcessClient(w http.ResponseWriter, r *http.Request) {
	var req ProcessOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err := ValidateRequest(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var n []int64
	var err error
	switch req.Operation {
	case "issue":
		n, err = h.acceptServ.ProcessClientIssue(r.Context(), req.UserID, req.OrderIDs)
	case "return":
		n, err = h.returnsServ.ProcessClientReturn(r.Context(), req.UserID, req.OrderIDs)
	default:
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	if err != nil && GetStatusCode(err) == http.StatusInternalServerError {
		w.WriteHeader(GetStatusCode(err))

		return
	}
	response := ProcessOrderResponce{
		Orders: n,
	}
	formattedJSON, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(formattedJSON); err != nil {
		log.Printf("ошибка при отправке ответа: %v\n", err)
	}
}

// ValidateRequest валидация запроса
func ValidateRequest(req ProcessOrderRequest) error {
	if req.UserID <= 0 {
		return errors.New("некорректный user ID")
	}

	if req.Operation == "" {
		return errors.New("некорректная операция")
	}

	if len(req.OrderIDs) == 0 {
		return errors.New("пустой список order IDs")
	}

	return nil
}
