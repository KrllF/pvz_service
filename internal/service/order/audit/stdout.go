package audite

import (
	"github.com/KrllF/pvz_service/internal/entity"
	"go.uber.org/zap"
)

// StdoutLog вывод логов в stdout
func (a *Audit) StdoutLog(lg interface{}) {
	switch v := lg.(type) {
	case entity.UpdateStatus:
		if v.NewStatus != a.Conf.AppConfig.Word {
			a.logger.Info("update status", zap.Any("log", v))
		}
	default:
		a.logger.Info("handler log", zap.Any("log", v))
	}
}
