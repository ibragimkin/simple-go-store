package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"payment-service/internal/application/repository"
	"payment-service/internal/domain"
)

// AccountDb реализует интерфейс AccountRepository
// и работает с таблицей accounts в PostgreSQL через pgxpool.
type AccountDb struct {
	db PgxPool
}

// NewAccountDb создаёт новый объект доступа к данным аккаунтов.
func NewAccountDb(db PgxPool) (repository.AccountRepository, error) {
	return AccountDb{db: db}, nil
}

// GetById возвращает аккаунт по его ID.
// Возвращает ошибку, если аккаунт не найден.
func (adb AccountDb) GetById(ctx context.Context, id int) (*domain.Account, error) {
	row := adb.db.QueryRow(ctx, `
SELECT id, user_id, balance, creation_date
FROM accounts
WHERE id=$1
`, id)

	var acc domain.Account
	err := row.Scan(&acc.Id, &acc.UserId, &acc.Balance, &acc.CreationDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

// Save сохраняет аккаунт в базу данных.
// Если аккаунт с таким id уже существует — обновляет баланс.
func (adb AccountDb) Save(ctx context.Context, account *domain.Account) error {
	_, err := adb.db.Exec(ctx, `
INSERT INTO accounts (id, user_id, balance, creation_date)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
    SET balance = EXCLUDED.balance
`, &account.Id, &account.UserId, &account.Balance, &account.CreationDate)
	return err
}

// GetByUserId возвращает аккаунт по user_id.
// Возвращает ошибку, если аккаунт не найден.
func (adb AccountDb) GetByUserId(ctx context.Context, userId int) (*domain.Account, error) {
	row := adb.db.QueryRow(ctx, `
SELECT id, user_id, balance, creation_date
FROM accounts
WHERE user_id=$1
`, userId)

	var acc domain.Account
	err := row.Scan(&acc.Id, &acc.UserId, &acc.Balance, &acc.CreationDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
