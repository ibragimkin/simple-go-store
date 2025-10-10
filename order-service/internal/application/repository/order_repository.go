package repository

import (
	"context"
	"order-service/internal/domain"
)

type OrderRepository interface {
	GetById(ctx context.Context, id int) (*domain.Order, error)
	Save(ctx context.Context, account *domain.Order) error
	GetUserOrders(ctx context.Context, userId int) ([]domain.Order, error)
}
