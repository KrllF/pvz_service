package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/KrllF/pvz_service/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

const (
	cursor = 10
)

// GetAll получить все заказы с курсором
func (r *Repo) GetAll(ctx context.Context, limit int64) ([]models.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetAll from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	ret := make([]models.Order, 0, capacity)

	query := `
				SELECT 
			orders.order_id AS orderid,
			orders.user_id AS userid,
			orders.order_weight AS weight,
			orders.order_price AS price,
			pack.pack_type AS packtype,
			pack.pack_type AS extrapack,
			statusorder.status_type AS status,
			orders.shelf_life AS shelflife,
			orders.two_days_of_life AS twodaysoflife,
			orders.created_at AS instoragefrom,
			orders.updated_at AS lastupdate
		FROM 
			orders
		LEFT JOIN 
			users ON orders.user_id = users.user_id
		LEFT JOIN 
			pack ON orders.pack_id = pack.pack_id
		LEFT JOIN 
			pack AS extrapack ON orders.extra_pack_id = extrapack.pack_id
		LEFT JOIN 
			statusorder ON orders.status_id = statusorder.status_id 
		WHERE order_id > $1
		ORDER BY order_id ASC
		LIMIT $2;
	`
	prevID := 0
	for {
		rows, err := r.tx.DB.Query(ctx, query, prevID, limit)

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			r.logger.Fatal("r.tx.DB.Query", zap.Error(err))

			return nil, fmt.Errorf("r.tx.DB.Query: %w", err)
		}
		ords, err := scanOrders(rows, cursor)
		if err != nil {
			r.logger.Error("scanOrders", zap.Error(err))

			return nil, fmt.Errorf("scanOrders: %w", err)
		}
		if len(ords) == 0 {
			break
		}
		ret = append(ret, ords...)
		prevID += int(ords[len(ords)-1].OrderID)
	}

	return ret, nil
}
