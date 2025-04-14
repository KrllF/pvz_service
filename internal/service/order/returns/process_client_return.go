package returns

import (
	"context"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// ProcessClientReturn обработка возврата заказов клиентом
func (s *Service) ProcessClientReturn(ctx context.Context, userID int64, orderIDs []int64) ([]int64, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ProcessClientReturn from returns-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	var errors error
	ret := make([]int64, 0, len(orderIDs))
	ok, err := s.repository.UserExist(ctx, userID)
	if err != nil {
		s.logger.Error("ошибка при получении информации о клиенте", zap.Int64("userID", userID), zap.Error(err))

		return nil, fmt.Errorf("ошибка при получении информации о клиенте: %w", err)
	}
	if !ok {
		return nil, errs.ErrUserNotFound
	}

	for _, orderID := range orderIDs {
		err := s.processReturn(ctx, userID, orderID)
		if err != nil {
			errors = multierr.Append(errors, err)

			continue
		}
		if err := s.cache.DeleteCache(orderID); err != nil {
			s.logger.Error("s.cache.DeleteCache", zap.Int64("orderID", orderID), zap.Error(err))
		}
		ret = append(ret, orderID)
	}

	return ret, errors
}

func (s *Service) processReturn(ctx context.Context, userID, orderID int64) error {
	err := s.txManager.RunSerializable(ctx, func(ctxTX context.Context) error {
		order, err := s.repository.GetOrderInfo(ctxTX, orderID)
		if err != nil {
			return fmt.Errorf("ошибка при получении информации о заказе: %w", err)
		}

		if order.UserID != userID {
			return errs.ErrInvalidData
		}

		switch order.Status {
		case consts.AcceptedSt:
			if time.Now().After(order.TwoDaysOfLife) {
				return errs.ErrShelfLifeExpired
			}

			err = s.repository.UpdateStatus(ctxTX, orderID, consts.ReturnedSt)
			if err != nil {
				return fmt.Errorf("ошибка при обновлении статуса: %w", err)
			}

		default:
			return errs.ErrInvalidOrderStatus
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("s.txManager.RunSerializable: %w", err)
	}

	return nil
}
