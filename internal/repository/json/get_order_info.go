package json

import (
	"fmt"

	"github.com/KrllF/pvz_service/internal/models"
)

func (s *repo) GetOrderInfo(orderID int64) (models.OrderInfo, error) {
	order, ok := s.Orders[orderID]
	if !ok {
		return models.OrderInfo{}, fmt.Errorf("нет заказа id%d", orderID)
	}

	return order, nil
}
