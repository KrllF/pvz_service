package show

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// ListReturns список заказов, которые вернул клиент
func (s *Service) ListReturns(ctx context.Context, limit, page int64, finder bool,
	orderFinder string,
) ([]models.Order, int, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListReturns from show-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	if finder {
		allOrders, err := s.repository.ListOrders(ctx, entity.WithStatus(consts.ReturnedSt),
			entity.WithPagination(limit, page), entity.WithFindByID(orderFinder))
		if err != nil {
			s.logger.Error("ошибка при получении списка заказов", zap.Error(err))

			return nil, 0, fmt.Errorf("ошибка при получении списка заказов: %w", err)
		}

		return allOrders, len(allOrders), nil
	}

	allOrders, err := s.repository.ListOrders(ctx, entity.WithStatus(consts.ReturnedSt),
		entity.WithPagination(limit, page))
	if err != nil {
		s.logger.Error("ошибка при получении списка заказов", zap.Error(err))

		return nil, 0, fmt.Errorf("ошибка при получении списка заказов: %w", err)
	}

	return allOrders, len(allOrders), nil
}
