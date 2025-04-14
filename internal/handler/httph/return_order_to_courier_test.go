//go:build unit

package httph

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestReturnOrderToCourier_validateReturnToCourier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    map[string]string
		wantID  int64
		wantErr error
	}{
		{
			name: "all good",
			args: map[string]string{
				"id": "1",
			},
			wantID:  1,
			wantErr: nil,
		},
		{
			name:    "without id",
			args:    map[string]string{},
			wantID:  0,
			wantErr: errors.New("нет параметра order_id"),
		},
		{
			name: "bad id",
			args: map[string]string{
				"id": "f",
			},
			wantID:  0,
			wantErr: strconv.ErrSyntax,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			orderID, err := validateReturnToCourier(tt.args)
			assert.Equal(t, tt.wantID, orderID)

			if tt.wantErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReturnOrderToCourier(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) *mocks.MockReturnsService
		url        string
		wantStatus int
	}{
		{
			name: "all good",
			mock: func(ctrl *gomock.Controller) *mocks.MockReturnsService {
				mockRet := mocks.NewMockReturnsService(ctrl)
				// тут не убрал, потому что gorilla/mux обрабатывает query
				mockRet.EXPECT().ReturnOrder(gomock.Any(), int64(1)).Return(nil)
				return mockRet
			},
			url:        "/return/1",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRet := tt.mock(ctrl)

			router := mux.NewRouter()
			serv := NewHandler(context.Background(), nil, mockRet, nil, nil)
			router.HandleFunc("/return/{id:[0-9]+}", serv.ReturnOrderToCourier).Methods(http.MethodDelete)

			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			assert.Equal(t, tt.wantStatus, res.Code)
		})
	}
}
