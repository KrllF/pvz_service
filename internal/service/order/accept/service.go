//go:generate mockgen -source=$GOFILE -destination=mocks/accept.go -package=mocks

package accept

import (
	"context"

	"github.com/KrllF/pvz_service/internal/models"
	"go.uber.org/zap"
)

type (
	// Repository определяет методы для взаимодействия с JSON-файлом.
	// Содержит 4 функций:
	// AddOrder - добавить заказ
	// GetOrderInfo - получить информацию о заказе (orderinfo)
	// ListOrders - получить все заказы
	// UpdateStatus - обновить заказы
	Repository interface {
		AddOrder(ctx context.Context, ord models.Order) error
		AddUser(ctx context.Context, userID int64) error
		GetOrderInfo(ctx context.Context, orderID int64) (models.Order, error)
		UserExist(ctx context.Context, userID int64) (bool, error)
		UpdateStatus(ctx context.Context, orderID int64, newStatus string) error
		SetTwoDaysOfLife(ctx context.Context, orderID int64) error
	}
	// Cache определяет методы для взаимодействия с кешом
	Cache interface {
		Put(key int64, value models.Order) error
		GetItem(key int64) (models.Order, error)
	}

	// TXManager интерфейс менеджера транзакций
	TXManager interface {
		RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
	}
	// Service accept
	Service struct {
		repository Repository
		cache      Cache
		txManager  TXManager
		logger     *zap.Logger
	}
)

// NewService создаёт новый экземляр serv
// Принимает интерфейс Repository - возвращает указатель на созданный serv
func NewService(storage Repository, cache Cache, tx TXManager, logger *zap.Logger) *Service {
	return &Service{
		repository: storage,
		cache:      cache,
		txManager:  tx,
		logger:     logger,
	}
}
