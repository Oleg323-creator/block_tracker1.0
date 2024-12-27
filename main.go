package main

import (
	"block_tracker1.0/internal/db"
	"block_tracker1.0/internal/runners"
	"block_tracker1.0/internal/usecases"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := db.ConnectionConfig{
		Host:     "localhost",
		Port:     "5431",
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	// DB -> REPOSITORY -> USECASES -> RUNNERS

	dbConn := db.NewDB(cfg)
	defer dbConn.Close()

	err = db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	url := fmt.Sprintf("wss://sepolia.infura.io/ws/v3/%s", os.Getenv("API_KEY"))
	repo := db.NewRepository(dbConn)

	usecase := usecases.NewUsecase(url, repo)
	defer usecase.Client.Close()

	runner := runners.NewRunner(usecase)

	err = runner.ListenBlockchain()
	if err != nil {
		log.Fatal(err)
	}
}
