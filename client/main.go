package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := LoadConfig()

	conn, err := grpc.NewClient(config.Client.Server, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	repoClient := pb.NewRepositoryServiceClient(conn)
	userClient := pb.NewUserServiceClient(conn)
	authClient := pb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if len(os.Args) < 2 {
		fmt.Println("Missing arguments")
		fmt.Println("Usage: ophelia-ci <service>")
		fmt.Println("Services: repo, user, auth")
		os.Exit(1)
	}

	service := os.Args[1]
	commandAndArgs := os.Args[2:]

	switch service {
	case "repo":
		handleRepoCommands(ctx, repoClient, commandAndArgs)
	case "user":
		handleUserCommands(ctx, userClient, commandAndArgs)
	case "auth":
		handleAuthCommands(ctx, authClient, commandAndArgs)
	default:
		fmt.Println("Invalid service. Use: repo, user, auth")
		os.Exit(1)
	}
}

func handleRepoCommands(ctx context.Context, client pb.RepositoryServiceClient, command string, args []string) {
	switch command {
	case "list":
		ListRepositories(ctx, client)
	case "show":
		getCmd := flag.NewFlagSet("show", flag.ExitOnError)
		getID := getCmd.String("id", "", "Repository ID")
		getName := getCmd.String("name", "", "Repository Name")
		getCmd.Parse(args)
		GetRepository(ctx, client, *getID, *getName)
	case "update":
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		updateID := updateCmd.String("id", "", "Repository ID")
		updateName := updateCmd.String("name", "", "Repository Name")
		updateDesc := updateCmd.String("desc", "", "Repository Description")
		updateCmd.Parse(args)
		UpdateRepository(ctx, client, *updateID, *updateName, *updateDesc)
	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		createName := createCmd.String("name", "", "Repository Name")
		createDesc := createCmd.String("desc", "", "Repository Description")
		createGitignore := createCmd.String("gitignore", "", "Repository Gitignore")
		createCmd.Parse(args)
		CreateRepository(ctx, client, *createName, *createDesc, *createGitignore)
	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		deleteID := deleteCmd.String("id", "", "Repository ID")
		deleteCmd.Parse(args)
		DeleteRepository(ctx, client, *deleteID)
	default:
		fmt.Println("Invalid repo command. Use: list, show, update, create, delete")
		os.Exit(1)
	}
}

func handleUserCommands(ctx context.Context, client pb.UserServiceClient, command string, args []string) {
	// Implement user service command handling
	switch command {
	case "list":
		ListUsers(ctx, client)
	case "show":
		// ...
	case "create":
		// ...
	case "delete":
		// ...
	default:
		fmt.Println("Invalid user command")
		os.Exit(1)
	}
}

func handleAuthCommands(ctx context.Context, client pb.AuthServiceClient, command string, args []string) {
	// Implement auth service command handling
	switch command {
	case "login":
		// ...
	case "unique":
		// ...
	default:
		fmt.Println("Invalid auth command")
		os.Exit(1)
	}
}
