package main

import (
	"log"
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
	server := NewAPISever(":5180", store)
	server.Run()
}
