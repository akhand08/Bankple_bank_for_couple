package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type DatabaseStore interface {
	CreateAccount()
}

type PgStore struct {
	Db *sql.DB
}

func NewPgStore() (*PgStore, error) {
	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		return nil, fmt.Errorf("environment variable DB_USER not set")
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, fmt.Errorf("environment variable DB_NAME not set")
	}

	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("environment variable DB_PASSWORD not set")
	}

	connectionString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable",
		dbUser, dbName, dbPassword)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Successfully connected to the database!")

	return &PgStore{
		Db: db,
	}, nil
}

func (pg *PgStore) CreateAccount() {
	fmt.Println("CreateAccount method called (implementation pending).")
}
