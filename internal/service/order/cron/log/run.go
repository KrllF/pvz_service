package log

import (
	"context"
	"fmt"
	"time"

	gocron "github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

// Run запустить крон
func (c *Cron) Run(ctx context.Context) error {
	_, err := c.Cron.NewJob(
		gocron.DurationJob(
			time.Minute,
		),
		gocron.NewTask(func(ctx context.Context) {
			_ = c.FillCache(ctx)
		},
			ctx,
		),
	)
	if err != nil {
		c.logger.Error("c.Cron.NewJob", zap.Error(err))

		return fmt.Errorf("c.Cron.NewJob: %w", err)
	}

	c.Cron.Start()

	<-ctx.Done()

	err = c.Cron.Shutdown()
	if err != nil {
		c.logger.Error("c.Cron.Shutdown", zap.Error(err))

		return fmt.Errorf("c.Cron.Shutdown: %w", err)
	}

	return nil
}
