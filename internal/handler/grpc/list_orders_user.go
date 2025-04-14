package grpc

import (
	"context"
	"errors"

	"github.com/KrllF/pvz_service/internal/converter"
	"github.com/KrllF/pvz_service/internal/models"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidateStructToListOrders структура для валидации
type ValidateStructToListOrders struct {
	n           int64
	inPVZ       bool
	limit       int64
	page        int64
	pag         bool
	finder      bool
	orderFinder string
}

// ListOrdersUser получить список заказов клиента
func (h *Handler) ListOrdersUser(ctx context.Context,
	req *desc.ListOrdersUserRequest,
) (*desc.ListOrdersUserResponse, error) {
	h.logger.Info("Processing request ListOrdersUser")

	userID := req.GetUserId()
	parseReq, err := validateParamsU(req)
	if err != nil {
		h.logger.Warn("validateParamsU", zap.Error(err))

		return nil, status.Errorf(codes.InvalidArgument, "validateParamsU: %v", err)
	}
	var ret []models.Order

	switch parseReq.inPVZ {
	case true:
		ret, err = h.showServ.ListOrdersUserPVZ(ctx, userID, parseReq.limit,
			parseReq.page, parseReq.pag, parseReq.finder, parseReq.orderFinder)
	case false:
		if parseReq.n == -1 {
			ret, err = h.showServ.ListOrdersUserAllN(ctx, userID, -1, parseReq.limit, parseReq.page,
				parseReq.pag, parseReq.finder, parseReq.orderFinder)
		} else {
			ret, err = h.showServ.ListOrdersUserAllN(ctx, userID, parseReq.n,
				parseReq.limit, parseReq.page, parseReq.pag, parseReq.finder, parseReq.orderFinder)
		}
	}
	if err != nil {
		h.logger.Error("h.showServ.ListOrdersUser", zap.Error(err))

		return nil, status.Errorf(GetGRPCStatusCode(err), "h.showServ.ListOrdersUser: %v", err)
	}

	conv := make([]*desc.Order, 0, len(ret))
	for _, val := range ret {
		conv = append(conv, converter.ConvertToProtoOrder(&val))
	}
	h.logger.Info("result ListOrdersUser OK")

	return &desc.ListOrdersUserResponse{Orders: conv}, nil
}

func validateParamsU(req *desc.ListOrdersUserRequest) (ValidateStructToListOrders, error) {
	n := nRet(req)

	inPVZ := inPVZRet(req)

	limit := limitRet(req)

	page := pageRet(req)

	pag := pagRet(req)

	finder, orderFinder, err := parseFinderOrder(req)
	if err != nil {
		return ValidateStructToListOrders{}, errors.New("incorrect parameters finder и order")
	}

	ret := ValidateStructToListOrders{
		n:           n,
		inPVZ:       inPVZ,
		limit:       limit,
		page:        page,
		pag:         pag,
		finder:      finder,
		orderFinder: orderFinder,
	}

	return ret, nil
}

func pagRet(req *desc.ListOrdersUserRequest) bool {
	return req.GetPag()
}

func limitRet(req *desc.ListOrdersUserRequest) int64 {
	limit := int64(-1)
	limitStr := req.GetLimit()
	if limitStr <= 0 {
		return limit
	}

	return limitStr
}

func pageRet(req *desc.ListOrdersUserRequest) int64 {
	page := int64(-1)
	pageStr := req.GetPage()
	if pageStr <= 0 {
		return page
	}

	return pageStr
}

func nRet(req *desc.ListOrdersUserRequest) int64 {
	n := int64(-1)
	nStr := req.GetN()
	if nStr <= 0 {
		return n
	}

	return nStr
}

func inPVZRet(req *desc.ListOrdersUserRequest) bool {
	return req.GetInPvz()
}
