package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func (u *UserServiceClient) GetUser(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Call gRPC service
	req := &user.GetUserRequest{
		Id: uint32(userID),
	}

	resp, err := u.client.GetUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"user":    resp.User,
		"message": resp.Message,
	})
}

func main() {
	// Create user service client
	userClient, err := NewUserServiceClient()
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	// Setup Gin router
	r := gin.Default()

	// User routes
	userRoutes := r.Group("/api/users")
	{
		userRoutes.GET("/:id", userClient.GetUser)
	}

	// Start HTTP server
	fmt.Println("API Gateway starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
