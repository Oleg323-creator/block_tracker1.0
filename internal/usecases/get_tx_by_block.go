package usecases

import (
	"block_tracker1.0/internal/db"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
)

func (u *Usecase) GetAllTxInfoByBlock(lastBlockInDb int64) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	block, err := u.Client.BlockByNumber(context.Background(), big.NewInt(lastBlockInDb))
	if err != nil {
		log.Fatalf("Failed to fetch block: %v", err)
	}

	// GETTING INFO ABOUT ALL TX IN BLOCK

	chainID, err := u.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Error getting chain ID: %v", err)
	}

	for _, tx := range block.Transactions() {

		//GETTING SENDER ADDRESS
		senderAddr, err := types.Sender(types.NewEIP155Signer(chainID), tx)
		if err != nil {
			senderAddr, err = types.Sender(types.HomesteadSigner{}, tx)
		}

		var toAddr string
		if tx.To() != nil {
			toAddr = tx.To().Hex()
		} else {
			toAddr = "Contract Creation"
		}

		//CHECKING ADDRESS
		if senderAddr.String() == os.Getenv("SENDER_ADDR") || toAddr == os.Getenv("CONTRACT_ADDRESS") {
			data := db.TxData{
				Hash:        tx.Hash().Hex(),
				FromAddr:    senderAddr.String(),
				ToAddr:      tx.To().Hex(),
				Value:       tx.Value().String(),
				BlockNumber: block.Number().Int64(),
			}
			err := u.Repository.SaveTxDataToDB(data)
			if err != nil {
				log.Fatalf("Error saiving data in db: %v", err)
			}
			fmt.Printf("From: %s\n", data.FromAddr)
			fmt.Printf("Transaction Hash: %s\n", data.Hash)
			fmt.Printf("To: %s\n", data.ToAddr)
			fmt.Printf("Value: %s\n", data.Value)
			fmt.Printf("Block number: %d", data.BlockNumber)
		}
	}
	return nil
}

/*
	signer := types.NewLondonSigner(chainID) // !!!!!!!!!!!!

	senderAddr, err := types.Sender(signer, tx)
	if err != nil {
		log.Fatalf("Error getting sennder address: %v", err)
	}

*/
