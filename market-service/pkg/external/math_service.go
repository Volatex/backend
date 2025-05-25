package external

import (
	"context"
	"fmt"

	mathservice "gitverse.ru/volatex/backend/market-service/api/proto"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultMathServiceEndpoint = "localhost:50055"
)

type MathServiceClient interface {
	CalculateVolatility(ctx context.Context, returns []float64) (float64, error)
}

type mathServiceClient struct {
	client mathservice.MathServiceClient
	logger logger.Interface
}

type MathServiceConfig struct {
	Endpoint string
}

func NewMathServiceClient(cfg MathServiceConfig, logger logger.Interface) (MathServiceClient, error) {
	if cfg.Endpoint == "" {
		cfg.Endpoint = defaultMathServiceEndpoint
	}

	conn, err := grpc.Dial(cfg.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to math service: %w", err)
	}

	client := mathservice.NewMathServiceClient(conn)
	return &mathServiceClient{
		client: client,
		logger: logger,
	}, nil
}

func (c *mathServiceClient) CalculateVolatility(ctx context.Context, returns []float64) (float64, error) {
	request := &mathservice.VolatilityRequest{
		Returns: returns,
	}

	response, err := c.client.CalculateVolatility(ctx, request)
	if err != nil {
		c.logger.Error(err, "Failed to calculate volatility")
		return 0, fmt.Errorf("failed to calculate volatility: %w", err)
	}

	return response.Volatility, nil
}
