package main

import (
	"context"
	"flag"
	"log"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
)

// handleSignals parses command line arguments for the signals command and makes the right call to the SignalsClient.
// The commands available are:
// - commit: Sends a commit signal to the server with the given commit hash, branch, repository and tag.
func handleSignals(ctx context.Context, client pb.SignalsClient, command string, args []string) {
	ctx = authenticateContext(ctx)
	switch command {
	case "commit":
		ensureArgsLength(args, 8, "Too many arguments\nUsage: ophelia-ci repo list")
		getCmd := flag.NewFlagSet("commit", flag.ExitOnError)
		getCommitHash := getCmd.String("hash", "", "Commit Hash")
		getBranch := getCmd.String("branch", "", "Branch")
		getRepositoryName := getCmd.String("repo", "", "Repository")
		getTag := getCmd.String("tag", "", "Tag")
		getCmd.Parse(args)
		SendCommitSignal(ctx, client, *getRepositoryName, *getCommitHash, *getBranch, *getTag)
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

// SendCommitSignal sends a commit signal to the server with the given commit hash, branch, repository and tag.
// If there is an error, it will log the error and exit.
func SendCommitSignal(ctx context.Context, client pb.SignalsClient, repo, hash, branch, tag string) {
	_, err := client.CommitSignal(ctx, &pb.CommitRequest{Repository: repo, CommitHash: hash, Branch: branch, Tag: tag})
	if err != nil {
		log.Fatalf("failed to send commit signal: %v", err)
	}
	
}
