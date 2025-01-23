package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MapRequest struct {
	UserID    string  `json:"userId"`
	Name      string  `json:"name"`
	ImageData string  `json:"imageData"`
	Matrix    [][]int `json:"matrix"` // Changed back to [][]int
	Width     int     `json:"width"`
	Height    int     `json:"height"`
}

type Map struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"userId" bson:"userId"`
	Name      string             `json:"name" bson:"name"`
	Width     int                `json:"width" bson:"width"`
	Height    int                `json:"height" bson:"height"`
	ImageData string             `json:"imageData" bson:"imageData"`
	Matrix    [][]int            `json:"matrix,omitempty" bson:"matrix,omitempty"` // Changed back to [][]int
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
