package helper

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
)

// ReadFile прочитать json файл
func ReadFile(path string) ([]models.Order, error) {
	f, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Файл не найден: %s", path)

			return nil, errs.ErrFileNotFound
		}
		log.Printf("Ошибка при чтении файла: %v", err)

		return nil, errs.ErrFileReadError
	}

	orders := make([]models.Order, 0)

	err = json.Unmarshal(f, &orders)
	if err != nil {
		return nil, errs.ErrInvalidData
	}

	return orders, nil
}
