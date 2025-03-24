package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

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