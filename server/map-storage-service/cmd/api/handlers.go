package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func handleSaveMap(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var mapData MapRequest
	if err := json.NewDecoder(r.Body).Decode(&mapData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate userID
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
		ImageData: []uint8(mapData.ImageData),
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

	var map_ Map
	err = db.Collection("maps").FindOne(context.Background(), bson.M{
		"_id":    id,
		"userId": userID,
	}).Decode(&map_)

	if err != nil {
		log.Printf("Error fetching map: %v", err)
		http.Error(w, "Map not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map_)
}

func handleUpdateMap(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var updateData MapRequest
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate userID
	if updateData.UserID == "" {
		log.Printf("No userID provided in request body")
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

	updateMap := Map{
		UserID:    updateData.UserID,
		Name:      updateData.Name,
		Width:     updateData.Width,
		Height:    updateData.Height,
		ImageData: []uint8(updateData.ImageData),
		UpdatedAt: time.Now(),
	}

	result, err := db.Collection("maps").UpdateOne(
		context.Background(),
		bson.M{"_id": id, "userId": updateData.UserID},
		bson.M{"$set": updateMap},
	)

	if err != nil {
		log.Printf("Error updating map: %v", err)
		http.Error(w, "Failed to update map", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		log.Printf("Map not found for user")
		http.Error(w, "Map not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
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
