package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/EdmilsonRodrigues/ophelia-ci"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = getSecret()
var uniqueKey = randomKey()

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	methodName := info.FullMethod

	if strings.Contains(methodName, "AuthenticationChallange") || strings.Contains(methodName, "Authentication") {
		return handler(ctx, req)
	}

	_, err := extractAndVerifyJWT(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func getSecret() string {
	config := LoadConfig()
	if config.Server.Secret != "" {
		return config.Server.Secret
	}
	return randomKey()
}

func randomKey() string {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		log.Fatalf("Failed to generate secret: %v", err)
	}
	return base64.StdEncoding.EncodeToString(secret)
}

func (s *server) AuthenticationChallange(ctx context.Context, req *pb.AuthenticationChallengeRequest) (*pb.AuthenticationChallengeResponse, error) {
	log.Printf("AuthenticationChallange with request: %v", req)
	challengeBytes := make([]byte, 32)
	_, err := rand.Read(challengeBytes)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	challenge := base64.StdEncoding.EncodeToString(challengeBytes)
	s.challenges.Store(req.Username, challenge)
	return &pb.AuthenticationChallengeResponse{Challenge: challenge}, nil
}

func (s *server) Authentication(ctx context.Context, req *pb.AuthenticationRequest) (response *pb.AuthenticationResponse, err error) {
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

	token, err := generateJWT(req.Username)
	if err != nil {
		return &pb.AuthenticationResponse{Authenticated: false}, err
	}

	return &pb.AuthenticationResponse{Authenticated: true, Token: token}, nil
}

func (s *server) UniqueKeyLogin(ctx context.Context, req *pb.UniqueKeyLoginRequest) (*pb.AuthenticationResponse, error) {
	log.Printf("UniqueKeyLogin with request: %v", req)
	if uniqueKey != "" && req.UniqueKey == uniqueKey {
		token, err := generateJWT(req.UniqueKey)
		if err != nil {
			return &pb.AuthenticationResponse{Authenticated: false}, err
		}
		uniqueKey = ""
		return &pb.AuthenticationResponse{Authenticated: true, Token: token}, nil
	}
	return &pb.AuthenticationResponse{Authenticated: false}, nil
}

func verifySignature(storedKey ssh.PublicKey, challengeBytes, signatureBytes []byte) bool {
	h := sha256.Sum256(challengeBytes)

	signature := &ssh.Signature{
		Format: storedKey.Type(),
		Blob:   signatureBytes,
	}

	err := storedKey.Verify(h[:], signature)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}

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

func verifyJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func extractAndVerifyJWT(ctx context.Context) (username string, err error) {
	tokenString, ok := extractTokenFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("no token found")
	}
	jwtMapClaims, err := verifyJWT(tokenString)
	if err != nil {
		return
	}
	username, ok = jwtMapClaims["username"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}
	return
}
