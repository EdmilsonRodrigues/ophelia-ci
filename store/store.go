package ophelia_ci_server

import (
	"database/sql"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci-server"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SQLRepositoryStore struct {
	db *sql.DB
}

type RepositoryStore interface {
	CreateTable() error
	CreateRepository(repo *pb.CreateRepositoryRequest) (pb.RepositoryResponse, error)
	GetRepository(id string) (*pb.RepositoryResponse, error)
	UpdateRepository(repo *pb.UpdateRepositoryRequest) (pb.RepositoryResponse, error)
	ListRepositories() (pb.ListRepositoryResponse, error)
	DeleteRepository(id string) error
}

func NewSQLRepositoryStore(db *sql.DB) *SQLRepositoryStore {
	return &SQLRepositoryStore{
		db: db,
	}
}

func (s *SQLRepositoryStore) CreateTable() error {
	log.Println("Creating repositories table...")
	query := `
		CREATE TABLE IF NOT EXISTS repositories (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
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
	query := "INSERT INTO repositories (id, name, description) VALUES (?, ?, ?)"
	_, err := s.db.Exec(query, id, repo.Name, repo.Description)
	log.Printf("Inserting repository %v with id %v into database...\n", repo.Name, id)
	if err != nil {
		log.Println("Error inserting repository:", err)
		return pb.RepositoryResponse{}, err
	}
	return pb.RepositoryResponse{
		Id:          id,
		Name:        repo.Name,
		Description: repo.Description,
	}, nil
}

func (s *SQLRepositoryStore) GetRepository(id string) (*pb.RepositoryResponse, error) {
	query := "SELECT id, name, description FROM repositories WHERE id = ?"
	var repo pb.RepositoryResponse
	row := s.db.QueryRow(query, id)
	err := row.Scan(repo.Id, repo.Name, repo.Description)
	log.Printf("Getting repository with id %v from database...\n", id)
	if err != nil {
		log.Println("Error getting repository:", err)
		return nil, err
	}
	return &repo, err
}

func (s *SQLRepositoryStore) UpdateRepository(repo *pb.UpdateRepositoryRequest) (pb.RepositoryResponse, error) {
	query := "UPDATE repositories SET name = ?, description = ? WHERE id = ?"
	_, err := s.db.Exec(query, repo.Name, repo.Description, repo.Id)
	log.Printf("Updating repository with id %v in database...\n", repo.Id)
	if err != nil {
		log.Println("Error updating repository:", err)
		return pb.RepositoryResponse{}, err
	}
	return pb.RepositoryResponse{
		Id:          repo.Id,
		Name:        repo.Name,
		Description: repo.Description,
	}, nil
}

func (s *SQLRepositoryStore) ListRepositories() (repos pb.ListRepositoryResponse, err error) {
	query := "SELECT id, name, description FROM repositories"
	rows, err := s.db.Query(query)
	log.Println("Getting all repositories from database...")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var repo pb.RepositoryResponse
		err = rows.Scan(&repo.Id, &repo.Name, &repo.Description)
		if err != nil {
			log.Println("Error scanning repository:", err)
			return
		}
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
