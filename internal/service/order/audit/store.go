package audite

import (
	"github.com/KrllF/pvz_service/internal/entity"
)

// StoreHTTP запись информации о запросе в канал
func (a *Audit) StoreHTTP(h entity.HTTPInfo) {
	a.HTTP <- h
}

// StoreStatus запись нового статуса в канал
func (a *Audit) StoreStatus(h entity.UpdateStatus) {
	a.Status <- h
}
