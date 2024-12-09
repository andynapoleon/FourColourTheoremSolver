package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type ColoringRequest struct {
	Image  map[string]uint8 `json:"image"`
	Width  int              `json:"width"`
	Height int              `json:"height"`
}

type AppConfig struct {
	Port            string
	ColoringService string
	AuthService     string
}

func loadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	return &AppConfig{
		Port:            getEnvOrDefault("PORT", "80"),
		ColoringService: getEnvOrDefault("COLORING_SERVICE_URL", "http://color-service"),
		AuthService:     getEnvOrDefault("AUTH_SERVICE_URL", "http://authentication-service"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupRoutes(router *mux.Router) {
	// Auth routes (unprotected)
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	router.HandleFunc("/api/v1/auth/register", handleRegister).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", handleLogin).Methods("POST")

	// Logout requires a valid token, but we'll handle it separately from other protected routes
	router.HandleFunc("/api/v1/auth/logout", handleLogout).Methods("POST")

	// Protected routes
	router.HandleFunc("/api/v1/maps/color", handleMapColoring).Methods("POST")
	// Add other protected routes here
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	router := mux.NewRouter()
	setupRoutes(router)

	// Add all middleware including auth
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)
	router.Use(authMiddleware) // Add this line to enable authentication

	log.Printf("Server starting on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, router); err != nil {
		log.Fatal(err)
	}
}
