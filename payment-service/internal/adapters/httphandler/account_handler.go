package httphandler

import (
	"context"
	"encoding/json"
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
// @Param amount query number true "amount"
// @Success 200 {object} interface{}
// @Router /accounts/{id} [patch]
func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id") // Extract path parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}
	err = h.accountService.Deposit(h.ctx, id, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CreateAccount
// @Param user_id query int true "user id"
// @Success 201 {object} interface{}
// @Router /accounts [post]
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")
	if userIdStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user id format", http.StatusBadRequest)
		return
	}
	account, err := h.accountService.CreateAccount(h.ctx, userId)
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
