package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"payment-service/internal/application/repository"
	"payment-service/internal/domain"
)

// TransactionDb реализует интерфейс repository.TransactionRepository
// и отвечает за работу с таблицей transactions в PostgreSQL.
type TransactionDb struct {
	db PgxPool
}

// NewTransactionDb создаёт новый экземпляр TransactionDb,
// принимая пул подключений к PostgreSQL.
func NewTransactionDb(db PgxPool) (repository.TransactionRepository, error) {
	return TransactionDb{db: db}, nil
}

// GetById возвращает транзакцию по её ID.
// Возвращает nil, nil если транзакция не найдена.
func (tdb TransactionDb) GetById(ctx context.Context, id int) (*domain.Transaction, error) {
	row := tdb.db.QueryRow(ctx, `
SELECT id, user_id, is_deposit, amount, date
FROM transactions
WHERE id = $1
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

// Save сохраняет новую транзакцию в базу данных.
// Если запись с таким ID уже существует — операция игнорируется.
func (tdb TransactionDb) Save(ctx context.Context, txn *domain.Transaction) error {
	_, err := tdb.db.Exec(ctx, `
INSERT INTO transactions (id, user_id, is_deposit, amount, date)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO NOTHING
`, &txn.Id, &txn.UserId, &txn.IsDeposit, &txn.Amount, &txn.Date)
	return err
}
