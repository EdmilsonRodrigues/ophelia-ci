package store

import (
	"database/sql"
	"log"
	"time"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserStore interface {
	CreateUser(user *pb.CreateUserRequest) (*pb.UserResponse, error)
	GetUser(id string) (*pb.UserResponse, error)
	GetUserByUsername(name string) (*pb.UserResponse, error)
	UpdateUser(user *pb.UpdateUserRequest) (*pb.UserResponse, error)
	ListUsers() (*pb.ListUserResponse, error)
	DeleteUser(id string) error
	GetPublicKeyByUsername(username string) (string, error)
}

type SQLUserStore struct {
	db *sql.DB
}

// NewSQLUserStore creates a new SQLUserStore given a database connection.
//
// If the users table does not exist in the database, it will be created.
//
// The function will log a fatal error if there is an issue creating the table.
func NewSQLUserStore(db *sql.DB) *SQLUserStore {
	store := &SQLUserStore{
		db: db,
	}
	err := store.CreateTable()
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
	return store
}

// CreateTable creates the users table in the SQLite database if it does not exist.
//
// The users table has the following columns:
// - id: the ID of the user, which is the primary key
// - username: the username of the user
// - public_key: the public key of the user
// - created_at: the timestamp when the user was created
// - updated_at: the timestamp when the user was last updated
//
// Returns an error if there is an issue creating the table.
func (s *SQLUserStore) CreateTable() error {
	log.Println("Creating users table...")
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            username TEXT NOT NULL,
            public_key TEXT,
			created_at INTEGER,
			updated_at INTEGER
        );
    `
	_, err := s.db.Exec(query)
	if err != nil {
		log.Println("Error creating users table:", err)
		return err
	}
	return nil
}

// CreateUser creates a new user with the given information.
//
// The request must contain the username and public key of the user to be created.
// The username is used to identify the user.
// The public key is used to store the user's public key.
//
// The response will contain the created user information.
//
// Parameters:
// - user: The request containing the username and public key.
//
// Returns:
// - *pb.UserResponse: The response containing the created user information.
// - error: An error if there is an issue creating the user.
func (s *SQLUserStore) CreateUser(user *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Adding user to database with request: %v", user)
	id := uuid.New().String()
	createdAt := timestamppb.Now()
	updatedAt := timestamppb.Now()
	query := "INSERT INTO users (id, username, public_key, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	_, err := s.db.Exec(query, id, user.Username, user.PublicKey, createdAt.AsTime().Unix(), updatedAt.AsTime().Unix())
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return &pb.UserResponse{}, err
	}
	return &pb.UserResponse{
		Id:       id,
		Username: user.Username,
	}, nil
}

// GetUser retrieves a user by ID from the database.
//
// Parameters:
// - id: The ID of the user to be retrieved.
//
// Returns:
// - *pb.UserResponse: The response containing the user information.
// - error: An error if there is an issue retrieving the user.
func (s *SQLUserStore) GetUser(id string) (*pb.UserResponse, error) {
	log.Printf("Getting user with id: %v", id)
	query := "SELECT id, username FROM users WHERE id = ?"
	row := s.db.QueryRow(query, id)
	user := &pb.UserResponse{}
	err := row.Scan(&user.Id, &user.Username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return &pb.UserResponse{}, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username from the database.
//
// Parameters:
// - username: The username of the user to be retrieved.
//
// Returns:
// - *pb.UserResponse: The response containing the user information.
// - error: An error if there is an issue retrieving the user.
func (s *SQLUserStore) GetUserByUsername(username string) (*pb.UserResponse, error) {
	log.Printf("Getting user with username: %v", username)
	query := "SELECT id, username FROM users WHERE username = ?"
	row := s.db.QueryRow(query, username)
	user := &pb.UserResponse{}
	err := row.Scan(&user.Id, &user.Username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return &pb.UserResponse{}, err
	}
	return user, nil
}

// UpdateUser updates an existing user with the given information.
//
// The request must contain the user ID, username and public key.
// The ID is used to identify the user to be updated.
// The username and public key are used to update the user information.
//
// The response will contain the user information.
//
// Parameters:
// - user: The request containing the user ID, username and public key.
//
// Returns:
// - *pb.UserResponse: The response containing the updated user information.
// - error: An error if there is an issue updating the user.
func (s *SQLUserStore) UpdateUser(user *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Updating user with request: %v", user)
	query := "UPDATE users SET username = ?, public_key = ?, updated_at = ? WHERE id = ?"
	_, err := s.db.Exec(query, user.Username, user.PublicKey, time.Now().Unix(), user.Id)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return &pb.UserResponse{}, err
	}
	return &pb.UserResponse{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}

// ListUsers lists all existing users in the database.
//
// The response will contain a list of existing users.
//
// Returns:
// - *pb.ListUserResponse: The response containing the list of users.
// - error: An error if there is an issue listing users.
func (s *SQLUserStore) ListUsers() (*pb.ListUserResponse, error) {
	log.Println("Listing users...")
	query := "SELECT id, username FROM users"
	rows, err := s.db.Query(query)
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return &pb.ListUserResponse{}, err
	}
	defer rows.Close()
	users := &pb.ListUserResponse{}
	for rows.Next() {
		user := &pb.UserResponse{}
		err := rows.Scan(&user.Id, &user.Username)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return &pb.ListUserResponse{}, err
		}
		users.Users = append(users.Users, user)
	}
	return users, nil
}

// DeleteUser deletes the user with the specified ID from the database.
//
// Parameters:
// - id: The ID of the user to be deleted.
//
// Returns:
// - error: An error if the user cannot be deleted.
func (s *SQLUserStore) DeleteUser(id string) error {
	log.Printf("Deleting user with id: %v", id)
	query := "DELETE FROM users WHERE id = ?"
	_, err := s.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}
	return nil
}

// GetPublicKeyByUsername retrieves the public key associated with the given username.
//
// This function queries the users table in the database to find the public key
// corresponding to the specified username. If the username does not exist or
// an error occurs during the query, an error is returned.
//
// Parameters:
// - username: The username for which the public key is to be retrieved.
//
// Returns:
// - string: The public key associated with the given username.
// - error: An error if the public key cannot be retrieved.

func (s *SQLUserStore) GetPublicKeyByUsername(username string) (string, error) {
	log.Printf("Getting public key with username: %v", username)
	query := "SELECT public_key FROM users WHERE username = ?"
	row := s.db.QueryRow(query, username)
	var publicKey string
	err := row.Scan(&publicKey)
	if err != nil {
		log.Printf("Error getting public key: %v", err)
		return "", err
	}
	return publicKey, nil
}
