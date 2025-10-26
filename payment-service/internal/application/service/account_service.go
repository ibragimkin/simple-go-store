package service

import (
	"context"
	"errors"
	"math/rand"
	"payment-service/internal/application/repository"
	"payment-service/internal/domain"
	"time"
)

// AccountService предоставляет бизнес-логику для работы со счетами пользователей.
type AccountService struct {
	accountDb repository.AccountRepository
}

// NewAccountService создаёт новый экземпляр AccountService.
func NewAccountService(accountDb repository.AccountRepository) *AccountService {
	return &AccountService{accountDb: accountDb}
}

// CreateAccount создаёт новый счёт для пользователя.
// Идентификатор генерируется случайно (временно), баланс устанавливается в 0.
// Если сохранение не удалось — возвращает ошибку.
func (as *AccountService) CreateAccount(ctx context.Context, userID int) (*domain.Account, error) {
	dupe, err := as.accountDb.GetByUserId(ctx, userID)
	if err == nil && dupe != nil {
		return nil, errors.New("account with that user_id already exists")
	}
	id := rand.Intn(2147483645) // TODO: change to go-uuid
	for range 5 {
		_, err := as.accountDb.GetById(ctx, id)
		if err != nil {
			break
		}
		id++
	}
	account := &domain.Account{
		Id:           id,
		UserId:       userID,
		Balance:      0,
		CreationDate: time.Now(),
	}
	err = as.accountDb.Save(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetAccount возвращает счёт по его идентификатору.
// Если счёт не найден, возвращает ошибку.
func (as *AccountService) GetAccount(ctx context.Context, id int) (*domain.Account, error) {
	account, err := as.accountDb.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	return account, nil
}

// GetUsersAccount возвращает счёт пользователя по его userId.
// Если счёт не найден — возвращает ошибку.
func (as *AccountService) GetUsersAccount(ctx context.Context, userId int) (*domain.Account, error) {
	account, err := as.accountDb.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	return account, nil
}

// Deposit пополняет баланс счёта на указанную сумму.
// Если счёт не найден или сумма отрицательная — возвращает ошибку.
func (as *AccountService) Deposit(ctx context.Context, id int, amount float64) error {
	account, err := as.accountDb.GetById(ctx, id)
	if err != nil {
		return err
	}
	if account == nil {
		return errors.New("account not found")
	}
	err = account.Deposit(amount)
	if err != nil {
		return err
	}
	err = as.accountDb.Save(ctx, account)
	return err
}

//func (as *AccountService) DeleteAccount(ctx context.Context, id int) error {} TODO
