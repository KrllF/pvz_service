package postgreslog

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

const (
	capacitySlice = 1000
	limit         = 10
)

// GetAllTask получить задачи с опциями
func (a *Repo) GetAllTask(ctx context.Context, opts ...entity.ListTaskOption) ([]entity.ResponceBDLog, error) {
	options := &entity.ListTasksOptions{}
	for _, opt := range opts {
		opt(options)
	}

	query := `
        WITH filtered_tasks AS (
            SELECT 
                id, 
                audit_log, 
                status_task, 
                attempts
            FROM taskslog
            WHERE id > $1
              AND ($2::text = '' OR status_task = (
                  SELECT id FROM statustask WHERE status_task = $2
              ))
            ORDER BY id ASC
            LIMIT $3
        )
        SELECT 
            ft.id,
            ft.audit_log,
            st.status_task AS status,
            ft.attempts
        FROM filtered_tasks ft
        LEFT JOIN statustask st ON ft.status_task = st.id;
    `

	args := []interface{}{0, options.Status, limit}
	ret := make([]entity.ResponceBDLog, 0, capacitySlice)
	prevID := 0

	for {
		args[0] = prevID

		if options.Count > 0 {
			remaining := int(options.Count) - len(ret)
			if remaining <= 0 {
				break
			}
			args[2] = min(limit, remaining)
		}

		rows, err := a.tx.DB.Query(ctx, query, args...)
		if err != nil {
			a.logger.Error("r.tx.DB.Query", zap.Error(err))

			return nil, fmt.Errorf("r.tx.DB.Query: %w", err)
		}
		defer rows.Close()

		lgs, err := scanOrders(rows, capacitySlice)
		if err != nil {
			a.logger.Error("scanOrders", zap.Error(err))

			return nil, fmt.Errorf("scanOrders: %w", err)
		}

		if len(lgs) == 0 {
			break
		}

		ret = append(ret, lgs...)

		prevID = int(lgs[len(lgs)-1].ID)

		if options.Count > 0 && len(ret) >= int(options.Count) {
			break
		}
	}

	return ret, nil
}

// GetStatusID получить id статуса
func (a *Repo) GetStatusID(ctx context.Context, statusName string) (int64, error) {
	queryStatus := `
	SELECT id FROM StatusTask WHERE status_task=$1
`
	var statusID int64
	if err := a.tx.GetDB(ctx).QueryRow(ctx, queryStatus, statusName).Scan(&statusID); err != nil {
		return statusID, fmt.Errorf("a.db.QueryRow: %w", err)
	}

	return statusID, nil
}

func scanOrders(rows pgx.Rows, capacity int64) ([]entity.ResponceBDLog, error) {
	lgs := make([]entity.ResponceBDLog, 0, capacity)
	for rows.Next() {
		var lg entity.ResponceBDLog
		err := rows.Scan(
			&lg.ID,
			&lg.LG,
			&lg.Status,
			&lg.Attempts,
		)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}

		lgs = append(lgs, lg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err(): %w", err)
	}

	return lgs, nil
}
