package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/application/service"
	"order-service/internal/infrastructure/kafka"
	"strconv"
)

type OrderHandler struct {
	orderService *service.OrderService
	messageBus   *kafka.MessageBus
	ctx          context.Context
}

func NewOrderHandler(ctx context.Context, orderService *service.OrderService, bus *kafka.MessageBus) *OrderHandler {
	return &OrderHandler{
		messageBus:   bus,
		ctx:          ctx,
		orderService: orderService,
	}
}

// CreateOrder godoc
// @Summary Create order
// @Description Creates a new order
// @Accept json
// @Produce json
// @Param order body CreateOrderRequest true "Order info"
// @Success 201 {object} interface{}
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	orderRequest := CreateOrderRequest{}
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	order, err := h.orderService.CreateOrder(h.ctx, orderRequest.UserID, orderRequest.ItemID, orderRequest.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetOrder godoc
// @Param id path int true "id"
// @Success 200 {object} interface{}
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := getIntPathValue(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order, err := h.orderService.GetById(h.ctx, orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PayOrder godoc
// @Param id path int true "id"
// @Success 200 {object} interface{}
// @Router /orders/{id} [patch]
func (h *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	id, err := getIntPathValue(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order, err := h.orderService.GetById(h.ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if order.IsPayed {
		http.Error(w, "order is payed", http.StatusBadRequest)
		return
	}

	txn := h.orderService.CreateTransaction(h.ctx, order)
	txnJson, err := json.Marshal(txn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	strUserId := []byte(strconv.Itoa(txn.UserId))
	message, err := h.messageBus.SendMessage(h.ctx, strUserId, txnJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if string(message) == "OK" {
		err = h.orderService.PayOrder(h.ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, string(message), http.StatusBadRequest)

}

// GetUserOrders godoc
// @Param id path int true "id"
// @Success 200 {object} interface{}
// @Router /users/{id}/orders [get]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userId, err := getIntPathValue(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orders, err := h.orderService.GetUserOrders(h.ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getIntPathValue(r *http.Request, key string) (int, error) {
	valueStr := r.PathValue(key)
	if valueStr == "" {
		return 0, fmt.Errorf("%s not provided", key)
	}
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format", key)
	}
	return valueInt, nil
}

func getIntQueryValue(r *http.Request, key string) (int, error) {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return 0, fmt.Errorf("%s not provided", key)
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format", key)
	}
	return value, nil
}

func getFloat64QueryValue(r *http.Request, key string) (float64, error) {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return 0, fmt.Errorf("%s not provided", key)
	}
	valueFloat, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format", key)
	}
	return valueFloat, nil
}

type CreateOrderRequest struct {
	UserID int     `json:"user_id"`
	ItemID int     `json:"item_id"`
	Price  float64 `json:"price"`
}
