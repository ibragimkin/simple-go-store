package service

import (
	"context"
	"fmt"
	"payment-service/internal/application/repository"
	"payment-service/internal/domain"
)

// PaymentService отвечает за обработку транзакций (пополнение и списание средств)
// и взаимодействие между счетами и историей транзакций.
type PaymentService struct {
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
}

// NewPaymentService создаёт новый экземпляр PaymentService.
// Возвращает ошибку, если один из репозиториев не инициализирован.
func NewPaymentService(accountsDb repository.AccountRepository, transactionsDb repository.TransactionRepository) (*PaymentService, error) {
	if accountsDb == nil || transactionsDb == nil {
		return nil, fmt.Errorf("nil repository")
	}
	return &PaymentService{accountRepository: accountsDb, transactionRepository: transactionsDb}, nil
}

// ProcessTransaction выбирает нужную операцию — Deposit или Withdraw —
// в зависимости от флага IsDeposit.
func (service *PaymentService) ProcessTransaction(ctx context.Context, transaction domain.Transaction) error {
	if transaction.IsDeposit {
		return service.Deposit(ctx, transaction)
	}
	return service.Withdraw(ctx, transaction)
}

// Withdraw выполняет списание средств со счёта пользователя.
// Проверяет, что транзакция не является депозитом, и что на балансе достаточно средств.
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
	return service.accountRepository.Save(ctx, account)
}

// Deposit выполняет пополнение счёта пользователя.
// Проверяет, что транзакция уникальна и не является списанием.
func (service *PaymentService) Deposit(ctx context.Context, transaction domain.Transaction) error {
	if !transaction.IsDeposit {
		return fmt.Errorf("transaction is not deposit")
	}

	// Проверяем, не существует ли уже транзакция с таким ID
	transactionFromDb, err := service.transactionRepository.GetById(ctx, transaction.Id)
	if err != nil || transactionFromDb != nil {
		return err
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
	return service.accountRepository.Save(ctx, account)
}
