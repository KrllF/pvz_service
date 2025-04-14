//go:generate mockgen -source=$GOFILE -destination=mocks/returns.go -package=mocks

package returns

import (
	"context"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	"go.uber.org/zap"
)

type (
	// Repository определяет методы для взаимодействия с JSON-файлом.
	// Содержит 4 функций:
	// RemoveOrder - удалить заказ
	// GetOrderInfo - получить информацию о заказе (orderinfo)
	// GetOrdersForUser - получить список id заказов пользователя
	// UpdateStatus - обновить заказы
	Repository interface {
		DeleteOrder(ctx context.Context, orderID int64) error
		DeleteUser(ctx context.Context, userID int64) error
		GetOrderInfo(ctx context.Context, orderID int64) (models.Order, error)
		UpdateStatus(ctx context.Context, orderID int64, newStatus string) error
		ListOrders(ctx context.Context, opts ...entity.ListOrdersOption) ([]models.Order, error)
		UserExist(ctx context.Context, userID int64) (bool, error)
	}
	// Cache определяет методы для взаимодействия с кешом
	Cache interface {
		DeleteCache(orderID int64) error
	}
	// TXManager интерфейс менеджера транзакций
	TXManager interface {
		RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
	}
	// Service returns
	Service struct {
		repository Repository
		cache      Cache
		txManager  TXManager
		logger     *zap.Logger
	}
)

// NewService создаёт новый экземляр returnsServ
// Принимает интерфейс Repository - возвращает указатель на созданный returnsServ
func NewService(storage Repository, cache Cache, tx TXManager, logger *zap.Logger) *Service {
	return &Service{
		repository: storage,
		cache:      cache,
		txManager:  tx,
		logger:     logger,
	}
}
