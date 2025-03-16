package main

import (
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

	db, err := sql.Open("sqlite3", "/var/lib/ophelia/repositories.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repoStore := store.NewSQLRepositoryStore(db)
	if err := repoStore.CreateTable(); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRepositoryServiceServer(s, &server{repositorieStore: repoStore})
	log.Println("Listening on port 50051")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
