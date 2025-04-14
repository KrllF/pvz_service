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

// ListOrdersUserPVZ список заказов юзера в ПВЗ
func (s *Service) ListOrdersUserPVZ(ctx context.Context, userID, limit, page int64,
	pag bool, finder bool, orderFinder string,
) ([]models.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListOrdersUserPVZ from show-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	var orders []models.Order
	if pag {
		orders, err = s.listordPVZpag(ctx, userID, limit, page, finder, orderFinder)
	} else {
		orders, err = s.listordPVZ(ctx, userID, finder, orderFinder)
	}
	if err != nil {
		s.logger.Error("ошибка при получении списка заказов", zap.Error(err))

		return nil, fmt.Errorf("ошибка при получении списка заказов: %w", err)
	}

	return orders, nil
}

func (s *Service) listordPVZpag(ctx context.Context, userID int64, limit, page int64,
	finder bool, orderFinder string,
) ([]models.Order, error) {
	if finder {
		return s.repository.ListOrders(ctx, entity.WithUserID(userID),
			entity.WithStatus(consts.DeliveredST, consts.ReturnedSt),
			entity.WithPagination(limit, page), entity.WithFindByID(orderFinder))
	}

	return s.repository.ListOrders(ctx, entity.WithUserID(userID),
		entity.WithStatus(consts.DeliveredST, consts.ReturnedSt), entity.WithPagination(limit, page))
}

func (s *Service) listordPVZ(ctx context.Context, userID int64, finder bool,
	orderFinder string,
) ([]models.Order, error) {
	if finder {
		return s.repository.ListOrders(ctx, entity.WithUserID(userID),
			entity.WithStatus(consts.DeliveredST, consts.ReturnedSt), entity.WithFindByID(orderFinder))
	}

	return s.repository.ListOrders(ctx, entity.WithUserID(userID),
		entity.WithStatus(consts.DeliveredST, consts.ReturnedSt))
}
