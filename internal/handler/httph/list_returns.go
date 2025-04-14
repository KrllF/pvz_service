package httph

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/KrllF/pvz_service/internal/models"
)

// ListReturnsResponse ответ на ListReturns
type ListReturnsResponse struct {
	Count   int64          `json:"count"`
	Returns []models.Order `json:"returns"`
}

// ListReturns получить список возвратов
func (h *Handler) ListReturns(w http.ResponseWriter, r *http.Request) {
	limit, page, err := parseLimitPage(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	finder, orderFinder, err := parseFinderOrder(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	returns, count, err := h.showServ.ListReturns(r.Context(), limit, page, finder, orderFinder)
	if err != nil {
		w.WriteHeader(GetStatusCode(err))

		return
	}

	response := ListReturnsResponse{
		Count:   int64(count),
		Returns: returns,
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

func parseLimitPage(queryParams url.Values) (int64, int64, error) {
	limitStr := queryParams.Get("limit")
	limit := int64(1)
	if limitStr != "" {
		var err error
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit < 1 {
			return 0, 0, errors.New("некорректный параметр 'limit'")
		}
	}

	pageStr := queryParams.Get("page")
	page := int64(1)
	if pageStr != "" {
		var err error
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil || page < 1 {
			return 0, 0, errors.New("некорректный параметр 'page'")
		}
	}

	return limit, page, nil
}
