package repository

import (
	"context"
	"payment-service/internal/domain"
)

// AccountRepository определяет интерфейс для работы с счетами пользователей.
type AccountRepository interface {
	// GetById возвращает счёт по его ID.
	GetById(ctx context.Context, id int) (*domain.Account, error)

	// GetByUserId возвращает счёт, принадлежащий конкретному пользователю.
	GetByUserId(ctx context.Context, userId int) (*domain.Account, error)

	// Save сохраняет или обновляет данные счёта.
	Save(ctx context.Context, account *domain.Account) error
}
