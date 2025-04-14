//go:build unit

package httph

import (
	"errors"
	"net/http"
	"testing"

	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/stretchr/testify/assert"
)

func TestGetStatusCode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args error
		want int
	}{
		{
			name: "ErrInvalidData",
			args: errs.ErrInvalidData,
			want: http.StatusBadRequest,
		},
		{
			name: "ErrInvalidSortField",
			args: errs.ErrInvalidSortField,
			want: http.StatusBadRequest,
		},
		{
			name: "ErrInvalidQueryParams",
			args: errs.ErrInvalidQueryParams,
			want: http.StatusBadRequest,
		},
		{
			name: "ErrInvalidOrderStatus",
			args: errs.ErrInvalidOrderStatus,
			want: http.StatusBadRequest,
		},

		{
			name: "ErrTimeNotExpired",
			args: errs.ErrTimeNotExpired,
			want: http.StatusConflict,
		},
		{
			name: "ErrShelfLifeExpired",
			args: errs.ErrShelfLifeExpired,
			want: http.StatusConflict,
		},

		{
			name: "ErrPackTypeNotFound",
			args: errs.ErrPackTypeNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "ErrStatusNotFound",
			args: errs.ErrStatusNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "ErrExtraPackNotFound",
			args: errs.ErrExtraPackNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "ErrOrderNotFound",
			args: errs.ErrOrderNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "ErrUserNotFound",
			args: errs.ErrUserNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "ErrFileNotFound",
			args: errs.ErrFileNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "default error",
			args: errors.New("goga"),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, GetStatusCode(tt.args))
		})
	}
}
