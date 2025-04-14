//nolint:all
package postgreslog

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/entity"
	"go.uber.org/zap"
)

func (a *Repo) AddStatusLog(ctx context.Context, lg entity.UpdateStatus) error {
	queryStatus := `
			INSERT INTO StatusLog(order_id, new_status) VALUES($1, $2)
		`
	_, err := a.tx.GetDB(ctx).Exec(ctx, queryStatus, lg.OrderID, lg.NewStatus)
	if err != nil {
		a.logger.Error("Ошибка при записи status лога в таблицу", zap.Error(err))

		return fmt.Errorf("ошибка при status лога в таблицу: %w", err)
	}

	return nil
}

func (a *Repo) AddHandlerLog(ctx context.Context, lg entity.HTTPInfo) error {
	queryHTTP := `
			INSERT INTO HandlerLog(method, request, code) VALUES($1, $2, $3)
		`

	_, err := a.tx.GetDB(ctx).Exec(ctx, queryHTTP, lg.Method, lg.Request, lg.ResponseCode)
	if err != nil {
		a.logger.Error("Ошибка при записи хэндлер лога в таблицу", zap.Error(err))

		return fmt.Errorf("ошибка при записи хэндлер лога в таблицу: %w", err)

	}

	return nil
}

func (a *Repo) AddTasks(ctx context.Context, lg []byte, statusID int64) error {
	queryTask := `
		INSERT INTO TasksLog(audit_log, status_task)
		VALUES($1, $2)
	`

	_, err := a.tx.GetDB(ctx).Exec(ctx, queryTask, lg, statusID)
	if err != nil {
		return fmt.Errorf("a.db.Exec: %w", err)
	}

	return nil
}
