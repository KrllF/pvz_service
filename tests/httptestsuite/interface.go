package httptestsuite

import (
	"context"
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DB interface
type DB interface {
	Close()
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	GetPool() *pgxpool.Pool
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
}

// Handler interface
type Handler interface {
	AcceptOrder(w http.ResponseWriter, r *http.Request)
	BulkAcceptOrders(w http.ResponseWriter, r *http.Request)
	DefaultHandler(w http.ResponseWriter, _ *http.Request)
	Init(cfg config.Config) *mux.Router
	ListHistory(w http.ResponseWriter, r *http.Request)
	ListOrdersUser(w http.ResponseWriter, r *http.Request)
	ListReturns(w http.ResponseWriter, r *http.Request)
	ProcessClient(w http.ResponseWriter, r *http.Request)
	ReturnOrderToCourier(w http.ResponseWriter, r *http.Request)
}

// AccService interface
type AccService interface {
	AcceptOrder(ctx context.Context, req entity.OrderEntry) error
	BulkAcceptOrders(ctx context.Context, path string) (int, error)
	ProcessClientIssue(ctx context.Context, userID int64, orderIDs []int64) (int64, error)
}

// ShowService interface
type ShowService interface {
	ListHistory(ctx context.Context) ([]models.Order, error)
	ListOrdersUserAllN(ctx context.Context, userID int64,
		opt int64, limit int64, page int64, pag bool) ([]models.Order, error)
	ListOrdersUserPVZ(ctx context.Context, userID int64, limit int64, page int64, pag bool) ([]models.Order, error)
	ListReturns(ctx context.Context, limit int64, page int64) ([]models.Order, int, error)
}

// RetService interface
type RetService interface {
	ProcessClientReturn(ctx context.Context, userID int64, orderIDs []int64) (int64, error)
	ReturnOrder(ctx context.Context, orderID int64) error
}

// Repo interface
type Repo interface {
	AddOrder(ctx context.Context, ord models.Order) error
	AddUser(ctx context.Context, userID int64) error
	DeleteOrder(ctx context.Context, orderID int64) error
	DeleteUser(ctx context.Context, userID int64) error
	GetOrderInfo(ctx context.Context, orderID int64) (models.Order, error)
	ListOrders(ctx context.Context, opts ...entity.ListOrdersOption) ([]models.Order, error)
	SetTwoDaysOfLife(ctx context.Context, orderID int64) error
	UpdateStatus(ctx context.Context, orderID int64, newStatus string) error
	UserExist(ctx context.Context, userID int64) (bool, error)
}

// Server interface
type Server interface {
	Close()
	Run() error
}
