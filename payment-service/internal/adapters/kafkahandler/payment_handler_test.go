package kafkahandler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"payment-service/internal/application/service"
	"payment-service/internal/domain"
	"testing"
	"time"
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

type mockTransactionRepository struct {
	data map[int]domain.Transaction
}

func (m *mockTransactionRepository) Save(ctx context.Context, transaction *domain.Transaction) error {
	m.data[transaction.Id] = *transaction
	return nil
}

func (m *mockTransactionRepository) GetById(ctx context.Context, id int) (*domain.Transaction, error) {
	tx, ok := m.data[id]
	if !ok {
		return nil, errors.New("id not found")
	}
	return &tx, nil
}

func setupTestEnv(t *testing.T) (context.Context, *service.PaymentService, *service.AccountService) {
	t.Helper()
	ctx := context.Background()
	accDb := &mockAccountRepository{data: make(map[int]domain.Account)}
	txDb := &mockTransactionRepository{data: make(map[int]domain.Transaction)}
	paymentService, _ := service.NewPaymentService(accDb, txDb)
	accService := service.NewAccountService(accDb)
	return ctx, paymentService, accService
}

func TestPaymentHandler_Success(t *testing.T) {
	ctx, paymentService, accService := setupTestEnv(t)
	acc, _ := accService.CreateAccount(ctx, 123)
	err := accService.Deposit(ctx, acc.Id, 10000)
	if err != nil {
		t.Errorf("error depositing account: %v", err)
	}
	handler := NewPaymentHandler(paymentService)
	tx := &domain.Transaction{
		Id:        1,
		UserId:    acc.UserId,
		IsDeposit: false,
		Amount:    999,
		Date:      time.Now(),
	}
	txJson, _ := json.Marshal(tx)
	msg := &kafka.Message{Key: []byte("123"), Value: txJson}
	res, err := handler(ctx, msg)
	if err != nil {
		t.Errorf("error processing transaction: %v", err)
	}
	if string(res.Value) != "OK" {
		t.Errorf("expected OK, got %s", string(res.Value))
	}
	return
}

func TestPaymentHandler_Fail(t *testing.T) {
	ctx, paymentService, accService := setupTestEnv(t)
	acc, _ := accService.CreateAccount(ctx, 123)
	err := accService.Deposit(ctx, acc.Id, 100)
	if err != nil {
		t.Errorf("error depositing account: %v", err)
	}
	handler := NewPaymentHandler(paymentService)
	tx := &domain.Transaction{
		Id:        1,
		UserId:    acc.UserId,
		IsDeposit: false,
		Amount:    125,
		Date:      time.Now(),
	}
	txJson, _ := json.Marshal(tx)
	msg := &kafka.Message{Key: []byte("123"), Value: txJson}
	_, err = handler(ctx, msg)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	acc, _ = accService.GetUsersAccount(ctx, 123)
	if acc.Balance != 100 {
		t.Errorf("account balance got changed")
	}
	return
}
