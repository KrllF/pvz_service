package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func (r *Repo) getPackID(ctx context.Context, packType string) (int64, error) {
	var packID int64
	query := `SELECT pack_id FROM pack WHERE pack_type = $1`
	err := r.tx.GetDB(ctx).QueryRow(ctx, query, packType).Scan(&packID)
	if errors.Is(err, pgx.ErrNoRows) {
		return packID, errs.ErrPackTypeNotFound
	}
	if err != nil {
		return packID, fmt.Errorf("не удалось получить pack_id, %w", err)
	}

	return packID, nil
}

func (r *Repo) getExtraPackID(ctx context.Context, extraPackType string) (pgtype.Int8, error) {
	var extraPackID pgtype.Int8
	query := `SELECT pack_id FROM pack WHERE pack_type = $1 AND is_extra=TRUE`
	err := r.tx.GetDB(ctx).QueryRow(ctx, query, extraPackType).Scan(&extraPackID)
	if errors.Is(err, pgx.ErrNoRows) {
		return pgtype.Int8{Status: pgtype.Null}, errs.ErrExtraPackNotFound
	}
	if err != nil {
		return pgtype.Int8{}, fmt.Errorf("не удалось получить extra_pack_id, %w", err)
	}
	extraPackID.Status = pgtype.Present

	return extraPackID, nil
}

func (r *Repo) getStatusID(ctx context.Context, statusName string) (int64, error) {
	var statusID int64
	query := `SELECT status_id FROM statusorder WHERE status_type = $1`
	err := r.tx.GetDB(ctx).QueryRow(ctx, query, statusName).Scan(&statusID)
	if errors.Is(err, pgx.ErrNoRows) {
		return statusID, errs.ErrStatusNotFound
	}
	if err != nil {
		return statusID, fmt.Errorf("не удалось получить status_id, %w", err)
	}

	return statusID, nil
}

func (r *Repo) existsOrder(ctx context.Context, orderID int64) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM orders WHERE order_id = $1);`
	var exists bool
	err := r.tx.GetDB(ctx).QueryRow(ctx, query, orderID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("не удалось проверить существование заказа с orderID %d: %w", orderID, err)
	}

	return exists, nil
}

func scanOrders(rows pgx.Rows, capacity int64) ([]models.Order, error) {
	orders := make([]models.Order, 0, capacity)
	for rows.Next() {
		var order models.Order
		var extraPackType pgtype.Text
		var twoDays pgtype.Timestamp
		err := rows.Scan(
			&order.OrderID,
			&order.UserID,
			&order.Weight,
			&order.Price,
			&order.Packaging.PackType,
			&extraPackType,
			&order.Status,
			&order.ShelfLife,
			&twoDays,
			&order.InStorageFrom,
			&order.LastUpdate,
		)
		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)

			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		order.TwoDaysOfLife = time.Time{}
		if twoDays.Status == pgtype.Present {
			order.TwoDaysOfLife = twoDays.Time
		}
		order.Packaging.ExtraPack = ""
		if extraPackType.Status == pgtype.Present {
			order.Packaging.ExtraPack = entity.PackType(extraPackType.String)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Ошибка при чтении строк: %v", err)

		return nil, fmt.Errorf("ошибка при чтении строк: %w", err)
	}

	return orders, nil
}
