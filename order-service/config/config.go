package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HttpPort           string
	DatabaseURL        string
	KafkaBrokers       []string
	KafkaRequestTopic  string
	KafkaResponseTopic string
	KafkaGroupID       string
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

	httpPort, err := mustGetEnv("HTTP_PORT")
	if err != nil {
		errs = append(errs, err.Error())
	}

	db, err := mustGetEnv("DATABASE_URL")
	if err != nil {
		errs = append(errs, err.Error())
	}

	brokers, err := mustGetEnv("KAFKA_URL")
	if err != nil {
		errs = append(errs, err.Error())
	}

	consumerTopic, err := mustGetEnv("KAFKA_REQUEST_TOPIC")
	if err != nil {
		errs = append(errs, err.Error())
	}

	producerTopic, err := mustGetEnv("KAFKA_RESPONSE_TOPIC")
	if err != nil {
		errs = append(errs, err.Error())
	}

	groupID, err := mustGetEnv("KAFKA_GROUP_ID")
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("config validation failed:\n  %s", strings.Join(errs, "\n  "))
	}
	return &Config{
		HttpPort:           httpPort,
		DatabaseURL:        db,
		KafkaBrokers:       strings.Split(brokers, ";"),
		KafkaRequestTopic:  consumerTopic,
		KafkaResponseTopic: producerTopic,
		KafkaGroupID:       groupID,
	}, nil
}
