package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

// handleRepoCommands is a function that parses command line arguments for the
// repository command and makes the right call to the RepositoryServiceClient.
// The commands available are:
// - list: Retrieves a list of all repositories
// - show: Retrieves a repository by ID or name
// - update: Updates a repository by ID
// - create: Creates a new repository
// - delete: Deletes a repository by ID
func handleRepoCommands(ctx context.Context, client pb.RepositoryServiceClient, command string, args []string) {
	ctx = authenticateContext(ctx)
	switch command {
	case "list":
		ensureArgsLength(args, 0, "Too many arguments\nUsage: ophelia-ci repo list")
		ListRepositories(ctx, client)
	case "show":
		ensureArgsLength(args, 2, "Wrong number of arguments\nUsage: ophelia-ci repo show --id <id>\nUsage: ophelia-ci repo show --name <name>")
		getCmd := flag.NewFlagSet("show", flag.ExitOnError)
		getID := getCmd.String("id", "", "Repository ID")
		getName := getCmd.String("name", "", "Repository Name")
		getCmd.Parse(args)
		GetRepository(ctx, client, *getID, *getName)
	case "update":
		ensureArgsLength(args, 6, "Wrong number of arguments\nUsage: ophelia-ci repo update --id <id> --name <name> --desc <desc>")
		updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
		updateID := updateCmd.String("id", "", "Repository ID")
		updateName := updateCmd.String("name", "", "Repository Name")
		updateDesc := updateCmd.String("desc", "", "Repository Description")
		updateCmd.Parse(args)
		UpdateRepository(ctx, client, *updateID, *updateName, *updateDesc)
	case "create":
		ensureArgsLength(args, 6, "Wrong number of arguments\nUsage: ophelia-ci repo create --name <name> --desc <desc> --gitignore <gitignore>")
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)
		createName := createCmd.String("name", "", "Repository Name")
		createDesc := createCmd.String("desc", "", "Repository Description")
		createGitignore := createCmd.String("gitignore", "", "Repository Gitignore")
		createCmd.Parse(args)
		CreateRepository(ctx, client, *createName, *createDesc, *createGitignore)
	case "delete":
		ensureArgsLength(args, 2, "Wrong number of arguments\nUsage: ophelia-ci repo delete --id <id>")
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		deleteID := deleteCmd.String("id", "", "Repository ID")
		deleteCmd.Parse(args)
		DeleteRepository(ctx, client, *deleteID)
	default:
		fmt.Println("Invalid repo command. Use: list, show, update, create, delete")
		os.Exit(1)
	}
}

// ListRepositories retrieves and prints a list of all repositories.
//
// This function sends a request to the RepositoryServiceClient to list all
// existing repositories. It displays each repository's ID, name, and
// description. If there is an error during the request, the function
// logs the error and terminates the program.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The RepositoryServiceClient used to access the repository service.
func ListRepositories(ctx context.Context, client pb.RepositoryServiceClient) {
	res, err := client.ListRepository(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("failed to list repositories: %v", err)
	}
	fmt.Println("Repositories:")
	for _, repo := range res.Repositories {
		fmt.Printf("ID: %s, Name: %s, Description: %s\n", repo.Id, repo.Name, repo.Description)
	}
	fmt.Println("")
}

// GetRepository retrieves a repository by either its ID or name.
//
// The request must contain either a non-empty ID or a non-empty name.
// The ID is used to identify the repository to be retrieved by ID.
// The name is used to identify the repository to be retrieved by name.
//
// The response will contain the repository information.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The RepositoryServiceClient used to access the repository service.
// - id: The ID of the repository to retrieve.
// - name: The name of the repository to retrieve.
func GetRepository(ctx context.Context, client pb.RepositoryServiceClient, id, name string) {
	if id == "" && name == "" {
		fmt.Println("Missing ID or Name")
		os.Exit(1)
		return
	}
	res, err := client.GetRepository(ctx, &pb.GetRepositoryRequest{Id: id, Name: name})
	if err != nil {
		log.Fatalf("failed to get repository: %v", err)
	}
	fmt.Println("Repository:")
	fmt.Printf("ID: %s, Name: %s, Description: %s\n\n", res.Id, res.Name, res.Description)
}

// UpdateRepository updates an existing repository with the given information.
//
// The request must contain the repository ID, name and description.
// The ID is used to identify the repository to be updated.
// The name and description are used to update the repository information.
//
// The response will contain the updated repository information.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The RepositoryServiceClient used to access the repository service.
// - id: The ID of the repository to update.
// - name: The name of the repository to update.
// - desc: The description of the repository to update.
func UpdateRepository(ctx context.Context, client pb.RepositoryServiceClient, id, name, desc string) {
	res, err := client.UpdateRepository(ctx, &pb.UpdateRepositoryRequest{Id: id, Name: name, Description: desc})
	if err != nil {
		log.Fatalf("failed to update repository: %v", err)
	}
	fmt.Printf("Updated Repository: ID: %s, Name: %s, Description: %s\n\n", res.Id, res.Name, res.Description)
}

// CreateRepository creates a new repository with the given information.
//
// The request must contain the repository name and description.
// The ID is generated by the server.
// The LastUpdate is set to the current timestamp by the server.
//
// The response will contain the created repository information.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The RepositoryServiceClient used to access the repository service.
// - name: The name of the repository to create.
// - desc: The description of the repository to create.
// - gitignore: The main language of the repository to create to generate the .gitignore file.
func CreateRepository(ctx context.Context, client pb.RepositoryServiceClient, name, desc, gitignore string) {
	if name == "" {
		fmt.Println("Missing Name")
		os.Exit(1)
		return
	}
	res, err := client.CreateRepository(ctx, &pb.CreateRepositoryRequest{Name: name, Description: desc, Gitignore: gitignore})
	if err != nil {
		log.Fatalf("failed to create repository: %v", err)
	}
	fmt.Printf("Created Repository: ID: %s, Name: %s, Description: %s\n\n", res.Id, res.Name, res.Description)
}

// DeleteRepository deletes a repository by its ID.
//
// This function sends a delete request to the RepositoryServiceClient using
// the provided ID. If the ID is empty, the function prints an error message
// and exits the program. If the deletion fails, it logs the error and
// terminates the program. Upon successful deletion, it prints a confirmation
// message with the repository ID.
//
// Parameters:
// - ctx: The context for the request, used for cancellation and deadlines.
// - client: The RepositoryServiceClient used to access the repository service.
// - id: The ID of the repository to be deleted.
func DeleteRepository(ctx context.Context, client pb.RepositoryServiceClient, id string) {
	if id == "" {
		fmt.Println("Missing ID")
		os.Exit(1)
		return
	}
	_, err := client.DeleteRepository(ctx, &pb.DeleteRepositoryRequest{Id: id})
	if err != nil {
		log.Fatalf("failed to delete repository: %v", err)
	}
	fmt.Printf("Deleted Repository with ID: %s\n", id)
}
