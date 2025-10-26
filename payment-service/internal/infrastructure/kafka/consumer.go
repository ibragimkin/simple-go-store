package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

// Consumer - обёртка над kafka.Reader, обеспечивающая чтение сообщений из топика.
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer создаёт нового Consumer с заданными брокерами, топиком и groupID.
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{reader: kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})}
}

// ReadMessage читает одно сообщение из Kafka.
// Возвращает ошибку, если чтение не удалось или контекст отменён.
func (c *Consumer) ReadMessage(ctx context.Context) (*kafka.Message, error) {
	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}
	return &message, nil
}

// Close закрывает Kafka reader и освобождает ресурсы.
func (c *Consumer) Close() error {
	err := c.reader.Close()
	if err != nil {
		return fmt.Errorf("error closing reader: %w", err)
	}
	return nil
}
