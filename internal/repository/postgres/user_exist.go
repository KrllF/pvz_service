package postgres

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// UserExist проверка наличия пользователя
func (r *Repo) UserExist(ctx context.Context, userID int64) (bool, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "UserExist from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	var ok bool
	query := `
	SELECT EXISTS(SELECT 1 from users WHERE user_id=$1)`
	err = r.tx.GetDB(ctx).QueryRow(ctx, query, userID).Scan(&ok)
	if err != nil {
		r.logger.Fatal("не удалось получить информацию о клиенте", zap.Error(err))

		return ok, fmt.Errorf("не удалось получить информацию о клиенте: %w", err)
	}

	return ok, nil
}
