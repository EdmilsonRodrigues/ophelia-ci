package main

import (
	"context"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	"github.com/EdmilsonRodrigues/ophelia-ci/server/git"
)

// CreateRepository creates a new repository with the given information.
//
// The request must contain the repository name, description and gitignore.
// The gitignore is used to generate the base .gitignore file for the repository.
//
// The response will contain the created repository information.
func (s *server) CreateRepository(ctx context.Context, req *pb.CreateRepositoryRequest) (*pb.RepositoryResponse, error) {
	log.Printf("Creating repository with request: %v", req)
	log.Printf("Creating git repository for %v", req.Name)
	err := git.CreateGitRepository(getRepoPath(req.Name), req.Gitignore)
	if err != nil {
		log.Printf("Error creating git repository: %v", err)
		return nil, err
	}
	log.Printf("Creating repository in database for %v", req.Name)
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
	log.Printf("Getting repository with id: %v", req.Id)
	old_repo, err := s.repositorieStore.GetRepository(req.Id)
	if err != nil {
		log.Printf("Error getting repository: %v", err)
		return nil, err
	}
	log.Printf("Updating git repository from %v to %v", getRepoPath(old_repo.Name), getRepoPath(req.Name))
	err = git.UpdateGitRepository(getRepoPath(old_repo.Name), getRepoPath(req.Name))
	if err != nil {
		log.Printf("Error updating git repository: %v", err)
		return nil, err
	}
	log.Printf("Updating repository in database: %v", getRepoPath(req.Name))
	response, err := s.repositorieStore.UpdateRepository(req)
	if err != nil {
		log.Printf("Error updating repository: %v", err)
		log.Printf("Rolling back git repository: %v", getRepoPath(req.Name))
		rollbackErr := git.UpdateGitRepository(getRepoPath(req.Name), getRepoPath(old_repo.Name))
		if rollbackErr != nil {
			log.Printf("Error rolling back git repository: %v", rollbackErr)
		}
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
	log.Printf("Getting repository with Id: %v", req.Id)
	old_repo, err := s.repositorieStore.GetRepository(req.Id)
	if err != nil {
		log.Printf("Error getting repository: %v", err)
		return nil, err
	}
	log.Printf("Deleting repository in database: %v", getRepoPath(old_repo.Name))
	err = s.repositorieStore.DeleteRepository(req.Id)
	if err != nil {
		log.Printf("Error deleting repository: %v", err)
		return nil, err
	}
	log.Printf("Deleting git repository: %v", getRepoPath(old_repo.Name))
	err = git.DeleteGitRepository(getRepoPath(old_repo.Name))
	if err != nil {
		log.Printf("Error deleting git repository: %v", err)
		log.Printf("Rolling back repository: %v", getRepoPath(old_repo.Name))
		_, rollbackErr := s.repositorieStore.CreateRepository(&pb.CreateRepositoryRequest{
			Name:        old_repo.Name,
			Description: old_repo.Description,
		})
		if rollbackErr != nil {
			log.Printf("Error rolling back repository: %v", rollbackErr)
		}
		return nil, err
	}

	return &pb.Empty{}, err
}

func getRepoPath(repoName string) string {
	return homePath + "/" + repoName + ".git"
}
