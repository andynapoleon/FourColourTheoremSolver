package main

// CreateUserRequest object
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginUserRequest object
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Image represents an image in the system
type Image struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// UserImage represents the many-to-many relationship between users and images
type UserImage struct {
	UserID  int
	ImageID int
}

// Constructor
func NewUser(name string, email string, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}
