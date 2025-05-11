package grpc

import (
	"gitverse.ru/volatex/backend/notification-service/internal/controller/grpc/middleware"
	v1 "gitverse.ru/volatex/backend/notification-service/internal/controller/grpc/v1"
	pb "gitverse.ru/volatex/backend/notification-service/pb/proto"
	"gitverse.ru/volatex/backend/notification-service/pkg/logger"
	"google.golang.org/grpc"
)

func NewServer(l logger.Interface, handler *v1.NotificationHandler) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			middleware.Recovery(l),
			middleware.Logger(l),
		),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterNotificationServiceServer(s, handler)
	return s
}
