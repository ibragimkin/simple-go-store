package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"order-service/internal/application/service"
	"order-service/internal/domain"
	"strconv"
	"testing"
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

func setupOrderTest(t *testing.T) (context.Context, *service.OrderService, *OrderHandler) {
	t.Helper()
	ctx := context.Background()
	orderDb := &mockAccountRepository{data: make(map[int]domain.Order)}
	orderService := service.NewOrderService(orderDb)
	handler := NewOrderHandler(ctx, orderService, nil)
	return ctx, orderService, handler
}

func TestCreateOrder_Success(t *testing.T) {
	_, svc, handler := setupOrderTest(t)

	body := bytes.NewBufferString(`{"user_id": 1, "item_id": 2, "price": 99.99}`)
	req := httptest.NewRequest(http.MethodPost, "/orders", body)
	w := httptest.NewRecorder()

	handler.CreateOrder(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 200 or 201, got %d", resp.StatusCode)
	}

	var order domain.Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if order.UserId != 1 || order.ItemId != 2 || order.Amount != 99.99 {
		t.Errorf("unexpected order data: %+v", order)
	}

	order2, err := svc.GetById(context.Background(), order.Id)
	if err != nil {
		t.Fatalf("get by id error: %v", err)
	}
	if order2.Id != order.Id {
		t.Errorf("unexpected order data: %+v", order2)
	}
}

func TestGetOrder_Success(t *testing.T) {
	ctx, svc, handler := setupOrderTest(t)

	order, _ := svc.CreateOrder(ctx, 1, 1, 50.0)
	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	req.SetPathValue("id", strconv.Itoa(order.Id))
	w := httptest.NewRecorder()

	handler.GetOrder(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var got domain.Order
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.Id != order.Id {
		t.Errorf("expected order id %d, got %d", order.Id, got.Id)
	}
}

func TestGetOrder_NotFound(t *testing.T) {
	_, _, handler := setupOrderTest(t)

	req := httptest.NewRequest(http.MethodGet, "/orders/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetOrder(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetUserOrders_Success(t *testing.T) {
	ctx, svc, handler := setupOrderTest(t)

	_, _ = svc.CreateOrder(ctx, 1, 10, 500)
	_, _ = svc.CreateOrder(ctx, 1, 20, 300)

	req := httptest.NewRequest(http.MethodGet, "/users/1/orders", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	handler.GetUserOrders(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var orders []domain.Order
	_ = json.NewDecoder(w.Body).Decode(&orders)
	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
}

func TestGetIntPathValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/10", nil)
	req.SetPathValue("id", "10")

	id, err := getIntPathValue(req, "id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 10 {
		t.Errorf("expected 10, got %v", id)
	}

	req.SetPathValue("id", "abc")
	_, err = getIntPathValue(req, "id")
	if err == nil {
		t.Error("expected error for invalid int format")
	}
}
