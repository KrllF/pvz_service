package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// ListOrders получить список заказов с параметрами
func (r *Repo) ListOrders(ctx context.Context, opts ...entity.ListOrdersOption) ([]models.Order, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ListOrders from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	options := &entity.ListOrdersOptions{}

	for _, opt := range opts {
		opt(options)
	}

	query, args, err := buildListOrdersQuery(options)
	if err != nil {
		r.logger.Error("ошибка при конструировании запроса", zap.Error(err))

		return nil, fmt.Errorf("ошибка при конструировании запроса: %w", err)
	}

	rows, err := r.tx.GetDB(ctx).Query(ctx, query, args...)
	if err != nil {
		r.logger.Fatal("Ошибка при выполнении запроса", zap.Error(err))

		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer rows.Close()

	orders, err := scanOrders(rows, capacity)
	if err != nil {
		r.logger.Error("ошибка при скане", zap.Error(err))

		return nil, fmt.Errorf("ошибка при скане: %w", err)
	}

	return orders, nil
}

func buildListOrdersQuery(options *entity.ListOrdersOptions) (string, []interface{}, error) {
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
			statusorder ON orders.status_id = statusorder.status_id    `

	var conditions []string
	var args []interface{}
	argIndex := 1

	if options.UserID != 0 {
		conditions = append(conditions, fmt.Sprintf("orders.user_id = $%d", argIndex))
		args = append(args, options.UserID)
		argIndex++
	}
	if options.OrderID != 0 {
		conditions = append(conditions, fmt.Sprintf("order_id = $%d", argIndex))
		args = append(args, options.OrderID)
		argIndex++
	}

	if len(options.Status) > 0 {
		var statusConditions []string
		for _, status := range options.Status {
			statusConditions = append(statusConditions, fmt.Sprintf("statusorder.status_type = $%d", argIndex))
			args = append(args, status)
			argIndex++
		}
		conditions = append(conditions, "("+strings.Join(statusConditions, " OR ")+")")
	}

	if options.FindByID != "" {
		conditions = append(conditions, fmt.Sprintf("orders.order_id::text LIKE $%d", argIndex))
		args = append(args, "%"+options.FindByID+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query, err := addSortingToQuery(query, options)
	if err != nil {
		return "", nil, fmt.Errorf("ошибка при добавлении условий сортировки: %w", err)
	}

	if options.Limit == 0 && options.Page == 0 {
		return query, args, nil
	}

	offset := (options.Page - 1) * options.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, options.Limit, offset)

	return query, args, nil
}

func addSortingToQuery(query string, options *entity.ListOrdersOptions) (string, error) {
	if options.SortBy == "" {
		return query, nil
	}

	SortField := map[string]struct{}{
		"order_id":     {},
		"user_id":      {},
		"order_weight": {},
		"order_price":  {},
		"shelf_life":   {},
		"created_at":   {},
		"updated_at":   {},
	}
	if _, ok := SortField[options.SortBy]; !ok {
		log.Printf("Недопустимое поле для сортировки: %s", options.SortBy)

		return "", errs.ErrInvalidSortField
	}
	query += fmt.Sprintf(" ORDER BY %s %s", options.SortBy, options.SortOrder)

	return query, nil
}
