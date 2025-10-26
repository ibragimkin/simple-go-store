package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"payment-service/internal/application/service"
	"payment-service/internal/domain"
	"strconv"
	"testing"
)

type mockAccountRepository struct {
	data map[int]domain.Account
}

func (m *mockAccountRepository) GetById(ctx context.Context, id int) (*domain.Account, error) {
	acc, ok := m.data[id]
	if !ok {
		return nil, errors.New("id not found")
	}
	return &acc, nil
}

func (m *mockAccountRepository) GetByUserId(ctx context.Context, userId int) (*domain.Account, error) {
	for _, acc := range m.data {
		if acc.UserId == userId {
			return &acc, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockAccountRepository) Save(ctx context.Context, account *domain.Account) error {
	m.data[account.Id] = *account
	return nil
}

// --- Тесты ---

func setupTestEnv(t *testing.T) (context.Context, *service.AccountService) {
	t.Helper()
	ctx := context.Background()
	accDb := &mockAccountRepository{data: make(map[int]domain.Account)}
	accService := service.NewAccountService(accDb)
	return ctx, accService
}

func TestGetAccount_Success(t *testing.T) {
	ctx, accService := setupTestEnv(t)
	acc, err := accService.CreateAccount(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	_ = accService.Deposit(ctx, acc.Id, 100)
	handler := NewAccountHandler(context.Background(), accService)
	req := httptest.NewRequest(http.MethodGet, "/accounts/", nil)
	req.SetPathValue("id", strconv.Itoa(acc.Id))
	w := httptest.NewRecorder()

	handler.GetAccount(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if acc.Balance != 100 {
		t.Errorf("expected balance 100, got %v", acc.Balance)
	}
}

func TestGetAccount_InvalidID(t *testing.T) {
	_, accService := setupTestEnv(t)
	handler := NewAccountHandler(context.Background(), accService)
	req := httptest.NewRequest(http.MethodGet, "/accounts/", nil)
	req.SetPathValue("id", "abc")
	w := httptest.NewRecorder()

	handler.GetAccount(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetAccount_NoAccount(t *testing.T) {
	_, accService := setupTestEnv(t)
	handler := NewAccountHandler(context.Background(), accService)
	req := httptest.NewRequest(http.MethodGet, "/accounts/", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.GetAccount(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetUsersAccount_Success(t *testing.T) {
	ctx, accService := setupTestEnv(t)
	_, err := accService.CreateAccount(ctx, 123)
	if err != nil {
		t.Errorf("error creating account: %v", err)
	}
	handler := NewAccountHandler(context.Background(), accService)
	req := httptest.NewRequest(http.MethodGet, "/users/{id}/account", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	handler.GetUsersAccount(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetUsersAccount_NoAccount(t *testing.T) {
	_, accService := setupTestEnv(t)
	handler := NewAccountHandler(context.Background(), accService)
	req := httptest.NewRequest(http.MethodGet, "/users/{id}/account", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()
	handler.GetUsersAccount(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeposit_Success(t *testing.T) {
	ctx, accService := setupTestEnv(t)
	account, err := accService.CreateAccount(ctx, 123)
	if err != nil {
		t.Errorf("error creating account: %v", err)
	}
	handler := NewAccountHandler(context.Background(), accService)

	body := bytes.NewBufferString(`{"amount": 25}`)
	req := httptest.NewRequest(http.MethodPatch, "/accounts/", body)
	req.SetPathValue("id", strconv.Itoa(account.Id))
	w := httptest.NewRecorder()

	handler.Deposit(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	account, err = accService.GetAccount(ctx, account.Id)
	if err != nil {
		t.Errorf("error getting account: %v", err)
	}
	if account.Balance != 25 {
		t.Errorf("expected balance 25, got %v", account.Balance)
	}
}

func TestDeposit_NoAccount(t *testing.T) {
	_, accService := setupTestEnv(t)
	handler := NewAccountHandler(context.Background(), accService)

	body := bytes.NewBufferString(`{"amount": 25}`)
	req := httptest.NewRequest(http.MethodPatch, "/accounts/", body)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	handler.Deposit(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestCreateAccount_Success(t *testing.T) {
	ctx, accService := setupTestEnv(t)
	acc, err := accService.CreateAccount(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	jsonReq, err := json.Marshal(CreateAccountRequest{UserId: 123})
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewBuffer(jsonReq)
	handler := NewAccountHandler(context.Background(), accService)
	req := httptest.NewRequest(http.MethodPost, "/accounts/", body)
	w := httptest.NewRecorder()
	handler.CreateAccount(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&acc); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if acc.UserId != 123 {
		t.Errorf("expected user_id = 123, got %v", acc.UserId)
	}
}

func TestCreateAccount_AlreadyExists(t *testing.T) {
	ctx, accService := setupTestEnv(t)
	acc, err := accService.CreateAccount(ctx, 123)
	if err != nil {
		t.Fatal(err)
	}
	_ = accService.Deposit(ctx, acc.Id, 100)
	jsonReq, err := json.Marshal(CreateAccountRequest{UserId: 123})
	if err != nil {
		t.Fatal(err)
	}
	body := bytes.NewBuffer(jsonReq)
	handler := NewAccountHandler(context.Background(), accService)
	req := httptest.NewRequest(http.MethodPost, "/accounts/", body)
	w := httptest.NewRecorder()
	handler.CreateAccount(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected 409, got %d", resp.StatusCode)
	}

	acc, err = accService.GetAccount(ctx, acc.Id)
	if err != nil {
		t.Fatal(err)
	}
	if acc.Balance != 100 {
		t.Errorf("expected balance = 100 (not changed), got %v", acc.Balance)
	}
}

func TestGetIntPathValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/accounts/10", nil)
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
