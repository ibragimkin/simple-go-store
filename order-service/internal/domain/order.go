package domain

import "time"

// Order представляет заказ, оформленный пользователем.
// Содержит информацию о товаре, пользователе, сумме и статусе оплаты.
type Order struct {
	Id           int        `json:"id"`            // Уникальный идентификатор заказа (в будущем заменить на UUID)
	UserId       int        `json:"user_id"`       // ID пользователя, оформившего заказ
	ItemId       int        `json:"item_id"`       // ID товара, на который оформлен заказ
	Amount       float64    `json:"amount"`        // Сумма заказа
	IsPayed      bool       `json:"is_payed"`      // Статус оплаты: true, если заказ оплачен
	CreationDate time.Time  `json:"creation_date"` // Дата создания заказа
	PaymentDate  *time.Time `json:"payment_date"`  // Дата оплаты (nil, если заказ ещё не оплачен)
}

// Pay помечает заказ как оплаченный и устанавливает текущую дату в PaymentDate.
func (o *Order) Pay() {
	o.IsPayed = true
	date := time.Now()
	o.PaymentDate = &date
}
