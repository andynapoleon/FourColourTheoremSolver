// handlers.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func handleMapColoring(w http.ResponseWriter, r *http.Request) {
	log.Printf("[MapColoring] Starting handler")

	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Read request body
	var coloringReq ColoringRequest
	if err := json.NewDecoder(r.Body).Decode(&coloringReq); err != nil {
		log.Printf("[MapColoring] Error parsing JSON: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Access the image data correctly
	imageData := coloringReq.Image.Data

	log.Printf("[MapColoring] Request decoded - Width: %d, Height: %d, Image pixels: %d",
		coloringReq.Width, coloringReq.Height, len(imageData))

	// Convert uint8 array to int array
	intImageData := make([]int, len(imageData))
	for i, v := range imageData {
		intImageData[i] = int(v)
	}

	config, err := loadConfig()
	if err != nil {
		log.Printf("[MapColoring] Error loading config: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	coloringURL := fmt.Sprintf("%s/api/solve", config.ColoringService)
	log.Printf("[MapColoring] Calling coloring service at: %s", coloringURL)

	// Create the request body
	serviceReq := map[string]interface{}{
		"image":  intImageData,
		"width":  coloringReq.Width,
		"height": coloringReq.Height,
	}

	jsonData, err := json.Marshal(serviceReq)
	if err != nil {
		log.Printf("[MapColoring] Error marshaling request: %v", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", coloringURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[MapColoring] Error creating request: %v", err)
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[MapColoring] Error calling coloring service: %v", err)
		http.Error(w, "Error communicating with coloring service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[MapColoring] Error reading response: %v", err)
		http.Error(w, "Error reading service response", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[MapColoring] Coloring service returned status: %d, body: %s",
			resp.StatusCode, string(responseBody))
		http.Error(w, string(responseBody), resp.StatusCode)
		return
	}

	log.Printf("[MapColoring] Successfully processed request")
	w.Write(responseBody)
}

func handleMapStorage(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		http.Error(w, "Failed to load configuration", http.StatusInternalServerError)
		return
	}

	// Forward the request to map storage service
	url := config.MapStorageService + r.URL.Path
	log.Printf("Full URL being called: %s", url)

	// Create new request
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	req.Header = r.Header
	log.Printf("Request headers: %v", req.Header)

	// Send request
	client := &http.Client{
		Timeout: 10 * time.Second, // Add timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error connecting to map storage service: %v", err)
		log.Printf("Attempted to connect to: %s", url)
		http.Error(w, "Failed to connect to map storage service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code and body
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response body: %v", err)
	}
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
	log.Printf("Forwarding request to auth service: %s", config.AuthService+"/auth/login")
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
