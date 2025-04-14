package accept

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/accept/helper"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// BulkAcceptOrders принять заказы из JSON файла
func (s *Service) BulkAcceptOrders(ctx context.Context, path string) (int, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "BulkAcceptOrders from accept-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	orders, err := helper.ReadFile(path)
	if err != nil {
		s.logger.Error("не удалось взаимодействие с файлом", zap.Error(err))

		return 0, fmt.Errorf("не удалось взаимодействие с файлом: %w", err)
	}

	var counter int

	for _, order := range orders {
		ord := models.Order{
			OrderID:       order.OrderID,
			UserID:        order.UserID,
			Weight:        order.Weight,
			Price:         order.Price,
			Packaging:     order.Packaging,
			ShelfLife:     order.ShelfLife,
			InStorageFrom: time.Now(),
			Status:        consts.DeliveredST,
			TwoDaysOfLife: time.Time{},
			LastUpdate:    time.Now(),
		}
		inc, err := s.bulkAdd(ctx, ord)
		if err != nil {
			s.logger.Error("ошибка при добавлении", zap.Int64("orderID", ord.OrderID), zap.Error(err))

			return counter, fmt.Errorf("ошибка при добавлении: %w", err)
		}
		counter += inc
	}

	return counter, nil
}

func (s *Service) bulkAdd(ctx context.Context, ord models.Order) (int, error) {
	_, err := s.repository.GetOrderInfo(ctx, ord.OrderID)
	if !errors.Is(err, errs.ErrOrderNotFound) && err != nil {
		return 0, fmt.Errorf("ошибка при обращении к базе: %w", err)
	}

	if err == nil {
		return 0, nil
	}

	if time.Now().After(ord.ShelfLife) {
		return 0, nil
	}

	err = s.txManager.RunSerializable(ctx, func(ctxTX context.Context) error {
		ok, err := s.repository.UserExist(ctxTX, ord.UserID)
		if err != nil {
			return fmt.Errorf("ошибка при выполнении UserExist: %w", err)
		}

		if !ok {
			err = s.repository.AddUser(ctxTX, ord.UserID)
			if err != nil {
				return fmt.Errorf("ошибка при добавлении пользователя: %w", err)
			}
		}

		err = s.repository.AddOrder(ctxTX, ord)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении пользователя: %w", err)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("s.txManager.RunSerializable: %w", err)
	}

	if err = s.cache.Put(ord.OrderID, ord); err != nil {
		s.logger.Error("s.cache.AddCache", zap.Int64("orderID", ord.OrderID), zap.Error(err))
	}

	return 1, nil
}
