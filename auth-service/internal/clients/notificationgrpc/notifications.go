package notificationgrpc

import (
	"context"

	pb "gitverse.ru/volatex/backend/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.NotificationServiceClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		client: pb.NewNotificationServiceClient(conn),
	}, nil
}

func (c *Client) SendVerificationCode(ctx context.Context, email string) error {
	_, err := c.client.SendVerificationCode(ctx, &pb.SendVerificationCodeRequest{Email: email})
	return err
}

func (c *Client) VerifyCode(ctx context.Context, email string, code string) (bool, error) {
	resp, err := c.client.VerifyCode(ctx, &pb.VerifyCodeRequest{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return false, err
	}
	return resp.IsValid, nil
}
