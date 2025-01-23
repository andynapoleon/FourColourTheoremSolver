package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getUserIDFromToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	// In a real implementation, you would decode and verify the JWT token
	// For now, we'll just use the token as the user ID
	return parts[1], nil
}

func handleSaveMap(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var map_ Map
	if err := json.NewDecoder(r.Body).Decode(&map_); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	map_.UserID = userID
	map_.CreatedAt = time.Now()
	map_.UpdatedAt = time.Now()

	result, err := db.Collection("maps").InsertOne(context.Background(), map_)
	if err != nil {
		log.Printf("Error saving map: %v", err)
		http.Error(w, "Failed to save map", http.StatusInternalServerError)
		return
	}

	map_.ID = result.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map_)
}

func handleGetMaps(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid map ID", http.StatusBadRequest)
		return
	}

	var map_ Map
	err = db.Collection("maps").FindOne(context.Background(), bson.M{
		"_id":    id,
		"userId": userID,
	}).Decode(&map_)

	if err != nil {
		http.Error(w, "Map not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map_)
}

func handleUpdateMap(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid map ID", http.StatusBadRequest)
		return
	}

	var updateMap Map
	if err := json.NewDecoder(r.Body).Decode(&updateMap); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updateMap.UpdatedAt = time.Now()

	result, err := db.Collection("maps").UpdateOne(
		context.Background(),
		bson.M{"_id": id, "userId": userID},
		bson.M{"$set": updateMap},
	)

	if err != nil {
		log.Printf("Error updating map: %v", err)
		http.Error(w, "Failed to update map", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Map not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDeleteMap(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
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
		http.Error(w, "Map not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
