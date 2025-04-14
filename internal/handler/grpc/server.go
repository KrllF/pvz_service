//go:generate mockgen -source=$GOFILE -destination=mocks/grpc.go -package=mocks

package grpc

import (
	"context"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"go.uber.org/zap"
)

type (
	// AcceptService определяет методы для управления заказами.
	// Содержит 3 функций:
	// AcceptOrder - принять заказ от курьера
	// BulkAcceptOrders - принять заказы от курьера через JSON-файл
	// ProcessClientIssue - выдать заказы клиента
	AcceptService interface {
		AcceptOrder(ctx context.Context, req entity.OrderEntry) error
		BulkAcceptOrders(ctx context.Context, path string) (int, error)
		ProcessClientIssue(ctx context.Context, userID int64, orderIDs []int64) ([]int64, error)
	}

	// ReturnsService определяет методы для управления заказами.
	// Содержит 2 функций:
	// ReturnOrder - вернуть заказ курьеру
	// ProcessClientReturn - принять возвраты клиента
	ReturnsService interface {
		ProcessClientReturn(ctx context.Context, userID int64, orderIDs []int64) ([]int64, error)
		ReturnOrder(ctx context.Context, orderID int64) error
	}

	// ShowService определяет методы для управления заказами.
	// Содержит 4 функций:
	// ListOrdersUserAllN - список заказов для клиента
	// ListOrdersUserPVZ - список заказов в пвз у клиента
	// ListReturns - список возвращённых заказов с пагинацией
	// ListHistory - история заказов
	ShowService interface {
		ListHistory(ctx context.Context) ([]models.Order, error)
		ListOrdersUserAllN(ctx context.Context, userID int64, opt, limit, page int64, pag bool,
			finder bool, orderFinder string) ([]models.Order, error)
		ListOrdersUserPVZ(ctx context.Context, userID, limit, page int64,
			pag bool, finder bool, orderFinder string) ([]models.Order, error)
		ListReturns(ctx context.Context, limit, page int64, finder bool, orderFinder string) ([]models.Order, int, error)
	}
	// Handler хэндлер ПВЗ
	Handler struct {
		desc.UnimplementedOrderV1Server
		acceptServ  AcceptService
		returnsServ ReturnsService
		showServ    ShowService
		logger      *zap.Logger
	}
)

// NewHandler конструктор хэндлера
func NewHandler(accept AcceptService, returns ReturnsService,
	show ShowService, logger *zap.Logger,
) *Handler {
	return &Handler{
		acceptServ:  accept,
		returnsServ: returns,
		showServ:    show,
		logger:      logger,
	}
}
