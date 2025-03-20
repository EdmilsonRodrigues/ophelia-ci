package main

import (
	"context"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Creating user with request: %v", req)
	response, err := s.userStore.CreateUser(req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}
	return response, err
}

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Updating user with request: %v", req)
	response, err := s.userStore.UpdateUser(req)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, err
	}
	return response, err
}

func (s *server) ListUser(ctx context.Context, req *pb.Empty) (*pb.ListUserResponse, error) {
	log.Printf("Listing users with request: %v", req)
	users, err := s.userStore.ListUsers()
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return nil, err
	}
	return users, err
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (response *pb.UserResponse, err error) {
	log.Printf("Getting user with request: %v", req)
	if req.Id == "" {
		response, err = s.userStore.GetUserByUsername(req.Username)
	} else {
		response, err = s.userStore.GetUser(req.Id)
	}
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, err
	}

	return response, err
}

func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	log.Printf("Deleting user with request: %v", req)
	err := s.userStore.DeleteUser(req.Id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return nil, err
	}
	return &pb.Empty{}, nil
}
