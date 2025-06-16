package router

import (
	"github.com/akhand08/Bankple_bank_for_couple/internal/db"
	"github.com/gorilla/mux"
)

func NewRouter(db db.DatabaseStore) *mux.Router {

	router := mux.NewRouter()

	return router

}
