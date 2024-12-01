package repository

import (
	"casino-service/internal/domain"
	"context"
	"database/sql"
	"errors"
)

type WalletRepository interface {
	GetBalance(ctx context.Context, userID int64) (*domain.Wallet, error)
	UpdateBalance(ctx context.Context, walletID int64, amount float64) error
	CreateTransaction(ctx context.Context, tx *domain.Transaction) error
	GetTransactionByID(ctx context.Context, txID int64) (*domain.Transaction, error)
}

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) GetBalance(ctx context.Context, userID int64) (*domain.Wallet, error) {
	wallet := &domain.Wallet{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, balance, currency, created_at, updated_at FROM wallets WHERE user_id = $1",
		userID,
	).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt, &wallet.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (r *walletRepository) UpdateBalance(ctx context.Context, walletID int64, amount float64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	result, err := tx.ExecContext(ctx,
		"UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND balance + $1 >= 0",
		amount, walletID)

	if err != nil {
		tx.Rollback()
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rows == 0 {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	return tx.Commit()
}

func (r *walletRepository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO transactions (wallet_id, amount, type, status, created_at) VALUES ($1, $2, $3, $4, NOW())",
		tx.WalletID, tx.Amount, tx.Type, tx.Status)
	return err
}

func (r *walletRepository) GetTransactionByID(ctx context.Context, txID int64) (*domain.Transaction, error) {
	tx := &domain.Transaction{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, wallet_id, amount, type, status, created_at FROM transactions WHERE id = $1",
		txID).Scan(&tx.ID, &tx.WalletID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
