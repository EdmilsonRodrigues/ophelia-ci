package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

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

func GetUser(ctx context.Context, client pb.UserServiceClient, id, name string) {
	res, err := client.GetUser(ctx, &pb.GetUserRequest{Id: id, Username: name})
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Println("User:")
	fmt.Printf("ID: %s, Username: %s\n\n", res.Id, res.Username)
}

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

func DeleteUser(ctx context.Context, client pb.UserServiceClient, id string) {
	_, err := client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	fmt.Printf("User with ID: %s successfully deleted\n\n", id)
}

func readPublicKey(path string) (string, error) {
	publicKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(publicKeyBytes), nil
}
