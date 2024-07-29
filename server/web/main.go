package main

import (
	"log"
	"os"
)

func main() {
	// Postgres initialization
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// Clear the users table
	// if _, err := store.db.Exec("DELETE FROM users"); err != nil {
	// 	log.Fatal(err)
	// }

	// Run server like a Go pro
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	server := NewAPISever(":"+port, store)
	server.Run()
}
