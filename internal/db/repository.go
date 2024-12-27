package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"log"
)

// TO CONNECT IT WITH DB
type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) SaveLastBlockToDB(block int64) error {

	queryBuilder := squirrel.Update("blocks").
		Set("block_number", block).
		Where(squirrel.Eq{"id": 1})

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, execErr := r.DB.ExecContext(context.Background(), query, args...)
	if execErr != nil {
		return fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	log.Println("Block number saved to db:", block)

	return nil
}

func (r *Repository) GetLastBlockFromDB() (int64, error) {
	queryBuilder := squirrel.Select("block_number").
		From("blocks").
		Where(squirrel.Eq{"id": 1})

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build SQL query: %v", err)
	}

	rows, execErr := r.DB.Query(query, args...)
	if execErr != nil {
		return 0, fmt.Errorf("failed to execute SQL query: %v", execErr)
	}

	var lastBlock int64
	if rows.Next() {
		if err = rows.Scan(&lastBlock); err != nil {
			return 0, fmt.Errorf("failed to scan result: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return lastBlock, nil
}

type TxData struct {
	Hash         string
	FromAddr     string
	ContractAddr string
	ToAddr       string
	Amount       string
	Value        string
	BlockNumber  int64
}

func (r *Repository) SaveTxDataToDB(data TxData) error {
	queryBuilder := squirrel.Insert("transactions").
		Columns("hash", "from_addr", "contract_addr", "to_addr", "amount", "value", "block_number").
		Values(data.Hash, data.FromAddr, data.ContractAddr, data.ToAddr, data.Amount, data.Value, data.BlockNumber).
		Suffix("ON CONFLICT (hash) DO NOTHING")

	query, args, err := queryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, execErr := r.DB.ExecContext(context.Background(), query, args...)
	if execErr != nil {
		return fmt.Errorf("failed to execute SQL query: %v", execErr)
	}
	log.Println("Tx data saved to DB")

	return nil
}

/*
func (r *Repository) SetLastBlockInDB(block int64) error {
	updateQueryBuilder := squirrel.Update("blocks").
		Set("block_number", block).
		Where(squirrel.Eq{"id": "1"})

	updateQuery, updateArgs, err := updateQueryBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("Failed to build SQL query: %v", err)
	}

	_, err = r.DB.Exec(updateQuery, updateArgs...)
	if err != nil {
		return fmt.Errorf("Failed to execute query: %v", err)
	}
	return nil
}
*/
/*
		if rows.Next() { // CHECKING IF THERE IS RESULT
			if err = rows.Scan(&lastBlock); err != nil {
				return 0, fmt.Errorf("failed to scan result: %v", err)
			}

			if err = rows.Err(); err != nil {
				log.Fatal(err)
			}

			return lastBlock, nil
		} else {
			lastReleasedBlock, err := usecases.GetLastReleasedBLock(r.URL)
			if err != nil {
				return 0, fmt.Errorf("Failed to fetch the latest block: %v", err)
			}

			err = r.SetLastBlockInDB(lastReleasedBlock) //ADDING LAST RELEASED BLOCK NUMBER AS DEFAULT VALUE IF NO VALUE IN DB
			if err != nil {
				return 0, fmt.Errorf("Failed to set the latest block: %v", err)
			}

			return lastReleasedBlock, nil
		}
	}
*/
