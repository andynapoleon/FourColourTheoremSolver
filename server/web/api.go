package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/crypto/bcrypt"
)

/*
ROUTES AND HANDLERS
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

	// Routes
	router.HandleFunc("/api/signup", makeHTTPHandlerFunc(s.handleSignUp))
	router.HandleFunc("/api/login", makeHTTPHandlerFunc(s.handleLogin))
	router.HandleFunc("/api/user", makeHTTPHandlerFunc(s.handleGetUser))

	// Create a CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(router)

	// Serve the port
	log.Println("JSON API Server running on port", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, handler); err != nil {
		log.Fatalf("Error starting server: %v", err) // Log and exit on error
	}
}

// Sign up
func (s *APIServer) handleSignUp(w http.ResponseWriter, r *http.Request) error {
	// Check if method is not allowed
	if r.Method != http.MethodPost {
		return writeJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "Method not allowed"})
	}

	// Decode the request body
	createUserReq := new(CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError, ApiError{Error: "Password too long, 60 characters maximum"})
	}

	// Create a new user in the database with hashed password
	user := NewUser(createUserReq.Name, createUserReq.Email, string(hashedPassword))
	if err := s.store.CreateUser(user); err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return writeJSON(w, http.StatusBadRequest, ApiError{Error: "Email already in use"})
		}
		if err.Error() == "name too long: maximum length is 60 characters" {
			return writeJSON(w, http.StatusBadRequest, ApiError{Error: "Name too long, 60 characters maximum"})
		} else if err.Error() == "email too long: maximum length is 60 characters" {
			return writeJSON(w, http.StatusBadRequest, ApiError{Error: "Email too long, 60 characters maximum"})
		} else if err.Error() == "password too long: maximum length is 60 characters" {
			return writeJSON(w, http.StatusBadRequest, ApiError{Error: "Password too long, 60 characters maximum"})
		}
	}

	return writeJSON(w, http.StatusOK, Message{Message: "New user created"})
}

// Login
func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	// Check if method is not allowed
	if r.Method != http.MethodPost {
		return writeJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "Method not allowed"})
	}

	// Decode the request body
	u := new(LoginUserRequest)
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		return err
	}

	// Grab user from the database
	dbUser, err := s.store.GetUserByEmail(u.Email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return writeJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid email"})
		}
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(u.Password))
	if err == nil {
		tokenString, err := CreateToken(u.Email)
		if err != nil {
			return writeJSON(w, http.StatusInternalServerError, ApiError{Error: "Error creating token"})
		}
		return writeJSON(w, http.StatusOK, LoginMessage{Token: tokenString, Name: dbUser.Name})
	} else {
		return writeJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid password"})
	}
}

// Get all users
func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) error {
	// Check if token is valid
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return writeJSON(w, http.StatusUnauthorized, ApiError{Error: "Missing authorization header"})
	}
	tokenString = tokenString[len("Bearer "):]
	err := VerifyToken(tokenString)
	if err != nil {
		return writeJSON(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
	}

	// Check if method is not allowed
	if r.Method != http.MethodGet {
		return writeJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "Method not allowed"})
	}

	users, err := s.store.GetUsers()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, &users)
}

/*
HELPER STRUCTS AND FUNCTIONS
*/
type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

type Message struct {
	Message string `json:"message"`
}

type LoginMessage struct {
	Token string `json:"token"`
	Name  string `json:"name"`
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
