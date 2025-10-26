package repository

import (
	"context"
	"payment-service/internal/domain"
)

// TransactionRepository определяет интерфейс для работы с транзакциями.
type TransactionRepository interface {
	// GetById возвращает транзакцию по её ID.
	GetById(ctx context.Context, id int) (*domain.Transaction, error)

	// Save сохраняет новую транзакцию в хранилище.
	Save(ctx context.Context, transaction *domain.Transaction) error
}
