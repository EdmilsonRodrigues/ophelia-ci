package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = getSecret()
	uniqueKey = randomKey()
	noAuthNeededFunctions = map[string]bool{
		"/user.AuthService/AuthenticationChallenge": true,
		"/user.AuthService/Authentication": true,
		"/user.AuthService/UniqueKeyLogin": true,
		"/health.HealthService/Health": true,
	}
)
const (
	uniqueKeyExpirationDays = 1
	
)

// AuthInterceptor is a gRPC interceptor that verifies the JWT token sent
// by the client in the Authorization header. It skips authentication for
// methods that are used for authentication.
//
// The interceptor is called by gRPC for each unary RPC received by the
// server. It extracts the JWT token from the context, verifies it and
// returns an error if the token is invalid or missing. If the token is
// valid, it calls the handler function to process the RPC.
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	methodName := info.FullMethod

	if noAuthNeededFunctions[methodName] {
		log.Println("Skipping authentication for method:", methodName)
		return handler(ctx, req)
	}
	log.Println("Authenticating method:", methodName)

	_, err := extractAndVerifyJWT(ctx)
	if err != nil {
		log.Println("Error extracting and verifying JWT:", err)
		return nil, err
	}
	return handler(ctx, req)
}

// getSecret retrieves the server secret from the configuration. If no secret
// is defined in the configuration, it generates and returns a random key.
// This function ensures that a consistent secret is used for operations
// requiring cryptographic signing or verification.
func getSecret() string {
	config := LoadConfig()
	if config.Server.Secret != "" {
		return config.Server.Secret
	}
	return randomKey()
}

// randomKey generates a random 32-byte secret and returns it as a base64 encoded string.
// It uses the crypto/rand package for secure random number generation.
// If the random number generation fails, it logs a fatal error and exits the application.
func randomKey() string {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		log.Fatalf("Failed to generate secret: %v", err)
	}
	return base64.StdEncoding.EncodeToString(secret)
}

// AuthenticationChallenge generates a random challenge for the provided username,
// stores it in the server's challenges map, and returns it in the response. This
// challenge can be used by the client to prove knowledge of a private key.
//
// Parameters:
// - ctx: The context for the request, which carries deadlines, cancellation signals,
//   and other request-scoped values.
// - req: The request containing the username for which the challenge is to be generated.
//
// Returns:
// - *pb.AuthenticationChallengeResponse: The response containing the generated challenge.
// - error: An error if the challenge generation fails.
func (s *server) AuthenticationChallenge(ctx context.Context, req *pb.AuthenticationChallengeRequest) (*pb.AuthenticationChallengeResponse, error) {
	log.Printf("AuthenticationChallenge with request: %v", req)
	challengeBytes := make([]byte, 32)
	_, err := rand.Read(challengeBytes)
	if err != nil {
		log.Println("Error generating challenge:", err)
		return nil, err
	}
	challenge := base64.StdEncoding.EncodeToString(challengeBytes)
	s.challenges.Store(req.Username, challenge)
	return &pb.AuthenticationChallengeResponse{Challenge: challenge}, nil
}

// Authentication verifies the user's challenge response and returns a JWT token
// if the authentication is successful.
//
// Parameters:
// - ctx: The context for the request, which carries deadlines, cancellation signals,
//   and other request-scoped values.
// - req: The request containing the username and challenge response.
//
// Returns:
// - *pb.AuthenticationResponse: The response containing the generated JWT token
//   if the authentication is successful, or an error if the authentication fails.
func (s *server) Authentication(ctx context.Context, req *pb.AuthenticationRequest) (response *pb.AuthenticationResponse, err error) {
	config := LoadConfig()
	log.Printf("Authentication with request: %v", req)
	response = &pb.AuthenticationResponse{Authenticated: false}
	storedChallenge, ok := s.challenges.Load(req.Username)
	if !ok {
		log.Println("challenge not found")
		return
	}
	s.challenges.Delete(req.Username)

	storedPublicKey, err := s.userStore.GetPublicKeyByUsername(req.Username)
	if err != nil {
		log.Println("Error getting public key:", err)
		return
	}

	storedKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(storedPublicKey))
	if err != nil {
		log.Println("Error parsing stored public key:", err)
		return
	}

	signature, err := base64.StdEncoding.DecodeString(req.Challenge)
	if err != nil {
		log.Println("Error decoding signature:", err)
		return
	}

	challengeBytes, err := base64.StdEncoding.DecodeString(storedChallenge.(string))
	if err != nil {
		log.Println("Error decoding challenge:", err)
		return
	}

	if !verifySignature(storedKey, challengeBytes, signature) {
		log.Println("Signature verification failed")
		return
	}

	token, err := generateJWT(req.Username, config.Server.ExpirationTime)
	if err != nil {
		log.Println("Error generating JWT:", err)
		return &pb.AuthenticationResponse{Authenticated: false}, err
	}

	return &pb.AuthenticationResponse{Authenticated: true, Token: token}, nil
}

