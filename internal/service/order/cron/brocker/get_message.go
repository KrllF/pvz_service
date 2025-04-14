package brocker

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
	"go.uber.org/zap"
)

// GetLog получить логи
func (c *Cron) GetLog(ctx context.Context, opts ...entity.ListTaskOption) ([]entity.ResponceBDLog, error) {
	var lg []entity.ResponceBDLog
	var err error
	err = c.txManager.RunReadCommitted(ctx, func(ctxTx context.Context) error {
		lg, err = c.Repo.GetAllTask(ctxTx, opts...)
		if err != nil {
			c.logger.Error("c.Repository.GetAllTask", zap.Error(err))

			return fmt.Errorf("c.Repository.GetAllTask: %w", err)
		}

		return nil
	})
	if err != nil {
		c.logger.Error("c.txManager.RunSerializable", zap.Error(err))

		return nil, fmt.Errorf("c.txManager.RunSerializable: %w", err)
	}
	ret := make([]entity.ResponceBDLog, 0, len(lg))
	for _, val := range lg {
		if (val.Status != consts.CompletedTask) && (val.Status != consts.NoAttemptsTask) {
			ret = append(ret, val)
		}
	}

	return ret, nil
}
