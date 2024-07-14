package main

import (
	"encoding/json"
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
}

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	return nil
}