// UniqueKeyLogin logs in a user using the unique key thst is generated when the server is started, 
// and returns a JWT token if the login is successful.
// This is used for the initial login when the server is started.
//
// Parameters:
// - ctx: The context for the request, which carries deadlines, cancellation signals,
//   and other request-scoped values.
// - req: The request containing the unique key.
//
// Returns:
// - *pb.AuthenticationResponse: The response containing the generated JWT token
//   if the login is successful, or an error if the login fails.
func (s *server) UniqueKeyLogin(ctx context.Context, req *pb.UniqueKeyLoginRequest) (*pb.AuthenticationResponse, error) {
	log.Printf("UniqueKeyLogin with request: %v", req)
	if uniqueKey != "" && req.UniqueKey == uniqueKey {
		token, err := generateJWT(req.UniqueKey, uniqueKeyExpirationDays)
		if err != nil {
			log.Println("Error generating JWT:", err)
			return &pb.AuthenticationResponse{Authenticated: false}, err
		}
		uniqueKey = ""
		return &pb.AuthenticationResponse{Authenticated: true, Token: token}, nil
	}
	log.Println("Invalid unique key")
	return &pb.AuthenticationResponse{Authenticated: false}, nil
}

// verifySignature verifies the signature of a challenge using the stored public key.
//
// Parameters:
// - storedKey: The stored public key.
// - challengeBytes: The challenge bytes.
// - signatureBytes: The signature bytes.
//
// Returns:
// - bool: True if the verification is successful, false otherwise.
func verifySignature(storedKey ssh.PublicKey, challengeBytes, signatureBytes []byte) bool {
	log.Printf("Public key type: %s", storedKey.Type())
	log.Printf("Signature: %x", signatureBytes)
	h := sha256.Sum256(challengeBytes)

	signature := &ssh.Signature{
		Format: storedKey.Type(),
		Blob:   signatureBytes,
	}

	err := storedKey.Verify(h[:], signature)

	if err != nil {
		log.Println("Error verifying signature:", err)
		return false
	}

	return true
}

// generateJWT generates a JWT token for the given username that expires in the given number of days.
//
// Parameters:
// - username: The username for which the token is to be generated.
// - expirationDays: The number of days until the token expires.
//
// Returns:
// - string: The generated JWT token.
// - error: An error if the token generation fails.
func generateJWT(username string, expirationDays int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * time.Duration(expirationDays)).Unix(),
	})
	return token.SignedString([]byte(jwtSecret))
}

// extractTokenFromContext extracts the JWT token from the authorization header of the given context.
//
// Parameters:
// - ctx: The context for which the token is to be extracted.
//
// Returns:
// - string: The JWT token if it is present in the context, or an empty string otherwise.
// - bool: true if a token is present, and false otherwise.
func extractTokenFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) <= 0 {
		return "", false
	}
	if len(authHeader[0]) > 7 && authHeader[0][:7] == "Bearer " {
		return authHeader[0][7:], true
	}
	return "", false
}

// verifyJWT verifies a JWT token and returns the claims if the token is valid.
//
// Parameters:
// - tokenString: The JWT token to be verified.
//
// Returns:
// - jwt.MapClaims: The claims of the token if it is valid, or nil otherwise.
// - error: An error if the token is invalid.
func verifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// extractAndVerifyJWT extracts a JWT token from the authorization header of the given context,
// verifies it, and returns the username if the token is valid.
//
// Parameters:
// - ctx: The context for which the token is to be extracted.
//
// Returns:
// - string: The username if the token is valid, or an empty string otherwise.
// - error: An error if the token is invalid or missing.
func extractAndVerifyJWT(ctx context.Context) (username string, err error) {
	tokenString, ok := extractTokenFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("no token found")
	}
	log.Println("Token found:", tokenString)
	jwtMapClaims, err := verifyJWT(tokenString)
	if err != nil {
		return
	}
	username, ok = jwtMapClaims["username"].(string)
	if !ok {
		err = fmt.Errorf("invalid token")
	}
	return
}
