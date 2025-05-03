package middleware

import (
	"context"
	"fmt"
	"gitverse.ru/volatex/backend/notification-service/pkg/logger"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func Logger(l logger.Interface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		clientIP := "unknown"
		if p, ok := peer.FromContext(ctx); ok {
			clientIP = p.Addr.String()
		}

		l.Info(fmt.Sprintf(
			"gRPC request - method: %s, client: %s, duration: %s",
			info.FullMethod, clientIP, duration,
		))

		return resp, err
	}
}
