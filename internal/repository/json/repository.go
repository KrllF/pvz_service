package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KrllF/pvz_service/internal/models"
)

// Структура repo - репозиторий для хранения данных о заказах и пользователях
// Она содержит 2 поля:
// - Orders: map, где ключ - ID заказа, а значением - информация о заказе (OrderInfo),
// - Users: map, где внешний ключ - ID пользователя, а внутренняя map содержит ID заказов,
//   связанных с этим пользователем. Значение во внутренней map - пустая структура (struct{}),
//	 важен только ключ,
// - path: путь к файлу JSON, из которого загружаются данные.

type repo struct {
	Orders map[int64]models.OrderInfo
	Users  map[int64]map[int64]struct{}
	path   string
}

// NewStorageRepository создаёт новый экземляр репозитория
// Принимает путь к JSON файлу - возвращает указатель на созданный репозиторий или ошибку
func NewStorageRepository(path string) (*repo, error) {
	if _, err := os.ReadFile(filepath.Clean(path)); os.IsNotExist(err) {
		return &repo{
			Orders: make(map[int64]models.OrderInfo),
			Users:  make(map[int64]map[int64]struct{}),
			path:   path,
		}, nil
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать: %w", err)
	}

	orders := make([]models.Order, 0)

	err = json.Unmarshal(f, &orders)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	mpOrd := make(map[int64]models.OrderInfo, len(orders))
	mpUs := make(map[int64]map[int64]struct{}, len(orders))
	for _, val := range orders {
		mpOrd[val.OrderID] = models.OrderInfo{
			UserID:        val.UserID,
			Weight:        val.Weight,
			Price:         val.Price,
			Packaging:     val.Packaging,
			InStorageFrom: val.InStorageFrom,
			ShelfLife:     val.ShelfLife,
			Status:        val.Status,
			TwoDaysOfLife: val.TwoDaysOfLife,
			LastUpdate:    val.LastUpdate,
		}
		if _, ok := mpUs[val.UserID]; !ok {
			mpUs[val.UserID] = make(map[int64]struct{})
		}
		mpUs[val.UserID][val.OrderID] = struct{}{}
	}

	return &repo{
		Orders: mpOrd,
		Users:  mpUs,
		path:   path,
	}, nil
}
