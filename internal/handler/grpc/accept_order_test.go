//go:build unit

package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/KrllF/pvz_service/pkg/logger"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var timeDate = time.Date(2026, time.April, 10, 10, 10, 10, 0, time.UTC)

func TestAcceptOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) *mocks.MockAcceptService
		req      *desc.AcceptRequest
		wantResp *desc.AcceptResponse
		wantErr  error
	}{
		{
			name: "all good",
			mock: func(ctrl *gomock.Controller) *mocks.MockAcceptService {
				mockAcc := mocks.NewMockAcceptService(ctrl)
				mockAcc.EXPECT().AcceptOrder(context.Background(), entity.OrderEntry{
					OrderID:   1,
					UserID:    1,
					Weight:    100,
					Price:     100,
					PackType:  entity.PackTypeFilm,
					ExtraPack: entity.PackTypeFilm,
					Shelflife: timeDate,
				}).Return(nil)

				return mockAcc
			},
			req: &desc.AcceptRequest{
				OrderId:   1,
				UserId:    1,
				Weight:    100,
				Price:     100,
				PackType:  desc.PackType_film,
				ExtraPack: desc.PackType_film.Enum(),
				ShelfLife: timestamppb.New(timeDate),
			},
			wantResp: &desc.AcceptResponse{Success: true},
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAcc := tt.mock(ctrl)
			logg, _ := logger.NewLogger(zapcore.DebugLevel)
			hand := Handler{desc.UnimplementedOrderV1Server{}, mockAcc, nil, nil, logg}
			resp, err := hand.AcceptOrder(context.Background(), tt.req)
			assert.Equal(t, tt.wantResp.GetSuccess(), resp.GetSuccess())
			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
