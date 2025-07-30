package main

import (
	"fmt"
	"gin/order_service/db"
	"gin/order_service/handler"
	"gin/order_service/repository"
	"gin/order_service/service"
	"gin/proto/generated/order"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	db.ConnectDatabase()

	// init layer
	orderRepo := repository.NewOrderRepository(db.DB)

	// Tạo gRPC client cho Product Service
	productClient, err := service.NewProductServiceClient("localhost:50052")
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}

	orderService := service.NewOrderService(orderRepo, productClient)
	orderHandler := handler.NewOrderHandler(orderService)

	// create gRPC server
	grpcServer := grpc.NewServer()
	// register order service
	order.RegisterOrderServiceServer(grpcServer, orderHandler)

	// Listen on port 50053
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("Order service gRPC server starting on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
