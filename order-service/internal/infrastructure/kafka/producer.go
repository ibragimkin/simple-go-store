package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

// Producer - обёртка над kafka.Writer, обеспечивающая отправку сообщений в топик.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer создаёт нового Producer с заданными брокерами и топиком.
func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{writer: &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}}
}

// SendMessage отправляет сообщение в Kafka.
// Возвращает ошибку при сбое записи.
func (p *Producer) SendMessage(ctx context.Context, message *kafka.Message) error {
	err := p.writer.WriteMessages(ctx, *message)
	if err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}
	return nil
}

// Close закрывает Kafka writer и освобождает ресурсы.
func (p *Producer) Close() error {
	err := p.writer.Close()
	if err != nil {
		return fmt.Errorf("closing writer: %w", err)
	}
	return nil
}
