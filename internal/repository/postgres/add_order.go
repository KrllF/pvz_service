package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/models"
	"github.com/jackc/pgtype"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// AddOrder добавить заказ
func (r *Repo) AddOrder(ctx context.Context, ord models.Order) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "AddOrder from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	queryAdd := `INSERT INTO orders(order_id,user_id,
									order_weight, order_price,
									pack_id, extra_pack_id,
									status_id, shelf_life,
									two_days_of_life, created_at, updated_at)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8, $9, $10, $11);`

	PackID, ExtraID, StatusID, err := r.getPackExtraStat(ctx, ord)
	if err != nil {
		r.logger.Error("ошибка при получении pack_id, extra_pack_id, status_id", zap.Error(err))

		return fmt.Errorf("ошибка при получении pack_id, extra_pack_id, status_id: %w", err)
	}
	twoDaysOfLife := pgtype.Timestamptz{
		Time:   time.Time{},
		Status: pgtype.Null,
	}
	_, err = r.tx.GetDB(ctx).Exec(ctx, queryAdd, ord.OrderID, ord.UserID, ord.Weight, ord.Price,
		PackID, ExtraID, StatusID, ord.ShelfLife, twoDaysOfLife, time.Now(), time.Now())
	if err != nil {
		r.logger.Fatal("Ошибка при добавлении заказа в таблицу", zap.Error(err))

		return fmt.Errorf("ошибка при добавлении заказа в таблицу: %w", err)
	}

	return nil
}

func (r *Repo) getPackExtraStat(ctx context.Context, ord models.Order) (int64, pgtype.Int8, int64, error) {
	PackID, err := r.getPackID(ctx, string(ord.Packaging.PackType))
	if err != nil {
		return 0, pgtype.Int8{}, 0, fmt.Errorf("не удалось получить pack_id: %w", err)
	}

	ExtraID := pgtype.Int8{Status: pgtype.Null}
	if ord.Packaging.ExtraPack != "" {
		ExtraID, err = r.getExtraPackID(ctx, string(ord.Packaging.ExtraPack))
		if err != nil {
			return 0, pgtype.Int8{}, 0, fmt.Errorf("не удалось получить extra_pack_id: %w", err)
		}
	}

	StatusID, err := r.getStatusID(ctx, ord.Status)
	if err != nil {
		return 0, pgtype.Int8{}, 0, fmt.Errorf("не удалось получить status_id: %w", err)
	}

	return PackID, ExtraID, StatusID, nil
}
