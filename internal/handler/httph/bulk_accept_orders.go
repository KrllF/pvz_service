package httph

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/KrllF/pvz_service/internal/errs"
)

// Bulk путь до json_file
type Bulk struct {
	PathJSON string `json:"path_json"`
}

// BulkAcceptOrders принять заказы из файла
func (h *Handler) BulkAcceptOrders(w http.ResponseWriter, r *http.Request) {
	var req Bulk
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if req.PathJSON == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	n, err := h.acceptServ.BulkAcceptOrders(r.Context(), req.PathJSON)
	if errors.Is(err, errs.ErrFileReadError) {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if err != nil {
		w.WriteHeader(GetStatusCode(err))

		return
	}

	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf("принято %d заказов", n)
	if _, err := w.Write([]byte(response)); err != nil {
		log.Printf("ошибка при отправке ответа: %v\n", err)
	}
}
