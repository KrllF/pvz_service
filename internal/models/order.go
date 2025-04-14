package models

import (
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
)

// Order содержит всю информацию о заказе
// Она содержит 7 полей:
// - OrderID: id заказа
// - UserID: id клиента, которому принадлежит заказ
// - InStorageFrom: время, с которого заказ хранится в ПВЗ (UTC+3:00)
// - ShelfLife: время, до которого хранится заказ (UTC+3:00)
// - Status: статус заказа ("accepted", "delivered", "returned")
// - TwoDayOfLife: время, до которого можно вернуть заказ клиенту, если он его забрал (UTC+3:00)
// - LastUpdate: время, когда у заказа последний раз изменился статус (UTC+3:00)
type Order struct {
	OrderID       int64       `json:"orderid"`
	UserID        int64       `json:"userid"`
	Weight        int64       `json:"weight"`
	Price         int64       `json:"price"`
	Packaging     entity.Pack `json:"packaging"`
	InStorageFrom time.Time   `json:"instoragefrom"`
	ShelfLife     time.Time   `json:"shelflife"`
	Status        string      `json:"status"`
	TwoDaysOfLife time.Time   `json:"twodaysoflife"`
	LastUpdate    time.Time   `json:"lastupdate"`
}

// OrderInfo содержит всю информацию о заказе
// Она содержит 6 полей:
// - UserID: id клиента, которому принадлежит заказ
// - InStorageFrom: время, с которого заказ хранится в ПВЗ (UTC+3:00)
// - ShelfLife: время, до которого хранится заказ (UTC+3:00)
// - Status: статус заказа ("accepted", "delivered", "returned")
// - TwoDayOfLife: время, до которого можно вернуть заказ клиенту, если он его забрал (UTC+3:00)
// - LastUpdate: время, когда у заказа последний раз изменился статус (UTC+3:00)
type OrderInfo struct {
	UserID        int64       `json:"userid"`
	Weight        int64       `json:"weight"`
	Price         int64       `json:"price"`
	Packaging     entity.Pack `json:"packaging"`
	InStorageFrom time.Time   `json:"instoragefrom"`
	ShelfLife     time.Time   `json:"shelflife"`
	Status        string      `json:"status"` // "accepted", "delivered", "returned"
	TwoDaysOfLife time.Time   `json:"twodaysoflife"`
	LastUpdate    time.Time   `json:"lastupdate"`
}

func (o Order) String() string {
	return fmt.Sprintf(`Order Details:
  Order ID:       %d
  User ID:        %d
  Weight:		  %d
  Price:		  %d
  PackType		  %s
  ExtraPack		  %v
  In Storage From:%s
  Shelf Life:     %s
  Status:         %s
  Two Days of Life: %s
  Last Update:    %s`,
		o.OrderID,
		o.UserID,
		o.Weight,
		o.Price,
		o.Packaging.PackType,
		o.Packaging.ExtraPack,
		o.InStorageFrom.Format("2006-01-02 15:04:05"),
		o.ShelfLife.Format("2006-01-02 15:04:05"),
		o.Status,
		o.TwoDaysOfLife.Format("2006-01-02 15:04:05"),
		o.LastUpdate.Format("2006-01-02 15:04:05"),
	)
}
