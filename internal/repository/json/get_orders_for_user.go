package json

import (
	"fmt"

	"github.com/KrllF/pvz_service/internal/models"
)

func (s *repo) GetOrdersForUser(userID int64) ([]models.Order, error) {
	if _, ok := s.Users[userID]; !ok {
		return nil, fmt.Errorf("нет пользователя id%d", userID)
	}
	result := make([]models.Order, 0, len(s.Users[userID]))
	for ord := range s.Users[userID] {
		order := models.Order{
			OrderID:       ord,
			UserID:        s.Orders[ord].UserID,
			Weight:        s.Orders[ord].Weight,
			Price:         s.Orders[ord].Price,
			Packaging:     s.Orders[ord].Packaging,
			InStorageFrom: s.Orders[ord].InStorageFrom,
			ShelfLife:     s.Orders[ord].ShelfLife,
			Status:        s.Orders[ord].Status,
			TwoDaysOfLife: s.Orders[ord].TwoDaysOfLife,
			LastUpdate:    s.Orders[ord].LastUpdate,
		}
		result = append(result, order)
	}

	return result, nil
}
