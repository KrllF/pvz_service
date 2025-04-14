//go:build unit

package returns

import (
	"context"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/returns/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestProcessClientReturn(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		userID   int64
		orderIDs []int64
		mock     func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager)
		wantNums []int64
		wantErr  error
	}{
		{
			name:     "all good",
			userID:   1,
			orderIDs: []int64{1, 2},
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					UserExist(context.Background(), int64(1)).
					Return(true, nil)
				mockRepo.EXPECT().
					GetOrderInfo(context.Background(), int64(1)).
					Return(models.Order{
						OrderID: 1, UserID: 1, Status: consts.AcceptedSt,
						TwoDaysOfLife: time.Now().Add(24 * time.Hour),
					}, nil)
				mockRepo.EXPECT().
					GetOrderInfo(context.Background(), int64(2)).
					Return(models.Order{
						OrderID: 2, UserID: 1, Status: consts.AcceptedSt,
						TwoDaysOfLife: time.Now().Add(24 * time.Hour),
					}, nil)
				mockRepo.EXPECT().
					UpdateStatus(context.Background(), int64(1), consts.ReturnedSt).
					Return(nil)
				mockRepo.EXPECT().
					UpdateStatus(context.Background(), int64(2), consts.ReturnedSt).
					Return(nil)
				mockTX := mocks.NewMockTXManager(ctrl)
				mockTX.EXPECT().
					RunSerializable(context.Background(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctxTX context.Context) error) error {
						return fn(ctx)
					}).Times(2)

				mockCache := mocks.NewMockCache(ctrl)
				mockCache.EXPECT().DeleteCache(int64(1)).Return(nil)
				mockCache.EXPECT().DeleteCache(int64(2)).Return(nil)
				return mockCache, mockRepo, mockTX
			},
			wantNums: []int64{1, 2},
			wantErr:  nil,
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

			nums, err := service.ProcessClientReturn(context.Background(), tt.userID, tt.orderIDs)
			assert.Equal(t, tt.wantNums, nums)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
