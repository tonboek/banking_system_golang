package repository

import (
	"banking_system_golang/internal/models"
	"context"
	"database/sql"
	"fmt"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, acc *models.Account) (int64, error) {
	q := `INSERT INTO accounts (owner, balance, created_at)
	          VALUES ($1, $2, NOW()) RETURNING id`
	err := r.db.QueryRowContext(ctx, q, acc.Owner, acc.Balance).Scan(&acc.ID)
	return acc.ID, err
}

func (r *AccountRepository) GetByID(ctx context.Context, id int) (*models.Account, error) {
	acc := &models.Account{}
	q := `SELECT id, owner, balance, created_at FROM accounts WHERE id=$1`
	err := r.db.QueryRowContext(ctx, q, id).Scan(&acc.ID, &acc.Owner, &acc.Balance, &acc.CreatedAt)
	if err != nil {
		return nil, err
	}

	// TODO: ** несущ. айди **

	return acc, nil
}

func (r *AccountRepository) AddMoneyToBalance(ctx context.Context, id int, amount float64) (float64, error) {
	q := `UPDATE accounts SET balance = balance + $1 WHERE id = $2 RETURNING balance`
	var newBal float64
	err := r.db.QueryRowContext(ctx, q, amount, id).Scan(&newBal)

	return newBal, err
}

func (r *AccountRepository) Transaction(ctx context.Context, fromID, toID int, amount float64) (*models.Account, *models.Account, *models.Transaction, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	defer tx.Rollback()

	var fromAcc models.Account
	queryFromID := `UPDATE accounts
		 SET balance = balance - $1
		 WHERE id = $2 AND balance >= $1
		 RETURNING id, owner, balance, created_at`
	err = tx.QueryRowContext(ctx, queryFromID, amount, fromID).Scan(&fromAcc.ID, &fromAcc.Owner, &fromAcc.Balance, &fromAcc.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil, fmt.Errorf("недостаточно средств на счету у ID: [%d], текущий баланс: %f", fromID, fromAcc.Balance)
	}

	if err != nil {
		return nil, nil, nil, err
	}

	var toAcc models.Account
	queryToID := `UPDATE accounts
		 SET balance = balance + $1
		 WHERE id = $2
		 RETURNING id, owner, balance, created_at`
	err = tx.QueryRowContext(ctx, queryToID, amount, toID).Scan(&toAcc.ID, &toAcc.Owner, &toAcc.Balance, &toAcc.CreatedAt)
	if err != nil {
		return nil, nil, nil, err
	}

	var t models.Transaction
	queryTx := `
		INSERT INTO transactions (from_id, to_id, amount, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, from_id, to_id, amount, created_at
	`
	err = tx.QueryRowContext(ctx, queryTx, fromID, toID, amount).Scan(&t.ID, &t.FromID, &t.ToID, &t.Amount, &t.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, err
	}

	return &fromAcc, &toAcc, &t, nil
}