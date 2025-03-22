package main

import (
	"context"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

// CreateRepository creates a new repository with the given information.
//
// The request must contain the repository name, description and gitignore.
// The gitignore is used to generate the base .gitignore file for the repository.
//
// The response will contain the created repository information.
func (s *server) CreateRepository(ctx context.Context, req *pb.CreateRepositoryRequest) (*pb.RepositoryResponse, error) {
	log.Printf("Creating repository with request: %v", req)
	response, err := s.repositorieStore.CreateRepository(req)
	if err != nil {
		log.Printf("Error creating repository: %v", err)
		return nil, err
	}
	return &response, err
}

// UpdateRepository updates an existing repository with the given information.
//
// The request must contain the repository ID, name and description.
// The ID is used to identify the repository to be updated.
// The name and description are used to update the repository information.
//
// The response will contain the updated repository information.
func (s *server) UpdateRepository(ctx context.Context, req *pb.UpdateRepositoryRequest) (*pb.RepositoryResponse, error) {
	log.Printf("Updating repository with request: %v", req)
	response, err := s.repositorieStore.UpdateRepository(req)
	if err != nil {
		log.Printf("Error updating repository: %v", err)
		return nil, err
	}
	return &response, err
}

// ListRepository lists all existing repositories.
//
// The request must contain an empty request message.
// The response will contain a list of existing repositories.
func (s *server) ListRepository(ctx context.Context, req *pb.Empty) (*pb.ListRepositoryResponse, error) {
	repos, err := s.repositorieStore.ListRepositories()
	if err != nil {
		log.Printf("Error listing repositories: %v", err)
		return nil, err
	}
	return &repos, err
}

// GetRepository gets a repository by either its ID or name.
//
// The request must contain either a non-empty ID or a non-empty name.
// The ID is used to identify the repository to be retrieved by ID.
// The name is used to identify the repository to be retrieved by name.
//
// The response will contain the repository information.
func (s *server) GetRepository(ctx context.Context, req *pb.GetRepositoryRequest) (response *pb.RepositoryResponse, err error) {
	log.Printf("Getting repository with request: %v", req)
	if req.Id == "" {
		response, err = s.repositorieStore.GetRepositoryByName(req.Name)
	} else {
		response, err = s.repositorieStore.GetRepository(req.Id)
	}
	if err != nil {
		log.Printf("Error getting repository: %v", err)
		return nil, err
	}
	return response, err
}

// DeleteRepository deletes an existing repository.
//
// The request must contain the ID of the repository to be deleted.
//
// The response will contain an empty message on success.
func (s *server) DeleteRepository(ctx context.Context, req *pb.DeleteRepositoryRequest) (*pb.Empty, error) {
	log.Printf("Deleting repository with request: %v", req)
	err := s.repositorieStore.DeleteRepository(req.Id)
	if err != nil {
		log.Printf("Error deleting repository: %v", err)
		return nil, err
	}
	return &pb.Empty{}, err
}
