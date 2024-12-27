package runners

import (
	"block_tracker1.0/internal/usecases"
	"context"
	"fmt"
	"log"
	"time"
)

type Runner struct {
	Ucase *usecases.Usecase
}

func NewRunner(usecase *usecases.Usecase) *Runner {
	return &Runner{usecase}
}

func (r *Runner) ListenBlockchain() error {
	lastReleasedBlock, err := r.Ucase.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Printf("Failed to fetch the latest block: %v", err)
	}
	log.Printf("Last released block: %d", lastReleasedBlock.Number().Int64())

	lastBlockInDB, err := r.Ucase.Repository.GetLastBlockFromDB()
	if err != nil {
		log.Printf("Error getting last block from DB: %v", err)
		err := r.Ucase.Repository.SaveLastBlockToDB(lastReleasedBlock.Number().Int64())
		if err != nil {
			log.Fatalf("Error setting last block into DB: %v", err)
		}
		lastBlockInDB = lastReleasedBlock.Number().Int64()
	}
	log.Printf("Last saved block in DB: %d", lastBlockInDB)

	for {
		for i := lastBlockInDB; i <= lastReleasedBlock.Number().Int64(); i++ {

			lastReleasedBlock, err = r.Ucase.Client.BlockByNumber(context.Background(), nil)
			if err != nil {
				fmt.Errorf("Failed to fetch the latest block: %v", err)
			}

			if lastBlockInDB == lastReleasedBlock.Number().Int64() {
				time.Sleep(3 * time.Second)
			}

			log.Printf("Last released block: %d", lastReleasedBlock.Number().Int64())

			err = r.Ucase.GetAllTxInfoByBlock(i)
			if err != nil {
				fmt.Errorf("Error getting tx data: %v", err)
			}
			log.Printf("All tx had been chacked in block with number: %d", i)

			err = r.Ucase.Repository.SaveLastBlockToDB(i)
			if err != nil {
				log.Printf("Error processing block %d: %v", i, err)
				continue
			}
		}
		log.Println("4")

		lastBlockInDB = lastReleasedBlock.Number().Int64()
		time.Sleep(2 * time.Second)
	}
	return nil
}
