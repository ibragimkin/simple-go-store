package service

import (
	"context"
	"errors"
	"payment-service/internal/domain"
	"testing"
	"time"
)

type mockAccountRepository struct {
	getByIdFunc     func(ctx context.Context, id int) (*domain.Account, error)
	getByUserIdFunc func(ctx context.Context, userId int) (*domain.Account, error)
	saveFunc        func(ctx context.Context, account *domain.Account) error
}

func (m *mockAccountRepository) GetById(ctx context.Context, id int) (*domain.Account, error) {
	if m.getByIdFunc != nil {
		return m.getByIdFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockAccountRepository) GetByUserId(ctx context.Context, userId int) (*domain.Account, error) {
	if m.getByUserIdFunc != nil {
		return m.getByUserIdFunc(ctx, userId)
	}
	return nil, nil
}

func (m *mockAccountRepository) Save(ctx context.Context, account *domain.Account) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, account)
	}
	return nil
}

func TestAccountService_Deposit(t *testing.T) {
	ctx := context.Background()
	account := &domain.Account{Id: 1, Balance: 100, CreationDate: time.Now()}

	tests := []struct {
		name       string
		amount     float64
		setupRepo  func() *mockAccountRepository
		wantErr    bool
		finalValue float64
	}{
		{
			name:   "успешное пополнение",
			amount: 50,
			setupRepo: func() *mockAccountRepository {
				return &mockAccountRepository{
					getByIdFunc: func(ctx context.Context, id int) (*domain.Account, error) {
						return account, nil
					},
					saveFunc: func(ctx context.Context, acc *domain.Account) error {
						if acc.Balance != 150 {
							t.Errorf("ожидался баланс 150, получен %.2f", acc.Balance)
						}
						return nil
					},
				}
			},
			wantErr:    false,
			finalValue: 150,
		},
		{
			name:   "отрицательная сумма — ошибка",
			amount: -10,
			setupRepo: func() *mockAccountRepository {
				return &mockAccountRepository{
					getByIdFunc: func(ctx context.Context, id int) (*domain.Account, error) {
						return account, nil
					},
				}
			},
			wantErr: true,
		},
		{
			name:   "счёт не найден",
			amount: 10,
			setupRepo: func() *mockAccountRepository {
				return &mockAccountRepository{
					getByIdFunc: func(ctx context.Context, id int) (*domain.Account, error) {
						return nil, nil
					},
				}
			},
			wantErr: true,
		},
		{
			name:   "ошибка при сохранении",
			amount: 10,
			setupRepo: func() *mockAccountRepository {
				return &mockAccountRepository{
					getByIdFunc: func(ctx context.Context, id int) (*domain.Account, error) {
						return account, nil
					},
					saveFunc: func(ctx context.Context, acc *domain.Account) error {
						return errors.New("save error")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAccountService(tt.setupRepo())
			err := svc.Deposit(ctx, 1, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ожидалась ошибка=%v, получено %v", tt.wantErr, err)
			}
		})
	}
}
