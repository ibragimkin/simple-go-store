package repository

import (
	"context"
	"order-service/internal/domain"
)

// OrderRepository определяет интерфейс для работы с хранилищем заказов.
// Позволяет получать, сохранять и извлекать заказы пользователя.
type OrderRepository interface {
	// GetById возвращает заказ по его ID.
	// Возвращает ошибку, если заказ не найден или произошла ошибка при чтении.
	GetById(ctx context.Context, id int) (*domain.Order, error)

	// Save сохраняет заказ в хранилище.
	// Если заказ с таким ID уже существует, он должен быть обновлён.
	Save(ctx context.Context, order *domain.Order) error

	// GetUserOrders возвращает список всех заказов пользователя по его ID.
	// В случае ошибки возвращает пустой срез и ошибку.
	GetUserOrders(ctx context.Context, userId int) ([]domain.Order, error)
}
