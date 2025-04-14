//go:build unit

package returns

import (
	"context"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/returns/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

var timeShelf = time.Date(2022, 10, 10, 10, 10, 10, 10, time.UTC)

func TestReturnOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args int64
		mock func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager)
		want error
	}{
		{
			name: "all good",
			args: int64(1),
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockCache := mocks.NewMockCache(ctrl)
				mockCache.EXPECT().DeleteCache(int64(1)).Return(nil)

				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetOrderInfo(context.Background(), int64(1)).
					Return(models.Order{UserID: int64(1), Status: consts.DeliveredST, ShelfLife: timeShelf}, nil)
				mockRepo.EXPECT().DeleteOrder(context.Background(), int64(1)).Return(nil)
				mockRepo.EXPECT().ListOrders(context.Background(), gomock.Any()).Return(make([]models.Order, 2), nil)
				mockTX := mocks.NewMockTXManager(ctrl)
				mockTX.EXPECT().RunSerializable(context.Background(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(ctxTX context.Context) error) error {
						return fn(ctx)
					},
				)

				return mockCache, mockRepo, mockTX
			},
			want: nil,
		},
		{
			name: "delete user",
			args: int64(1),
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockCache := mocks.NewMockCache(ctrl)
				mockCache.EXPECT().DeleteCache(int64(1)).Return(nil)

				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetOrderInfo(context.Background(), int64(1)).
					Return(models.Order{UserID: int64(1), Status: consts.DeliveredST, ShelfLife: timeShelf}, nil)
				mockRepo.EXPECT().DeleteOrder(context.Background(), int64(1)).Return(nil)
				mockRepo.EXPECT().ListOrders(context.Background(), gomock.Any()).Return(make([]models.Order, 0), nil)
				mockRepo.EXPECT().DeleteUser(context.Background(), int64(1))
				mockTX := mocks.NewMockTXManager(ctrl)
				mockTX.EXPECT().RunSerializable(context.Background(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(ctxTX context.Context) error) error {
						return fn(ctx)
					},
				)

				return mockCache, mockRepo, mockTX
			},
			want: nil,
		},
		{
			name: "bad status",
			args: int64(1),
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetOrderInfo(context.Background(), int64(1)).
					Return(models.Order{UserID: int64(1), Status: consts.AcceptedSt, ShelfLife: timeShelf}, nil)
				mockTX := mocks.NewMockTXManager(ctrl)
				mockTX.EXPECT().RunSerializable(context.Background(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, fn func(ctxTX context.Context) error) error {
						return fn(ctx)
					},
				)

				return nil, mockRepo, mockTX
			},
			want: errs.ErrInvalidOrderStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCache, mockRepo, mockTX := tt.mock(ctrl)
			logg, _ := logger.NewLogger(zapcore.DebugLevel)
			service := &Service{
				repository: mockRepo,
				cache:      mockCache,
				txManager:  mockTX,
				logger:     logg,
			}
			if tt.want != nil {
				assert.ErrorIs(t, service.ReturnOrder(context.Background(), tt.args), tt.want)
			} else {
				assert.NoError(t, service.ReturnOrder(context.Background(), tt.args))
			}
		})
	}
}
