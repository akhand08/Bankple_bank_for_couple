package utils

import (
	"github.com/akhand08/Bankple_bank_for_couple/internal/models"
	"github.com/gorilla/mux"
)

func NewAPIServer(listenAddr string, router *mux.Router) *models.APIServer {
	return &models.APIServer{ListenAddress: listenAddr, Router: router}

}
