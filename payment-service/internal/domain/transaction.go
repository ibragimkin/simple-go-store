package domain

import "time"

// Transaction описывает операцию пополнения или снятия средств.
type Transaction struct {
	Id        int       `json:"id"`         // Уникальный идентификатор транзакции
	UserId    int       `json:"user_id"`    // Идентификатор пользователя, связанного с операцией
	IsDeposit bool      `json:"is_deposit"` // Тип операции: true — пополнение, false — снятие
	Amount    float64   `json:"amount"`     // Сумма операции
	Date      time.Time `json:"date"`       // Дата выполнения транзакции
}
