// handlers.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func handleMapColoring(w http.ResponseWriter, r *http.Request) {

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

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
	log.Println("Forwarding request to coloring service:", coloringURL)

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

	// Add "Bearer" prefix to the token
	req.Header.Set("Authorization", "Bearer "+token)

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

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Read and validate the request body
	var registrationReq struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registrationReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Forward the request to auth service
	jsonData, err := json.Marshal(registrationReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to process request",
		})
		return
	}

	resp, err := http.Post(
		config.AuthService+"/auth/register",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to connect to auth service",
		})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to read auth service response",
		})
		return
	}

	// Check if the response body is valid JSON
	var jsonResponse interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid response from auth service",
		})
		return
	}

	// Set status code and write response
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
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
