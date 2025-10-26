package config

import (
	"fmt"
	"os"
	"strings"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	HttpPort           string   // Порт для HTTP сервера
	DatabaseURL        string   // URL для подключения к базе данных
	KafkaBrokers       []string // Список брокеров Kafka
	KafkaConsumerTopic string   // Топик для потребления сообщений
	KafkaProducerTopic string   // Топик для производства сообщений
	KafkaGroupID       string   // Group ID для Kafka consumer
}

// mustGetEnv получает значение обязательной переменной окружения или возвращает ошибку если она пустая
func mustGetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

// LoadConfig загружает конфигурацию из переменных окружения и возвращает Config
// Возвращает ошибку если какие-то обязательные переменные не установлены
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
		KafkaConsumerTopic: consumerTopic,
		KafkaProducerTopic: producerTopic,
		KafkaGroupID:       groupID,
	}, nil
}
