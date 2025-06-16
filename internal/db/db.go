package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DatabaseStore interface {
	CreateAccount()
}

type PgStore struct {
	Db *sql.DB
}

func NewPgStore() (*PgStore, error) {

	connectionString := "user=bankple dbname=bankple-db password=bankple-hidden sslmode=disable"
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return &PgStore{
		Db: db,
	}, nil

}

func (pg *PgStore) CreateAccount() {

}
