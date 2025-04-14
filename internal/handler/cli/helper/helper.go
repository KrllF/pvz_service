package helper

import (
	"fmt"
	"strconv"
)

// ParseID кастит string к int64
func ParseID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат id: %w", err)
	}

	return id, nil
}

// ParseIDs кастит []string к []int64
func ParseIDs(idStrs []string) ([]int64, error) {
	ids := make([]int64, len(idStrs))
	for i, idStr := range idStrs {
		id, err := ParseID(idStr)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}

	return ids, nil
}

// ParseOptionalInt испульзуется для каста параметр paramName, значения string, к int64
func ParseOptionalInt(value, paramName string) (int64, error) {
	if value == "" {
		return 0, nil
	}
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат параметра %s: %w", paramName, err)
	}

	return num, nil
}
