package main

import (
	"context"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

func (s *server) CreateRepository(ctx context.Context, req *pb.CreateRepositoryRequest) (*pb.RepositoryResponse, error) {
	log.Printf("Creating repository with request: %v", req)
	response, err := s.repositorieStore.CreateRepository(req)
	if err != nil {
		log.Printf("Error creating repository: %v", err)
		return nil, err
	}
	return &response, err
}

func (s *server) UpdateRepository(ctx context.Context, req *pb.UpdateRepositoryRequest) (*pb.RepositoryResponse, error) {
	log.Printf("Updating repository with request: %v", req)
	response, err := s.repositorieStore.UpdateRepository(req)
	if err != nil {
		log.Printf("Error updating repository: %v", err)
		return nil, err
	}
	return &response, err
}

func (s *server) ListRepository(ctx context.Context, req *pb.Empty) (*pb.ListRepositoryResponse, error) {
	log.Printf("Listing repositories with request: %v", req)
	repos, err := s.repositorieStore.ListRepositories()
	if err != nil {
		log.Printf("Error listing repositories: %v", err)
		return nil, err
	}
	return &repos, err
}

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

func (s *server) DeleteRepository(ctx context.Context, req *pb.DeleteRepositoryRequest) (*pb.Empty, error) {
	log.Printf("Deleting repository with request: %v", req)
	err := s.repositorieStore.DeleteRepository(req.Id)
	if err != nil {
		log.Printf("Error deleting repository: %v", err)
		return nil, err
	}
	return &pb.Empty{}, err
}
