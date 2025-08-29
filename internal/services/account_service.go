package services

import (
	"banking_system_golang/internal/models"
	"banking_system_golang/internal/repository"
	"context"
	"errors"
)

type AccountService struct {
	repo *repository.AccountRepository
}

func NewAccountService(repo *repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, owner string) (*models.Account, error) {
	acc := &models.Account{
		Owner:   owner,
		Balance: 0,
	}
	_, err := s.repo.Create(ctx, acc)
	return acc, err
}

func (s *AccountService) GetAccountByID(ctx context.Context, id int) (*models.Account, error) {
	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return acc, err
}

func (s *AccountService) AddMoneyToBalance(ctx context.Context, id int, amount float64) (*models.Account, error) {
	if amount < 0 {
		return nil, errors.New("сумма должна быть больше 0")
	}

	newBal, err := s.repo.AddMoneyToBalance(ctx, id, amount)
	if err != nil {
		return nil, err
	}

	acc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	acc.Balance = newBal

	return acc, err
}

func (s *AccountService) Transaction(ctx context.Context, fromID, toID int, amount float64) (*models.TransactionResult, error) {
	if amount < 0 {
		return nil, errors.New("сумма должна быть больше 0")
	}

	fromAcc, toAcc, err := s.repo.Transaction(ctx, fromID, toID, amount)
	if err != nil {
		return nil, err
	}

	return &models.TransactionResult{FromAcc: fromAcc, ToAcc: toAcc}, nil
}
