package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Verify all required environment variables are set
	requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			log.Fatalf("Required environment variable %s is not set", env)
		}
	}

	db, err := initDB()
	if err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.Close()

	app := &App{
		db: db,
	}

	router := mux.NewRouter()
	setupRoutes(router, app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Auth service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func setupRoutes(router *mux.Router, app *App) {
	router.HandleFunc("/auth/register", app.handleRegister).Methods("POST")
	router.HandleFunc("/auth/login", app.handleLogin).Methods("POST")
	router.HandleFunc("/auth/verify", app.handleVerifyToken).Methods("POST")
	router.HandleFunc("/auth/refresh", app.handleRefreshToken).Methods("POST")
	router.HandleFunc("/auth/logout", app.handleLogout).Methods("POST")
}
