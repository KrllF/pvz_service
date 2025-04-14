package httph

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ReturnOrderToCourierResponse ответ на ReturnOrderToCourier
type ReturnOrderToCourierResponse struct {
	OrderDel int64 `json:"orderdel"`
}

// ReturnOrderToCourier вернуть заказы курьеру
func (h *Handler) ReturnOrderToCourier(w http.ResponseWriter, r *http.Request) {
	orderID, err := validateReturnToCourier(mux.Vars(r))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.returnsServ.ReturnOrder(r.Context(), orderID)
	if err != nil {
		w.WriteHeader(GetStatusCode(err))

		return
	}

	response := ReturnOrderToCourierResponse{
		OrderDel: orderID,
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

func validateReturnToCourier(mp map[string]string) (int64, error) {
	key, ok := mp["id"]
	if !ok {
		return 0, errors.New("нет параметра order_id")
	}

	orderID, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
