package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB CONFIG
type ConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// CONNECTION TO DB
func NewDB(cfg ConnectionConfig) *sql.DB {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.DBName, cfg.SSLMode)
	log.Printf("Connecting to the database with connection string: %s", connString)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Error opening database connection:", err)

	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database:", err)
	}

	log.Println("Successfully connected to the database")
	return db
}
