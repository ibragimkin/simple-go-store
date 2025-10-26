package postgres

import (
	"context"
	"github.com/pashagolub/pgxmock/v3"
	"payment-service/internal/domain"
	"testing"
	"time"
)

func TestAccountDb_GetById(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("ошибка создания мока: %v", err)
	}
	defer mock.Close()

	account := domain.Account{
		Id:           1,
		UserId:       42,
		Balance:      100.5,
		CreationDate: time.Now(),
	}

	rows := pgxmock.NewRows([]string{"id", "user_id", "balance", "creation_date"}).
		AddRow(account.Id, account.UserId, account.Balance, account.CreationDate)
	mock.ExpectQuery(`SELECT id, user_id, balance, creation_date FROM accounts WHERE id=`).
		WithArgs(1).
		WillReturnRows(rows)

	db := AccountDb{db: mock}
	got, err := db.GetById(context.Background(), 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got.Id != account.Id || got.Balance != account.Balance {
		t.Errorf("ожидалось %+v, получено %+v", account, got)
	}
}

func TestAccountDb_Save(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("ошибка создания мока: %v", err)
	}
	defer mock.Close()

	account := domain.Account{
		Id:           2,
		UserId:       11,
		Balance:      500,
		CreationDate: time.Now(),
	}

	mock.ExpectExec(`INSERT INTO accounts`).
		WithArgs(&account.Id, &account.UserId, &account.Balance, &account.CreationDate).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	db := AccountDb{db: mock}
	err = db.Save(context.Background(), &account)
	if err != nil {
		t.Fatalf("ошибка при Save: %v", err)
	}
}

func TestAccountDb_GetByUserId(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("ошибка создания мока: %v", err)
	}
	defer mock.Close()

	account := domain.Account{
		Id:           3,
		UserId:       77,
		Balance:      250.25,
		CreationDate: time.Now(),
	}

	rows := pgxmock.NewRows([]string{"id", "user_id", "balance", "creation_date"}).
		AddRow(account.Id, account.UserId, account.Balance, account.CreationDate)
	mock.ExpectQuery(`SELECT id, user_id, balance, creation_date FROM accounts WHERE user_id=`).
		WithArgs(77).
		WillReturnRows(rows)

	db, _ := NewAccountDb(mock)
	got, err := db.GetByUserId(context.Background(), 77)
	if err != nil {
		t.Fatalf("ошибка при GetByUserId: %v", err)
	}
	if got.UserId != account.UserId {
		t.Errorf("ожидался user_id %d, получен %d", account.UserId, got.UserId)
	}
}
