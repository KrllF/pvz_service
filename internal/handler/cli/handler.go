package cli

import (
	"context"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
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
		ProcessClientIssue(ctx context.Context, userID int64, orderIDs []int64) (int64, error)
	}

	// ReturnsService определяет методы для управления заказами.
	// Содержит 2 функций:
	// ReturnOrder - вернуть заказ курьеру
	// ProcessClientReturn - принять возвраты клиента
	ReturnsService interface {
		ProcessClientReturn(ctx context.Context, userID int64, orderIDs []int64) (int64, error)
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
		ListOrdersUserAllN(ctx context.Context, userID int64, opt int64) ([]models.Order, error)
		ListOrdersUserPVZ(ctx context.Context, userID int64) ([]models.Order, error)
		ListReturns(ctx context.Context, limit, page int64) ([]models.Order, int, error)
	}

	clihandler struct {
		acceptServ  AcceptService
		returnsServ ReturnsService
		showServ    ShowService
	}
)

// NewCliHandler создаёт новый экземляр CLI-хэндлера
// Принимает интерфейс Service - возвращает указатель на созданный CLI-хэндлер
func NewCliHandler(accept AcceptService, returns ReturnsService, show ShowService) *clihandler {
	return &clihandler{acceptServ: accept, returnsServ: returns, showServ: show}
}
