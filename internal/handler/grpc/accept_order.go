package grpc

import (
	"context"
	"errors"

	"github.com/KrllF/pvz_service/internal/converter"
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/pkg/monitoring"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AcceptOrder принять заказ
func (h *Handler) AcceptOrder(ctx context.Context, req *desc.AcceptRequest) (*desc.AcceptResponse, error) {
	h.logger.Info("Processing request AcceptOrder")

	var extraPack entity.PackType

	if packValue := req.GetExtraPack(); packValue != 0 {
		extraPack = converter.ConvertPackType(int64(packValue))
	}

	orderEntry := entity.OrderEntry{
		OrderID:   req.GetOrderId(),
		UserID:    req.GetUserId(),
		Weight:    req.GetWeight(),
		Price:     req.GetPrice(),
		PackType:  converter.ConvertPackType(int64(req.GetPackType())),
		ExtraPack: extraPack,
		Shelflife: req.GetShelfLife().AsTime(),
	}

	err := ValidateReq(orderEntry)
	if err != nil {
		h.logger.Warn("error during validation", zap.Error(err))

		return nil, status.Errorf(codes.InvalidArgument, "error during validation: %v", err)
	}

	if err := h.acceptServ.AcceptOrder(ctx, orderEntry); err != nil {
		h.logger.Error("h.acceptServ.AcceptOrder", zap.Error(err))

		return nil, status.Errorf(GetGRPCStatusCode(err), "h.acceptServ.AcceptOrder: %v", err)
	}
	monitoring.OrderPrice(orderEntry.Price)
	monitoring.OrderWeight(orderEntry.Weight)
	monitoring.SetCounterMainPackaging(string(orderEntry.PackType))

	h.logger.Info("result AcceptOrder OK")

	return &desc.AcceptResponse{Success: true}, nil
}

// ValidateReq валидация заказа
func ValidateReq(ord entity.OrderEntry) error {
	if ord.PackType == "" {
		return errors.New("bad PackType")
	}
	if ord.Price < 0 {
		return errors.New("bad Price")
	}
	if ord.Weight < 0 {
		return errors.New("bad Weight")
	}
	if ord.UserID < 0 {
		return errors.New("bad UserID")
	}
	if ord.OrderID < 0 {
		return errors.New("bad OrderID")
	}

	return nil
}
