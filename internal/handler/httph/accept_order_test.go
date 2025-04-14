//go:build unit

package httph

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAcceptOrder_validation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args entity.OrderEntry
		want bool
	}{
		{
			name: "wrong pack_type",
			args: entity.OrderEntry{
				OrderID: 1, UserID: 1,
				Weight: 1, Price: 1,
				PackType: "", ExtraPack: "",
				Shelflife: time.Now(),
			},
			want: false,
		},
		{
			name: "wrong price",
			args: entity.OrderEntry{
				OrderID: 1, UserID: 1,
				Weight: 1, Price: -1,
				PackType: "box", ExtraPack: "",
				Shelflife: time.Now(),
			},
			want: false,
		},
		{
			name: "wrong weight",
			args: entity.OrderEntry{
				OrderID: 1, UserID: 1,
				Weight: -1, Price: 1,
				PackType: "bag", ExtraPack: "",
				Shelflife: time.Now(),
			},
			want: false,
		},
		{
			name: "wrong user_id",
			args: entity.OrderEntry{
				OrderID: 1, UserID: -1,
				Weight: 1, Price: 1,
				PackType: "film", ExtraPack: "",
				Shelflife: time.Now(),
			},
			want: false,
		},
		{
			name: "wrong order_id",
			args: entity.OrderEntry{
				OrderID: -1, UserID: 1,
				Weight: 1, Price: 1,
				PackType: "bag", ExtraPack: "film",
				Shelflife: time.Now(),
			},
			want: false,
		},
		{
			name: "all good",
			args: entity.OrderEntry{
				OrderID: 1, UserID: 1,
				Weight: 1, Price: 1,
				PackType: "bag", ExtraPack: "",
				Shelflife: time.Now(),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, ValidateReq(tt.args))
		})
	}
}

func TestAcceptOrder(t *testing.T) {
	t.Parallel()

	validShelflife, err := time.Parse(time.RFC3339, "2026-10-10T15:15:10Z")
	if err != nil {
		t.Fatalf("failed to parse valid Shelflife: %v", err)
	}

	tests := []struct {
		name          string
		orderEntry    entity.OrderEntry
		mock          func(ctrl *gomock.Controller) *mocks.MockAcceptService
		wantResStatus int
		wantResBody   []byte
	}{
		{
			name: "nil",
			orderEntry: entity.OrderEntry{
				OrderID:   1,
				UserID:    2,
				Weight:    100,
				Price:     200,
				PackType:  "box",
				ExtraPack: "",
				Shelflife: validShelflife,
			},
			mock: func(ctrl *gomock.Controller) *mocks.MockAcceptService {
				mockAcc := mocks.NewMockAcceptService(ctrl)
				mockAcc.EXPECT().AcceptOrder(context.Background(), entity.OrderEntry{
					OrderID:   1,
					UserID:    2,
					Weight:    100,
					Price:     200,
					PackType:  "box",
					ExtraPack: "",
					Shelflife: validShelflife,
				}).Return(nil)
				return mockAcc
			},
			wantResStatus: http.StatusCreated,
			wantResBody:   []byte("Заказ успешно принят"),
		},
		{
			name: "bad pack_type",
			orderEntry: entity.OrderEntry{
				OrderID:   1,
				UserID:    2,
				Weight:    100,
				Price:     200,
				PackType:  "",
				ExtraPack: "",
				Shelflife: validShelflife,
			},
			mock:          func(ctrl *gomock.Controller) *mocks.MockAcceptService { return nil },
			wantResStatus: http.StatusBadRequest,
			wantResBody:   []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAcc := tt.mock(ctrl)
			serv := NewHandler(context.Background(), mockAcc, nil, nil, nil)
			reqBody, err := json.Marshal(tt.orderEntry)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/accept", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			serv.AcceptOrder(rec, req)
			assert.Equal(t, tt.wantResStatus, rec.Code)
			if len(tt.wantResBody) != 0 {
				assert.Equal(t, tt.wantResBody, rec.Body.Bytes())
			}
		})
	}
}
