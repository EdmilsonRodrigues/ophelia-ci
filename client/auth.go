package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc/metadata"
)


// handleAuthCommands handles authentication commands.
//
// The available commands are "login" and "unique". The "login" command
// takes a username and private key as arguments, and performs a login
// operation with the server. The "unique" command takes a unique key as
// an argument, and performs a login operation with the unique key.
//
// If the login is successful, the function sets the
// OPHELIA_CI_CLIENT_TOKEN environment variable to the token returned by the
// server, and prints a success message. If the login fails, the function
// prints an error message and exits with a non-zero status code.
func handleAuthCommands(ctx context.Context, client pb.AuthServiceClient, command string, args []string) {
	var token string
	var err error
	switch command {
	case "--help":
		printAuthHelp()
		os.Exit(0)
	case "login":
		ensureArgsLength(args, 4, "Wrong number of arguments\nUsage: ophelia-ci auth login --username <username> --private-key <private-key>")

		loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
		username := loginCmd.String("username", "", "Username")
		privateKey := loginCmd.String("private-key", "", "Private Key")
		loginCmd.Parse(args)

		token, err = login(ctx, client, *username, *privateKey)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}

	case "unique":
		ensureArgsLength(args, 2, "Wrong number of arguments\nUsage: ophelia-ci auth unique --key <unique-key>")

		uniqueCmd := flag.NewFlagSet("unique", flag.ExitOnError)
		uniqueKey := uniqueCmd.String("key", "", "Unique Key")
		uniqueCmd.Parse(args)

		token, err = uniqueKeyLogin(ctx, client, *uniqueKey)
		if err != nil {
			log.Fatalf("Unique key login failed: %v", err)
		}

	default:
		log.Fatalln("Invalid auth command. Use 'login' or 'unique'.")
	}
	if token == "" {
		log.Fatalf("Login failed")
	}
	log.Printf("Token: %s\n", token)

	setToken(token)
	fmt.Println("Login successful")
}

func printAuthHelp() {
	fmt.Println("Usage: ophelia-ci auth <command>")
	fmt.Println("Commands:")
	fmt.Println("	login	Authenticate a user using their username and private key")
	fmt.Println("	unique	Authenticate a user using the server unique key")
}

func setToken(token string) {
	config := LoadConfig()
	config.Client.AuthToken = token
	SaveConfig(config)
}

// login authenticates a user using their username and private key, and returns a JWT token if the authentication is successful.
//
// Parameters:
//   - ctx: The context for the request, which carries deadlines, cancellation signals,
//     and other request-scoped values.
//   - client: The client to use for the authentication request.
//   - username: The username of the user to authenticate.
//   - privateKey: The path to the private key file of the user to authenticate.
//
// Returns:
// - string: The JWT token if the authentication is successful, or an empty string if the authentication fails.
// - error: An error if the authentication fails.
func login(ctx context.Context, client pb.AuthServiceClient, username, privateKey string) (string, error) {
	if username == "" || privateKey == "" {
		return "", fmt.Errorf("username and private key are required")
	}

	privateKeyBytes, err := os.ReadFile(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to read private key: %w", err)
	}

	privateKeyObj, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	challengeResponse, err := client.AuthenticationChallenge(ctx, &pb.AuthenticationChallengeRequest{Username: username})
	if err != nil {
		return "", fmt.Errorf("failed to get authentication challenge: %w", err)
	}

	challengeBytes, err := base64.StdEncoding.DecodeString(challengeResponse.Challenge)
	if err != nil {
		return "", fmt.Errorf("failed to decode challenge: %w", err)
	}

	h := sha256.Sum256(challengeBytes)
	signature, err := privateKeyObj.Sign(rand.Reader, h[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign challenge: %w", err)
	}
	signatureBase64 := base64.StdEncoding.EncodeToString(signature.Blob)

	authResponse, err := client.Authentication(ctx, &pb.AuthenticationRequest{Username: username, Challenge: signatureBase64})
	if err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	if !authResponse.Authenticated {
		return "", fmt.Errorf("authentication failed")
	}

	return authResponse.Token, nil
}

// uniqueKeyLogin authenticates a user using a unique key, and returns a JWT token if the authentication is successful.
//
// Parameters:
//   - ctx: The context for the request, which carries deadlines, cancellation signals,
//     and other request-scoped values.
//   - client: The client to use for the authentication request.
//   - uniqueKey: The unique key for authentication.
//
// Returns:
// - string: The JWT token if the authentication is successful, or an empty string if the authentication fails.
// - error: An error if the authentication fails.
func uniqueKeyLogin(ctx context.Context, client pb.AuthServiceClient, uniqueKey string) (string, error) {
	if uniqueKey == "" {
		return "", fmt.Errorf("unique key is required")
	}

	authResponse, err := client.UniqueKeyLogin(ctx, &pb.UniqueKeyLoginRequest{UniqueKey: uniqueKey})
	if err != nil {
		return "", fmt.Errorf("unique key login failed: %w", err)
	}

	if !authResponse.Authenticated {
		return "", fmt.Errorf("unique key authentication failed")
	}

	return authResponse.Token, nil
}

// getToken retrieves the JWT token from the environment variable defined by tokenVariable.
//
// Returns:
// - string: The JWT token if it is present in the environment variable, or an empty string if it is not.
// - bool: true if a token is present, and false otherwise.
func getToken() (string, bool) {
	config := LoadConfig()
	if config.Client.AuthToken == "" {
		return "", false
	}
	return config.Client.AuthToken, true
}

// authenticateContext adds a JWT token to the outgoing context's metadata for authorization.
//
// Parameters:
//   - ctx: The context to which the JWT token is to be added. It carries deadlines,
//     cancellation signals, and other request-scoped values.
//
// Returns:
// - context.Context: The context with the JWT token added to its metadata.
//
// If no token is found in the environment variable, the function prints an error
// message and exits the program.
func authenticateContext(ctx context.Context) context.Context {
	token, ok := getToken()
	if !ok {
		fmt.Println("No token found")
		os.Exit(1)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
	return ctx
}
