package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccountByID(int) (*Account, error)
}

type PgStore struct {
	db *sql.DB
}

func NewPgStore() (*PgStore, error) {

	connectionString := "user=bankple dbname=bankple-db password=bankple-hidden sslmode=disable"
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return &PgStore{
		db: db,
	}, nil

}

func (pg *PgStore) CreateTable(tableName string) error {

	sqlQuery := fmt.Sprintf(`
			create table if not exists %s (
				id serial primary key,
				first_name text,
				last_name text,
				email text,
				phone text,
				balance float,
				created_at timestamp
			)`, tableName)

	_, err := pg.db.Exec(sqlQuery)

	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("error creating table '%s': %w", tableName, err)
	}

	return nil

}

func (pg *PgStore) CreateAccount(account *Account) error {

	pg.CreateTable("accounts")

	sqlQuery := `
		insert into accounts (first_name, last_name, email, phone, balance, created_at)
		values ($1, $2, $3, $4, $5, $6)
	`

	_, err := pg.db.Exec(
		sqlQuery,
		account.FirstName,
		account.LastName,
		account.Email,
		account.Phone,
		account.Balance,
		account.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PgStore) UpdateAccount(account *Account) error {
	return nil
}

func (pg *PgStore) DeleteAccount(id int) error {
	return nil
}

func (pg *PgStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
