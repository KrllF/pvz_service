package interceptor

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// SpanInterceptor interceptor
func SpanInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Создаем новый спан
	span, ctx := opentracing.StartSpanFromContext(ctx, info.FullMethod)
	defer span.Finish()
	// Вызываем обработчик gRPC
	resp, err := handler(ctx, req)
	// Логируем ошибки
	if err != nil {
		span.SetTag("error", true)
		span.LogKV("error", err.Error())
	}

	return resp, err
}
