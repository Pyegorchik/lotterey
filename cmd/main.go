package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"lottery/internal/handler"
	"lottery/internal/repository"
	"lottery/internal/service"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	db *sql.DB
}

func main() {

	ctx := context.Background()

	// Инициализация базы данных
	db, err := initDB(ctx)
	if err != nil {
		log.Fatal("DB init failed:", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("Database initialized successfully!")

	repo := repository.NewRepository(db)
	lotteryService := service.NewLotteryService(repo)
	lotteryHandler := handler.NewLotteryHandler(lotteryService)

	// Настройка маршрутов
	router := mux.NewRouter()

	// API эндпоинты
	router.HandleFunc("/draws", lotteryHandler.CreateDraw).Methods("POST")
	router.HandleFunc("/tickets", lotteryHandler.CreateTicket).Methods("POST")
	router.HandleFunc("/draws/{draw_id}/close", lotteryHandler.CloseDraw).Methods("POST")
	router.HandleFunc("/draws/{draw_id}/results", lotteryHandler.GetResults).Methods("GET")

	// Дополнительные эндпоинты
	router.HandleFunc("/draws", lotteryHandler.GetDraws).Methods("GET")
	router.HandleFunc("/draws/{draw_id}", lotteryHandler.GetDraw).Methods("GET")

	fmt.Println("Lottery API server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initDB(ctx context.Context) (*sql.DB, error) {
	var err error
	db, err := sql.Open("sqlite3", "lottery.db")
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile("cmd/sql/init.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read SQL file: %w", err)
	}

	queries := strings.Split(string(content), ";")

	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}

		if _, err := db.ExecContext(ctx, q); err != nil {
			return nil, fmt.Errorf("query failed: %s\nerror: %w", q, err)
		}
	}

	return db, nil
}
