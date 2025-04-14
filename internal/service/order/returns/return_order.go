package returns

import (
	"context"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// ReturnOrder возврат заказа курьеру
func (s *Service) ReturnOrder(ctx context.Context, orderID int64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "ReturnOrder from returns-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	err = s.txManager.RunSerializable(ctx, func(ctxTX context.Context) error {
		k, err := s.repository.GetOrderInfo(ctxTX, orderID)
		if err != nil {
			s.logger.Error("ошибка при обращении к базе", zap.Error(err))

			return fmt.Errorf("ошибка при получении информации о заказе: %w", err)
		}

		if k.Status == consts.AcceptedSt {
			s.logger.Warn("недопустимый статус заказа")

			return errs.ErrInvalidOrderStatus
		}

		if time.Now().Before(k.ShelfLife) {
			s.logger.Warn("некорректное время")

			return errs.ErrTimeNotExpired
		}

		err = s.repository.DeleteOrder(ctxTX, orderID)
		if err != nil {
			s.logger.Error("ошибка при удалении заказа",
				zap.Int64("orderID", orderID), zap.Error(err))

			return fmt.Errorf("ошибка при удалении заказа: %w", err)
		}

		ret, err := s.repository.ListOrders(ctxTX, entity.WithUserID(k.UserID))
		if err != nil {
			s.logger.Error("ошибка при получении информации о пользователе",
				zap.Int64("userID", k.UserID), zap.Error(err))

			return fmt.Errorf("ошибка при получении информации о пользователе: %w", err)
		}

		if len(ret) != 0 {
			return nil
		}

		err = s.repository.DeleteUser(ctxTX, k.UserID)
		if err != nil {
			s.logger.Error("ошибка при удалении пользователя",
				zap.Int64("userID", k.UserID), zap.Error(err))

			return fmt.Errorf("ошибка при удалении пользователя: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("s.txManager.RunSerializable: %w", err)
	}

	if err = s.cache.DeleteCache(orderID); err != nil {
		s.logger.Error("s.cache.DeleteCache", zap.Int64("orderID", orderID), zap.Error(err))
	}

	return nil
}
