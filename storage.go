package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	UpdateUser(*User) error
	GetUsers() ([]*User, error)
	GetUserById(int) (*User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=goweblikepeterparker sslmode=disable"

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

func (s *PostgresStore) Init() error {
	return s.CreateUserTable()
}

func (s *PostgresStore) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL, 
		password VARCHAR(100) NOT NULL
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateUser(user *User) error {
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

func (s *PostgresStore) UpdateUser(user *User) error {
	return nil
}

func (s *PostgresStore) GetUserById(id int) (*User, error) {
	return nil, nil
}

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
