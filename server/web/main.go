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

	// // Drop the table users first
	// err = store.ClearPostgresStore()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Run server
	server := NewAPISever(":5180", store)
	server.Run()
}
