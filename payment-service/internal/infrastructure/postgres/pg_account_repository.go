package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"payment-service/internal/application/repository"
	"payment-service/internal/domain"
)

type AccountDb struct {
	db *pgxpool.Pool
}

func NewAccountDb(db *pgxpool.Pool) (repository.AccountRepository, error) {
	return AccountDb{db: db}, nil
}

func (adb AccountDb) GetById(ctx context.Context, id int) (*domain.Account, error) {
	row := adb.db.QueryRow(ctx, `
SELECT id, user_id, balance, creation_date
FROM accounts
WHERE id=$1
`, id)
	var acc domain.Account = domain.Account{}
	err := row.Scan(&acc.Id, &acc.UserId, &acc.Balance, &acc.CreationDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (adb AccountDb) Save(ctx context.Context, account *domain.Account) error {
	_, err := adb.db.Exec(ctx, `
INSERT INTO accounts (id, user_id, balance, creation_date)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
                    SET BALANCE = EXCLUDED.balance
`, &account.Id, &account.UserId, &account.Balance, &account.CreationDate)
	return err
}

func (adb AccountDb) GetByUserId(ctx context.Context, userId int) (*domain.Account, error) {
	sql := `
SELECT id, user_id, balance, creation_date
FROM accounts
WHERE user_id=$1`
	row := adb.db.QueryRow(ctx, sql, userId)
	var acc domain.Account = domain.Account{}
	err := row.Scan(&acc.Id, &acc.UserId, &acc.Balance, &acc.CreationDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("account not found")
		}
		return nil, err
	}
	return &acc, nil
}
