package brocker

import (
	"context"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/entity"
	gocron "github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

const (
	timeFailed  = 2     // s
	timeCreated = 10000 // ms
	countFix    = 10
)

// Run запустить крон
func (c *Cron) Run(ctx context.Context) error {
	_, err := c.CronFailed.NewJob(
		gocron.DurationJob(
			timeFailed*time.Second,
		),
		gocron.NewTask(func(ctx context.Context) {
			_ = c.SendMessage(ctx, entity.WithTaskStatus(consts.FailedTask), entity.WithTaskCount(countFix))
		},
		),
	)
	if err != nil {
		c.logger.Error("c.Cron.NewJob", zap.Error(err))

		return fmt.Errorf("c.Cron.NewJob: %w", err)
	}
	_, err = c.CronCreated.NewJob(
		gocron.DurationJob(
			timeCreated*time.Microsecond,
		),
		gocron.NewTask(func(ctx context.Context) {
			_ = c.SendMessage(ctx, entity.WithTaskStatus(consts.CreatedTask), entity.WithTaskCount(countFix))
		},
			ctx,
		),
	)
	if err != nil {
		c.logger.Error("c.Cron.NewJob", zap.Error(err))

		return fmt.Errorf("c.Cron.NewJob: %w", err)
	}

	c.CronCreated.Start()
	c.CronFailed.Start()

	<-ctx.Done()

	err = c.CronCreated.Shutdown()
	if err != nil {
		c.logger.Error("c.Cron.Shutdown", zap.Error(err))

		return fmt.Errorf("c.Cron.Shutdown: %w", err)
	}

	err = c.CronFailed.Shutdown()
	if err != nil {
		c.logger.Error("c.Cron.Shutdown", zap.Error(err))

		return fmt.Errorf("c.Cron.Shutdown: %w", err)
	}

	return nil
}
