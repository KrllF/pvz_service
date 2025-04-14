package audite

import (
	"context"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/entity"
	"go.uber.org/zap"
)

const (
	sizeChan = 100
)

type (
	repo interface {
		AddStatusLog(ctx context.Context, lg entity.UpdateStatus) error
		AddHandlerLog(ctx context.Context, lg entity.HTTPInfo) error
		AddTasks(ctx context.Context, lg []byte, statusID int64) error
		GetStatusID(ctx context.Context, statusName string) (int64, error)
	}
	// TXManager интерфейс менеджера транзакций
	TXManager interface {
		RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
	}
	// Audit аудит сервис
	Audit struct {
		Repo   repo
		TX     TXManager
		Conf   config.Config
		HTTP   chan entity.HTTPInfo
		Status chan entity.UpdateStatus
		DB     chan interface{}
		Stdout chan interface{}
		logger *zap.Logger
	}
)

// NewAudite конструктор аудит сервиса
func NewAudite(db repo, tx TXManager, conf config.Config, logg *zap.Logger) *Audit {
	return &Audit{
		Repo: db, TX: tx, Conf: conf,
		HTTP: make(chan entity.HTTPInfo, sizeChan), Status: make(chan entity.UpdateStatus, sizeChan),
		DB: make(chan interface{}, sizeChan), Stdout: make(chan interface{}, sizeChan),
		logger: logg,
	}
}
