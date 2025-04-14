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
	limitReturns       = int64(10)
	pageReturns        = int64(1)
	finderReturns      = false
	orderFinderReturns = ""
)

func TestListReturns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      *desc.ListReturnsRequest
		mock     func(ctrl *gomock.Controller) *mocks.MockShowService
		wantResp *desc.ListReturnsResponse
		wantErr  bool
	}{
		{
			name: "all good",
			req: &desc.ListReturnsRequest{
				Limit:       &limitReturns,
				Page:        &pageReturns,
				Finder:      &finderReturns,
				OrderFinder: &orderFinderReturns,
			},
			mock: func(ctrl *gomock.Controller) *mocks.MockShowService {
				mockShow := mocks.NewMockShowService(ctrl)
				mockShow.EXPECT().ListReturns(
					context.Background(),
					int64(10),
					int64(1),
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
					1, // Count
					nil,
				)
				return mockShow
			},
			wantResp: &desc.ListReturnsResponse{
				Count: 1,
				Returns: []*desc.Order{
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
			wantErr: false,
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

			resp, err := hand.ListReturns(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantResp.Count, resp.Count)
				assert.Equal(t, tt.wantResp.Returns[0].GetOrderid(), resp.Returns[0].GetOrderid())
				assert.Equal(t, tt.wantResp.Returns[0].GetUserid(), resp.Returns[0].GetUserid())
				assert.Equal(t, tt.wantResp.Returns[0].GetWeight(), resp.Returns[0].GetWeight())
				assert.Equal(t, tt.wantResp.Returns[0].GetPrice(), resp.Returns[0].GetPrice())
				assert.Equal(t, tt.wantResp.Returns[0].Packaging.GetPacktype(), resp.Returns[0].Packaging.GetPacktype())
				assert.Equal(t, tt.wantResp.Returns[0].GetInstoragefrom(), resp.Returns[0].GetInstoragefrom())
				assert.Equal(t, tt.wantResp.Returns[0].GetStatus(), resp.Returns[0].GetStatus())
			}
		})
	}
}
