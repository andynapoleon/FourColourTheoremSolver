package main

import "math/rand"

// User represents a user in the system
type User struct {
	ID       int    `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // omit password in JSON responses for security
}

// Image represents an image in the system
type Image struct {
	ID  int    `json:"id" gorm:"primary_key"`
	URL string `json:"url"`
}

// UserImage represents the many-to-many relationship between users and images
type UserImage struct {
	UserID  int `gorm:"primaryKey"`
	ImageID int `gorm:"primaryKey"`
}

// Constructor
func NewUser(name string, email string, password string) *User {
	return &User{
		ID:       rand.Intn(10000),
		Name:     name,
		Email:    email,
		Password: password,
	}
}
