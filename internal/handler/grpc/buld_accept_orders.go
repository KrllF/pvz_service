package grpc

import (
	"context"

	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

// BulkAcceptOrders принять заказы из файла
func (h *Handler) BulkAcceptOrders(ctx context.Context,
	req *desc.BulkRequest,
) (*desc.BulkResponse, error) {
	h.logger.Info("Processing request BulkAcceptOrders")

	path := req.GetPath()
	n, err := h.acceptServ.BulkAcceptOrders(ctx, path)
	if err != nil {
		h.logger.Error("h.acceptServ.BulkAcceptOrders", zap.Error(err))

		return nil, status.Errorf(GetGRPCStatusCode(err), "h.acceptServ.BulkAcceptOrders: %v", err)
	}
	h.logger.Info("result BulkAcceptOrders OK")

	return &desc.BulkResponse{Count: int64(n)}, nil
}
