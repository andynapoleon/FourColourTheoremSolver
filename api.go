package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
MAIN FUNCTIONS
*/

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPISever(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPHandlerFunc(s.handleUser))

	router.HandleFunc("/user/{id}", makeHTTPHandlerFunc(s.handleGetUserById))

	// serve the port
	log.Println("JSON API Server running on port", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, router); err != nil {
		log.Fatalf("Error starting server: %v", err) // Log and exit on error
	}
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUser(w, r)
	} else if r.Method == "POST" {
		return s.handleCreateUser(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	// Decode the request body
	createUserReq := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	// Create a new user in the database
	user := NewUser(createUserReq.Name, createUserReq.Email, createUserReq.Password)
	if err := s.store.CreateUser(user); err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, "New user created!")
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, &users)
}

func (s *APIServer) handleGetUserById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"] // get all vars sent with the request
	fmt.Println(id)
	return writeJSON(w, http.StatusOK, &User{})
}

/*
HELPER STRUCTS AND FUNCTIONS
*/
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

// Sends response in JSON format
func writeJSON(w http.ResponseWriter, status int, v any) error {
	// v is data to be encoded into json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// Converts apiFunc to handlerFunc
func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error
			writeJSON(w, http.StatusBadRequest, ApiError{
				Error: err.Error(),
			})
		}
	}
}
