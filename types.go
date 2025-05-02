package main

import "math/rand/v2"

type Account struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Id        int32
	Password  string
}

func CreateNewAccount() *Account {

	return &Account{
		FirstName: "President",
		LastName:  "XYZ",
		Email:     "president.xyz@email.com",
		Phone:     "0123456789",
		Id:        rand.Int32(),
		Password:  "hidden",
	}
}
