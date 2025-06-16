package models

import (
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddress string
	Router        *mux.Router
}

func (server APIServer) Run() {

	http.ListenAndServe(server.ListenAddress, server.Router)

}
