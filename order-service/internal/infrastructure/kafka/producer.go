package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{writer: &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic}}
}

func (p *Producer) SendMessage(ctx context.Context, message *kafka.Message) error {
	err := p.writer.WriteMessages(ctx, *message)
	if err != nil {
		return fmt.Errorf("error writing message: %s\n", err)
	}
	return nil
}

func (p *Producer) Close() error {
	err := p.writer.Close()
	if err != nil {
		return fmt.Errorf("closing writer: %w", err)
	}
	return nil
}
