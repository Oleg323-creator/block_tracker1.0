package main

import (
	"block_tracker1.0/internal/db"
	"block_tracker1.0/internal/usecases"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
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

	dbConn := db.NewDB(cfg)
	defer dbConn.Close()

	url := fmt.Sprintf("wss://sepolia.infura.io/ws/v3/%s", os.Getenv("API_KEY"))
	repo := db.NewRepository(dbConn)

	usecase := usecases.NewUsecase(url, repo)
	defer usecase.Client.Close()

	// RUN MIGRATIONS
	err = db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	lastReleasedBlock, err := usecase.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Printf("Failed to fetch the latest block: %v", err)
	}
	log.Printf("Last released block: %d", lastReleasedBlock.Number().Int64())

	lastBlockInDB, err := repo.GetLastBlockFromDB()
	if err != nil {
		log.Printf("Error getting last block from DB: %v", err)
		err := repo.SaveLastBlockToDB(lastReleasedBlock.Number().Int64())
		if err != nil {
			log.Fatalf("Error setting last block into DB: %v", err)
		}
		lastBlockInDB = lastReleasedBlock.Number().Int64()
	}
	log.Printf("Last saved block in DB: %d", lastBlockInDB)

	for {
		for i := lastBlockInDB; i <= lastReleasedBlock.Number().Int64(); i++ {

			lastReleasedBlock, err = usecase.Client.BlockByNumber(context.Background(), nil)
			if err != nil {
				log.Printf("Failed to fetch the latest block: %v", err)
			}
			if lastBlockInDB == lastReleasedBlock.Number().Int64() {
				time.Sleep(3 * time.Second)
			}
			log.Printf("Last released block: %d", lastReleasedBlock.Number().Int64())

			err = usecase.GetAllTxInfoByBlock(i)
			if err != nil {
				log.Fatalf("Error getting tx data: %v", err)
			}
			log.Printf("All tx had been chacked in block with number: %d", i)

			err = repo.SaveLastBlockToDB(i)
			if err != nil {
				log.Printf("Error processing block %d: %v", i, err)
				continue
			}
		}
		log.Println("4")

		lastBlockInDB = lastReleasedBlock.Number().Int64()
		time.Sleep(2 * time.Second)

	}
}
