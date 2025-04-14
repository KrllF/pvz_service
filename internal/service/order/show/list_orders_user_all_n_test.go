//go:build unit

package show

import (
	"context"
	"errors"
	"testing"

	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/show/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestListOrdersUserAllN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		userID      int64
		opt         int64
		limit       int64
		page        int64
		pag         bool
		finder      bool
		orderFinder string
		mock        func(ctrl *gomock.Controller) (*mocks.MockRepository, error)
		wantOrders  []models.Order
		wantErr     error
	}{
		{
			name:        "all good with opt < 0 and pagination",
			userID:      1,
			opt:         -1,
			limit:       10,
			page:        1,
			pag:         true,
			finder:      false,
			orderFinder: "",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]models.Order{
						{OrderID: 1, UserID: 1},
						{OrderID: 2, UserID: 1},
					}, nil)
				return mockRepo, nil
			},
			wantOrders: []models.Order{
				{OrderID: 1, UserID: 1},
				{OrderID: 2, UserID: 1},
			},
			wantErr: nil,
		},
		{
			name:        "all good with opt < 0 without pagination",
			userID:      1,
			opt:         -1,
			limit:       0,
			page:        0,
			pag:         false,
			finder:      false,
			orderFinder: "",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(gomock.Any(),
						gomock.Any(),
					).
					Return([]models.Order{
						{OrderID: 1, UserID: 1},
						{OrderID: 2, UserID: 1},
					}, nil)
				return mockRepo, nil
			},
			wantOrders: []models.Order{
				{OrderID: 1, UserID: 1},
				{OrderID: 2, UserID: 1},
			},
			wantErr: nil,
		},
		{
			name:        "all good with opt >= 0 and finder",
			userID:      1,
			opt:         5,
			limit:       0,
			page:        0,
			pag:         false,
			finder:      true,
			orderFinder: "123",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(gomock.Any(),
						gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]models.Order{
						{OrderID: 123, UserID: 1},
					}, nil)
				return mockRepo, nil
			},
			wantOrders: []models.Order{
				{OrderID: 123, UserID: 1},
			},
			wantErr: nil,
		},
		{
			name:        "empty list",
			userID:      1,
			opt:         5,
			limit:       0,
			page:        0,
			pag:         false,
			finder:      false,
			orderFinder: "",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return([]models.Order{}, nil)
				return mockRepo, nil
			},
			wantOrders: []models.Order{},
			wantErr:    nil,
		},
		{
			name:        "repository error",
			userID:      1,
			opt:         5,
			limit:       0,
			page:        0,
			pag:         false,
			finder:      false,
			orderFinder: "",
			mock: func(ctrl *gomock.Controller) (*mocks.MockRepository, error) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					ListOrders(gomock.Any(),
						gomock.Any(),
						gomock.Any(),
					).
					Return(nil, errors.New("repository error"))
				return mockRepo, nil
			},
			wantOrders: nil,
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

			orders, err := service.ListOrdersUserAllN(context.Background(), tt.userID, tt.opt, tt.limit, tt.page, tt.pag, tt.finder, tt.orderFinder)

			assert.Equal(t, tt.wantOrders, orders)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
