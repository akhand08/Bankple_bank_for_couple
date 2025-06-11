package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {

	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleGetAccountByID))
	router.HandleFunc("/money", makeHTTPHandlerFunc(s.handleMoney))
	router.HandleFunc("/transfer-money", makeHTTPHandlerFunc(s.HandleTransferMoney))
	log.Println("The server is listening on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)

}

func (s *APIServer) handleMoney(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "POST" {
		return s.HandleDepositMoney(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)

}

func (s *APIServer) HandleTransferMoney(w http.ResponseWriter, r *http.Request) error {

	transferMoneyRequest := new(TransferMoney)

	err := json.NewDecoder(r.Body).Decode(transferMoneyRequest)

	if err != nil {
		return err
	}

	response, err := s.store.TransferMoney(transferMoneyRequest)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, response)

}

func (s *APIServer) HandleDepositMoney(w http.ResponseWriter, r *http.Request) error {

	depositMoneyRequestBody := new(DepositMoney)

	err := json.NewDecoder(r.Body).Decode(depositMoneyRequestBody)

	if err != nil {
		return err
	}

	userAccount, err := s.store.DepositMoney(depositMoneyRequestBody)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, userAccount)

}
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	allAccounts, err := s.store.GetAccount()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, allAccounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	strID := mux.Vars(r)["id"]
	accountID, err := strconv.Atoi(strID)

	if err != nil {
		return err
	}

	userAccount, err := s.store.GetAccountByID(accountID)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, userAccount)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	accountRequestBody := new(CreateAccountRequestBody)

	err := json.NewDecoder(r.Body).Decode(accountRequestBody)
	if err != nil {
		return err
	}

	newAccount := CreateNewAccount(accountRequestBody.FirstName, accountRequestBody.LastName, accountRequestBody.Email, accountRequestBody.Phone)

	err = s.store.CreateAccount(newAccount)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, *newAccount)

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	strID := mux.Vars(r)["id"]
	accountID, err := strconv.Atoi(strID)

	if err != nil {
		return err
	}

	err = s.store.DeleteAccount(accountID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accountID)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {

	updateAccount := new(Account)

	err := json.NewDecoder(r.Body).Decode(updateAccount)

	if err != nil {
		return nil
	}

	err = s.store.UpdateAccount(updateAccount)

	if err != nil {
		fmt.Println(err)
		return err

	}

	return WriteJSON(w, http.StatusOK, *updateAccount)

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application-json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}
