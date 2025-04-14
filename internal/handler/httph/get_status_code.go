package httph

import (
	"errors"
	"net/http"

	"github.com/KrllF/pvz_service/internal/errs"
)

// GetStatusCode получить статус код, если случилась ошибка
func GetStatusCode(err error) int {
	switch {
	case errors.Is(err, errs.ErrPackTypeNotFound),
		errors.Is(err, errs.ErrStatusNotFound),
		errors.Is(err, errs.ErrExtraPackNotFound),
		errors.Is(err, errs.ErrOrderNotFound),
		errors.Is(err, errs.ErrUserNotFound),
		errors.Is(err, errs.ErrFileNotFound):
		return http.StatusNotFound

	case errors.Is(err, errs.ErrInvalidData),
		errors.Is(err, errs.ErrInvalidSortField),
		errors.Is(err, errs.ErrInvalidQueryParams),
		errors.Is(err, errs.ErrInvalidOrderStatus):
		return http.StatusBadRequest

	case errors.Is(err, errs.ErrTimeNotExpired),
		errors.Is(err, errs.ErrShelfLifeExpired):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
