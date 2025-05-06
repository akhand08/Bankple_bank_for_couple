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
	GetAccount() ([]*Account, error)
	GetAccountByID(int) (*Account, error)

	DepositMoney(*DepositMoney) (*Account, error)
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
		returning id
	`
	var userId int
	err := pg.db.QueryRow(
		sqlQuery,
		account.FirstName,
		account.LastName,
		account.Email,
		account.Phone,
		account.Balance,
		account.CreatedAt,
	).Scan(&userId)

	if err != nil {
		return err
	}

	account.ID = userId

	return nil
}

func (pg *PgStore) UpdateAccount(account *Account) error {

	sqlQuery := `UPDATE accounts
	SET first_name = $1,
		last_name = $2,
		email = $3,
		phone = $4,
		balance = $5
	WHERE id = $6
	RETURNING id, first_name, last_name, email, phone, balance, created_at`

	row := pg.db.QueryRow(sqlQuery,
		account.FirstName,
		account.LastName,
		account.Email,
		account.Phone,
		account.Balance,
		account.ID,
	)

	err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Email, &account.Phone, &account.Balance, &account.CreatedAt)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil

}

func (pg *PgStore) DeleteAccount(id int) error {

	sqlQuery := `delete from accounts where id = $1`

	_, err := pg.db.Exec(sqlQuery, id)

	if err != nil {
		return err
	}

	return nil
}

func (pg *PgStore) GetAccount() ([]*Account, error) {

	sqlQuery := `select * from accounts`

	rows, err := pg.db.Query(sqlQuery)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var allAccounts []*Account
	for rows.Next() {

		account := new(Account)

		rows.Scan(&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Email,
			&account.Phone,
			&account.Balance,
			&account.CreatedAt)
		allAccounts = append(allAccounts, account)

	}

	return allAccounts, nil

}

func (pg *PgStore) GetAccountByID(id int) (*Account, error) {

	sqlQuery := `select * from accounts where id = $1`
	userAccount := new(Account)

	row := pg.db.QueryRow(sqlQuery, id)

	err := row.Scan(&userAccount.ID, &userAccount.FirstName, &userAccount.LastName, &userAccount.Email, &userAccount.Phone, &userAccount.Balance, &userAccount.CreatedAt)

	if err != nil {
		return nil, err
	}

	return userAccount, nil

}

func (pg *PgStore) DepositMoney(depositMoneyRequest *DepositMoney) (*Account, error) {
	return nil, nil
}
