package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"gin/proto/generated/user"
	"gin/user_service/db"
	"gin/user_service/email"
	"gin/user_service/handler"
	userredis "gin/user_service/redis"
	"gin/user_service/repository"
	"gin/user_service/service"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
	// Connect to database
	db.ConnectDatabase()

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-19112.c52.us-east-1-4.ec2.redns.redis-cloud.com:19112",
		Username: "default",
		Password: "pA4GVyJVQTLUeCXNBsKnauUAtKQND7Jl",
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")

	// Initialize email service
	emailService := email.NewEmailService()

	// Initialize Redis subscriber
	subscriber := userredis.NewRedisSubscriber(redisClient, emailService)

	// Start Redis subscriber in goroutine
	go func() {
		subscriber.SubscribeToOrderEvents(ctx)
	}()

	// Initialize layers
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register user service
	user.RegisterUserServiceServer(grpcServer, userHandler)

	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("User service gRPC server starting on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
