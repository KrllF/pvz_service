package json

import (
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/consts"
)

const (
	hourInTwoDays = 48
)

func (s *repo) UpdateStatus(orderID int64, newStatus string) error {
	if _, ok := s.Orders[orderID]; !ok {
		return fmt.Errorf("нет заказа id %d", orderID)
	}
	if newStatus == consts.AcceptedSt || newStatus == consts.ReturnedSt {
		ord := s.Orders[orderID]
		ord.Status = newStatus
		ord.LastUpdate = time.Now()
		s.Orders[orderID] = ord
		if err := s.saveOrderToFile(); err != nil {
			return fmt.Errorf("не удалось сохранить изменения заказа id%d: %w", orderID, err)
		}
	}
	if newStatus == consts.AcceptedSt {
		ord := s.Orders[orderID]
		ord.TwoDaysOfLife = time.Now().Add(hourInTwoDays * time.Hour)
		s.Orders[orderID] = ord
	}

	return nil
}
