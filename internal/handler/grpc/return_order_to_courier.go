package grpc

import (
	"context"

	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReturnOrderToCourier вернуть заказы курьеру
func (h *Handler) ReturnOrderToCourier(ctx context.Context,
	req *desc.ReturnOrderToCourierRequest,
) (*desc.ReturnOrderToCourierResponse, error) {
	h.logger.Info("Processing request ReturnOrderToCourier")

	orderID := req.GetOrderId()
	if orderID <= 0 {
		h.logger.Warn("validation bad")

		return nil, status.Error(codes.InvalidArgument, "validation bad")
	}
	err := h.returnsServ.ReturnOrder(ctx, orderID)
	if err != nil {
		h.logger.Error("h.returnsServ.ReturnOrder", zap.Error(err))

		return nil, status.Errorf(GetGRPCStatusCode(err), "h.returnsServ.ReturnOrder: %v", err)
	}
	h.logger.Info("result ReturnOrderToCourier OK")

	return &desc.ReturnOrderToCourierResponse{Success: true}, nil
}
