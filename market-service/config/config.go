package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		PG      PG
		Swagger Swagger
		GRPC    GRPC
		JWT     JWT
		Math    Math
	}

	App struct {
		Name    string `env:"APP_NAME,required"`
		Version string `env:"APP_VERSION,required"`
	}

	HTTP struct {
		Port           string `env:"HTTP_PORT,required"`
		UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}

	Swagger struct {
		Enabled bool `env:"SWAGGER_ENABLED" envDefault:"true"`
	}

	GRPC struct {
		NotificationAddr string `env:"GRPC_NOTIFICATION_ADDR,required"`
	}

	JWT struct {
		Secret     string `env:"JWT_SECRET,required"`
		TTL        string `env:"JWT_TTL,required"`
		RefreshTTL string `env:"REFRESH_TTL,required"`
	}

	Math struct {
		Endpoint string `env:"MATH_SERVICE_ENDPOINT" envDefault:"math_service:50055"`
	}
)

// NewConfig - возвращает конфигурацию приложения.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
