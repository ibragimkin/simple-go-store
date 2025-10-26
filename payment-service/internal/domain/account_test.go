package domain

import (
	"testing"
)

func TestAccountDeposit(t *testing.T) {
	account := &Account{Balance: 100.0}

	err := account.Deposit(50.0)
	if err != nil {
		t.Errorf("Deposit failed: %v", err)
	}
	if account.Balance != 150.0 {
		t.Errorf("Expected 150.0, got %v", account.Balance)
	}

	err = account.Deposit(-10.0)
	if err == nil {
		t.Error("Expected error for negative amount")
	}
}

func TestAccountWithdraw(t *testing.T) {
	account := &Account{Balance: 100.0}

	err := account.Withdraw(30.0)
	if err != nil {
		t.Errorf("Withdraw failed: %v", err)
	}
	if account.Balance != 70.0 {
		t.Errorf("Expected 70.0, got %v", account.Balance)
	}

	err = account.Withdraw(100.0)
	if err == nil {
		t.Error("Expected error for insufficient funds")
	}

	err = account.Withdraw(-10.0)
	if err == nil {
		t.Error("Expected error for negative amount")
	}
}

func TestAccountEdgeCases(t *testing.T) {
	account := &Account{Balance: 0.0}

	err := account.Withdraw(10.0)
	if err == nil {
		t.Error("Expected error when withdrawing from zero balance")
	}

	err = account.Deposit(0.0)
	if err != nil {
		t.Errorf("Deposit zero should not error: %v", err)
	}
}
