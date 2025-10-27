package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// MessageBus — это высокоуровневая обёртка над Kafka Producer и Consumer,
// обеспечивающая двустороннюю коммуникацию между сервисами.
//
// Основная идея — реализовать механизм корреляции сообщений,
// чтобы можно было отправить запрос (через Producer) и получить
// конкретный ответ (через Consumer) по тому же ключу (correlation key).
type MessageBus struct {
	consumer       *Consumer              // Kafka consumer для чтения сообщений
	producer       *Producer              // Kafka producer для отправки сообщений
	correlationMap map[string]chan []byte // Карта ключей корреляции -> каналы для передачи ответов
}

// NewMessageBus создаёт новый экземпляр MessageBus с заданными Consumer и Producer.
// correlationMap инициализируется пустой map для отслеживания ожидаемых ответов.
func NewMessageBus(consumer *Consumer, producer *Producer) *MessageBus {
	return &MessageBus{
		consumer:       consumer,
		producer:       producer,
		correlationMap: make(map[string]chan []byte),
	}
}

// SendMessage отправляет сообщение в Kafka с заданным ключом и значением,
// затем блокирующе ожидает ответа с тем же ключом через ReceiveMessage.
//
// Возвращает тело ответа (msg.Value) либо ошибку, если:
//   - не удалось отправить сообщение в Kafka,
//   - не пришёл ответ в течение таймаута (60 секунд),
//   - или контекст был отменён.
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

// StartReading запускает бесконечный цикл чтения сообщений из Kafka.
// Для каждого прочитанного сообщения проверяет наличие соответствующего
// канала в correlationMap и пересылает значение сообщения туда.
//
// Если ключ не найден, выводит предупреждение в лог.
//
// Цикл завершается только при закрытии контекста (ctx.Done()) или ошибке чтения.
func (mb *MessageBus) StartReading(ctx context.Context) {
	for {
		msg, err := mb.consumer.ReadMessage(ctx)
		if err != nil {
			if err == io.EOF || errors.Is(err, context.Canceled) {
				log.Printf("Reader stopped gracefully")
				return
			}
			log.Printf("Error reading message: %s\n", err)
			continue
		}

		strKey := string(msg.Key)
		ch, ok := mb.correlationMap[strKey]
		if ok {
			ch <- msg.Value
		} else {
			log.Printf("Key not found in correlation map: %s", strKey)
		}
	}
}

// ReceiveMessage ожидает получение ответа по заданному ключу корреляции.
//
// Поведение:
//   - Если контекст завершён — возвращает ошибку контекста.
//   - Если в течение 60 секунд не поступило сообщение — возвращает timeout-ошибку.
//   - Если канал закрыт — возвращает ошибку io.EOF.
//   - Иначе возвращает полученное сообщение.
//
// Используется внутри SendMessage для получения ответа на запрос.
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
