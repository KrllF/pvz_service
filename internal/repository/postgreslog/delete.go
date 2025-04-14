package postgreslog

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// DeleteTask удалить задачу
func (a *Repo) DeleteTask(ctx context.Context, id int64) error {
	query := `
			DELETE FROM taskslog WHERE id=$1 returning id
		`
	var ID int64
	err := a.tx.GetDB(ctx).QueryRow(ctx, query, id).Scan(&ID)
	if err != nil {
		a.logger.Error("a.tx.GetDB(ctx).QueryRow", zap.Error(err))

		return fmt.Errorf("a.tx.GetDB(ctx).QueryRow: %w", err)
	}

	return nil
}
