package domain

import (
	"time"
)

type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Transaction struct {
	ID          int64     `json:"id"`
	WalletID    int64     `json:"wallet_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`   // deposit, withdraw, bet, win
	Status      string    `json:"status"` // pending, completed, failed
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}
