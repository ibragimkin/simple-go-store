package repository

import (
	"context"
	"payment-service/internal/domain"
)

type TransactionRepository interface {
	GetById(ctx context.Context, id int) (*domain.Transaction, error)
	Save(ctx context.Context, transaction *domain.Transaction) error
}
