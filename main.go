package main

import (
	"block_tracker1.0/internal/db"
	"block_tracker1.0/internal/usecases"
	"context"
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

	dbConn := db.NewDB(cfg)

	url := fmt.Sprintf("wss://sepolia.infura.io/ws/v3/%s", os.Getenv("API_KEY"))
	repo := db.NewRepository(dbConn)

	usecase := usecases.NewUsecase(url, repo)
	defer usecase.Client.Close()

	// RUN MIGRATIONS
	err = db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	for {
		lastReleasedBlock, err := usecase.Client.BlockByNumber(context.Background(), nil)
		if err != nil {
			log.Fatalf("Failed to fetch block: %v", err)
		}
		log.Printf("Last released block: %d", lastReleasedBlock.Number().Int64())

		lastBlockInDB, err := repo.GetLastBlockFromDB()
		if err != nil {
			log.Printf("Error: %v", err)
			err := repo.SaveLastBlockToDB(lastReleasedBlock.Number().Int64())
			if err != nil {
				log.Fatalf("Error setting last block into DB: %v", err)
			}
			lastBlockInDB = lastReleasedBlock.Number().Int64()
		}
		log.Printf("Last saved block in DB: %d", lastBlockInDB)

		for i := lastBlockInDB; i < lastReleasedBlock.Number().Int64(); i++ {
			err = repo.SaveLastBlockToDB(i)
			if err != nil {
				return
			}

			err = usecase.GetAllTxInfoByBlock(lastBlockInDB)
			if err != nil {
				log.Fatalf("Error getting tx data: %v", err)
			}
			log.Printf("All tx had been chacked in block with number: %d", i)
		}
		log.Println("4")

		lastBlockInDB = lastReleasedBlock.Number().Int64()

	}
}

/*
package main

import (
	"block_tracker1.0/internal/db"
	"block_tracker1.0/internal/usecases"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

	// DB CONNECT
	dbConn := db.NewDB(cfg)

	// RUN MIGRATIONS
	err = db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	//MIGRATIONS CONNECT
	url := fmt.Sprintf("wss://sepolia.infura.io/ws/v3/%s", os.Getenv("API_KEY"))
	repo := db.NewRepository(dbConn, url)

	for {
		lastReleasedBlock, err := usecases.GetLastReleasedBLock(url)
		if err != nil {
			log.Fatal(err)
		}

		lastBlockInDB, err := repo.GetLastBlockFromDB()
		if err != nil {
			log.Fatal(err)
		}

		for i := lastBlockInDB; i < lastReleasedBlock; i++ {
			for _, tx := range block.Transactions() {
				chainID, err := client.NetworkID(context.Background())
				if err != nil {
					continue
				}
				from, err := types.Sender(types.NewLondonSigner(chainID), tx)
				if err != nil {
					continue
				}

				if from.String() == "0x886577048713f65d6e26e61e82597A523887645B" {
					fmt.Printf("From: %s\n", from.Hex())
					fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
					fmt.Printf("To: %s\n", tx.To().Hex())
					fmt.Printf("Value: %s\n", tx.Value().String())
				}
			}
		}
		lastBlockInDB = lastReleasedBlock
		lastReleasedBlock = block.Number().Int64()
	}
}

*/

/*
err := repo.SaveLastBlockToDB(header.Number.Int64())
if err != nil {
log.Printf("Failed to save block number %d: %v\n", header.Number.Int64(), err)
continue*/
/*
for {
	for i := lastBlockInDB; i < lastReleasedBock; i++{
		for _, tx := range block.Transactions() {
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			continue
		}
		from, err := types.Sender(types.NewLondonSigner(chainID), tx)
		if err != nil {
			continue
		}

		if from.String() == "0x886577048713f65d6e26e61e82597A523887645B" {
			fmt.Printf("From: %s\n", from.Hex())
			fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
			fmt.Printf("To: %s\n", tx.To().Hex())
			fmt.Printf("Value: %s\n", tx.Value().String())
		}
	}
}
*/
