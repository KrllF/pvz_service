package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KrllF/pvz_service/internal/models"
)

const filePermissions = 0o600

func (s *repo) saveOrderToFile() error {
	f, err := os.OpenFile(filepath.Clean(s.path), os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePermissions)
	if err != nil {
		return fmt.Errorf("ошибка при открытие файла: %w", err)
	}
	defer f.Close()

	ordersToSave := make([]*models.Order, 0, len(s.Orders))
	for orderID, orderInfo := range s.Orders {
		order := &models.Order{
			OrderID:       orderID,
			UserID:        orderInfo.UserID,
			Weight:        orderInfo.Weight,
			Price:         orderInfo.Price,
			Packaging:     orderInfo.Packaging,
			ShelfLife:     orderInfo.ShelfLife,
			InStorageFrom: orderInfo.InStorageFrom,
			Status:        orderInfo.Status,
			TwoDaysOfLife: orderInfo.TwoDaysOfLife,
			LastUpdate:    orderInfo.LastUpdate,
		}
		ordersToSave = append(ordersToSave, order)
	}

	b, err := json.MarshalIndent(ordersToSave, "", "\t")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %w", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("f.Write: %w", err)
	}

	return nil
}
