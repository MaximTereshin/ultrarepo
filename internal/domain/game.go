package domain

type Game struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	MinBet      float64 `json:"min_bet"`
	MaxBet      float64 `json:"max_bet"`
	IsActive    bool    `json:"is_active"`
}

type GameSession struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"user_id"`
	GameID    int64   `json:"game_id"`
	BetAmount float64 `json:"bet_amount"`
	WinAmount float64 `json:"win_amount"`
	Status    string  `json:"status"` // active, completed, failed
	Result    string  `json:"result"` // выигрышная комбинация или результат
}
