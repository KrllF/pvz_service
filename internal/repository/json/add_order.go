package json

import (
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/models"
)

func (s *repo) AddOrder(ord models.Order) error {
	if _, ok := s.Orders[ord.OrderID]; ok {
		return fmt.Errorf("нельзя дважды принять один заказ id%d", ord.OrderID)
	}
	s.Orders[ord.OrderID] = models.OrderInfo{
		UserID:        ord.UserID,
		Weight:        ord.Weight,
		Price:         ord.Price,
		Packaging:     ord.Packaging,
		InStorageFrom: time.Now(),
		ShelfLife:     ord.ShelfLife,
		Status:        ord.Status,
		TwoDaysOfLife: time.Time{},
		LastUpdate:    ord.LastUpdate,
	}
	if _, ok := s.Users[ord.UserID]; !ok {
		s.Users[ord.UserID] = make(map[int64]struct{})
	}
	s.Users[ord.UserID][ord.OrderID] = struct{}{}

	if err := s.saveOrderToFile(); err != nil {
		return fmt.Errorf("не удалось добавь заказ в файл: %w", err)
	}

	return nil
}
