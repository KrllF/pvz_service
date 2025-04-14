//go:build unit

package httph

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListReturns_parseLimitPage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		args     url.Values
		wantlim  int64
		wantpage int64
		err      error
	}{
		{
			name: "bad limit",
			args: url.Values{
				"limit": []string{"-1"},
				"page":  []string{"1"},
			},
			wantlim:  0,
			wantpage: 0,
			err:      errors.New("некорректный параметр 'limit'"),
		},
		{
			name: "bad page",
			args: url.Values{
				"limit": []string{"1"},
				"page":  []string{"-1"},
			},
			wantlim:  0,
			wantpage: 0,
			err:      errors.New("некорректный параметр 'page'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			limit, page, err := parseLimitPage(tt.args)
			assert.Equal(t, tt.wantlim, limit)
			assert.Equal(t, tt.wantpage, page)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestPvzServerListHistory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mock          func(ctrl *gomock.Controller) *mocks.MockShowService
		wantResStatus int
	}{
		{
			name: "all good",
			mock: func(ctrl *gomock.Controller) *mocks.MockShowService {
				mockShow := mocks.NewMockShowService(ctrl)
				mockShow.EXPECT().ListReturns(context.Background(), int64(1), int64(1), false, "").Return([]models.Order{
					{
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
					},
				}, 1, nil)
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
			req := httptest.NewRequest(http.MethodGet, "/history/returns", nil)
			res := httptest.NewRecorder()
			serv.ListReturns(res, req)
			assert.Equal(t, tt.wantResStatus, res.Code)
		})
	}
}
