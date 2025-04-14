//go:build unit

package grpc

import (
	"context"
	"testing"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/handler/mocks"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/pkg/logger"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	N      = int64(1)
	pvz    = true
	limit  = int64(1)
	page   = int64(1)
	pag    = true
	find   = false
	finers = ""
)

func TestListOrdersUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      *desc.ListOrdersUserRequest
		mock     func(ctrl *gomock.Controller) *mocks.MockShowService
		wantResp *desc.ListOrdersUserResponse
		wantErr  error
	}{
		{
			name: "all good",
			req: &desc.ListOrdersUserRequest{
				UserId:      int64(1),
				N:           &N,
				InPvz:       &pvz,
				Limit:       &limit,
				Page:        &page,
				Pag:         &pag,
				Finder:      &find,
				OrderFinder: &finers,
			},
			mock: func(ctrl *gomock.Controller) *mocks.MockShowService {
				mockShow := mocks.NewMockShowService(ctrl)
				mockShow.EXPECT().ListOrdersUserPVZ(
					context.Background(),
					int64(1),
					int64(1),
					int64(1),
					true,
					false,
					"",
				).Return(
					[]models.Order{
						{
							OrderID: 1,
							UserID:  1,
							Weight:  100,
							Price:   100,
							Packaging: entity.Pack{
								PackType: entity.PackTypeFilm,
							},
							InStorageFrom: timeStorage,
							Status:        consts.DeliveredST,
						},
					},
					nil,
				)
				return mockShow
			},
			wantResp: &desc.ListOrdersUserResponse{
				Orders: []*desc.Order{
					{
						Orderid: 1,
						Userid:  1,
						Weight:  100,
						Price:   100,
						Packaging: &desc.Pack{
							Packtype: desc.PackType_film,
						},
						Instoragefrom: timestamppb.New(timeStorage),
						Status:        desc.OrderStatus_delivered.String(),
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockShow := tt.mock(ctrl)
			logg, _ := logger.NewLogger(zapcore.DebugLevel)
			hand := Handler{desc.UnimplementedOrderV1Server{}, nil, nil, mockShow, logg}

			resp, err := hand.ListOrdersUser(context.Background(), tt.req)

			assert.Equal(t, tt.wantResp.Orders[0].GetOrderid(), resp.Orders[0].GetOrderid())
			assert.Equal(t, tt.wantResp.Orders[0].Userid, resp.Orders[0].GetUserid())
			assert.Equal(t, tt.wantResp.Orders[0].Weight, resp.Orders[0].GetWeight())
			assert.Equal(t, tt.wantResp.Orders[0].Price, resp.Orders[0].GetPrice())
			assert.Equal(t, tt.wantResp.Orders[0].Packaging.GetPacktype(), resp.Orders[0].Packaging.GetPacktype())
			assert.Equal(t, tt.wantResp.Orders[0].GetInstoragefrom(), resp.Orders[0].GetInstoragefrom())
			assert.Equal(t, tt.wantResp.Orders[0].GetStatus(), resp.Orders[0].GetStatus())

			if tt.wantErr != nil {
				assert.ErrorIs(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
