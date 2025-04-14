package grpc

import (
	"context"

	"github.com/KrllF/pvz_service/internal/converter"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	retSize = 1000
)

// ListHistory получить историю заказов
func (h *Handler) ListHistory(ctx context.Context, _ *emptypb.Empty) (*desc.ListHistoryResponse, error) {
	h.logger.Info("Processing request ListHistory")

	history, err := h.showServ.ListHistory(ctx)
	if err != nil {
		h.logger.Error("h.showServ.ListHistory", zap.Error(err))

		return nil, status.Errorf(GetGRPCStatusCode(err), "h.showServ.ListHistory: %v", err)
	}
	ret := make([]*desc.Order, 0, retSize)
	for _, val := range history {
		ret = append(ret, converter.ConvertToProtoOrder(&val))
	}
	h.logger.Info("result ListHistory OK")

	return &desc.ListHistoryResponse{Orders: ret}, nil
}
