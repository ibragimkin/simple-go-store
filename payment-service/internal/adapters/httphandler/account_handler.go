package httphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"payment-service/internal/application/service"
	"strconv"
)

type AccountHandler struct {
	accountService *service.AccountService
	ctx            context.Context
}

func NewAccountHandler(ctx context.Context, accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService, ctx: ctx}
}

// GetAccount godoc
// @Summary      Получить аккаунт
// @Description  Возвращает информацию об аккаунте по ID
// @Tags         accounts
// @Param        id   path      int  true  "Account ID"
// @Produce      json
// @Success 200 {object} interface{}
// @Failure      404  {object} interface{}
// @Router       /accounts/{id} [get]
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // Extract path parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.GetAccount(h.ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Printf("Failed to encode account to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Deposit godoc
// @Param id path int true "account id"
// @Param amount body DepositRequest true "Deposit amount"
// @Success 200 {object} interface{}
// @Router /accounts/{id} [patch]
func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	depositRequest := DepositRequest{}
	err = json.NewDecoder(r.Body).Decode(&depositRequest)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}

	err = h.accountService.Deposit(h.ctx, id, depositRequest.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CreateAccount godoc
// @Summary Create account
// @Description Creates a new account
// @Param data body CreateAccountRequest true "Account info"
// @Success 201 {object} interface{}
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	createRequest := CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		http.Error(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	account, err := h.accountService.CreateAccount(h.ctx, createRequest.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Printf("Failed to encode account to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetUsersAccount
// @Summary Get users account by id
// @Param        id   path      int  true  "User ID"
// @Param data body GetByUserIdRequest true "Get by users id"
// @Success 201 {object} interface{}
// @Router /users/{id}/account [get]
func (h *AccountHandler) GetUsersAccount(w http.ResponseWriter, r *http.Request) {
	userId, err := getIntPathValue(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	account, err := h.accountService.GetUsersAccount(h.ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Printf("Failed to encode account to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

type CreateAccountRequest struct {
	UserId int `json:"user_id"`
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type GetByUserIdRequest struct {
	UserId int `json:"user_id"`
}
