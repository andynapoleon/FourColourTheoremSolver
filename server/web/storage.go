package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Storage interface
type Storage interface {
	CreateUser(*User) error
	GetUsers() ([]*User, error)
	GetUserByEmail(string) (*User, error)
}

// PostgresStore struct
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore creates a new PostgresStore
func NewPostgresStore() (*PostgresStore, error) {
	connStr := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

// Init initializes the PostgresStore
func (s *PostgresStore) Init() error {
	return s.CreateUserTable()
}

// CreateUserTable creates the users table
func (s *PostgresStore) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name VARCHAR(60) NOT NULL,
		email VARCHAR(60) UNIQUE NOT NULL, 
		password VARCHAR(60) NOT NULL
	);`

	_, err := s.db.Exec(query)
	return err
}

// CreateUser creates a new user
func (s *PostgresStore) CreateUser(user *User) error {
	if len(user.Name) > 60 {
		return fmt.Errorf("name too long: maximum length is 60 characters")
	}
	if len(user.Email) > 60 {
		return fmt.Errorf("email too long: maximum length is 60 characters")
	}
	if len(user.Password) > 60 {
		return fmt.Errorf("password too long: maximum length is 60 characters")
	}
	query := `
        INSERT INTO users (name, email, password) 
        VALUES ($1, $2, $3)
    `
	resp, err := s.db.Query(
		query,
		user.Name,
		user.Email,
		user.Password,
	)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

// GetUsers gets all users
func (s *PostgresStore) GetUsers() ([]*User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for rows.Next() {
		user := new(User)
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByEmail gets a user by email
func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	user := new(User)
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return user, nil
}
