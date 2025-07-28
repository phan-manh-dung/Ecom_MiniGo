package main

import (
	"fmt"
	"log"

	"gin/api-gateway/handler"
	"gin/api-gateway/router"
	"gin/proto/generated/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	client user.UserServiceClient
}

func NewUserServiceClient() (*UserServiceClient, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	client := user.NewUserServiceClient(conn)
	return &UserServiceClient{client: client}, nil
}

func main() {
	// Create gRPC client
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user service: %v", err)
	}
	userGrpcClient := user.NewUserServiceClient(conn)
	userHandler := handler.NewUserServiceClient(userGrpcClient)

	r := gin.Default()

	// Register user routes
	router.RegisterUserRoutes(r, userHandler)

	// Start HTTP server
	fmt.Println("API Gateway starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
