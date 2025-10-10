package domain

import (
	"fmt"
	"time"
)

type Account struct {
	Id           int       `json:"id"`      // TODO: change to go-uuid
	UserId       int       `json:"user_id"` // TODO: change to go-uuid
	Balance      float64   `json:"balance"`
	CreationDate time.Time `json:"creation_date"`
}

func (a *Account) Deposit(amount float64) error {
	if amount < 0 {
		return fmt.Errorf("amount must be not negative")
	}
	a.Balance += amount
	return nil
}

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
