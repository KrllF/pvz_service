package brocker

import (
	"context"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	gocron "github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type (
	// Repository интерфейс
	Repository interface {
		GetAllTask(ctx context.Context, opts ...entity.ListTaskOption) ([]entity.ResponceBDLog, error)
		UpdateStatusTask(ctx context.Context, id int64, newStatus string, completedAt time.Time) error
		UpdateAttempts(ctx context.Context, id int64) error
	}
	// Brocker интерфейс
	Brocker interface {
		SendMessage(message []byte, topicname string) error
	}
	// TXManager интерфейс менеджера транзакций
	TXManager interface {
		RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
		RunReadCommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	}
	// Cron структура
	Cron struct {
		Repo        Repository
		Brock       Brocker
		TopicName   string
		txManager   TXManager
		CronCreated gocron.Scheduler
		CronFailed  gocron.Scheduler
		logger      *zap.Logger
	}
)

// NewCron конструктор крона
func NewCron(repo Repository, brocker Brocker, tx TXManager, topicName string, logg *zap.Logger) (*Cron, error) {
	cronFailed, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	cronCreated, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	cr := Cron{
		Repo:        repo,
		Brock:       brocker,
		TopicName:   topicName,
		txManager:   tx,
		CronCreated: cronCreated,
		CronFailed:  cronFailed,
		logger:      logg,
	}

	return &cr, nil
}
