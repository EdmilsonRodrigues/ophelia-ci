package main

import (
	"context"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

// CommitSignal is a gRPC service method that sends a signal to the server when a commit is pushed to a repository.
//
// Parameters:
//   - ctx: The context for the request, which carries deadlines, cancellation signals,
//     and other request-scoped values.
//   - req: The request containing the repository name and commit hash.
//
// Returns:
//   - *pb.Empty: An empty response message indicating the signal was sent successfully.
//   - error: An error if there is an issue sending the signal.
func (s *server) CommitSignal(ctx context.Context, req *pb.CommitRequest) (*pb.Empty, error) {
	// Improve this method in roadmap
	log.Printf("Commit signal with request: %v", req)
	repo, err := s.repositorieStore.GetRepositoryByName(req.Repository)
	if err != nil {
		log.Printf("Error getting repository: %v", err)
		return nil, err
	}

	_, err = s.repositorieStore.UpdateRepository(&pb.UpdateRepositoryRequest{
		Id:          repo.Id,
		Name:        repo.Name,
		Description: repo.Description,
	})

	if err != nil {
		log.Printf("Error updating repository: %v", err)
		return nil, err
	}

	return &pb.Empty{}, nil
}
