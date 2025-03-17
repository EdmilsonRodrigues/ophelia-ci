package main

import (
	"fmt"
	"log"
	"net"

	"database/sql"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	store "github.com/EdmilsonRodrigues/ophelia-ci/server/store"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedRepositoryServiceServer
	repositorieStore store.RepositoryStore
}

func main() {
	log.Println("Ophelia CI Server Service started!")

	config := LoadConfig()

	repoStore, err := openRepositoryStore()
	if err != nil {
		log.Fatalf("Failed to open repository store: %v", err)
	}
	defer repoStore.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRepositoryServiceServer(s, &server{repositorieStore: repoStore})
	log.Printf("Listening on port %d\n", config.Server.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func openRepositoryStore() (*store.SQLRepositoryStore, error) {
	db, err := sql.Open("sqlite3", "/var/lib/ophelia/repositories.db")
	if err != nil {
		return nil, err
	}
	return store.NewSQLRepositoryStore(db), nil
}
