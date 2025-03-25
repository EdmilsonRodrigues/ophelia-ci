package main

import (
	"context"
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
	command := os.Args[2]
	args := os.Args[3:]

	switch service {
	case "repo":
		handleRepoCommands(ctx, repoClient, command, args)
	case "user":
		handleUserCommands(ctx, userClient, command, args)
	case "auth":
		handleAuthCommands(ctx, authClient, command, args)
	default:
		fmt.Println("Invalid service. Use: repo, user, auth")
		os.Exit(1)
	}
}

func ensureArgsLength(args []string, length int, message string) {
	if len(args) < length {
		fmt.Println(message)
		os.Exit(1)
	}
}
