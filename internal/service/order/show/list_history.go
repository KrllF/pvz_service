package show

import (
	"context"
	"fmt"
	"sort"

	"github.com/KrllF/pvz_service/internal/models"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// ListHistory история заказов
func (s *Service) ListHistory(ctx context.Context) ([]models.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListHistory from show-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	lst, err := s.cache.ListItems(ctx)
	if err != nil {
		s.logger.Error("s.cache.ListOrders()", zap.Error(err))

		return nil, fmt.Errorf("s.cache.ListOrders(): %w", err)
	}
	sort.Slice(lst, func(i, j int) bool {
		return lst[i].LastUpdate.Before(lst[j].LastUpdate)
	})

	return lst, nil
}
