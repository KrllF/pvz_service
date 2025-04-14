package entity

import (
	"time"
)

// OrderEntry структура для передачи иформации о заказе, чтобы добавить в базу
type OrderEntry struct {
	OrderID   int64     `json:"order_id"`
	UserID    int64     `json:"user_id"`
	Weight    int64     `json:"weight"`
	Price     int64     `json:"price"`
	PackType  PackType  `json:"pack_type"`
	ExtraPack PackType  `json:"extra_pack,omitempty"`
	Shelflife time.Time `json:"shelf_life"`
}
