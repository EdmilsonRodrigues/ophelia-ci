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

func (s *SQLUserStore) CreateUser(user *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Creating user with request: %v", user)
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

func (s *SQLUserStore) GetPublicKeyByUsername(id string) (string, error) {
	log.Printf("Getting private key with id: %v", id)
	query := "SELECT public_key FROM users WHERE id = ?"
	row := s.db.QueryRow(query, id)
	var privateKey string
	err := row.Scan(&privateKey)
	if err != nil {
		log.Printf("Error getting private key: %v", err)
		return "", err
	}
	return privateKey, nil
}
