package repository

import (
	"context"
	"payment-service/internal/domain"
)

type AccountRepository interface {
	GetById(ctx context.Context, id int) (*domain.Account, error)
	GetByUserId(ctx context.Context, userId int) (*domain.Account, error)
	Save(ctx context.Context, account *domain.Account) error
}
