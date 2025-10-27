package service

import (
	"context"
	"fmt"
	"math/rand"
	"order-service/internal/application/repository"
	"order-service/internal/domain"
	"time"
)

// OrderService отвечает за бизнес-логику, связанную с заказами.
// Он использует репозиторий для сохранения и получения данных о заказах.
type OrderService struct {
	orderRepository repository.OrderRepository
}

// NewOrderService создаёт новый экземпляр OrderService.
func NewOrderService(orderRepository repository.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

// GetById возвращает заказ по его ID.
// Возвращает ошибку, если заказ не найден или произошла ошибка при обращении к репозиторию.
func (os *OrderService) GetById(ctx context.Context, id int) (*domain.Order, error) {
	order, err := os.orderRepository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order by id %d: %w", id, err)
	}
	return order, nil
}

// Save сохраняет заказ в репозитории.
// Возвращает ошибку, если операция не удалась.
func (os *OrderService) Save(ctx context.Context, order *domain.Order) error {
	err := os.orderRepository.Save(ctx, order)
	if err != nil {
		return fmt.Errorf("error saving order: %w", err)
	}
	return nil
}

// CreateOrder создаёт новый заказ с указанными параметрами.
// Генерирует случайный ID.
// Возвращает созданный заказ или ошибку при сохранении.
func (os *OrderService) CreateOrder(ctx context.Context, userId int, itemId int, amount float64) (*domain.Order, error) {
	id := rand.Intn(2147483645)
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

// GetUserOrders возвращает все заказы пользователя по его ID.
// Возвращает ошибку, если произошёл сбой при получении данных.
func (os *OrderService) GetUserOrders(ctx context.Context, userId int) ([]domain.Order, error) {
	orders, err := os.orderRepository.GetUserOrders(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user orders: %w", err)
	}
	return orders, nil
}

// PayOrder помечает заказ как оплаченный, устанавливая дату оплаты.
// Возвращает ошибку, если заказ не найден или уже оплачен.
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

// CreateTransaction создаёт транзакцию для оплаты заказа.
// Генерирует случайный ID.
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
