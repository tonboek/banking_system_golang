package models

import (
	"time"
)

// пользователь
type User struct {
	ID    int64 		`json:"id" db:"id"`
	Name string 		`json:"name" db:"name"`
	Username string 	`json:"username" db:"username"`
	Email string 		`json:"email" db:"email"`
	Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// банковский счет пользователя
type Account struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Owner     string    `json:"owner" db:"owner"`
	Balance   float64   `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// операция (банк. перевод)
type Transaction struct {
	ID        int64 	`json:"id" db:"id"`
	FromID    int64 	`json:"from_id" db:"from_id"`
	ToID      int64 	`json:"to_id" db:"to_id"`
	Amount    float64 	`json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// результат перевода
type TransactionResult struct {
	FromAcc *Account `json:"from"`
	ToAcc *Account 	`json:"to"`
	Transaction *Transaction `json:"transaction"`
}