package grpc

import (
	"context"

	"github.com/KrllF/pvz_service/internal/converter"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListReturns получить список возвратов
func (h *Handler) ListReturns(ctx context.Context,
	req *desc.ListReturnsRequest,
) (*desc.ListReturnsResponse, error) {
	h.logger.Info("Processing request ListReturns")

	limit, page := parseLimitPage(req)
	finder, orderFinder, err := parseFinderOrder(req)
	if err != nil {
		h.logger.Warn("parseLimitPage", zap.Error(err))

		return nil, status.Errorf(codes.InvalidArgument, "parseLimitPage: %v", err)
	}
	returns, count, err := h.showServ.ListReturns(ctx, limit, page, finder, orderFinder)
	if err != nil {
		h.logger.Error("h.showServ.ListReturns", zap.Error(err))

		return nil, status.Errorf(codes.InvalidArgument, "h.showServ.ListReturns: %v", err)
	}
	conv := make([]*desc.Order, 0, count)
	for _, val := range returns {
		conv = append(conv, converter.ConvertToProtoOrder(&val))
	}
	h.logger.Info("result ListReturns OK")

	return &desc.ListReturnsResponse{Count: int64(count), Returns: conv}, nil
}

func parseLimitPage(req *desc.ListReturnsRequest) (int64, int64) {
	limitStr := req.GetLimit()
	pageStr := req.GetPage()
	limit := int64(1)
	page := int64(1)
	if limitStr <= 0 {
		limitStr = limit
	}
	if page <= 0 {
		pageStr = page
	}

	return limitStr, pageStr
}
