package external

import (
	"context"
	"fmt"

	"github.com/tinkoff/invest-api-go-sdk/investgo"
	investapi "github.com/tinkoff/invest-api-go-sdk/proto"
	"gitverse.ru/volatex/backend/market-service/pkg/logger"
)

const (
	defaultEndpoint = "sandbox-invest-public-api.tinkoff.ru:443"
)

type TinkoffClient interface {
	GetLastPrice(ctx context.Context, identifier string, isTicker bool) (float64, error)
	BuyAsset(ctx context.Context, ticker string, quantity int64, accountID string) error
}

type tinkoffSDKWrapper struct {
	client *investgo.Client
	logger logger.Interface
}

// TinkoffConfig represents configuration for Tinkoff client
type TinkoffConfig struct {
	Token    string
	EndPoint string
}

// NewTinkoffClient creates a new Tinkoff client with the given configuration
func NewTinkoffClient(cfg TinkoffConfig, logger logger.Interface) (TinkoffClient, error) {
	if cfg.EndPoint == "" {
		cfg.EndPoint = defaultEndpoint
	}

	config := investgo.Config{
		Token:    cfg.Token,
		EndPoint: cfg.EndPoint,
	}

	sdkClient, err := investgo.NewClient(context.Background(), config, &sdkLogger{logger: logger})
	if err != nil {
		return nil, fmt.Errorf("failed to create Tinkoff client: %w", err)
	}
	return &tinkoffSDKWrapper{
		client: sdkClient,
		logger: logger,
	}, nil
}

func (t *tinkoffSDKWrapper) GetLastPrice(ctx context.Context, identifier string, isTicker bool) (float64, error) {
	if t.client == nil {
		return 0, fmt.Errorf("Tinkoff client is not initialized")
	}

	var figi string
	if isTicker {
		instrumentsService := t.client.NewInstrumentsServiceClient()
		if instrumentsService == nil {
			return 0, fmt.Errorf("Instruments service is not initialized")
		}

		instrumentResp, err := instrumentsService.InstrumentByTicker(identifier, "TQBR")
		if err != nil {
			t.logger.Error(fmt.Errorf("error getting instrument by ticker: %w", err), "Error getting instrument by ticker", "ticker", identifier)
			return 0, err
		}

		instrument := instrumentResp.GetInstrument()
		if instrument == nil {
			return 0, fmt.Errorf("No instrument found for ticker: %s", identifier)
		}

		figi = instrument.Figi
		if figi == "" {
			return 0, fmt.Errorf("Empty FIGI for ticker: %s", identifier)
		}
	} else {
		figi = identifier
	}

	marketDataService := t.client.NewMarketDataServiceClient()
	if marketDataService == nil {
		return 0, fmt.Errorf("Market data service is not initialized")
	}

	lastPricesResp, err := marketDataService.GetLastPrices([]string{figi})
	if err != nil {
		t.logger.Error(fmt.Errorf("error getting latest prices: %w", err), "Error getting latest prices")
		return 0, err
	}

	for _, price := range lastPricesResp.GetLastPrices() {
		priceValue := float64(price.GetPrice().Units) + float64(price.GetPrice().Nano)/1e9
		t.logger.Info("Asset price received", "price", priceValue)
		return priceValue, nil
	}

	return 0, fmt.Errorf("No price found for FIGI: %s", figi)
}

func (t *tinkoffSDKWrapper) BuyAsset(ctx context.Context, ticker string, quantity int64, accountID string) error {
	if t.client == nil {
		return fmt.Errorf("Tinkoff client is not initialized")
	}

	instrumentsService := t.client.NewInstrumentsServiceClient()
	if instrumentsService == nil {
		return fmt.Errorf("Instruments service is not initialized")
	}
	instrumentResp, err := instrumentsService.InstrumentByTicker(ticker, "TQBR")
	if err != nil {
		t.logger.Error(fmt.Errorf("error getting instrument by ticker: %w", err), "Error getting instrument by ticker", "ticker", ticker)
		return err
	}
	instrument := instrumentResp.GetInstrument()
	if instrument == nil {
		return fmt.Errorf("No instrument found for ticker: %s", ticker)
	}
	figi := instrument.Figi
	if figi == "" {
		return fmt.Errorf("Empty FIGI for ticker: %s", ticker)
	}

	ordersService := t.client.NewOrdersServiceClient()
	if ordersService == nil {
		return fmt.Errorf("Orders service is not initialized")
	}
	orderResponse, err := ordersService.PostOrder(&investgo.PostOrderRequest{
		Quantity:     quantity,
		Price:        nil,
		Direction:    investapi.OrderDirection_ORDER_DIRECTION_BUY,
		AccountId:    accountID,
		OrderType:    investapi.OrderType_ORDER_TYPE_MARKET,
		InstrumentId: figi,
	})
	if err != nil {
		t.logger.Error(fmt.Errorf("error creating buy order: %w", err), "Error creating buy order", "ticker", ticker, "quantity", quantity)
		return err
	}

	t.logger.Info("Buy order created successfully",
		"order_id", orderResponse.GetOrderId(),
		"status", orderResponse.GetExecutionReportStatus(),
		"lots_requested", orderResponse.GetLotsRequested(),
		"lots_executed", orderResponse.GetLotsExecuted())

	return nil
}

type sdkLogger struct {
	logger logger.Interface
}

func (l *sdkLogger) Infof(template string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(template, args...))
}

func (l *sdkLogger) Errorf(template string, args ...interface{}) {
	l.logger.Error(fmt.Errorf(template, args...))
}

func (l *sdkLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal(fmt.Errorf(template, args...))
}
