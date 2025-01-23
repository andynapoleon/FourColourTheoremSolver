package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AppConfig struct {
	Port     string
	MongoURI string
	Database string
}

var db *mongo.Database

func loadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	return &AppConfig{
		Port:     getEnvOrDefault("PORT", "8083"),
		MongoURI: getEnvOrDefault("MONGO_URI", "mongodb://mongodb:27017"),
		Database: getEnvOrDefault("MONGO_DB", "mapstore"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func connectDB(config *AppConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	db = client.Database(config.Database)
	return nil
}

func setupRoutes(router *mux.Router) {
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	router.HandleFunc("/api/v1/maps", handleSaveMap).Methods("POST")
	router.HandleFunc("/api/v1/maps", handleGetMaps).Methods("GET")
	router.HandleFunc("/api/v1/maps/{id}", handleGetMap).Methods("GET")
	// router.HandleFunc("/api/v1/maps/{id}", handleUpdateMap).Methods("PUT")
	router.HandleFunc("/api/v1/maps/{id}", handleDeleteMap).Methods("DELETE")
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	if err := connectDB(config); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	router := mux.NewRouter()
	setupRoutes(router)

	log.Printf("Map Storage Service starting on port %s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, router); err != nil {
		log.Fatal(err)
	}
}
