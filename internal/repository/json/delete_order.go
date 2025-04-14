package json

import "fmt"

func (s *repo) DeleteOrder(orderID int64) error {
	order, ok := s.Orders[orderID]
	if !ok {
		return fmt.Errorf("нет заказа id %d", orderID)
	}

	delete(s.Users[order.UserID], orderID)

	if len(s.Users[order.UserID]) == 0 {
		delete(s.Users, order.UserID)
	}

	delete(s.Orders, orderID)

	err := s.saveOrderToFile()
	if err != nil {
		return fmt.Errorf("не удалось сохранить изменения: %w", err)
	}

	return nil
}
