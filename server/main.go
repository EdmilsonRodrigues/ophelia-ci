package ophelia_ci_server

import (
	"context"
	"database/sql"
	"log"
	"net"

	"google.golang.org/grpc"
        pb "github.com/EdmilsonRodrigues/ophelia-ci-server"
        store "github.com/EdmilsonRodrigues/ophelia-ci-server/store"
)

type server struct {
        pb.UnimplementedRepositoryServiceServer
        repositorieStore store.RepositoryStore
}

func newServer() *server {
        db, err := sql.Open("sqlite3", "/var/lib/git/repositories.db")
        if err != nil {
                log.Fatal(err)
        }
        return &server{
                repositorieStore: store.NewSQLRepositoryStore(db),
        }
}

func (s *server) CreateRepository(ctx context.Context, req *pb.CreateRepositoryRequest) (*pb.RepositoryResponse, error) {
        log.Printf("CreateRepository called with request: %v", req)
        return nil, nil
}

func (s *server) UpdateRepository(ctx context.Context, req *pb.UpdateRepositoryRequest) (*pb.RepositoryResponse, error) {
        log.Printf("UpdateRepository called with request: %v", req)
        return nil, nil
}


func (s *server) ListRepository(ctx context.Context, req *pb.Empty) (*pb.ListRepositoryResponse, error) {
        log.Printf("ListRepository called with request: %v", req)
        return nil, nil
}

func (s *server) GetRepository(ctx context.Context, req *pb.GetRepositoryRequest) (*pb.RepositoryResponse, error) {
        log.Printf("GetRepository called with request: %v", req)
        return nil, nil
}

func (s *server) DeleteRepository(ctx context.Context, req *pb.DeleteRepositoryRequest) (*pb.Empty, error) {
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
        pb.RegisterRepositoryServiceServer(s, newServer())
        log.Println("Listening on port 50051")
        if err := s.Serve(lis); err != nil {
                log.Fatalf("Failed to serve: %v", err)
        }
}
