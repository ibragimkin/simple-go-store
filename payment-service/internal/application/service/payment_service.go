package service

import (
	"context"
	"fmt"
	"payment-service/internal/application/repository"
	"payment-service/internal/domain"
)

type PaymentService struct {
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
}

func NewPaymentService(accountsDb repository.AccountRepository, transactionsDb repository.TransactionRepository) (*PaymentService, error) {
	if accountsDb == nil || transactionsDb == nil {
		return nil, fmt.Errorf("nil repository")
	}
	return &PaymentService{accountRepository: accountsDb, transactionRepository: transactionsDb}, nil
}

func (service *PaymentService) ProcessTransaction(ctx context.Context, transaction domain.Transaction) error {
	if transaction.IsDeposit {
		return service.Deposit(ctx, transaction)
	}
	return service.Withdraw(ctx, transaction)
}

func (service *PaymentService) Withdraw(ctx context.Context, transaction domain.Transaction) error {
	if transaction.IsDeposit {
		return fmt.Errorf("transaction is not withdrawal")
	}

	account, err := service.accountRepository.GetByUserId(ctx, transaction.UserId)
	if err != nil || account == nil {
		return err
	}
	err = account.Withdraw(transaction.Amount)
	if err != nil {
		return err
	}
	err = service.transactionRepository.Save(ctx, &transaction)
	if err != nil {
		return err
	}
	err = service.accountRepository.Save(ctx, account)
	return err
}

func (service *PaymentService) Deposit(ctx context.Context, transaction domain.Transaction) error {
	if !transaction.IsDeposit {
		return fmt.Errorf("transaction is not deposit")
	}
	transactionFromDb, err := service.transactionRepository.GetById(ctx, transaction.Id)
	if err != nil || transactionFromDb != nil {
		return nil
	}
	account, err := service.accountRepository.GetByUserId(ctx, transaction.UserId)
	if err != nil {
		return err
	}
	err = account.Deposit(transaction.Amount)
	if err != nil {
		return err
	}
	err = service.transactionRepository.Save(ctx, &transaction)
	if err != nil {
		return err
	}
	err = service.accountRepository.Save(ctx, account)
	return err
}
