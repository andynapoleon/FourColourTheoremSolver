// handlers.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func handleMapColoring(w http.ResponseWriter, r *http.Request) {
	// Read request body
	var coloringReq ColoringRequest
	if err := json.NewDecoder(r.Body).Decode(&coloringReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if coloringReq.Width <= 0 || coloringReq.Height <= 0 || len(coloringReq.Image) == 0 {
		http.Error(w, "Invalid dimensions or empty image", http.StatusBadRequest)
		return
	}

	// Forward request to Map Coloring Service
	jsonData, err := json.Marshal(coloringReq)
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	config, _ := loadConfig()
	coloringURL := fmt.Sprintf("%s/api/solve", config.ColoringService)

	resp, err := http.Post(coloringURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Error communicating with coloring service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error from coloring service", resp.StatusCode)
		return
	}

	// Forward the response back to client
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func verifyToken(token string) bool {
	config, _ := loadConfig()

	req, err := http.NewRequest("POST", config.AuthService+"/auth/verify", nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	config, _ := loadConfig()

	// Forward the request to auth service
	resp, err := http.Post(
		config.AuthService+"/auth/register",
		"application/json",
		r.Body,
	)
	if err != nil {
		http.Error(w, "Failed to connect to auth service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy status code and response body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	config, _ := loadConfig()

	// Forward the request to auth service
	resp, err := http.Post(
		config.AuthService+"/auth/login",
		"application/json",
		r.Body,
	)
	if err != nil {
		http.Error(w, "Failed to connect to auth service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy status code and response body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	config, _ := loadConfig()

	// Create new request to auth service
	req, err := http.NewRequest("POST", config.AuthService+"/auth/logout", nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Forward the authorization token
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "No authorization token provided", http.StatusUnauthorized)
		return
	}
	req.Header.Set("Authorization", token)

	// Send request to auth service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to connect to auth service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy status code and response body
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
