package entity

// HTTPInfo хранит информацию об HTTP-запросе
type HTTPInfo struct {
	Method       string `json:"method"`
	Request      string `json:"request"`
	ResponseCode int    `json:"code"`
}

// UpdateStatus хранит информацию об изменении статуса заказа
type UpdateStatus struct {
	OrderID   int64  `json:"orderid"`
	NewStatus string `json:"status"`
}
