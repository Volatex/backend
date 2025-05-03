package middleware

import (
	"context"
	"fmt"
	"gitverse.ru/volatex/backend/notification-service/pkg/logger"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(l logger.Interface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				l.Error(fmt.Sprintf(
					"PANIC in method: %s\nError: %v\nStack:\n%s",
					info.FullMethod, r, string(debug.Stack()),
				))
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
