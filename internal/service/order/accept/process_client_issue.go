package accept

import (
	"context"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/errs"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const (
	twoDays = 2
)

// ProcessClientIssue отдать заказы клиенту
func (s *Service) ProcessClientIssue(ctx context.Context, userID int64, orderIDs []int64) ([]int64, error) {
	ret := make([]int64, 0, len(orderIDs))
	ok, err := s.repository.UserExist(ctx, userID)
	if err != nil {
		s.logger.Error("ошибка при получении информации о клиенте", zap.Int64("userID", userID), zap.Error(err))

		return nil, fmt.Errorf("ошибка при получении информации о клиенте: %w", err)
	}
	if !ok {
		return nil, errs.ErrUserNotFound
	}
	var errors error

	for _, orderID := range orderIDs {
		err := s.processIssue(ctx, userID, orderID)
		if err != nil {
			errors = multierr.Append(errors, err)

			continue
		}
		k, err := s.cache.GetItem(orderID)
		if err != nil {
			s.logger.Error("s.cache.GetItem", zap.Int64("orderID", orderID), zap.Error(err))
		}
		k.Status = consts.AcceptedSt
		k.TwoDaysOfLife = time.Now().Add(twoDays * time.Hour)
		err = s.cache.Put(orderID, k)
		if err != nil {
			s.logger.Error("s.cache.AddCache", zap.Int64("orderID", orderID), zap.Error(err))
		}
		ret = append(ret, orderID)
	}

	return ret, errors
}

func (s *Service) processIssue(ctx context.Context, userID, orderID int64) error {
	err := s.txManager.RunSerializable(ctx, func(ctxTX context.Context) error {
		order, err := s.repository.GetOrderInfo(ctxTX, orderID)
		if err != nil {
			return err
		}
		if order.UserID != userID {
			return errs.ErrInvalidData
		}

		switch order.Status {
		case consts.DeliveredST:
			if time.Now().After(order.ShelfLife) {
				return errs.ErrShelfLifeExpired
			}

			err = s.repository.UpdateStatus(ctxTX, orderID, consts.AcceptedSt)
			if err != nil {
				return fmt.Errorf("ошибка при обновлении статуса заказа: %w", err)
			}
			err = s.repository.SetTwoDaysOfLife(ctxTX, orderID)
			if err != nil {
				return fmt.Errorf("ошибка при обновлении времени хранения заказа у клиента: %w", err)
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
