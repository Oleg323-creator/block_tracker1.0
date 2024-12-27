package usecases

import (
	"block_tracker1.0/internal/db"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"strings"
)

var erc20ABI = `[{"constant":false,"inputs":[{"name":"recipient","type":"address"},
{"name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],
"payable":false,"stateMutability":"nonpayable","type":"function"}]`

func (u *Usecase) GetAllTxInfoByBlock(lastBlockInDb int64) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	block, err := u.Client.BlockByNumber(context.Background(), big.NewInt(lastBlockInDb))
	if err != nil {
		log.Fatalf("Failed to fetch block: %v", err)
	}

	chainID, err := u.Client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("Error getting chain ID: %v", err)
	}

	tokenABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		log.Fatalf("Error loading ABI: %v", err)
	}

	// GETTING INFO ABOUT ALL TX IN BLOCK
	for _, tx := range block.Transactions() {

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
		if senderAddr.String() == os.Getenv("SENDER_ADDR") || toAddr == os.Getenv("CONTRACT_ADDRESS") {
			//CHECKING ADDRESS
			toAddress, amount, err := CheckTxData(tx, tokenABI)
			if err != nil {
				continue
			}

			data := db.TxData{
				Hash:         tx.Hash().Hex(),
				FromAddr:     senderAddr.String(),
				ContractAddr: toAddr,
				ToAddr:       toAddress.Hex(),
				Amount:       amount.String(),
				Value:        tx.Value().String(),
				BlockNumber:  block.Number().Int64(),
			}
			err = u.Repository.SaveTxDataToDB(data)
			if err != nil {
				log.Fatalf("Error saiving data in db: %v", err)
			}
			fmt.Printf("From: %s\n", data.FromAddr)
			fmt.Printf("Transaction Hash: %s\n", data.Hash)
			fmt.Printf("To: %s\n", data.ContractAddr)
			fmt.Printf("Recipient: %s\n", toAddress.Hex())
			fmt.Printf("Amount: %s tokens\n", amount.String())
			fmt.Printf("Value: %s\n", data.Value)
			fmt.Printf("Block number: %d", data.BlockNumber)
		}
	}
	return nil
}

func CheckTxData(tx *types.Transaction, tokenABI abi.ABI) (common.Address, *big.Int, error) {
	var emptyAddress common.Address
	info := tx.Data()
	if len(info) < 4 {
		var emptyAddress common.Address
		return emptyAddress, big.NewInt(0), fmt.Errorf("Got simple transaction")
	}

	method, err := tokenABI.MethodById(info[:4])
	if err != nil || method.Name != "transfer" {
		return emptyAddress, big.NewInt(0), fmt.Errorf("There is no transfer method")
	}

	args := make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(args, info[4:])
	if err != nil {
		return emptyAddress, big.NewInt(0), fmt.Errorf("Error unpacking transaction data: %v", err)
	}

	toAddress, ok1 := args["recipient"].(common.Address)
	amount, ok2 := args["amount"].(*big.Int)
	if !ok1 || !ok2 {
		return emptyAddress, big.NewInt(0), fmt.Errorf("Error extracting data from unpacked map")
	}

	return toAddress, amount, nil
}
