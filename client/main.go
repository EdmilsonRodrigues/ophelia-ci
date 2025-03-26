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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if len(os.Args) < 2 {
		printOpheliaHelp()
		os.Exit(1)
	}

	service := os.Args[1]

	if len(os.Args) < 3 {
		handleCommands(ctx, conn, service, "--help", []string{})
		os.Exit(1)
	}

	handleCommands(ctx, conn, service, os.Args[2], os.Args[3:])
}

func printOpheliaHelp() {
	fmt.Println("Usage: ophelia-ci <service> <command> [arguments]")
	fmt.Println("Services:")
	fmt.Println("	repo	Repository service")
	fmt.Println("	user	User service")
	fmt.Println("	auth	Authentication service")
}

func printHelp(service string) {
	switch service {
	case "repo":
		printRepoHelp()
	case "user":
		printUserHelp()
	case "auth":
		printAuthHelp()
	default:
		printOpheliaHelp()
	}
}

func handleCommands(ctx context.Context, conn *grpc.ClientConn, service, command string, args []string) {
	repoClient := pb.NewRepositoryServiceClient(conn)
	userClient := pb.NewUserServiceClient(conn)
	authClient := pb.NewAuthServiceClient(conn)
	signalClient := pb.NewSignalsClient(conn)

	switch service {
	case "--help":
		printOpheliaHelp()
	case "help":
		printHelp(command)
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
		fmt.Println(message + "\n")
		os.Exit(1)
	}
}
