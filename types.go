package main

import (
	"time"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Balance   float32   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateAccountRequestBody struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func CreateNewAccount(firstName string, lastName string, email string, phone string) *Account {

	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Balance:   0.00,
		CreatedAt: time.Now().UTC(),
	}
}
