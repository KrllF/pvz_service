//go:build unit

package httph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBulkAcceptOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		path          Bulk
		mock          func(ctrl *gomock.Controller) *mocks.MockAcceptService
		wantResStatus int
		wantResBody   []byte
	}{
		{
			name: "kot",
			path: Bulk{PathJSON: "asdfa/fasdfadsf/asdf"},
			mock: func(ctrl *gomock.Controller) *mocks.MockAcceptService {
				mockAcc := mocks.NewMockAcceptService(ctrl)
				mockAcc.EXPECT().BulkAcceptOrders(context.Background(), "asdfa/fasdfadsf/asdf").Return(4, nil)
				return mockAcc
			},
			wantResStatus: http.StatusOK,
			wantResBody:   []byte(fmt.Sprintf("принято %d заказов", 4)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAcc := tt.mock(ctrl)
			serv := NewHandler(context.Background(), mockAcc, nil, nil, nil)
			reqBody, err := json.Marshal(tt.path)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/accept/bulk", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()
			serv.BulkAcceptOrders(res, req)
			assert.Equal(t, tt.wantResStatus, res.Code)
			assert.Equal(t, tt.wantResBody, res.Body.Bytes())
		})
	}
}
