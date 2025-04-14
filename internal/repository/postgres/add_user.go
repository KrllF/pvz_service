package postgres

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

// AddUser добавить клиента
func (r *Repo) AddUser(ctx context.Context, userID int64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "AddUser from postgres-repository")
	defer span.Finish()
	var err error

	defer func() {
		if err != nil {
			span.SetTag("error", true)
			span.LogKV("error", err.Error())
		}
	}()

	_, err = r.tx.GetDB(ctx).Exec(ctx, `INSERT INTO users(user_id) VALUES($1)`, userID)
	if err != nil {
		r.logger.Fatal("ошибка при добавлении пользователя", zap.Error(err))

		return fmt.Errorf("ошибка при добавлении пользователя: %w", err)
	}

	return nil
}
