package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func handleSaveMap(w http.ResponseWriter, r *http.Request) {
	// Read and log the raw request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Log the raw request
	log.Printf("Raw request body: %s", string(body))

	// Create a new reader with the body data
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// First try to decode into a map to inspect the data
	var rawData map[string]interface{}
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&rawData); err != nil {
		log.Printf("Error decoding raw data: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log the matrix type and structure
	if matrix, ok := rawData["matrix"]; ok {
		log.Printf("Matrix type: %T", matrix)
		log.Printf("Matrix value: %v", matrix)
	}

	// Reset the body reader
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Try to decode into the actual struct
	var mapData MapRequest
	if err := json.NewDecoder(r.Body).Decode(&mapData); err != nil {
		log.Printf("Error decoding into MapRequest: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if mapData.UserID == "" {
		log.Printf("No userID provided in request body")
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	// Create new map
	newMap := Map{
		UserID:    mapData.UserID,
		Name:      mapData.Name,
		Width:     mapData.Width,
		Height:    mapData.Height,
		ImageData: mapData.ImageData,
		Matrix:    mapData.Matrix,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert into database
	result, err := db.Collection("maps").InsertOne(context.Background(), newMap)
	if err != nil {
		log.Printf("Error saving map: %v", err)
		http.Error(w, "Failed to save map", http.StatusInternalServerError)
		return
	}

	// Get the inserted ID
	newMap.ID = result.InsertedID.(primitive.ObjectID)

	// Return the saved map
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMap)
}

func handleGetMaps(w http.ResponseWriter, r *http.Request) {
	// Get userID from query parameter
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		log.Printf("No userID provided in query parameters")
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	cursor, err := db.Collection("maps").Find(context.Background(), bson.M{"userId": userID})
	if err != nil {
		log.Printf("Error fetching maps: %v", err)
		http.Error(w, "Failed to fetch maps", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var maps []Map
	if err = cursor.All(context.Background(), &maps); err != nil {
		log.Printf("Error decoding maps: %v", err)
		http.Error(w, "Failed to decode maps", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(maps)
}

func handleGetMap(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mapID := vars["id"]

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(mapID)
	if err != nil {
		log.Printf("Invalid map ID format: %v", err)
		http.Error(w, "Invalid map ID", http.StatusBadRequest)
		return
	}

	// Find the map in the database
	var mapData Map
	err = db.Collection("maps").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&mapData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Map not found", http.StatusNotFound)
			return
		}
		log.Printf("Error finding map: %v", err)
		http.Error(w, "Failed to retrieve map", http.StatusInternalServerError)
		return
	}

	// Return the map data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapData)
}

func handleDeleteMap(w http.ResponseWriter, r *http.Request) {
	// Get userID from query parameter
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		log.Printf("No userID provided in query parameters")
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Printf("Invalid map ID: %v", err)
		http.Error(w, "Invalid map ID", http.StatusBadRequest)
		return
	}

	result, err := db.Collection("maps").DeleteOne(context.Background(), bson.M{
		"_id":    id,
		"userId": userID,
	})

	if err != nil {
		log.Printf("Error deleting map: %v", err)
		http.Error(w, "Failed to delete map", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		log.Printf("Map not found for user")
		http.Error(w, "Map not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
