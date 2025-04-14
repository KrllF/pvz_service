package interceptor

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/KrllF/pvz_service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	needSize = 2
)

// BasicAuthGRPCInterceptor Basic Auth для GRPC
func BasicAuthGRPCInterceptor(config config.Config) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{},
		_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "miss metadata")
		}

		authHeaders := meta.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing header")
		}
		authSlc := strings.Split(authHeaders[0], " ")
		if authSlc[0] != "Basic" {
			return nil, status.Error(codes.Unauthenticated, "invalid header format")
		}
		if len(authSlc) != needSize {
			return nil, status.Error(codes.Unauthenticated, "invalid header format")
		}
		logNpass, err := base64.StdEncoding.DecodeString(authSlc[1])
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "failed to decode")
		}

		authData := strings.Split(string(logNpass), ":")
		if len(authData) != needSize || authData[0] != config.Login || authData[1] != config.Password {
			return nil, status.Error(codes.Unauthenticated, "invalid auth")
		}

		return handler(ctx, req)
	}
}
