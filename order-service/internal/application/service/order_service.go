package service

import (
	"context"
	"fmt"
	"math/rand"
	"order-service/internal/application/repository"
	"order-service/internal/domain"
	"time"
)

type OrderService struct {
	orderRepository repository.OrderRepository
}

func NewOrderService(orderRepository repository.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

func (os *OrderService) GetById(ctx context.Context, id int) (*domain.Order, error) {
	order, err := os.orderRepository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id %d: %w", id, err)
	}
	return order, nil
}

func (os *OrderService) Save(ctx context.Context, account *domain.Order) error {
	err := os.orderRepository.Save(ctx, account)
	if err != nil {
		return fmt.Errorf("error saving order: %w", err)
	}
	return nil
}

func (os *OrderService) CreateOrder(ctx context.Context, userId int, itemId int, amount float64) (*domain.Order, error) {
	id := rand.Intn(2147483645) // TODO: change to go-uuid
	order := &domain.Order{
		Id:           id,
		UserId:       userId,
		ItemId:       itemId,
		Amount:       amount,
		IsPayed:      false,
		CreationDate: time.Now(),
		PaymentDate:  nil,
	}
	err := os.orderRepository.Save(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("error saving order: %w", err)
	}
	return order, nil
}

func (os *OrderService) GetUserOrders(ctx context.Context, userId int) ([]domain.Order, error) {
	orders, err := os.orderRepository.GetUserOrders(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user orders: %w", err)
	}
	return orders, nil
}

func (os *OrderService) PayOrder(ctx context.Context, id int) error {
	order, err := os.orderRepository.GetById(ctx, id)
	if err != nil {
		return err
	}
	if order.IsPayed {
		return fmt.Errorf("order is already payed")
	}
	order.IsPayed = true
	paymentDate := time.Now()
	order.PaymentDate = &paymentDate
	err = os.Save(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (os *OrderService) CreateTransaction(ctx context.Context, order *domain.Order) *domain.Transaction {
	id := rand.Intn(2147483645)
	return &domain.Transaction{
		Id:        id,
		UserId:    order.UserId,
		IsDeposit: false,
		Amount:    order.Amount,
		Date:      time.Now(),
	}
}
