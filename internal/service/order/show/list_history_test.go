//go:build unit

package show

import (
	"context"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/show/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestListHistory(t *testing.T) {
	t.Parallel()
	fixedTime := time.Now()

	tests := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository)
		wantOrders []models.Order
		wantErr    error
	}{
		{
			name: "all good",
			mock: func(ctrl *gomock.Controller) (*mocks.MockCache, *mocks.MockRepository) {
				mockCache := mocks.NewMockCache(ctrl)
				mockCache.EXPECT().ListItems(context.Background()).
					Return([]models.Order{
						{LastUpdate: fixedTime.Add(time.Hour)},
						{LastUpdate: fixedTime},
					}, nil)
				return mockCache, nil
			},
			wantOrders: []models.Order{
				{LastUpdate: fixedTime},
				{LastUpdate: fixedTime.Add(time.Hour)},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache, mockRepo := tt.mock(ctrl)
			logg, _ := logger.NewLogger(zapcore.DebugLevel)
			service := &Service{
				repository: mockRepo,
				cache:      mockCache,
				logger:     logg,
			}
			ords, err := service.ListHistory(context.Background())
			assert.Equal(t, tt.wantOrders, ords)
			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
