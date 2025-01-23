package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Map struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string             `json:"userId" bson:"userId"`
	Name      string             `json:"name" bson:"name"`
	Width     int                `json:"width" bson:"width"`
	Height    int                `json:"height" bson:"height"`
	ImageData []uint8            `json:"imageData" bson:"imageData"`
	Solution  [][]int            `json:"solution,omitempty" bson:"solution,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
