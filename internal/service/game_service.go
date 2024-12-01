package service

import (
	"casino-service/internal/domain"
	"casino-service/internal/repository"
	"casino-service/pkg/utils"
	"context"
	"errors"
	"fmt"
	"log"
)

type GameService interface {
	PlaceBet(ctx context.Context, userID int64, gameID int64, amount float64) (*domain.GameSession, error)
	ProcessGameResult(ctx context.Context, sessionID int64) (*domain.GameSession, error)
	GetBalance(ctx context.Context, userID int64) (*domain.Wallet, error)
}

type gameService struct {
	walletRepo repository.WalletRepository
	random     utils.RandomGenerator
}

func NewGameService(walletRepo repository.WalletRepository) GameService {
	return &gameService{
		walletRepo: walletRepo,
		random:     utils.NewRandomGenerator(),
	}
}

func (s *gameService) PlaceBet(ctx context.Context, userID int64, gameID int64, amount float64) (*domain.GameSession, error) {
	log.Printf("Placing bet: userID=%d, gameID=%d, amount=%.2f", userID, gameID, amount)

	// Проверяем баланс пользователя
	wallet, err := s.walletRepo.GetBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	if wallet.Balance < amount {
		return nil, errors.New("insufficient funds")
	}

	// Создаем игровую сессию
	session := &domain.GameSession{
		UserID:    userID,
		GameID:    gameID,
		BetAmount: amount,
		Status:    "active",
	}

	// Списываем ставку
	err = s.walletRepo.UpdateBalance(ctx, userID, -amount)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	return session, nil
}

func (s *gameService) ProcessGameResult(ctx context.Context, sessionID int64) (*domain.GameSession, error) {
	// Создаем новую сессию с переданным ID
	session := &domain.GameSession{
		ID:     sessionID,
		Status: "active",
	}

	// Генерируем результат игры
	result, winAmount, err := s.generateSlotResult(session.BetAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to generate result: %w", err)
	}

	session.Result = result
	session.WinAmount = winAmount
	session.Status = "completed"

	// Если есть выигрыш, начисляем его на баланс
	if winAmount > 0 {
		err = s.walletRepo.UpdateBalance(ctx, session.UserID, winAmount)
		if err != nil {
			session.Status = "failed"
			return nil, fmt.Errorf("failed to update balance: %w", err)
		}
	}

	return session, nil
}

func (s *gameService) generateSlotResult(betAmount float64) (string, float64, error) {
	num, err := s.random.GenerateNumber(1, 100)
	if err != nil {
		return "", 0, fmt.Errorf("failed to generate random number: %w", err)
	}

	switch {
	case num <= 40: // 40% шанс маленького выигрыша
		return "small_win", betAmount * 1.5, nil
	case num <= 70: // 30% шанс среднего выигрыша
		return "medium_win", betAmount * 2, nil
	case num <= 85: // 15% шанс большого выигрыша
		return "big_win", betAmount * 5, nil
	case num <= 87: // 2% шанс джекпота
		return "jackpot", betAmount * 50, nil
	default: // 13% шанс проигрыша
		return "no_win", 0, nil
	}
}

func (s *gameService) GetBalance(ctx context.Context, userID int64) (*domain.Wallet, error) {
	return s.walletRepo.GetBalance(ctx, userID)
}
