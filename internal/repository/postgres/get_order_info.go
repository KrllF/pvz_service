package postgres

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/jackc/pgtype"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// GetOrderInfo получить информацию о заказе
func (r *Repo) GetOrderInfo(ctx context.Context, orderID int64) (models.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetOrderInfo from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	var m models.Order
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
		WHERE orders.order_id = $1
	`
	ok, err := r.existsOrder(ctx, orderID)
	if err != nil {
		r.logger.Error("не удалось проверить существование заказа", zap.Int64("orderID", orderID))

		return models.Order{}, fmt.Errorf("не удалось проверить существование заказа с orderID %d: %w", orderID, err)
	}
	if !ok {
		r.logger.Warn("заказ не найден", zap.Error(errs.ErrOrderNotFound))

		return models.Order{}, errs.ErrOrderNotFound
	}

	var extraPack pgtype.Text
	var twoDays pgtype.Timestamp
	err = r.tx.GetDB(ctx).QueryRow(ctx, query, orderID).
		Scan(&m.OrderID, &m.UserID,
			&m.Weight, &m.Price,
			&m.Packaging.PackType, &extraPack,
			&m.Status, &m.ShelfLife, &twoDays,
			&m.InStorageFrom, &m.LastUpdate)
	if err != nil {
		r.logger.Fatal("не удалось получить информацию о заказе", zap.Int64("orderID", orderID), zap.Error(err))

		return models.Order{}, fmt.Errorf("не удалось получить информацию о заказе с orderID %d: %w", orderID, err)
	}
	if extraPack.Status == pgtype.Present {
		m.Packaging.ExtraPack = entity.PackType(extraPack.String)
	}
	if twoDays.Status == pgtype.Present {
		m.TwoDaysOfLife = twoDays.Time
	}

	return m, nil
}
