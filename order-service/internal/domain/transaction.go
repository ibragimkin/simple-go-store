package domain

import "time"

// Transaction представляет транзакцию - операцию по списанию или пополнению счёта.
type Transaction struct {
	Id        int       `json:"id"`         // Уникальный идентификатор транзакции
	UserId    int       `json:"user_id"`    // ID пользователя, к которому относится транзакция
	IsDeposit bool      `json:"is_deposit"` // true - если это пополнение, false - если списание
	Amount    float64   `json:"amount"`     // Сумма транзакции
	Date      time.Time `json:"date"`       // Дата и время проведения транзакции
}
