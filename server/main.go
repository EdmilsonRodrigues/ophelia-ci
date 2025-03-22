package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"

	"database/sql"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	"github.com/EdmilsonRodrigues/ophelia-ci/server/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.UnimplementedRepositoryServiceServer
	pb.UnimplementedUserServiceServer
	pb.UnimplementedAuthServiceServer
	pb.UnimplementedHealthServiceServer

	userStore        store.UserStore
	repositorieStore store.RepositoryStore
	challenges       sync.Map
}

// Main starts the Ophelia CI Server Service.
func main() {
	log.Println("Ophelia CI Server Service started!")

	config := LoadConfig()

	db, err := sql.Open("sqlite3", config.Server.DBPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repoStore := store.NewSQLRepositoryStore(db)
	userStore := store.NewSQLUserStore(db)

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.UnaryInterceptor(AuthInterceptor)}

	if config.SSL.CertFile != "" && config.SSL.KeyFile != "" {
		log.Println("Using SSL")
		cert, err := tls.LoadX509KeyPair(config.SSL.CertFile, config.SSL.KeyFile)
		if err != nil {
			log.Fatalf("Failed to load credentials: %v", err)
		}

		opts = append(opts, grpc.Creds(
			credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
			})))
	}

	s := grpc.NewServer(opts...)

	mainServer := &server{
		repositorieStore: repoStore,
		userStore:        userStore,
	}
	pb.RegisterRepositoryServiceServer(s, mainServer)
	pb.RegisterUserServiceServer(s, mainServer)
	pb.RegisterAuthServiceServer(s, mainServer)
	pb.RegisterHealthServiceServer(s, mainServer)
	log.Printf("Listening on port %d\n", config.Server.Port)
	log.Printf("For logging in for the first time, use the following key: %v", uniqueKey)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Health is a gRPC service method that provides a simple health check.
// It returns an empty response to indicate the server is running and reachable.
//
// Parameters:
// - ctx: The context for the request, which carries deadlines, cancellation signals,
//   and other request-scoped values.
// - req: An empty request message as defined in the protobuf service definition.
//
// Returns:
// - *pb.Empty: An empty response message indicating the server is healthy.
// - error: Always returns nil since this operation is expected to succeed.
func (s *server) Health(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
