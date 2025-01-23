// app.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type App struct {
	db *sql.DB
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type RefreshResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

func (app *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Print out the request body
	fmt.Println("REQUEST BODY: ", req.Email, req.Password, req.Name)

	// Validate input
	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Email and password are required",
		})
		return
	}

	// Check if database connection is alive
	if err := app.db.Ping(); err != nil {
		log.Printf("Database connection error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Database connection error",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error processing password",
		})
		return
	}

	// Insert user
	var userID int
	err = app.db.QueryRow(
		`INSERT INTO users (email, password_hash, name) 
         VALUES ($1, $2, $3) 
         RETURNING id`,
		req.Email,
		string(hashedPassword),
		req.Name,
	).Scan(&userID)

	if err != nil {
		log.Printf("Database error during user insertion: %v", err)

		if strings.Contains(err.Error(), "unique constraint") ||
			strings.Contains(err.Error(), "duplicate key") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "User with this email already exists",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
		"userId":  strconv.Itoa(userID),
	})
}

func (app *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user User
	err := app.db.QueryRow(
		"SELECT id, email, password_hash FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, expiresAt, err := generateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Store session
	_, err = app.db.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		user.ID,
		token,
		expiresAt,
	)

	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	})
}

func (app *App) handleVerifyToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	// Verify token
	claims, err := verifyToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if token exists in sessions
	var exists bool
	err = app.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM sessions WHERE token = $1 AND expires_at > NOW())",
		token,
	).Scan(&exists)

	if err != nil || !exists {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"user_id": claims.UserID,
	})
}

func (app *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	// Delete session
	_, err := app.db.Exec("DELETE FROM sessions WHERE token = $1", token)
	if err != nil {
		http.Error(w, "Error processing logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

func (app *App) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	oldToken := r.Header.Get("Authorization")
	if oldToken == "" {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	// Verify old token
	claims, err := verifyToken(oldToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check if old token exists in sessions and is still valid
	var exists bool
	err = app.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM sessions WHERE token = $1 AND expires_at > NOW())",
		oldToken,
	).Scan(&exists)

	if err != nil || !exists {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Generate new token
	newToken, expiresAt, err := generateToken(claims.UserID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Start transaction
	tx, err := app.db.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Delete old session
	_, err = tx.Exec("DELETE FROM sessions WHERE token = $1", oldToken)
	if err != nil {
		http.Error(w, "Error updating session", http.StatusInternalServerError)
		return
	}

	// Create new session
	_, err = tx.Exec(
		"INSERT INTO sessions (user_id, token, expires_at) VALUES ($1, $2, $3)",
		claims.UserID,
		newToken,
		expiresAt,
	)
	if err != nil {
		http.Error(w, "Error creating new session", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(RefreshResponse{
		Token:     newToken,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	})
}
