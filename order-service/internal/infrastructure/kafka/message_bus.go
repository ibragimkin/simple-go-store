package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"time"
)

type MessageBus struct {
	consumer       *Consumer
	producer       *Producer
	correlationMap map[string]chan []byte
}

func NewMessageBus(consumer *Consumer, producer *Producer) *MessageBus {
	return &MessageBus{consumer: consumer, producer: producer, correlationMap: make(map[string]chan []byte)}
}

func (mb *MessageBus) SendMessage(ctx context.Context, key []byte, value []byte) ([]byte, error) {
	err := mb.producer.SendMessage(ctx, &kafka.Message{Key: key, Value: value})
	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}
	strKey := string(key)
	mb.correlationMap[strKey] = make(chan []byte)
	kafkaMsg, err := mb.ReceiveMessage(ctx, strKey)
	if err != nil {
		return nil, fmt.Errorf("error receiving message: %w", err)
	}
	return kafkaMsg.Value, nil
}

func (mb *MessageBus) StartReading(ctx context.Context) {
	for {
		msg, err := mb.consumer.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF || errors.Is(err, context.Canceled) {
				log.Printf("ERROR AAAAA")
			}
			log.Printf("Error reading message: %s\n", err)
		}
		strKey := string(msg.Key)
		_, ok := mb.correlationMap[strKey]
		if ok {
			mb.correlationMap[strKey] <- msg.Value
		} else {
			log.Printf("Key not found in map: %s", strKey)
		}
	}
}

func (mb *MessageBus) ReceiveMessage(ctx context.Context, key string) (*kafka.Message, error) {
	ch, ok := mb.correlationMap[key]
	if !ok {
		return nil, fmt.Errorf("key not found in correlation map: %s", key)
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(60 * time.Second):
		return nil, fmt.Errorf("timeout waiting for message: %s", key)
	case msg, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("message channel closed: %w", io.EOF)
		}
		return &kafka.Message{Value: msg}, nil
	}
}
