package postgres

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// DeleteUser удалить клиента
func (r *Repo) DeleteUser(ctx context.Context, userID int64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "DeleteUser from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	count, err := r.tx.GetDB(ctx).Exec(ctx, `DELETE FROM users WHERE user_id = $1`, userID)
	if err != nil {
		r.logger.Fatal("ошибка при удалении пользователя", zap.Error(err))

		return fmt.Errorf("ошибка при удалении пользователя: %w", err)
	}
	if count.RowsAffected() == 0 {
		r.logger.Error("Пользователя нет", zap.Int64("userID", userID))

		return errs.ErrUserNotFound
	}
	r.logger.Info("Пользователь успешно удален", zap.Int64("userID", userID))

	return nil
}
