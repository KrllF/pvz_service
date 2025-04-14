package json

import "github.com/KrllF/pvz_service/internal/models"

func (s *repo) GetAllOrders() ([]models.Order, error) {
	result := make([]models.Order, 0, len(s.Orders))
	for orderID, orderInfo := range s.Orders {
		order := models.Order{
			OrderID:       orderID,
			UserID:        orderInfo.UserID,
			Weight:        orderInfo.Weight,
			Price:         orderInfo.Price,
			Packaging:     orderInfo.Packaging,
			InStorageFrom: orderInfo.InStorageFrom,
			ShelfLife:     orderInfo.ShelfLife,
			Status:        orderInfo.Status,
			TwoDaysOfLife: orderInfo.TwoDaysOfLife,
			LastUpdate:    orderInfo.LastUpdate,
		}
		result = append(result, order)
	}

	return result, nil
}
