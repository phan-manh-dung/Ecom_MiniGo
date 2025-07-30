package main

import (
	"fmt"
	"gin/product_service/db"
	"gin/product_service/handler"
	"gin/product_service/repository"
	"gin/product_service/service"
	"gin/proto/generated/product"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	db.ConnectDatabase()

	// init layer
	productRepo := repository.NewUserRepository(db.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// create gRPC server
	grpcServer := grpc.NewServer()
	// register product service
	product.RegisterProductServiceServer(grpcServer, productHandler)

	// Listen on port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("User service gRPC server starting on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
