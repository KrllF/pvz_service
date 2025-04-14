//go:build unit

package accept

import (
	"context"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/accept/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestAcceptOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args entity.OrderEntry
		mock func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager)
		want error
	}{
		{
			name: "all good",
			args: entity.OrderEntry{
				OrderID:   int64(1),
				UserID:    int64(1),
				Weight:    100,
				Price:     100,
				PackType:  "bag",
				ExtraPack: "",
				Shelflife: time.Now().Add(24 * time.Hour),
			},
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockCache := mocks.NewMockCache(ctrl)
				mockCache.EXPECT().Put(int64(1), gomock.Any()).Return(nil)

				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetOrderInfo(context.Background(), int64(1)).Return(models.Order{}, errs.ErrOrderNotFound)
				mockRepo.EXPECT().UserExist(context.Background(), int64(1)).Return(false, nil)
				mockRepo.EXPECT().AddUser(context.Background(), int64(1)).Return(nil)
				mockRepo.EXPECT().AddOrder(context.Background(), gomock.Any()).Return(nil)

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
			name: "order exists",
			args: entity.OrderEntry{
				OrderID:   1,
				UserID:    1,
				Weight:    100,
				Price:     100,
				PackType:  "bag",
				ExtraPack: "",
				Shelflife: time.Now().Add(24 * time.Hour),
			},
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetOrderInfo(context.Background(), int64(1)).Return(models.Order{}, nil)

				return nil, mockRepo, nil
			},
			want: errs.ErrInvalidData,
		},
		{
			name: "bad shelf life",
			args: entity.OrderEntry{
				OrderID:   1,
				UserID:    1,
				Weight:    100,
				Price:     100,
				PackType:  "bag",
				ExtraPack: "",
				Shelflife: time.Now().Add(-1 * time.Hour),
			},
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().GetOrderInfo(context.Background(), int64(1)).Return(models.Order{}, errs.ErrOrderNotFound)

				return nil, mockRepo, nil
			},
			want: errs.ErrInvalidData,
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

			assert.Equal(t, tt.want, service.AcceptOrder(context.Background(), tt.args))
		})
	}
}
