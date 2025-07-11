package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) error
	GetAccount() ([]*Account, error)
	GetAccountByID(int) (*Account, error)

	DepositMoney(*DepositMoney) (*Account, error)
	TransferMoney(*TransferMoney) (string, error)
}

type PgStore struct {
	db *sql.DB
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

	sqlQueryToGetAccount := `select * from accounts where id = $1`
	userAccount := new(Account)

	row := pg.db.QueryRow(sqlQueryToGetAccount, depositMoneyRequest.ID)
	err := row.Scan(&userAccount.ID, &userAccount.FirstName, &userAccount.LastName, &userAccount.Email, &userAccount.Phone, &userAccount.Balance, &userAccount.CreatedAt)

	sqlQueryToDeposit := `UPDATE accounts
		SET balance = $1
		WHERE id = $2
		RETURNING id, first_name, last_name, email, phone, balance, created_at`

	row = pg.db.QueryRow(sqlQueryToDeposit,
		userAccount.Balance+depositMoneyRequest.Amount,
		userAccount.ID,
	)

	err = row.Scan(&userAccount.ID, &userAccount.FirstName, &userAccount.LastName, &userAccount.Email, &userAccount.Phone, &userAccount.Balance, &userAccount.CreatedAt)

	if err != nil {
		return nil, err
	}

	return userAccount, nil

}

func (pg *PgStore) TransferMoney(transferMoneyRequest *TransferMoney) (string, error) {

	senderAccount, err := pg.GetAccountByID(transferMoneyRequest.SenderId)

	if err != nil {
		return "", err
	}

	receiverAccount, err := pg.GetAccountByID(transferMoneyRequest.ReceiverId)

	if err != nil {
		return "", err
	}

	sqlQueryToDeduct := `UPDATE accounts
		SET balance = $1
		WHERE id = $2
		RETURNING id, first_name, last_name, email, phone, balance, created_at`

	row := pg.db.QueryRow(sqlQueryToDeduct,
		senderAccount.Balance-transferMoneyRequest.Amount,
		senderAccount.ID,
	)

	err = row.Scan(&senderAccount.ID, &senderAccount.FirstName, &senderAccount.LastName, &senderAccount.Email, &senderAccount.Phone, &senderAccount.Balance, &senderAccount.CreatedAt)

	if err != nil {
		return "", err
	}

	sqlQueryToAdd := `UPDATE accounts
		SET balance = $1
		WHERE id = $2
		RETURNING id, first_name, last_name, email, phone, balance, created_at`

	row = pg.db.QueryRow(sqlQueryToAdd,
		receiverAccount.Balance+transferMoneyRequest.Amount,
		receiverAccount.ID,
	)

	err = row.Scan(&receiverAccount.ID, &receiverAccount.FirstName, &receiverAccount.LastName, &receiverAccount.Email, &receiverAccount.Phone, &receiverAccount.Balance, &receiverAccount.CreatedAt)

	if err != nil {
		return "", err
	}

	return "Money trasfer has been successful", nil

}
