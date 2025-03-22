package main

import (
	"context"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

// CreateUser creates a new user with the given information.
//
// The request must contain the username and public key of the user to be created.
// The username is used to identify the user.
// The public key is used to store the user's public key.
//
// The response will contain the created user information.
//
// Parameters:
// - ctx: The context for the request, which carries deadlines, cancellation signals,
//   and other request-scoped values.
// - req: The request containing the username and public key.
//
// Returns:
// - *pb.UserResponse: The response containing the created user information.
// - error: An error if there is an issue creating the user.
func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Creating user with request: %v", req)
	response, err := s.userStore.CreateUser(req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}
	return response, err
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
// - ctx: The context for the request, which carries deadlines, cancellation signals,
//   and other request-scoped values.
// - req: The request containing the user ID, username and public key.
//
// Returns:
// - *pb.UserResponse: The response containing the updated user information.
// - error: An error if there is an issue updating the user.
func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	log.Printf("Updating user with request: %v", req)
	response, err := s.userStore.UpdateUser(req)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, err
	}
	return response, err
}

// ListUser retrieves a list of all users.
//
// The request must contain an empty request message.
// The response will contain a list of users.
//
// Parameters:
// - ctx: The context for the request, which carries deadlines, cancellation signals,
//   and other request-scoped values.
// - req: An empty request message as defined in the protobuf service definition.
//
// Returns:
// - *pb.ListUserResponse: The response containing the list of users.
// - error: An error if there is an issue retrieving the user list.
func (s *server) ListUser(ctx context.Context, req *pb.Empty) (*pb.ListUserResponse, error) {
	log.Printf("Listing users with request: %v", req)
	users, err := s.userStore.ListUsers()
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return nil, err
	}
	return users, err
}

// GetUser gets a user by either its ID or username.
//
// The request must contain either a non-empty ID or a non-empty username.
// The ID is used to identify the user to be retrieved by ID.
// The username is used to identify the user to be retrieved by username.
//
// The response will contain the user information.
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

// DeleteUser deletes a user by ID.
//
// The request must contain the ID of the user to be deleted.
//
// The response will contain an empty message on success.
func (s *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	log.Printf("Deleting user with request: %v", req)
	err := s.userStore.DeleteUser(req.Id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return nil, err
	}
	return &pb.Empty{}, nil
}
