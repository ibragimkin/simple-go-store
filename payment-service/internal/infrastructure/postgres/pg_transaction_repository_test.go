package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"payment-service/internal/domain"
	"payment-service/internal/infrastructure/postgres"
)

// TestTransactionDb_GetById проверяет корректность получения транзакции по ID.
func TestTransactionDb_GetById(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock: %v", err)
	}
	defer mock.Close()

	db, _ := postgres.NewTransactionDb(mock)
	ctx := context.Background()

	rows := pgxmock.NewRows([]string{"id", "user_id", "is_deposit", "amount", "date"}).
		AddRow(1, 10, true, 100.0, time.Now())

	mock.ExpectQuery(`SELECT id, user_id, is_deposit, amount, date FROM transactions WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	txn, err := db.GetById(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if txn == nil || txn.Id != 1 || txn.UserId != 10 {
		t.Errorf("unexpected result: %+v", txn)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

// TestTransactionDb_Save проверяет корректность сохранения транзакции.
func TestTransactionDb_Save(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock: %v", err)
	}
	defer mock.Close()

	db, _ := postgres.NewTransactionDb(mock)
	ctx := context.Background()

	txn := &domain.Transaction{
		Id:        2,
		UserId:    42,
		IsDeposit: true,
		Amount:    250.5,
		Date:      time.Now(),
	}

	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(&txn.Id, &txn.UserId, &txn.IsDeposit, &txn.Amount, &txn.Date).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = db.Save(ctx, txn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
