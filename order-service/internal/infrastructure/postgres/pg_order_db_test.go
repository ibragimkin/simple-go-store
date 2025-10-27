package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
	"order-service/internal/domain"
	"testing"
	"time"
)

func TestNewPgOrderDb(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	db, err := NewPgOrderDb(mock)
	require.NoError(t, err)
	require.NotNil(t, db)
}

func TestPgOrderDb_GetById_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	order := domain.Order{
		Id:           1,
		UserId:       2,
		ItemId:       3,
		Amount:       100,
		IsPayed:      false,
		CreationDate: time.Now(),
		PaymentDate:  nil,
	}

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "item_id", "amount", "is_payed", "creation_date", "payment_date",
	}).AddRow(order.Id, order.UserId, order.ItemId, order.Amount,
		order.IsPayed, order.CreationDate, order.PaymentDate)

	mock.ExpectQuery("SELECT id, user_id, item_id").
		WithArgs(&order.Id).
		WillReturnRows(rows)

	db, _ := NewPgOrderDb(mock)
	result, err := db.GetById(context.Background(), order.Id)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, order.Id, result.Id)
	require.Equal(t, order.UserId, result.UserId)
}

func TestPgOrderDb_GetById_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	id := 1

	mock.ExpectQuery("SELECT id, user_id, item_id").
		WithArgs(&id).
		WillReturnError(pgx.ErrNoRows)

	db, _ := NewPgOrderDb(mock)
	res, err := db.GetById(context.Background(), id)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestPgOrderDb_Save_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	order := domain.Order{
		Id:           1,
		UserId:       2,
		ItemId:       3,
		Amount:       50,
		IsPayed:      false,
		CreationDate: time.Now(),
	}

	mock.ExpectExec("INSERT INTO orders").
		WithArgs(&order.Id, &order.UserId, &order.ItemId,
			&order.Amount, &order.IsPayed, &order.CreationDate, order.PaymentDate).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	db, _ := NewPgOrderDb(mock)
	err = db.Save(context.Background(), &order)
	require.NoError(t, err)
}

func TestPgOrderDb_Save_Error(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	order := domain.Order{Id: 1}
	mock.ExpectExec("INSERT INTO orders").
		WithArgs(&order.Id, &order.UserId, &order.ItemId,
			&order.Amount, &order.IsPayed, &order.CreationDate, order.PaymentDate).
		WillReturnError(errors.New("insert failed"))

	db, _ := NewPgOrderDb(mock)
	err = db.Save(context.Background(), &order)
	require.Error(t, err)
}

func TestPgOrderDb_GetUserOrders_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	userId := 10
	order1 := domain.Order{Id: 1, UserId: userId, ItemId: 2, Amount: 30, IsPayed: false, CreationDate: time.Now(), PaymentDate: nil}
	order2 := domain.Order{Id: 2, UserId: userId, ItemId: 3, Amount: 50, IsPayed: false, CreationDate: time.Now(), PaymentDate: nil}

	rows := pgxmock.NewRows([]string{
		"id", "user_id", "item_id", "amount", "is_payed", "creation_date", "payment_date",
	}).AddRow(order1.Id, order1.UserId, order1.ItemId, order1.Amount, order1.IsPayed, order1.CreationDate, order1.PaymentDate).
		AddRow(order2.Id, order2.UserId, order2.ItemId, order2.Amount, order2.IsPayed, order2.CreationDate, order2.PaymentDate)

	mock.ExpectQuery("SELECT id, user_id, item_id, amount, is_payed, creation_date, payment_date").
		WithArgs(&userId).
		WillReturnRows(rows)

	db, _ := NewPgOrderDb(mock)
	orders, err := db.GetUserOrders(context.Background(), userId)
	require.NoError(t, err)
	require.Len(t, orders, 2)
	require.Equal(t, order1.Id, orders[0].Id)
	require.Equal(t, order2.Id, orders[1].Id)
}

func TestPgOrderDb_GetUserOrders_QueryError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	userId := 99
	mock.ExpectQuery("SELECT id, user_id, item_id").
		WithArgs(&userId).
		WillReturnError(errors.New("query failed"))

	db, _ := NewPgOrderDb(mock)
	orders, err := db.GetUserOrders(context.Background(), userId)
	require.Error(t, err)
	require.Nil(t, orders)
}
