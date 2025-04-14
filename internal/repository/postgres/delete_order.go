package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// DeleteOrder удалить заказ
func (r *Repo) DeleteOrder(ctx context.Context, orderID int64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "DeleteOrder from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	query := `DELETE FROM orders WHERE order_id = $1 RETURNING user_id`

	var userID int64
	err = r.tx.GetDB(ctx).QueryRow(ctx, query, orderID).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Warn("Заказа нет", zap.Int64("orderID", orderID))

		return errs.ErrOrderNotFound
	}
	if err != nil {
		r.logger.Fatal("Ошибка при удалении заказа", zap.Int64("orderID", orderID), zap.Error(err))

		return fmt.Errorf("ошибка при удалении заказа %w", err)
	}

	r.logger.Info("заказ успешно удален", zap.Int64("orderID", orderID))

	return nil
}
