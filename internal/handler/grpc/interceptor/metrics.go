package interceptor

import (
	"context"

	"github.com/KrllF/pvz_service/pkg/monitoring"
	"google.golang.org/grpc"
)

// MetricsInterceptor interceptor для логирования
func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		monitoring.SetCounterReq(info.FullMethod)

		resp, err := handler(ctx, req)
		if err != nil {
			monitoring.SetcounterOKERR("err")

			return resp, err
		}
		monitoring.SetcounterOKERR("ok")

		return resp, err
	}
}
