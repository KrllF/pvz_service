package grpc

import (
	"context"
	"errors"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/pkg/monitoring"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProcessClient принять или вернуть заказы клиента
func (h *Handler) ProcessClient(ctx context.Context,
	req *desc.ProcessClientRequest,
) (*desc.ProcessClientResponse, error) {
	h.logger.Info("Processing request ProcessClient")

	if err := ValidateRequest(req); err != nil {
		h.logger.Warn("ValidateRequest", zap.Error(err))

		return nil, status.Errorf(codes.InvalidArgument, "ValidateRequest: %v", err)
	}
	var n []int64
	var err error
	switch req.GetOperation().String() {
	case consts.IssueAction:
		n, err = h.acceptServ.ProcessClientIssue(ctx, req.GetUserId(), req.GetOrderIds())
	case consts.ReturnAction:
		n, err = h.returnsServ.ProcessClientReturn(ctx, req.GetUserId(), req.GetOrderIds())
	default:
		h.logger.Warn("bad operation")

		return nil, status.Errorf(codes.InvalidArgument,
			"bad operation")
	}

	if err != nil && GetGRPCStatusCode(err) == codes.Internal {
		return nil, status.Errorf(GetGRPCStatusCode(err), "h.returnsServ.ProcessClient: %v", err)
	}
	switch req.GetOperation().String() {
	case consts.IssueAction:
		monitoring.SetCountAcceptOrder(int64(len(n)))
	case consts.ReturnAction:
		monitoring.SetCountReturnOrder(int64(len(n)))
	}
	h.logger.Info("result ProcessClient OK")

	return &desc.ProcessClientResponse{OrderIds: n}, nil
}

// ValidateRequest валидация запроса
func ValidateRequest(req *desc.ProcessClientRequest) error {
	if req.GetUserId() <= 0 {
		return errors.New("bad user ID")
	}

	if req.GetOperation().String() == "" {
		return errors.New("bad operation")
	}

	if len(req.GetOrderIds()) == 0 {
		return errors.New("empty order_IDs")
	}

	return nil
}
