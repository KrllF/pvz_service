package grpc

import (
	"errors"

	"github.com/KrllF/pvz_service/internal/errs"
	"google.golang.org/grpc/codes"
)

// GetGRPCStatusCode получить GRPC статус-код по ошибке
func GetGRPCStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, errs.ErrPackTypeNotFound),
		errors.Is(err, errs.ErrStatusNotFound),
		errors.Is(err, errs.ErrExtraPackNotFound),
		errors.Is(err, errs.ErrOrderNotFound),
		errors.Is(err, errs.ErrUserNotFound),
		errors.Is(err, errs.ErrFileNotFound):
		return codes.NotFound

	case errors.Is(err, errs.ErrInvalidData),
		errors.Is(err, errs.ErrInvalidSortField),
		errors.Is(err, errs.ErrInvalidQueryParams),
		errors.Is(err, errs.ErrInvalidOrderStatus):
		return codes.InvalidArgument

	case errors.Is(err, errs.ErrTimeNotExpired),
		errors.Is(err, errs.ErrShelfLifeExpired):
		return codes.FailedPrecondition

	default:
		return codes.Internal
	}
}
