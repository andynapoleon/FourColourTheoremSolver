package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listenAddr string
}

type ApiError struct {
	Error string
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	// v is data to be encoded into json
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

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

func NewAPISever(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPHandlerFunc(s.handleUser))

	router.HandleFunc("/user/{id}", makeHTTPHandlerFunc(s.handleGetUser))

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
	user := NewUser("Andy Tran", "aqtran@ualberta.ca", "123456")
	return writeJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"] // get all vars sent with the request
	fmt.Println(id)
	return writeJSON(w, http.StatusOK, &User{})
}
