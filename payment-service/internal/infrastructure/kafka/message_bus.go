package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
)

// MessageBus объединяет Kafka consumer и producer, обеспечивая двустороннюю обработку сообщений.
type MessageBus struct {
	consumer *Consumer
	producer *Producer
}

// NewMessageBus создаёт новый экземпляр MessageBus с указанными топиками и groupID.
func NewMessageBus(brokers []string, consumerTopic, producerTopic, groupID string) *MessageBus {
	return &MessageBus{
		consumer: NewConsumer(brokers, consumerTopic, groupID),
		producer: NewProducer(brokers, producerTopic),
	}
}

// Start запускает бесконечный цикл чтения сообщений из consumer-топика,
// передаёт каждое сообщение обработчику handler,
// и отправляет результат обратно через producer.
//
// Обработка сообщений выполняется асинхронно (в отдельных горутинах).
func (mb *MessageBus) Start(ctx context.Context, handler func(ctx context.Context, message *kafka.Message) (*kafka.Message, error)) {
	for {
		select {
		case <-ctx.Done():
			err := mb.consumer.Close()
			if err != nil {
				log.Printf("Error closing consumer: %s\n", err)
			}
			err = mb.producer.Close()
			if err != nil {
				log.Printf("Error closing producer: %s\n", err)
			}
			return
		default:
		}

		m, err := mb.consumer.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			log.Printf("Error reading message: %s\n", err)
			continue
		}

		go func() {
			response, err := handler(ctx, m)
			if err != nil {
				log.Printf("Error processing message: %s\n", err)
				return
			}
			err = mb.producer.SendMessage(ctx, response)
			if err != nil {
				log.Printf("Error sending message: %s", err)
			}
		}()
	}

	err := mb.producer.Close()
	if err != nil {
		log.Printf("Error closing producer: %s\n", err)
	}
	err = mb.consumer.Close()
	if err != nil {
		log.Printf("Error closing consumer: %s\n", err)
	}
}
