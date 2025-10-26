package domain

import (
	"fmt"
	"time"
)

// Account представляет счёт пользователя.
// Хранит информацию о текущем балансе и дате создания.
type Account struct {
	Id           int       `json:"id"`            // Уникальный идентификатор счёта
	UserId       int       `json:"user_id"`       // Идентификатор пользователя, которому принадлежит счёт
	Balance      float64   `json:"balance"`       // Текущий баланс счёта
	CreationDate time.Time `json:"creation_date"` // Дата создания счёта
}

// Deposit увеличивает баланс счёта на указанную сумму.
// Возвращает ошибку, если сумма отрицательная.
func (a *Account) Deposit(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("amount must be not negative")
	}
	a.Balance += amount
	return nil
}

// Withdraw уменьшает баланс счёта на указанную сумму.
// Возвращает ошибку, если сумма отрицательная или средств недостаточно.
func (a *Account) Withdraw(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("amount must be not negative")
	}
	if a.Balance-amount < 0 {
		return fmt.Errorf("not enough balance for withdraw")
	}
	a.Balance -= amount
	return nil
}
