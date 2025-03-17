package main

import (
        "context"
        "flag"
        "fmt"
        "log"
        "os"
        "time"

        "google.golang.org/grpc"
        "google.golang.org/grpc/credentials/insecure"
        pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

const (
        address = "localhost:50051"
)

func main() {
        listCmd := flag.NewFlagSet("list", flag.ExitOnError)
        getCmd := flag.NewFlagSet("show", flag.ExitOnError)
        updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
        createCmd := flag.NewFlagSet("create", flag.ExitOnError)

        // Get command flags
        getID := getCmd.String("id", "", "Repository ID")
        getName := getCmd.String("name", "", "Repository Name")

        // Update command flags
        updateID := updateCmd.String("id", "", "Repository ID")
        updateName := updateCmd.String("name", "", "Repository Name")
        updateDesc := updateCmd.String("desc", "", "Repository Description")

        // Create command flags
        createName := createCmd.String("name", "", "Repository Name")
        createDesc := createCmd.String("desc", "", "Repository Description")
        createGitignore := createCmd.String("gitignore", "", "Repository Gitignore")

        if len(os.Args) < 2 {
                fmt.Println("Usage: cli <command> [arguments]")
                fmt.Println("Commands: list, get, update, create")
                os.Exit(1)
        }

        conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
                log.Fatalf("did not connect: %v", err)
        }
        defer conn.Close()
        client := pb.NewRepositoryServiceClient(conn)
        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        switch os.Args[1] {
        case "list":
                listCmd.Parse(os.Args[2:])
                ListRepositories(ctx, client)
        case "show":
                getCmd.Parse(os.Args[2:])
                GetRepository(ctx, client, *getID, *getName)
        case "update":
                updateCmd.Parse(os.Args[2:])
                UpdateRepository(ctx, client, *updateID, *updateName, *updateDesc)
        case "create":
                createCmd.Parse(os.Args[2:])
                CreateRepository(ctx, client, *createName, *createDesc, *createGitignore)
        default:
                fmt.Println("Invalid command. Use: list, get, update, create")
                os.Exit(1)
        }
}

func ListRepositories(ctx context.Context, client pb.RepositoryServiceClient) {
        res, err := client.ListRepository(ctx, &pb.Empty{})
        if err != nil {
                log.Fatalf("failed to list repositories: %v", err)
        }
        for _, repo := range res.Repositories {
                fmt.Printf("ID: %s, Name: %s, Description: %s\n", repo.Id, repo.Name, repo.Description)
        }
}

func GetRepository(ctx context.Context, client pb.RepositoryServiceClient, id, name string) {
        res, err := client.GetRepository(ctx, &pb.GetRepositoryRequest{Id: id, Name: name})
        if err != nil {
                log.Fatalf("failed to get repository: %v", err)
        }
        fmt.Printf("ID: %s, Name: %s, Description: %s\n", res.Id, res.Name, res.Description)
}

func UpdateRepository(ctx context.Context, client pb.RepositoryServiceClient, id, name, desc string) {
        res, err := client.UpdateRepository(ctx, &pb.UpdateRepositoryRequest{Id: id, Name: name, Description: desc})
        if err != nil {
                log.Fatalf("failed to update repository: %v", err)
        }
        fmt.Printf("Updated Repository: ID: %s, Name: %s, Description: %s\n", res.Id, res.Name, res.Description)
}

func CreateRepository(ctx context.Context, client pb.RepositoryServiceClient, name, desc, gitignore string) {
        res, err := client.CreateRepository(ctx, &pb.CreateRepositoryRequest{Name: name, Description: desc, Gitignore: gitignore})
        if err != nil {
                log.Fatalf("failed to create repository: %v", err)
        }
        fmt.Printf("Created Repository: ID: %s, Name: %s, Description: %s\n", res.Id, res.Name, res.Description)
}