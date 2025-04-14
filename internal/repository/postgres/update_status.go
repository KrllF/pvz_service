package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// UpdateStatus обновить статус заказа
func (r *Repo) UpdateStatus(ctx context.Context, orderID int64, newStatus string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UpdateStatus from postgres-repository")
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
		r.logger.Fatal("не удалось проверить существование заказа", zap.Int64("orderID", orderID), zap.Error(err))

		return fmt.Errorf("не удалось проверить существование заказа с orderID %d: %w", orderID, err)
	}
	if !ok {
		r.logger.Warn("заказ не найден", zap.Error(errs.ErrOrderNotFound))

		return errs.ErrOrderNotFound
	}

	query := `UPDATE orders SET status_id = $1, updated_at = $2 WHERE order_id = $3;`
	NewStatusID, err := r.getStatusID(ctx, newStatus)
	if err != nil {
		r.logger.Fatal("не удалось получить статус", zap.Error(err))

		return fmt.Errorf("не удалось получить статус: %w", err)
	}

	ret, err := r.tx.GetDB(ctx).Exec(ctx, query, NewStatusID, time.Now(), orderID)
	if err != nil {
		r.logger.Fatal("не удалось обновить", zap.Error(err))

		return fmt.Errorf("не удалось обновить: %w", err)
	}

	if ret.RowsAffected() == 0 {
		r.logger.Warn("невалидные данные", zap.Error(err))

		return errs.ErrInvalidData
	}

	return nil
}
