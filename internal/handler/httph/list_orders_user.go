package httph

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/KrllF/pvz_service/internal/models"
	"github.com/gorilla/mux"
)

// ListOrdersResponse ответ на ListOrdersUser
type ListOrdersResponse struct {
	Orders []models.Order `json:"orders"`
}

// ValidateStructToListOrders структура для валидации
type ValidateStructToListOrders struct {
	n           int64
	inPVZ       bool
	limit       int64
	page        int64
	pag         bool
	finder      bool
	orderFinder string
}

// ListOrdersUser получить список заказов клиента
func (h *Handler) ListOrdersUser(w http.ResponseWriter, r *http.Request) {
	key, ok := mux.Vars(r)[queryParamKey]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	userID, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	parseReq, err := validateQueryParamsU(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	var ret []models.Order

	switch parseReq.inPVZ {
	case true:
		ret, err = h.showServ.ListOrdersUserPVZ(r.Context(), userID, parseReq.limit,
			parseReq.page, parseReq.pag, parseReq.finder, parseReq.orderFinder)
	case false:
		if parseReq.n == -1 {
			ret, err = h.showServ.ListOrdersUserAllN(r.Context(), userID, -1, parseReq.limit, parseReq.page,
				parseReq.pag, parseReq.finder, parseReq.orderFinder)
		} else {
			ret, err = h.showServ.ListOrdersUserAllN(r.Context(), userID, parseReq.n,
				parseReq.limit, parseReq.page, parseReq.pag, parseReq.finder, parseReq.orderFinder)
		}
	}

	if err != nil {
		w.WriteHeader(GetStatusCode(err))

		return
	}

	response := ListOrdersResponse{
		Orders: ret,
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

func validateQueryParamsU(queryParams url.Values) (ValidateStructToListOrders, error) {
	n, err := nRet(queryParams)
	if err != nil {
		return ValidateStructToListOrders{}, errors.New("некорректный параметр n")
	}

	inPVZ, err := inPVZRet(queryParams)
	if err != nil {
		return ValidateStructToListOrders{}, errors.New("некорректный параметр in_pvz")
	}

	limit, err := limitRet(queryParams)
	if err != nil {
		return ValidateStructToListOrders{}, fmt.Errorf("некорректный параметр limit: %w", err)
	}

	page, err := pageRet(queryParams)
	if err != nil {
		return ValidateStructToListOrders{}, errors.New("некорректный параметр page")
	}

	pag, err := pagRet(queryParams)
	if err != nil {
		return ValidateStructToListOrders{}, errors.New("некорректный параметр pag")
	}
	finder, orderFinder, err := parseFinderOrder(queryParams)
	if err != nil {
		return ValidateStructToListOrders{}, errors.New("некорректный параметры finder и order")
	}

	ret := ValidateStructToListOrders{
		n:           n,
		inPVZ:       inPVZ,
		limit:       limit,
		page:        page,
		pag:         pag,
		finder:      finder,
		orderFinder: orderFinder,
	}

	return ret, nil
}

func pagRet(q url.Values) (bool, error) {
	var err error
	pag := false
	pagStr := q.Get("pag")
	if pagStr != "" {
		pag, err = strconv.ParseBool(pagStr)
		if err != nil {
			return pag, errors.New("некорректный параметр pag")
		}
	}

	return pag, nil
}

func limitRet(q url.Values) (int64, error) {
	var err error
	limit := int64(-1)
	limitStr := q.Get("limit")
	if limitStr == "" {
		return limit, nil
	}
	limit, err = strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return -1, errors.New("некорректный параметр limit")
	}
	if limit <= 0 {
		return limit, errors.New("отрицательный параметр limit")
	}

	return limit, nil
}

func pageRet(q url.Values) (int64, error) {
	var err error
	page := int64(-1)
	pageStr := q.Get("page")
	if pageStr == "" {
		return page, nil
	}
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		return -1, errors.New("некорректный параметр page")
	}
	if page <= 0 {
		return page, errors.New("отрицательный параметр page")
	}

	return page, nil
}

func nRet(q url.Values) (int64, error) {
	var err error
	n := int64(-1)
	nStr := q.Get("n")
	if nStr == "" {
		return n, nil
	}

	n, err = strconv.ParseInt(nStr, 10, 64)
	if err != nil {
		return -1, errors.New("некорректный параметр n")
	}
	if n <= 0 {
		return n, errors.New("отрицательный параметр n")
	}

	return n, nil
}

func inPVZRet(q url.Values) (bool, error) {
	var err error
	inPVZ := false
	inPVZStr := q.Get("in_pvz")
	if inPVZStr == "" {
		return inPVZ, nil
	}
	inPVZ, err = strconv.ParseBool(inPVZStr)
	if err != nil {
		return inPVZ, errors.New("некорректный параметр in_pvz")
	}

	return inPVZ, nil
}
