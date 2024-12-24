package usecases

import (
	"block_tracker1.0/internal/db"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type Usecase struct {
	URL        string
	Repository *db.Repository
	Client     *ethclient.Client
}

func NewUsecase(url string, rep *db.Repository) *Usecase {
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}

	return &Usecase{
		URL:        url,
		Repository: rep,
		Client:     client,
	}
}
