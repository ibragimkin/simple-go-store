package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"order-service/internal/domain"
)

// PgOrderDb реализует интерфейс OrderRepository,
// обеспечивая доступ к данным заказов через PostgreSQL.
type PgOrderDb struct {
	db PgxPool
}

// NewPgOrderDb создаёт новый экземпляр PgOrderDb,
// используя переданный пул соединений PostgreSQL.
func NewPgOrderDb(pool PgxPool) (*PgOrderDb, error) {
	return &PgOrderDb{db: pool}, nil
}

// GetById возвращает заказ по его ID из базы данных.
// Если заказ не найден — возвращает ошибку с pgx.ErrNoRows.
// При других ошибках возвращает ошибку выполнения SQL-запроса.
func (p *PgOrderDb) GetById(ctx context.Context, id int) (*domain.Order, error) {
	sql := `
		SELECT id, user_id, item_id, amount, is_payed, creation_date, payment_date 
		FROM orders 
		WHERE id = $1`
	row := p.db.QueryRow(ctx, sql, &id)

	var order domain.Order
	err := row.Scan(&order.Id, &order.UserId, &order.ItemId, &order.Amount,
		&order.IsPayed, &order.CreationDate, &order.PaymentDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, fmt.Errorf("error getting order: %w", err)
	}
	return &order, nil
}

// Save сохраняет заказ в базу данных.
// Если заказ с таким ID уже существует — обновляет флаги оплаты и дату платежа.
// При ошибках выполнения SQL-запроса возвращает подробную ошибку.
func (p *PgOrderDb) Save(ctx context.Context, order *domain.Order) error {
	sql := `
		INSERT INTO orders(id, user_id, item_id, amount, is_payed, creation_date, payment_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE 
		SET is_payed = EXCLUDED.is_payed,
		    payment_date = EXCLUDED.payment_date;`

	_, err := p.db.Exec(ctx, sql,
		&order.Id, &order.UserId, &order.ItemId,
		&order.Amount, &order.IsPayed, &order.CreationDate, order.PaymentDate)
	if err != nil {
		return fmt.Errorf("error inserting order: %w", err)
	}
	return nil
}

// GetUserOrders возвращает все заказы, принадлежащие пользователю с указанным ID.
// Если при запросе или чтении данных возникает ошибка — возвращает её.
func (p *PgOrderDb) GetUserOrders(ctx context.Context, userId int) ([]domain.Order, error) {
	sql := `
		SELECT id, user_id, item_id, amount, is_payed, creation_date, payment_date 
		FROM orders 
		WHERE user_id = $1`

	rows, err := p.db.Query(ctx, sql, &userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(&order.Id, &order.UserId, &order.ItemId,
			&order.Amount, &order.IsPayed, &order.CreationDate, &order.PaymentDate)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return orders, nil
}
