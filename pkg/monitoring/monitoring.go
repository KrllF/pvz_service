package monitoring

import (
	"log"
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	orderPrice = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "order_price_without_packaging",
		Help:    "price histogram of new orders",
		Buckets: []float64{1000.0, 1500.0, 2000.0, 2500.0, 3000.0},
	})
	orderWeight = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "order_weight_without_packaging",
		Help:    "weight histogram of new orders",
		Buckets: []float64{100.0, 200.0, 300.0, 400.0, 500.0, 1000.0, 2000.0},
	})

	counterMainPackaging = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "packaging_total",
		Help: "the number of packages of different types",
	}, []string{"pack_type"})
	counterAcceptOrder = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "accept_orders_total",
		Help: "number of orders accepted",
	})
	counterReturnOrder = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "return_orders_total",
		Help: "number of refunds",
	})
	counterReq = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_total",
		Help: "number of requests",
	}, []string{"handler"})
	counterOKERR = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_okerr_total",
		Help: "number of ok or err requests",
	}, []string{"status"})
)

// OrderPrice добавить
func OrderPrice(price int64) {
	orderPrice.Observe(float64(price))
}

// OrderWeight добавить
func OrderWeight(weight int64) {
	orderWeight.Observe(float64(weight))
}

// SetCounterMainPackaging увеличить счётчик
func SetCounterMainPackaging(pack string) {
	counterMainPackaging.With(prometheus.Labels{"pack_type": pack}).Inc()
}

// SetCountAcceptOrder увеличить счётчик
func SetCountAcceptOrder(counter int64) {
	counterAcceptOrder.Add(float64(counter))
}

// SetCountReturnOrder увеличить счётчик
func SetCountReturnOrder(counter int64) {
	counterReturnOrder.Add(float64(counter))
}

// SetCounterReq увеличить счётчик
func SetCounterReq(handler string) {
	counterReq.With(prometheus.Labels{"handler": handler}).Inc()
}

// SetcounterOKERR увеличить счётчик
func SetcounterOKERR(status string) {
	counterOKERR.With(prometheus.Labels{"status": status}).Inc()
}

// nolint
func init() {
	prometheus.MustRegister(orderPrice, counterMainPackaging, counterAcceptOrder,
		counterReturnOrder, counterReq, counterOKERR)
}

// StartMetricsServer запустить сервер Prometheus
func StartMetricsServer(conf config.Config) {
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":" + conf.ConnectConfig.PROMETHEUS_PORT,
		ReadTimeout:  conf.AppConfig.Read,
		WriteTimeout: conf.AppConfig.Write,
		IdleTimeout:  conf.AppConfig.Idle,
	}

	go func() {
		log.Println("Сервер Prometheus прослушивается...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()
}
