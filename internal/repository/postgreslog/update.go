package postgreslog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// UpdateStatusTask обновить статус задачи
func (a *Repo) UpdateStatusTask(ctx context.Context, id int64, newStatus string, completedAt time.Time) error {
	query := `
	UPDATE taskslog 
	SET status_task = $1, updated_at = NOW(), completed_at = $2
	WHERE id = $3;
`
	NewStatusID, err := a.getStatusID(ctx, newStatus)
	if err != nil {
		a.logger.Error("a.getStatusID", zap.Error(err))

		return fmt.Errorf("не удалось получить статус: %w", err)
	}

	_, err = a.tx.GetDB(ctx).Exec(ctx, query, NewStatusID, completedAt, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.logger.Info("задача не найдена", zap.Int64("orderID", id), zap.Error(err))

			return fmt.Errorf("задача с id=%d не найдена", id)
		}
		a.logger.Error("не удалось обновить задачу", zap.Error(err))

		return fmt.Errorf("не удалось обновить задачу: %w", err)
	}

	return nil
}

// UpdateAttempts обновить попытку
func (a *Repo) UpdateAttempts(ctx context.Context, id int64) error {
	query := `
        UPDATE TasksLog
        SET attempts = attempts + 1
        WHERE id = $1
    `

	result, err := a.tx.GetDB(ctx).Exec(ctx, query, id)
	if err != nil {
		a.logger.Error("a.tx.GetDB(ctx).Exec", zap.Error(err))

		return fmt.Errorf("a.tx.GetDB(ctx).Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("result.RowsAffected == 0")
	}

	return nil
}

func (a *Repo) getStatusID(ctx context.Context, statusName string) (int64, error) {
	var statusID int64
	query := `SELECT id FROM statustask WHERE status_task = $1`
	err := a.tx.GetDB(ctx).QueryRow(ctx, query, statusName).Scan(&statusID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.logger.Info("стату не найден", zap.String("status", statusName), zap.Error(err))

			return 0, fmt.Errorf("статус '%s' не найден", statusName)
		}
		a.logger.Error("ошибка при получении ID статуса", zap.String("status", statusName), zap.Error(err))

		return 0, fmt.Errorf("ошибка при получении ID статуса '%s': %w", statusName, err)
	}

	return statusID, nil
}
