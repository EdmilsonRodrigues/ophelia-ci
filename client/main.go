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

// Main is the entry point of the command-line client for Ophelia CI. It parses the
// command-line arguments and calls the corresponding service method.
//
// The client takes three arguments:
//
//   1. The service name (one of "repo", "user", "auth", or "signal").
//   2. The command name (service-specific).
//   3. The command arguments (service-specific).
//
// The client will print the result of the command as a YAML object and exit with a
// non-zero status code if the command fails.
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
	signalClient := pb.NewSignalsClient(conn)

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
	case "signal":
		handleSignals(ctx, signalClient, command, args)
	default:
		fmt.Println("Invalid service. Use: repo, user, auth")
		os.Exit(1)
	}
}

// ensureArgsLength checks if the provided arguments slice contains at least
// the specified number of elements. If not, it prints the provided message
// and exits the program with a non-zero status code.
func ensureArgsLength(args []string, length int, message string) {
	if len(args) < length {
		fmt.Println(message)
		os.Exit(1)
	}
}
