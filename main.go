package ophelia_ci_server

import (
	"context"
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
        UnimplementedRepositoryServiceServer
        repositorieStore RepositoryStore
}

func newServer() *server {
        db, err := sql.Open("sqlite3", "/var/lib/git/repositories.db")
        if err != nil {
                log.Fatal(err)
        }
        return &server{
                repositorieStore: NewSQLRepositoryStore(db),
        }
}

func (s *server) CreateRepository(ctx context.Context, req *CreateRepositoryRequest) (*RepositoryResponse, error) {
        log.Printf("CreateRepository called with request: %v", req)
        return nil, nil
}

func (s *server) UpdateRepository(ctx context.Context, req *UpdateRepositoryRequest) (*RepositoryResponse, error) {
        log.Printf("UpdateRepository called with request: %v", req)
        return nil, nil
}


func (s *server) ListRepository(ctx context.Context, req *Empty) (*ListRepositoryResponse, error) {
        log.Printf("ListRepository called with request: %v", req)
        return nil, nil
}

func (s *server) GetRepository(ctx context.Context, req *GetRepositoryRequest) (*RepositoryResponse, error) {
        log.Printf("GetRepository called with request: %v", req)
        return nil, nil
}

func (s *server) DeleteRepository(ctx context.Context, req *DeleteRepositoryRequest) (*Empty, error) {
        log.Printf("DeleteRepository called with request: %v", req)
        return nil, nil
}

func main() {
        log.Println("Ophelia CI Server Service started!")
        lis, err := net.Listen("tcp", ":50051")
        if err != nil {
                log.Fatalf("Failed to listen: %v", err)
        }
        s := grpc.NewServer()
        RegisterRepositoryServiceServer(s, newServer())
        log.Println("Listening on port 50051")
        if err := s.Serve(lis); err != nil {
                log.Fatalf("Failed to serve: %v", err)
        }
}
