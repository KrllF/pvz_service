package entity

import "strings"

const (
	sortDesc = "desc"
	sortAsc  = "asc"
)

// ListOrdersOptions - паттерн опций
type ListOrdersOptions struct {
	UserID    int64
	OrderID   int64
	Status    []string
	Limit     int64
	Page      int64
	SortBy    string
	SortOrder string
	FindByID  string
}

// ListOrdersOption функции опций
type ListOrdersOption func(*ListOrdersOptions)

// WithUserID заказы по userID
func WithUserID(userID int64) ListOrdersOption {
	return func(opts *ListOrdersOptions) {
		opts.UserID = userID
	}
}

// WithOrderID заказы по orderID
func WithOrderID(orderID int64) ListOrdersOption {
	return func(opts *ListOrdersOptions) {
		opts.OrderID = orderID
	}
}

// WithPagination заказы с пагинацией
func WithPagination(limit, page int64) ListOrdersOption {
	return func(opts *ListOrdersOptions) {
		opts.Limit = limit
		opts.Page = page
	}
}

// WithStatus заказы по status
func WithStatus(status ...string) ListOrdersOption {
	return func(opts *ListOrdersOptions) {
		opts.Status = append(opts.Status, status...)
	}
}

// WithSorting заказы отсортированные по столбцу sortBy по направлению sortOrder
func WithSorting(sortBy, sortOrder string) ListOrdersOption {
	return func(opts *ListOrdersOptions) {
		opts.SortBy = strings.ToLower(sortBy)
		opts.SortOrder = sortAsc
		if strings.ToLower(sortOrder) == sortDesc {
			opts.SortOrder = sortDesc
		}
	}
}

// WithFindByID заказы поиск по order_id
func WithFindByID(orderID string) ListOrdersOption {
	return func(opts *ListOrdersOptions) {
		opts.FindByID = orderID
	}
}
