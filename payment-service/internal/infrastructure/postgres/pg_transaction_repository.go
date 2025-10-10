package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"payment-service/internal/application/repository"
	"payment-service/internal/domain"
)

type TransactionDb struct {
	db *pgxpool.Pool
}

func NewTransactionDb(db *pgxpool.Pool) (repository.TransactionRepository, error) {
	return TransactionDb{db: db}, nil
}

func (tdb TransactionDb) GetById(ctx context.Context, id int) (*domain.Transaction, error) {
	row := tdb.db.QueryRow(ctx, `
SELECT (id, user_id, is_deposit, amount, date)
FROM transactions
WHERE id=$1
`, id)
	txn := domain.Transaction{}
	err := row.Scan(&txn.Id, &txn.UserId, &txn.IsDeposit, &txn.Amount, &txn.Date)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &txn, nil

}

func (tdb TransactionDb) Save(ctx context.Context, txn *domain.Transaction) error {
	_, err := tdb.db.Exec(ctx, `
INSERT INTO transactions (id, user_id, is_deposit, amount, date )
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO NOTHING
`, &txn.Id, &txn.UserId, &txn.IsDeposit, &txn.Amount, &txn.Date)
	return err
}
