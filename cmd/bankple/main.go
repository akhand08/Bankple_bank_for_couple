package main

import (
	"fmt"
	"log"

	"github.com/akhand08/Bankple_bank_for_couple/internal/db"
	"github.com/akhand08/Bankple_bank_for_couple/internal/router"

	"github.com/akhand08/Bankple_bank_for_couple/pkg/utils"
)

func main() {

	db, err := db.NewPgStore()
	fmt.Printf("%+v\n", db)
	router := router.NewRouter(db)

	if err != nil {
		log.Fatal(err)
	}

	server := utils.NewAPIServer(":3000", router)
	server.Run()

}
