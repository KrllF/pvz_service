package postgreslog

import (
	"github.com/KrllF/pvz_service/internal/repository/txmanager"
	"go.uber.org/zap"
)

// Repo структура репо слоя
type Repo struct {
	tx     *txmanager.TxManager
	logger *zap.Logger
}

// NewRepository новый Repo
func NewRepository(tx *txmanager.TxManager, logg *zap.Logger) (*Repo, error) {
	return &Repo{tx: tx, logger: logg}, nil
}
