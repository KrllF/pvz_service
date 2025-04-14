package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/KrllF/pvz_service/internal/config"
	"github.com/KrllF/pvz_service/internal/entity"
	g "github.com/KrllF/pvz_service/internal/handler/grpc"
	intercep "github.com/KrllF/pvz_service/internal/handler/grpc/interceptor"
	"github.com/KrllF/pvz_service/pkg/monitoring"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type (
	// HTTP сервер
	HTTP struct {
		conf       config.Config
		httpServer *http.Server
	}
	// GRPC сервер
	GRPC struct {
		conf       config.Config
		grpcServer *grpc.Server
	}
	// Aud сервис
	Aud interface {
		StoreStatus(h entity.UpdateStatus)
		StoreHTTP(h entity.HTTPInfo)
	}
)

// NewServer конструктор сервера
func NewServer(cfg config.Config, handler http.Handler) *HTTP {
	return &HTTP{
		conf: cfg,
		httpServer: &http.Server{
			Addr:              cfg.ConnectConfig.HTTP_HOST + ":" + cfg.ConnectConfig.HTTP_PORT,
			Handler:           handler,
			ReadTimeout:       cfg.AppConfig.Read,
			WriteTimeout:      cfg.AppConfig.Write,
			IdleTimeout:       cfg.AppConfig.Idle,
			ReadHeaderTimeout: cfg.AppConfig.ReadHeader,
		},
	}
}

// Run запуск сервера
func (s *HTTP) Run() error {
	return s.httpServer.ListenAndServe()
}

// Close остановка HTTP сервера
func (s *HTTP) Close() {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Println("не удалось закрыть сервер")
	}
}

// NewServerGRPC конструктор сервера
func NewServerGRPC(cfg config.Config, srv *g.Handler, audit Aud) *GRPC {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(intercep.BasicAuthGRPCInterceptor(cfg),
			intercep.LoggingInterceptor(audit), intercep.MetricsInterceptor(),
		),
		grpc.UnaryInterceptor(intercep.SpanInterceptor),
	)
	reflection.Register(s)
	desc.RegisterOrderV1Server(s, srv)

	return &GRPC{
		conf:       cfg,
		grpcServer: s,
	}
}

// Run запуск сервера
func (s *GRPC) Run() error {
	monitoring.StartMetricsServer(s.conf)

	lis, err := net.Listen("tcp", s.conf.ConnectConfig.GRPC_HOST+":"+s.conf.ConnectConfig.GRPC_PORT)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	return s.grpcServer.Serve(lis)
}

// Close остановка GRPC сервера
func (s *GRPC) Close() {
	s.grpcServer.Stop()
}
