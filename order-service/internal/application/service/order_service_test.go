package service

import (
	"context"
	"errors"
	"order-service/internal/application/repository"
	"order-service/internal/domain"
	"testing"
	"time"
)

type mockAccountRepository struct {
	data map[int]domain.Order
}

func (m *mockAccountRepository) GetById(ctx context.Context, id int) (*domain.Order, error) {
	order, ok := m.data[id]
	if !ok {
		return nil, errors.New("id not found")
	}
	return &order, nil
}

func (m *mockAccountRepository) Save(ctx context.Context, order *domain.Order) error {
	m.data[order.Id] = *order
	return nil
}

func (m *mockAccountRepository) GetUserOrders(ctx context.Context, userId int) ([]domain.Order, error) {
	orders := make([]domain.Order, 0)
	for _, order := range m.data {
		if order.UserId == userId {
			orders = append(orders, order)
		}
	}
	if len(orders) == 0 {
		return nil, errors.New("user not found")
	}
	return orders, nil
}

func setupTestEnv(t *testing.T) (context.Context, repository.OrderRepository, *OrderService) {
	t.Helper()
	ctx := context.Background()
	orderDb := &mockAccountRepository{data: make(map[int]domain.Order)}
	orderService := NewOrderService(orderDb)
	return ctx, orderDb, orderService
}

func TestOrderService_CreateOrder(t *testing.T) {
	ctx, db, orderService := setupTestEnv(t)
	order, err := orderService.CreateOrder(ctx, 1, 2, 300)
	if err != nil {
		t.Errorf("error creating order: %v", err)
	}
	order2, err := db.GetById(ctx, order.Id)
	if err != nil {
		t.Errorf("error getting order: %v", err)
	}
	if order.Id != order2.Id {
		t.Errorf("order.Id != order.Id")
	}
}

func TestOrderService_Save(t *testing.T) {
	ctx, _, orderService := setupTestEnv(t)
	order := domain.Order{
		Id:           11,
		UserId:       123,
		ItemId:       144,
		Amount:       1000,
		IsPayed:      false,
		CreationDate: time.Time{},
		PaymentDate:  nil,
	}
	err := orderService.Save(ctx, &order)
	if err != nil {
		t.Errorf("error saving order: %v", err)
	}
	order2, err := orderService.GetById(ctx, order.Id)
	if err != nil {
		t.Errorf("error getting order: %v", err)
	}
	if order.Id != order2.Id {
		t.Errorf("order.Id != order.Id")
	}
	return
}

func TestOrderService_GetById_Fail(t *testing.T) {
	ctx, _, orderService := setupTestEnv(t)
	order, err := orderService.GetById(ctx, 1)
	if err == nil {
		t.Errorf("expected error getting order, got nil")
	}
	if order != nil {
		t.Errorf("expected order to be nil, got %v", order.Id)
	}
	return
}

func TestOrderService_GetUserOrders(t *testing.T) {
	ctx, _, orderService := setupTestEnv(t)
	_, err := orderService.CreateOrder(ctx, 1, 2, 300)
	if err != nil {
		t.Errorf("error creating order: %v", err)
	}
	_, err = orderService.CreateOrder(ctx, 1, 10, 500)
	if err != nil {
		t.Errorf("error creating order: %v", err)
	}
	orders, err := orderService.GetUserOrders(ctx, 1)
	if err != nil {
		t.Errorf("error getting user orders: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("len(orders) != 2")
	}
	if orders[0].ItemId != 2 && orders[1].ItemId != 10 {
		t.Errorf("item_id mismatch, expected 2 and 10 got %v and %v", orders[0].ItemId, orders[1].ItemId)
	}
	return
}

func TestOrderService_GetUserOrders_Fail(t *testing.T) {
	ctx, _, orderService := setupTestEnv(t)
	orders, err := orderService.GetUserOrders(ctx, 1)
	if err == nil {
		t.Errorf("expected error getting user orders, got nil")
	}
	if orders != nil {
		t.Errorf("expected order to be nil, got %v", orders)
	}
	return
}

func TestPayOrder_Success(t *testing.T) {
	ctx, db, svc := setupTestEnv(t)

	order := domain.Order{
		Id:           1,
		UserId:       10,
		ItemId:       100,
		Amount:       500,
		IsPayed:      false,
		CreationDate: time.Now(),
	}
	_ = svc.Save(ctx, &order)

	err := svc.PayOrder(ctx, order.Id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := db.GetById(ctx, 1)
	if !updated.IsPayed {
		t.Errorf("expected order to be paid")
	}
	if updated.PaymentDate == nil {
		t.Errorf("expected payment date to be set")
	}
}

func TestPayOrder_AlreadyPaid(t *testing.T) {
	ctx, db, svc := setupTestEnv(t)

	now := time.Now()
	order := domain.Order{
		Id:           2,
		UserId:       10,
		ItemId:       100,
		Amount:       700,
		IsPayed:      true,
		CreationDate: now,
		PaymentDate:  &now,
	}
	_ = db.Save(ctx, &order)

	err := svc.PayOrder(ctx, 2)
	if err == nil {
		t.Errorf("expected error for already paid order, got nil")
	}
}

func TestCreateTransaction(t *testing.T) {
	ctx, _, svc := setupTestEnv(t)

	order := &domain.Order{
		Id:      3,
		UserId:  42,
		ItemId:  5,
		Amount:  1200.50,
		IsPayed: false,
	}

	tx := svc.CreateTransaction(ctx, order)

	if tx.UserId != order.UserId {
		t.Errorf("expected UserId %d, got %d", order.UserId, tx.UserId)
	}
	if tx.Amount != order.Amount {
		t.Errorf("expected Amount %.2f, got %.2f", order.Amount, tx.Amount)
	}
	if tx.IsDeposit {
		t.Errorf("expected IsDeposit = false for order transaction")
	}
}
