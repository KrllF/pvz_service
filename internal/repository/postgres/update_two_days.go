package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const (
	hourInTwoDays = 48
)

// SetTwoDaysOfLife установить два дня жизни
func (r *Repo) SetTwoDaysOfLife(ctx context.Context, orderID int64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "SetTwoDaysOfLife from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	ok, err := r.existsOrder(ctx, orderID)
	if err != nil {
		r.logger.Warn("не удалось проверить существование заказа", zap.Int64("orderID", orderID), zap.Error(err))

		return fmt.Errorf("не удалось проверить существование заказа с orderID %d: %w", orderID, err)
	}
	if !ok {
		r.logger.Warn("заказ не найден", zap.Error(errs.ErrOrderNotFound))

		return errs.ErrOrderNotFound
	}

	query := `UPDATE orders SET two_days_of_life = $1, updated_at = $2 WHERE order_id = $3;`
	ret, err := r.tx.GetDB(ctx).Exec(ctx, query, time.Now().Add(hourInTwoDays*time.Hour), time.Now(), orderID)
	if err != nil {
		r.logger.Fatal("не удалось обновить заказ", zap.Error(err))

		return fmt.Errorf("не удалось обновить заказ: %w", err)
	}
	if ret.RowsAffected() == 0 {
		r.logger.Warn("невалидные данные", zap.Error(err))

		return errs.ErrInvalidData
	}

	return nil
}
