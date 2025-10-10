package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HttpPort          string
	OrderServiceURL   string
	PaymentServiceURL string
}

func mustGetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func LoadConfig() (*Config, error) {
	errs := make([]string, 0)

	orderService, err := mustGetEnv("ORDER_SERVICE_URL")
	if err != nil {
		errs = append(errs, err.Error())
	}

	paymentService, err := mustGetEnv("PAYMENT_SERVICE_URL")
	if err != nil {
		errs = append(errs, err.Error())
	}

	httpPort, err := mustGetEnv("HTTP_PORT")
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("config validation failed:\n  %s", strings.Join(errs, "\n  "))
	}
	return &Config{
		HttpPort:          httpPort,
		OrderServiceURL:   orderService,
		PaymentServiceURL: paymentService,
	}, nil
}
