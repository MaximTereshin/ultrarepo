package main

import (
	"casino-service/internal/handler"
	"casino-service/internal/repository"
	"casino-service/internal/service"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Получаем параметры подключения к БД из переменных окружения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://casino_user:casino_password@localhost:5433/casino_db?sslmode=disable"
	}

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Инициализируем репозитории
	walletRepo := repository.NewWalletRepository(db)

	// Инициализируем сервисы
	gameService := service.NewGameService(walletRepo)

	// Инициализируем хендлеры
	gameHandler := handler.NewGameHandler(gameService)

	// Создаем новый мультиплексор
	mux := http.NewServeMux()

	// Регистрируем обработчики с конкретными методами
	mux.HandleFunc("/bet", gameHandler.PlaceBet)
	mux.HandleFunc("/result", gameHandler.ProcessGameResult)
	mux.HandleFunc("/balance", gameHandler.GetBalance)

	// Запускаем HTTP сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
