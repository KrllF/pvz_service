//go:generate mockgen -source=$GOFILE -destination=mocks/show.go -package=mocks

package show

import (
	"context"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	"go.uber.org/zap"
)

type (
	// Repository определяет методы для взаимодействия с JSON-файлом.
	// Содержит 7 функций:
	// AddOrder - добавить заказ
	// RemoveOrder - удалить заказ
	// GetOrderInfo - получить информацию о заказе (orderinfo)
	// GetOrdersForUser - получить список id заказов пользователя
	// GetAllOrders - получить все заказы
	// UpdateStatus - обновить заказы
	// ReadFile - прочитать файл
	Repository interface {
		ListOrders(ctx context.Context, opts ...entity.ListOrdersOption) ([]models.Order, error)
	}
	// Cache определяет методы для взаимодействия с кэшом
	Cache interface {
		ListItems(ctx context.Context) ([]models.Order, error)
	}
	// Service show
	Service struct {
		repository Repository
		cache      Cache
		logger     *zap.Logger
	}
)

// NewService создаёт новый экземляр serv
// Принимает интерфейс Repository - возвращает указатель на созданный serv
func NewService(storage Repository, cache Cache, logger *zap.Logger) (*Service, error) {
	return &Service{
		repository: storage,
		cache:      cache,
		logger:     logger,
	}, nil
}
