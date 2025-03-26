package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

// handleUserCommands parses command line arguments for the user command and makes the right call to the UserServiceClient.
// The commands available are:
// - list: Retrieves a list of all users
// - show: Retrieves a user by ID or username
// - create: Creates a new user
// - delete: Deletes a user by ID
func handleUserCommands(ctx context.Context, client pb.UserServiceClient, command string, args []string) {
	ctx = authenticateContext(ctx)
	switch command {
	case "list":
		ensureArgsLength(args, 0, "Too many arguments\nUsage: ophelia-ci user list")
		ListUsers(ctx, client)
	case "show":
		ensureArgsLength(args, 4, "Wrong number of arguments\nUsage: ophelia-ci user show --id <id>\nUsage: ophelia-ci user show --username <username>")
		getCmd := flag.NewFlagSet("show", flag.ExitOnError)
		getID := getCmd.String("id", "", "User ID")
		getUsername := getCmd.String("username", "", "User Username")
		getCmd.Parse(args)
		GetUser(ctx, client, *getID, *getUsername)
	case "create":
		ensureArgsLength(args, 4, "Wrong number of arguments\nUsage: ophelia-ci user create --username <username> --public-key <public-key>")
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		createUsername := createCmd.String("username", "", "User Username")
		createPublicKey := createCmd.String("public-key", "", "User Public Key")
		createCmd.Parse(args)
		CreateUser(ctx, client, *createUsername, *createPublicKey)
	case "delete":
		ensureArgsLength(args, 2, "Wrong number of arguments\nUsage: ophelia-ci user delete --id <id>")
	default:
		fmt.Println("Invalid user command")
		os.Exit(1)
	}
}

// ListUsers retrieves and prints a list of all users.
//
// This function sends a request to the UserServiceClient to list all
// existing users. It displays each user's ID and username. If there is an
// error during the request, the function logs the error and terminates the
// program.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The UserServiceClient used to access the user service.
func ListUsers(ctx context.Context, client pb.UserServiceClient) {
	res, err := client.ListUser(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Failed to list users: %v", err)
	}
	fmt.Println("Users:")
	for _, user := range res.Users {
		fmt.Printf("ID: %s, Username: %s\n", user.Id, user.Username)
	}
	fmt.Println("")
}

// GetUser retrieves and prints a user by either its ID or username.
//
// The request must contain either a non-empty ID or a non-empty username.
// The ID is used to identify the user to be retrieved by ID.
// The username is used to identify the user to be retrieved by username.
//
// The response will contain the user information.
// If there is an error during the request, the function logs the error and terminates the program.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The UserServiceClient used to access the user service.
// - id: The ID of the user to be retrieved.
// - name: The username of the user to be retrieved.
func GetUser(ctx context.Context, client pb.UserServiceClient, id, name string) {
	res, err := client.GetUser(ctx, &pb.GetUserRequest{Id: id, Username: name})
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Println("User:")
	fmt.Printf("ID: %s, Username: %s\n\n", res.Id, res.Username)
}

// CreateUser creates a new user with the given information.
//
// The request must contain the username and public key of the user to be created.
// The username is used to identify the user.
// The public key is used to store the user's public key.
//
// The response will contain the created user information.
// If there is an error during the request, the function logs the error and terminates the program.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The UserServiceClient used to access the user service.
// - username: The username of the user to be created.
// - publicKey: The path to the public key file of the user to be created.
func CreateUser(ctx context.Context, client pb.UserServiceClient, username, publicKey string) {
	publicKeyString, err := readPublicKey(publicKey)
	if err != nil {
		log.Fatalf("Failed to read public key: %v", err)
	}
	res, err := client.CreateUser(ctx, &pb.CreateUserRequest{Username: username, PublicKey: publicKeyString})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Println("User created:")
	fmt.Printf("ID: %s, Username: %s\n\n", res.Id, res.Username)
}

// UpdateUser updates a user with the given information.
//
// The request must contain the ID, username and public key of the user to be updated.
// The ID is used to identify the user to be updated.
// The username and public key are used to update the user information.
//
// The response will contain the updated user information.
// If there is an error during the request, the function logs the error and terminates the program.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The UserServiceClient used to access the user service.
// - id: The ID of the user to be updated.
// - username: The username of the user to be updated.
// - publicKey: The path to the public key file of the user to be updated.
func UpdateUser(ctx context.Context, client pb.UserServiceClient, id, username, publicKey string) {
	publicKeyString, err := readPublicKey(publicKey)
	if err != nil {
		log.Fatalf("Failed to read public key: %v", err)
	}
	res, err := client.UpdateUser(ctx, &pb.UpdateUserRequest{Id: id, Username: username, PublicKey: publicKeyString})
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
	fmt.Println("User updated:")
	fmt.Printf("ID: %s, Username: %s\n\n", res.Id, res.Username)
}

// DeleteUser deletes a user by ID.
//
// The request must contain the ID of the user to be deleted.
// If there is an error during the request, the function logs the error and terminates the program.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The UserServiceClient used to access the user service.
// - id: The ID of the user to be deleted.
func DeleteUser(ctx context.Context, client pb.UserServiceClient, id string) {
	_, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	fmt.Printf("User with ID: %s successfully deleted\n\n", id)
}

// readPublicKey reads the content of a public key file and returns it as a string.
//
// Parameters:
// - path: The path to the public key file.
//
// Returns:
// - string: The content of the public key file as a string.
// - error: An error if the file cannot be read.
func readPublicKey(path string) (string, error) {
	publicKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(publicKeyBytes), nil
}
