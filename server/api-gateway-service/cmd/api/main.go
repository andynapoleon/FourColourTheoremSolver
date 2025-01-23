package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type ColoringRequest struct {
	Image struct {
		Data []uint8 `json:"data"`
	} `json:"image"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AppConfig struct {
	Port              string
	ColoringService   string
	AuthService       string
	MapStorageService string
}

func loadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	return &AppConfig{
		Port:              getEnvOrDefault("PORT", "80"),
		ColoringService:   getEnvOrDefault("COLORING_SERVICE_URL", "http://solver-service"),
		AuthService:       getEnvOrDefault("AUTHENTICATION_SERVICE_URL", "http://authentication-service"),
		MapStorageService: getEnvOrDefault("MAP_STORAGE_SERVICE_URL", "http://map-storage-service"),
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
	}).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/auth/register", handleRegister).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/auth/login", handleLogin).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/auth/logout", handleLogout).Methods("POST", "OPTIONS")

	// Map solver routes (protected)
	router.HandleFunc("/api/v1/maps/color", handleMapColoring).Methods("POST", "OPTIONS")

	// Map storage routes (protected)
	router.HandleFunc("/api/v1/maps", handleMapStorage).Methods("POST", "GET", "OPTIONS")
	router.HandleFunc("/api/v1/maps/{id}", handleMapStorage).Methods("GET", "PUT", "DELETE", "OPTIONS")
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	router := mux.NewRouter()

	// Apply CORS middleware first
	router.Use(corsMiddleware)

	// Setup routes
	setupRoutes(router)

	// Other middleware
	router.Use(loggingMiddleware)
	router.Use(authMiddleware)

	log.Printf("Server starting on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, router); err != nil {
		log.Fatal(err)
	}
}
