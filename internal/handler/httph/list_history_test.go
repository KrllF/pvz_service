//go:build unit

package httph

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListHistory(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string

		mock          func(ctrl *gomock.Controller) *mocks.MockShowService
		wantResStatus int
	}{
		{
			name: "all ok",
			mock: func(ctrl *gomock.Controller) *mocks.MockShowService {
				mockShow := mocks.NewMockShowService(ctrl)
				mockShow.EXPECT().ListHistory(context.Background()).Return([]models.Order{{
					OrderID:       1,
					UserID:        1,
					Weight:        1,
					Price:         1,
					Packaging:     entity.Pack{PackType: "box", ExtraPack: ""},
					InStorageFrom: time.Time{},
					ShelfLife:     time.Time{},
					Status:        "accepted",
					TwoDaysOfLife: time.Time{},
					LastUpdate:    time.Time{},
				}}, nil)
				return mockShow
			},
			wantResStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockShow := tt.mock(ctrl)
			serv := NewHandler(context.Background(), nil, nil, mockShow, nil)
			req := httptest.NewRequest(http.MethodGet, "/history/", nil)
			res := httptest.NewRecorder()
			serv.ListHistory(res, req)
			assert.Equal(t, tt.wantResStatus, res.Code)
		})
	}
}
