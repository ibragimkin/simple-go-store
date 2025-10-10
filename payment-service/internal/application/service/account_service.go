package service

import (
	"context"
	"errors"
	"math/rand"
	"payment-service/internal/application/repository"
	"payment-service/internal/domain"
	"time"
)

type AccountService struct {
	accountDb repository.AccountRepository
}

func NewAccountService(accountDb repository.AccountRepository) *AccountService {
	return &AccountService{accountDb: accountDb}
}

func (as *AccountService) CreateAccount(ctx context.Context, userID int) (*domain.Account, error) {
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
	err := as.accountDb.Save(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

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
