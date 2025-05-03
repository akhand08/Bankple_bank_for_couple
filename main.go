package main

import (
	"fmt"
	"log"
)

func main() {

	dbStore, err := NewPgStore()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", dbStore)

	// server := NewAPIServer(":3000", dbStore)
	// server.Run()

	fmt.Println("Welcome to Bankple")

}
