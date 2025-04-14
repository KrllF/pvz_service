package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KrllF/pvz_service/internal/models"
)

func (s *repo) ReadFile(path string) ([]models.Order, error) {
	f, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать: %w", err)
	}

	orders := make([]models.Order, 0)

	err = json.Unmarshal(f, &orders)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return orders, nil
}
