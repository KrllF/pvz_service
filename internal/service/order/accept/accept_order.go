package accept

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/internal/service/order/accept/helper"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// AcceptOrder принять заказ у курьера
func (s *Service) AcceptOrder(ctx context.Context, req entity.OrderEntry) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "AcceptOrder from accept-service")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	_, err = s.repository.GetOrderInfo(ctx, req.OrderID)
	if !errors.Is(err, errs.ErrOrderNotFound) && err != nil {
		s.logger.Error("ошибка при обращении к базе", zap.Error(err))

		return fmt.Errorf("ошибка при обращении к базе: %w", err)
	}
	if err == nil {
		s.logger.Warn("заказ уже есть", zap.Int64("orderID", req.OrderID))

		return errs.ErrInvalidData
	}

	if time.Now().After(req.Shelflife) {
		s.logger.Warn("некорректное время")

		return errs.ErrInvalidData
	}

	ord := models.Order{
		OrderID:       req.OrderID,
		UserID:        req.UserID,
		Weight:        req.Weight,
		Price:         req.Price,
		Packaging:     entity.NewPack(req.PackType, req.ExtraPack),
		ShelfLife:     req.Shelflife,
		InStorageFrom: time.Now(),
		Status:        consts.DeliveredST,
		TwoDaysOfLife: time.Time{},
		LastUpdate:    time.Now(),
	}

	if err = helper.Validate(ord); err != nil {
		return fmt.Errorf("ошибка при валидации заказа: %w", err)
	}
	ord.Price = helper.CalcPrice(ord)

	err = s.txManager.RunSerializable(ctx, func(ctxTX context.Context) error {
		ok, err := s.repository.UserExist(ctxTX, ord.UserID)
		if err != nil {
			return fmt.Errorf("ошибка при выполнении UserExist: %w", err)
		}
		if !ok {
			err = s.repository.AddUser(ctxTX, req.UserID)
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
		s.logger.Error("s.txManager.RunSerializable", zap.Error(err))

		return fmt.Errorf("s.txManager.RunSerializable: %w", err)
	}

	if err = s.cache.Put(ord.OrderID, ord); err != nil {
		s.logger.Error("s.cache.Put", zap.Int64("orderID", ord.OrderID), zap.Error(err))
	}

	return nil
}
