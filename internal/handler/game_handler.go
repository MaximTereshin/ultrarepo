package handler

import (
	"casino-service/internal/domain"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type GameService interface {
	PlaceBet(ctx context.Context, userID int64, gameID int64, amount float64) (*domain.GameSession, error)
	ProcessGameResult(ctx context.Context, sessionID int64) (*domain.GameSession, error)
	GetBalance(ctx context.Context, userID int64) (*domain.Wallet, error)
}

type GameHandler struct {
	gameService GameService
}

func NewGameHandler(gameService GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

type PlaceBetRequest struct {
	GameID int64   `json:"game_id"`
	Amount float64 `json:"amount"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path == "/bet" {
			h.PlaceBet(w, r)
		} else if r.URL.Path == "/result" {
			h.ProcessGameResult(w, r)
		}
	case http.MethodGet:
		if r.URL.Path == "/balance" {
			h.GetBalance(w, r)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GameHandler) PlaceBet(w http.ResponseWriter, r *http.Request) {
	var req PlaceBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{"invalid request format"})
		return
	}

	if req.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{"amount must be greater than 0"})
		return
	}

	// В реальном приложении userID должен браться из JWT токена
	userID := int64(1) // Временное решение

	session, err := h.gameService.PlaceBet(r.Context(), userID, req.GameID, req.Amount)
	if err != nil {
		switch err.Error() {
		case "insufficient funds":
			writeJSON(w, http.StatusBadRequest, ErrorResponse{"insufficient funds"})
		case "game is not active":
			writeJSON(w, http.StatusBadRequest, ErrorResponse{"game is not active"})
		default:
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{"internal server error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, session)
}

func (h *GameHandler) ProcessGameResult(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{"session id is required"})
		return
	}

	id, err := strconv.ParseInt(sessionID, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{"invalid session id"})
		return
	}

	result, err := h.gameService.ProcessGameResult(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *GameHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// В реальном приложении userID должен браться из JWT токена
	userID := int64(1) // Временное решение

	wallet, err := h.gameService.GetBalance(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, wallet)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
