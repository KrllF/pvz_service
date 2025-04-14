//go:generate mockgen -source=$GOFILE -destination=../mocks/mocks.go -package=mocks

package httph

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	"github.com/KrllF/pvz_service/middleware"
	"github.com/gorilla/mux"
)

const (
	queryParamKey = "id"
)

type (
	// AcceptService определяет методы для управления заказами.
	// Содержит 3 функций:
	// AcceptOrder - принять заказ от курьера
	// BulkAcceptOrders - принять заказы от курьера через JSON-файл
	// ProcessClientIssue - выдать заказы клиента
	AcceptService interface {
		AcceptOrder(ctx context.Context, req entity.OrderEntry) error
		BulkAcceptOrders(ctx context.Context, path string) (int, error)
		ProcessClientIssue(ctx context.Context, userID int64, orderIDs []int64) ([]int64, error)
	}

	// ReturnsService определяет методы для управления заказами.
	// Содержит 2 функций:
	// ReturnOrder - вернуть заказ курьеру
	// ProcessClientReturn - принять возвраты клиента
	ReturnsService interface {
		ProcessClientReturn(ctx context.Context, userID int64, orderIDs []int64) ([]int64, error)
		ReturnOrder(ctx context.Context, orderID int64) error
	}

	// ShowService определяет методы для управления заказами.
	// Содержит 4 функций:
	// ListOrdersUserAllN - список заказов для клиента
	// ListOrdersUserPVZ - список заказов в пвз у клиента
	// ListReturns - список возвращённых заказов с пагинацией
	// ListHistory - история заказов
	ShowService interface {
		ListHistory(ctx context.Context) ([]models.Order, error)
		ListOrdersUserAllN(ctx context.Context, userID int64, opt, limit, page int64, pag bool,
			finder bool, orderFinder string) ([]models.Order, error)
		ListOrdersUserPVZ(ctx context.Context, userID, limit, page int64,
			pag bool, finder bool, orderFinder string) ([]models.Order, error)
		ListReturns(ctx context.Context, limit, page int64, finder bool, orderFinder string) ([]models.Order, int, error)
	}
	// AuditService определяет методы для аудит-логов.
	// Содержит 3 функции:
	// Run - запускает воркеры
	// StorageHttp - сохраняет логи хедлеров
	// StorageStatus - сохраняет логи обновления статуса
	AuditService interface {
		StoreHTTP(h entity.HTTPInfo)
		StoreStatus(h entity.UpdateStatus)
	}
	// Handler хэндлер ПВЗ
	Handler struct {
		ctx         context.Context
		acceptServ  AcceptService
		returnsServ ReturnsService
		showServ    ShowService
		auditServ   AuditService
	}
)

// NewHandler конструктор нового хедлера
func NewHandler(ctx context.Context, accept AcceptService,
	returns ReturnsService, show ShowService, audit AuditService,
) *Handler {
	return &Handler{ctx: ctx, acceptServ: accept, returnsServ: returns, showServ: show, auditServ: audit}
}

// Init инициализация роутера
func (h *Handler) Init(cfg config.Config) *mux.Router {
	handler := mux.NewRouter()
	handler.Use(middleware.LogHandlerMiddleware(h.auditServ))
	handler.Use(middleware.BasicAuthMiddleware(cfg.ConnectConfig))
	handler.Use(middleware.SpanMiddleware())

	handler.HandleFunc("/accept/", h.AcceptOrder).Methods(http.MethodPost)
	handler.HandleFunc("/accept/bulk", h.BulkAcceptOrders).Methods(http.MethodPost)

	processHandler := handler.PathPrefix("/process").Subrouter()
	processHandler.Use(middleware.LogStatusUpdateMiddleware(h.auditServ))
	processHandler.HandleFunc("", h.ProcessClient).Methods(http.MethodPut)

	handler.HandleFunc(fmt.Sprintf("/return/{%s:[0-9]+}", queryParamKey), h.ReturnOrderToCourier).
		Methods(http.MethodDelete)
	handler.HandleFunc("/history/", h.ListHistory).Methods(http.MethodGet)
	handler.HandleFunc("/history/returns", h.ListReturns).Methods(http.MethodGet)
	handler.HandleFunc(fmt.Sprintf("/history/orders/{%s:[0-9]+}", queryParamKey), h.ListOrdersUser).Methods(http.MethodGet)
	handler.HandleFunc("/", h.DefaultHandler)

	return handler
}
