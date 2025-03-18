package ophelia_ci_server

import (
	"database/sql"
	"log"
	"time"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SQLRepositoryStore struct {
	db *sql.DB
}

type RepositoryStore interface {
	CreateTable() error
	CreateRepository(repo *pb.CreateRepositoryRequest) (pb.RepositoryResponse, error)
	GetRepository(id string) (*pb.RepositoryResponse, error)
	GetRepositoryByName(name string) (*pb.RepositoryResponse, error)
	UpdateRepository(repo *pb.UpdateRepositoryRequest) (pb.RepositoryResponse, error)
	ListRepositories() (pb.ListRepositoryResponse, error)
	DeleteRepository(id string) error
	Close() error
}

func NewSQLRepositoryStore(db *sql.DB) *SQLRepositoryStore {
	store := &SQLRepositoryStore{
		db: db,
	}
	err := store.CreateTable()
	if err != nil {
		log.Fatalf("Failed to create repositories table: %v", err)
	}
	return store
}

func (s *SQLRepositoryStore) CreateTable() error {
	log.Println("Creating repositories table...")
	query := `
        CREATE TABLE IF NOT EXISTS repositories (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT,
            last_update INTEGER
        );
    `
	_, err := s.db.Exec(query)
	if err != nil {
		log.Println("Error creating repositories table:", err)
		return err
	}
	return nil
}

func (s *SQLRepositoryStore) CreateRepository(repo *pb.CreateRepositoryRequest) (pb.RepositoryResponse, error) {
	id := uuid.New().String()
	now := timestamppb.Now()
	query := "INSERT INTO repositories (id, name, description, last_update) VALUES (?, ?, ?, ?)"
	_, err := s.db.Exec(query, id, repo.Name, repo.Description, now.Seconds)
	log.Printf("Inserting repository %v with id %v into database...\n", repo.Name, id)
	if err != nil {
		log.Println("Error inserting repository:", err)
		return pb.RepositoryResponse{}, err
	}
	return pb.RepositoryResponse{
		Id:          id,
		Name:        repo.Name,
		Description: repo.Description,
		LastUpdate:  now,
	}, nil
}

func (s *SQLRepositoryStore) GetRepository(id string) (*pb.RepositoryResponse, error) {
	query := "SELECT id, name, description, last_update FROM repositories WHERE id = ?"
	var repo pb.RepositoryResponse
	var lastUpdateSeconds int64
	row := s.db.QueryRow(query, id)
	err := row.Scan(&repo.Id, &repo.Name, &repo.Description, &lastUpdateSeconds)
	log.Printf("Getting repository with id %v from database...\n", id)
	if err != nil {
		log.Println("Error getting repository:", err)
		return nil, err
	}
	repo.LastUpdate = timestamppb.New(time.Unix(lastUpdateSeconds, 0))
	return &repo, err
}

func (s *SQLRepositoryStore) GetRepositoryByName(name string) (*pb.RepositoryResponse, error) {
	query := "SELECT id, name, description, last_update FROM repositories WHERE name = ?"
	var repo pb.RepositoryResponse
	var lastUpdateSeconds int64
	row := s.db.QueryRow(query, name)
	err := row.Scan(&repo.Id, &repo.Name, &repo.Description, &lastUpdateSeconds)
	if err != nil {
		return nil, err
	}
	repo.LastUpdate = timestamppb.New(time.Unix(lastUpdateSeconds, 0))
	return &repo, nil
}

func (s *SQLRepositoryStore) UpdateRepository(repo *pb.UpdateRepositoryRequest) (pb.RepositoryResponse, error) {
	now := timestamppb.Now()
	query := "UPDATE repositories SET name = ?, description = ?, last_update = ? WHERE id = ?"
	_, err := s.db.Exec(query, repo.Name, repo.Description, now.Seconds, repo.Id)
	log.Printf("Updating repository with id %v in database...\n", repo.Id)
	if err != nil {
		log.Println("Error updating repository:", err)
		return pb.RepositoryResponse{}, err
	}
	return pb.RepositoryResponse{
		Id:          repo.Id,
		Name:        repo.Name,
		Description: repo.Description,
		LastUpdate:  now,
	}, nil
}

func (s *SQLRepositoryStore) ListRepositories() (repos pb.ListRepositoryResponse, err error) {
	query := "SELECT id, name, description, last_update FROM repositories"
	rows, err := s.db.Query(query)
	log.Println("Getting all repositories from database...")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var repo pb.RepositoryResponse
		var lastUpdateSeconds int64
		err = rows.Scan(&repo.Id, &repo.Name, &repo.Description, &lastUpdateSeconds)
		if err != nil {
			log.Println("Error scanning repository:", err)
			return
		}
		repo.LastUpdate = timestamppb.New(time.Unix(lastUpdateSeconds, 0))
		repos.Repositories = append(repos.Repositories, &repo)
	}
	return
}

func (s *SQLRepositoryStore) DeleteRepository(id string) (err error) {
	query := "DELETE FROM repositories WHERE id = ?"
	_, err = s.db.Exec(query, id)
	log.Printf("Deleting repository with id %v from database...\n", id)
	if err != nil {
		log.Println("Error deleting repository:", err)
		return
	}
	return
}

func (s *SQLRepositoryStore) Close() error {
	return s.db.Close()
}
