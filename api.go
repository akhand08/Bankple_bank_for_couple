package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleAccount))
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

	// if r.Method == "DELETE" {
	// 	return s.handleDeleteAccount(w, r)
	// }

	return fmt.Errorf("Method not allowed %s", r.Method)

}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	// vars := mux.Vars(r)
	// fmt.Println(vars["id"])
	// userAccount := CreateNewAccount()

	// return WriteJSON(w, http.StatusOK, userAccount)
	return nil
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

	WriteJSON(w, http.StatusOK, *newAccount)

	return nil
}

// func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application-json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}
