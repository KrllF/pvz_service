package log

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/models"
	"go.uber.org/zap"
)

const (
	limit = 10
)

// FillCache заполнить кэш
func (c *Cron) FillCache(ctx context.Context) error {
	orders, err := c.Repository.GetAll(ctx, limit)
	if err != nil {
		c.logger.Error("c.Repository.ListOrders", zap.Error(err))

		return fmt.Errorf("c.Repository.ListOrders: %w", err)
	}
	if err = c.Cache.UpdateAllCache(orders, func(ord models.Order) int64 { return ord.OrderID }); err != nil {
		c.logger.Error("c.Cache.UpdateAllCache", zap.Error(err))

		return fmt.Errorf("c.Cache.UpdateAllCache: %w", err)
	}

	return nil
}
