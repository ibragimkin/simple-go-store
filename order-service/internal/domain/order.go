package domain

import "time"

type Order struct {
	Id           int        `json:"id"`
	UserId       int        `json:"user_id"`
	ItemId       int        `json:"item_id"`
	Amount       float64    `json:"amount"`
	IsPayed      bool       `json:"is_payed"`
	CreationDate time.Time  `json:"creation_date"`
	PaymentDate  *time.Time `json:"payment_date"`
}
