package domain

import "time"

type Transaction struct {
	Id        int       `json:"id"`      // TODO: change to go-uuid
	UserId    int       `json:"user_id"` // TODO: change to go-uuid
	IsDeposit bool      `json:"is_deposit"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
}
