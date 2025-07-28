package main

import (
	"fmt"
	"log"
	"net"

	"gin/proto/generated/user"
	"gin/user_service/db"
	"gin/user_service/handler"

	"google.golang.org/grpc"
)

func main() {
	// Connect to database
	db.ConnectDatabase()

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register user service
	userHandler := handler.NewUserHandler(db.DB)
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
