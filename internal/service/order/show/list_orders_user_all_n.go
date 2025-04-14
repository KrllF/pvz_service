package show

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// ListOrdersUserAllN список заказов user
func (s *Service) ListOrdersUserAllN(ctx context.Context, userID int64,
	opt, limit, page int64, pag bool, finder bool,
	orderFinder string,
) ([]models.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListOrdersUserAllN from show-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	var orders []models.Order
	if opt < 0 {
		if pag {
			orders, err = s.listord(ctx, userID, limit, page, finder, orderFinder)
		} else {
			orders, err = s.repository.ListOrders(ctx, entity.WithUserID(userID))
		}
		if err != nil {
			s.logger.Error("ошибка при получении списка заказов", zap.Error(err))

			return nil, fmt.Errorf("ошибка при получении списка заказов: %w", err)
		}

		return orders, nil
	}

	if finder {
		orders, err = s.repository.ListOrders(ctx, entity.WithUserID(userID),
			entity.WithPagination(opt, 1), entity.WithFindByID(orderFinder))
	} else {
		orders, err = s.repository.ListOrders(ctx, entity.WithUserID(userID), entity.WithPagination(opt, 1))
	}
	if err != nil {
		s.logger.Error("ошибка при получении списка заказов", zap.Error(err))

		return nil, fmt.Errorf("ошибка при получении списка заказов: %w", err)
	}

	return orders, nil
}

func (s *Service) listord(ctx context.Context, userID int64,
	limit, page int64, finder bool,
	orderFinder string,
) ([]models.Order, error) {
	if finder {
		return s.repository.ListOrders(ctx, entity.WithUserID(userID), entity.WithPagination(limit, page),
			entity.WithFindByID(orderFinder))
	}

	return s.repository.ListOrders(ctx, entity.WithUserID(userID), entity.WithPagination(limit, page))
}
