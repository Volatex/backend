package v1

import (
	"context"
	"gitverse.ru/volatex/backend/notification-service/internal/usecase"
	pb "gitverse.ru/volatex/backend/notification-service/pb/proto"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	uc usecase.NotificationUseCase
}

func NewNotificationHandler(uc usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

func (h *NotificationHandler) SendVerificationCode(ctx context.Context, req *pb.SendVerificationCodeRequest) (*pb.SendVerificationCodeResponse, error) {
	err := h.uc.SendVerificationCode(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &pb.SendVerificationCodeResponse{Message: "Verification code sent"}, nil
}
