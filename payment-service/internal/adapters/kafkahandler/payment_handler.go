package kafkahandler

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"payment-service/internal/application/service"
	"payment-service/internal/domain"
)

// NewPaymentHandler возвращает функцию-обработчик Kafka-сообщений,
// которая десериализует JSON-тело сообщения в структуру domain.Transaction,
// передаёт её в PaymentService для обработки,
// и возвращает Kafka-ответ с результатом ("OK" или текст ошибки).
//
// В случае ошибки десериализации или бизнес-логики сервис возвращает
// сообщение с описанием проблемы и соответствующую ошибку.
func NewPaymentHandler(service *service.PaymentService) func(ctx context.Context, message *kafka.Message) (*kafka.Message, error) {
	return func(ctx context.Context, message *kafka.Message) (*kafka.Message, error) {
		var tx domain.Transaction
		err := json.Unmarshal(message.Value, &tx)
		if err != nil {
			resp := "invalid JSON: " + err.Error()
			return &kafka.Message{Key: message.Key, Value: []byte(resp)}, err
		}
		err = service.ProcessTransaction(ctx, tx)
		if err != nil {
			response := "Error processing transaction: " + err.Error()
			return &kafka.Message{Key: message.Key, Value: []byte(response)}, err
		}
		return &kafka.Message{Key: message.Key, Value: []byte("OK")}, err
	}
}
