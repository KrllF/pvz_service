package interceptor

import (
	"context"
	"encoding/json"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Aud интерфейс сервиса аудит лога
type Aud interface {
	StoreStatus(h entity.UpdateStatus)
	StoreHTTP(h entity.HTTPInfo)
}

// LoggingInterceptor  interceptor для логирования
func LoggingInterceptor(aud Aud) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		reqJSON, _ := json.Marshal(req)

		resp, err := handler(ctx, req)

		grpcInfo := entity.HTTPInfo{
			Method:       info.FullMethod,
			Request:      string(reqJSON),
			ResponseCode: int(codes.OK),
		}
		if err != nil {
			grpcInfo.ResponseCode = int(status.Code(err))
		}
		aud.StoreHTTP(grpcInfo)
		if err == nil {
			logStatusChanges(req, aud)
		}

		return resp, err
	}
}

func logStatusChanges(req interface{}, aud Aud) {
	switch r := req.(type) {
	case *desc.ProcessClientRequest:
		switch r.GetOperation().String() {
		case consts.IssueAction:
			for _, orderID := range r.GetOrderIds() {
				aud.StoreStatus(entity.UpdateStatus{OrderID: orderID, NewStatus: consts.AcceptedSt})
			}
		case consts.ReturnAction:
			for _, orderID := range r.GetOrderIds() {
				aud.StoreStatus(entity.UpdateStatus{OrderID: orderID, NewStatus: consts.ReturnedSt})
			}
		}
	default:
	}
}
