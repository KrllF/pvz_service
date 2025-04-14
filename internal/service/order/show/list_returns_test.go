//go:build unit

package show

import (
	"context"
	"errors"
	"testing"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/show/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestListReturns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		limit       int64
		page        int64
		finder      bool
		orderFinder string
		mock        func(ctrl *gomock.Controller) (*mocks.MockRepository, error)
		wantOrders  []models.Order
		wantCount   int
		wantErr     error
	}{
		{
			name:        "all good without finder",
			limit:       10,
			page:        1,
			finder:      false,
			orderFinder: "",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(context.Background(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]models.Order{
						{OrderID: 1, Status: consts.ReturnedSt},
						{OrderID: 2, Status: consts.ReturnedSt},
					}, nil)
				return mockRepo, nil
			},
			wantOrders: []models.Order{
				{OrderID: 1, Status: consts.ReturnedSt},
				{OrderID: 2, Status: consts.ReturnedSt},
			},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name:        "all good with finder",
			limit:       10,
			page:        1,
			finder:      true,
			orderFinder: "123",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(context.Background(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]models.Order{
						{OrderID: 123, Status: consts.ReturnedSt},
					}, nil)
				return mockRepo, nil
			},
			wantOrders: []models.Order{
				{OrderID: 123, Status: consts.ReturnedSt},
			},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name:        "empty list",
			limit:       10,
			page:        1,
			finder:      false,
			orderFinder: "",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(context.Background(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]models.Order{}, nil)
				return mockRepo, nil
			},
			wantOrders: []models.Order{},
			wantCount:  0,
			wantErr:    nil,
		},
		{
			name:        "repository error",
			limit:       10,
			page:        1,
			finder:      false,
			orderFinder: "",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(context.Background(),
						gomock.Any(),
						gomock.Any(),
					).
					Return(nil, errors.New("repository error"))
				return mockRepo, nil
			},
			wantOrders: nil,
			wantCount:  0,
			wantErr:    errors.New("ошибка при получении списка заказов: repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo, _ := tt.mock(ctrl)

			logg, _ := logger.NewLogger(zapcore.DebugLevel)
			service := &Service{
				repository: mockRepo,
				logger:     logg,
			}

			orders, count, err := service.ListReturns(context.Background(), tt.limit, tt.page, tt.finder, tt.orderFinder)

			assert.Equal(t, tt.wantOrders, orders)
			assert.Equal(t, tt.wantCount, count)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
