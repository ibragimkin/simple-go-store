package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{reader: kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})}
}

func (c *Consumer) ReadMessage(ctx context.Context) (*kafka.Message, error) {
	message, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}
	return &message, nil
}

func (c *Consumer) Close() error {
	err := c.reader.Close()
	if err != nil {
		return fmt.Errorf("error closing reader: %w", err)
	}
	return nil
}

//func (c *Consumer) ReadMessage(ctx context.Context, brokers []string, topic string, groupID string, handler func(message *kafkahandler.Message) error) {
//	reader := kafkahandler.NewReader(kafkahandler.ReaderConfig{
//		Brokers: brokers,
//		Topic:   topic,
//		GroupID: groupID,
//	})
//	defer func() {
//		err := reader.Close()
//		if err != nil {
//			log.Printf("Error closing reader: %s\n", err)
//		}
//		err = c.producer.Close()
//		if err != nil {
//			log.Printf("Error closing writer: %s\n", err)
//		}
//	}()
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		default:
//		}
//		m, err := reader.ReadMessage(ctx)
//		if err != nil {
//			if errors.Is(err, context.Canceled) {
//				return
//			}
//			log.Printf("Error reading message: %s\n", err)
//			continue
//		}
//		go func() {
//			err := handler(&m)
//			if err != nil {
//				log.Printf("Error handling message: %s\n", err)
//				c.producer.SendMessage(ctx, m.Key, []byte(err.Error()))
//			}
//		}()
//	}
//}
