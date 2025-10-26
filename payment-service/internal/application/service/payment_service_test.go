package service

import (
	"context"
	"errors"
	"payment-service/internal/domain"
	"testing"
)

type mockAccountRepo struct {
	getByUserIdFunc func(ctx context.Context, id int) (*domain.Account, error)
	saveFunc        func(ctx context.Context, acc *domain.Account) error
}

func (m *mockAccountRepo) GetByUserId(ctx context.Context, id int) (*domain.Account, error) {
	if m.getByUserIdFunc != nil {
		return m.getByUserIdFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockAccountRepo) GetById(ctx context.Context, id int) (*domain.Account, error) {
	return nil, nil
}
func (m *mockAccountRepo) Save(ctx context.Context, acc *domain.Account) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, acc)
	}
	return nil
}

type mockTransactionRepo struct {
	getByIdFunc func(ctx context.Context, id int) (*domain.Transaction, error)
	saveFunc    func(ctx context.Context, tx *domain.Transaction) error
}

func (m *mockTransactionRepo) GetById(ctx context.Context, id int) (*domain.Transaction, error) {
	if m.getByIdFunc != nil {
		return m.getByIdFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockTransactionRepo) Save(ctx context.Context, tx *domain.Transaction) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, tx)
	}
	return nil
}

func TestPaymentService_Deposit(t *testing.T) {
	ctx := context.Background()
	account := &domain.Account{Id: 1, UserId: 10, Balance: 100}

	tests := []struct {
		name      string
		tx        domain.Transaction
		setupMock func() (*mockAccountRepo, *mockTransactionRepo)
		wantErr   bool
	}{
		{
			name: "успешное пополнение",
			tx:   domain.Transaction{Id: 1, UserId: 10, IsDeposit: true, Amount: 50},
			setupMock: func() (*mockAccountRepo, *mockTransactionRepo) {
				return &mockAccountRepo{
						getByUserIdFunc: func(ctx context.Context, id int) (*domain.Account, error) {
							return account, nil
						},
						saveFunc: func(ctx context.Context, acc *domain.Account) error {
							if acc.Balance != 150 {
								t.Errorf("ожидался баланс 150, получен %.2f", acc.Balance)
							}
							return nil
						},
					}, &mockTransactionRepo{
						getByIdFunc: func(ctx context.Context, id int) (*domain.Transaction, error) {
							return nil, nil
						},
					}
			},
			wantErr: false,
		},
		{
			name: "ошибка при сохранении транзакции",
			tx:   domain.Transaction{Id: 1, UserId: 10, IsDeposit: true, Amount: 50},
			setupMock: func() (*mockAccountRepo, *mockTransactionRepo) {
				return &mockAccountRepo{
						getByUserIdFunc: func(ctx context.Context, id int) (*domain.Account, error) {
							return account, nil
						},
					}, &mockTransactionRepo{
						getByIdFunc: func(ctx context.Context, id int) (*domain.Transaction, error) {
							return nil, nil
						},
						saveFunc: func(ctx context.Context, tx *domain.Transaction) error {
							return errors.New("db error")
						},
					}
			},
			wantErr: true,
		},
		{
			name: "некорректный тип транзакции",
			tx:   domain.Transaction{Id: 1, UserId: 10, IsDeposit: false, Amount: 50},
			setupMock: func() (*mockAccountRepo, *mockTransactionRepo) {
				return &mockAccountRepo{}, &mockTransactionRepo{}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accRepo, txRepo := tt.setupMock()
			svc, _ := NewPaymentService(accRepo, txRepo)
			err := svc.Deposit(ctx, tt.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ожидалась ошибка=%v, получено %v", tt.wantErr, err)
			}
		})
	}
}
