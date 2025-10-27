package domain

import (
	"testing"
	"time"
)

func TestOrder_Pay(t *testing.T) {
	order := Order{
		Id:           123,
		UserId:       11,
		ItemId:       50,
		Amount:       1000,
		IsPayed:      false,
		CreationDate: time.Now(),
		PaymentDate:  nil,
	}
	order.Pay()
	if !order.IsPayed {
		t.Errorf("Order is not payed")
	}
	if order.PaymentDate == nil {
		t.Errorf("Payment date is nil")
	}
	return
}
