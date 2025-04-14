//go:build unit

package accept

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/accept/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

var timeShelf = time.Date(2030, 10, 10, 10, 10, 10, 10, time.UTC)

func TestProcessClientIssue(t *testing.T) {
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
					UpdateStatus(context.Background(), int64(1), consts.AcceptedSt).
					Return(nil)
				mockRepo.EXPECT().
					UpdateStatus(context.Background(), int64(2), consts.AcceptedSt).
					Return(nil)

				mockRepo.EXPECT().
					SetTwoDaysOfLife(context.Background(), int64(1)).
					Return(nil)
				mockRepo.EXPECT().
					SetTwoDaysOfLife(context.Background(), int64(2)).
					Return(nil)

				mockRepo.EXPECT().
					GetOrderInfo(context.Background(), int64(1)).
					Return(models.Order{
						UserID: int64(1), ShelfLife: timeShelf,
						Status: consts.DeliveredST,
					}, nil)
				mockRepo.EXPECT().
					GetOrderInfo(context.Background(), int64(2)).
					Return(models.Order{
						UserID: int64(1), ShelfLife: timeShelf,
						Status: consts.DeliveredST,
					}, nil)

				mockTX := mocks.NewMockTXManager(ctrl)
				mockTX.EXPECT().
					RunSerializable(context.Background(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(ctxTX context.Context) error) error {
						return fn(ctx)
					}).Times(2)

				mockCache := mocks.NewMockCache(ctrl)
				mockCache.EXPECT().
					GetItem(int64(1)).
					Return(models.Order{OrderID: 1, Status: consts.DeliveredST, ShelfLife: timeShelf}, nil)
				mockCache.EXPECT().
					GetItem(int64(2)).
					Return(models.Order{OrderID: 2, Status: consts.DeliveredST, ShelfLife: timeShelf}, nil)
				mockCache.EXPECT().
					Put(int64(1), gomock.Any()).
					Return(nil)
				mockCache.EXPECT().
					Put(int64(2), gomock.Any()).
					Return(nil)

				return mockCache, mockRepo, mockTX
			},
			wantNums: []int64{1, 2},
			wantErr:  nil,
		},
		{
			name:     "user not found",
			userID:   999,
			orderIDs: []int64{1, 2, 3},
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					UserExist(context.Background(), int64(999)).
					Return(false, nil)

				return nil, mockRepo, nil
			},
			wantNums: nil,
			wantErr:  errs.ErrUserNotFound,
		},
		{
			name:     "bad UserExist",
			userID:   12,
			orderIDs: []int64{1, 2, 3},
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository, *mocks.MockTXManager) {
				mockRepo := mocks.NewMockRepository(ctrl)
				mockRepo.EXPECT().
					UserExist(context.Background(), int64(12)).
					Return(false, errors.New("error"))

				return nil, mockRepo, nil
			},
			wantNums: nil,
			wantErr:  errors.New("error"),
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

			nums, err := service.ProcessClientIssue(context.Background(), tt.userID, tt.orderIDs)
			assert.Equal(t, tt.wantNums, nums)
			if tt.wantErr != nil {
				assert.ErrorContains(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
