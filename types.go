package main

import "math/rand/v2"

type Account struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Id        int32  `json:"id"`
	Password  string `json:"password"`
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
