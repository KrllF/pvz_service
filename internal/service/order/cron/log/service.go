package log

import (
	"context"

	"github.com/KrllF/pvz_service/internal/models"
	gocron "github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type (
	// Repository интерфейс
	Repository interface {
		GetAll(ctx context.Context, limit int64) ([]models.Order, error)
	}
	// Cache интерфейс
	Cache interface {
		UpdateAllCache(items []models.Order, getKey func(ord models.Order) int64) error
	}
	// Cron структура
	Cron struct {
		Repository
		Cache
		Cron   gocron.Scheduler
		logger *zap.Logger
	}
)

// NewCron конструктор крона
func NewCron(repository Repository, c Cache, logg *zap.Logger) (*Cron, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	cr := Cron{
		Repository: repository,
		Cache:      c,
		Cron:       s,
		logger:     logg,
	}

	return &cr, nil
}
