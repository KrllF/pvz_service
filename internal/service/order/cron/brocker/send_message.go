package brocker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
	"go.uber.org/zap"
)

const (
	lastAttempt = 2
	batchSize   = 5
)

// SendMessage отправить сообщение в брокер
func (c *Cron) SendMessage(ctx context.Context, opts ...entity.ListTaskOption) error {
	var wg sync.WaitGroup
	lg, err := c.GetLog(ctx, opts...)
	if err != nil {
		c.logger.Error("c.GetAllTask", zap.Error(err))

		return fmt.Errorf("c.GetAllTask: %w", err)
	}
	if len(lg) == 0 {
		return errors.New("len(ret) == 0")
	}
	for i := 0; i < len(lg); i += batchSize {
		end := i + batchSize
		if end > len(lg) {
			end = len(lg)
		}

		batch := lg[i:end]
		wg.Add(1)
		go func(batch []entity.ResponceBDLog) {
			defer wg.Done()
			for _, val := range batch {
				select {
				case <-ctx.Done():
					return
				default:
					_ = c.executor(ctx, val)
				}
			}
		}(batch)
	}
	wg.Wait()

	return nil
}

func (c *Cron) executor(ctx context.Context, val entity.ResponceBDLog) error {
	if err := c.Repo.UpdateAttempts(ctx, val.ID); err != nil {
		c.logger.Error("c.Repo.UpdateAttempts", zap.Error(err))

		return fmt.Errorf("c.Repo.UpdateAttempts: %w", err)
	}

	if err := c.Repo.UpdateStatusTask(ctx, val.ID, consts.ProcessingTask, time.Time{}); err != nil {
		c.logger.Error("c.Repo.UpdateStatusTask", zap.Error(err))

		return fmt.Errorf("c.Repo.UpdateStatusTask: %w", err)
	}
	err := c.Brock.SendMessage(val.LG, c.TopicName)
	if err == nil {
		if err := c.Repo.UpdateStatusTask(ctx, val.ID, consts.CompletedTask, time.Now().UTC()); err != nil {
			c.logger.Error("c.Repo.UpdateStatusTask", zap.Error(err))

			return fmt.Errorf("c.Repo.UpdateStatusTask: %w", err)
		}

		return nil
	}

	if val.Attempts == lastAttempt {
		if errStatus := c.Repo.UpdateStatusTask(ctx, val.ID, consts.NoAttemptsTask, time.Time{}); errStatus != nil {
			c.logger.Error("c.Repo.UpdateStatusTask", zap.Error(errStatus))
			c.logger.Error("c.Brock.SendMessage", zap.Error(err))

			return fmt.Errorf("c.Repo.UpdateStatusTask: %w, c.Brock.SendMessage: %w", errStatus, err)
		}

		return fmt.Errorf("c.Brock.SendMessage: %w", err)
	}
	if errStatus := c.Repo.UpdateStatusTask(ctx, val.ID, consts.FailedTask, time.Time{}); errStatus != nil {
		c.logger.Error("c.Repo.UpdateStatusTask", zap.Error(errStatus))
		c.logger.Error("c.Brock.SendMessage", zap.Error(err))

		return fmt.Errorf("c.Repo.UpdateStatusTask: %w, c.Brock.SendMessage: %w", errStatus, err)
	}

	return fmt.Errorf("c.Brock.SendMessage: %w", err)
}
