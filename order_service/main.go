package main

import (
	"fmt"
	"gin/order_service/db"
	"gin/order_service/handler"
	"gin/order_service/repository"
	"gin/order_service/service"
	"gin/proto/generated/order"
	"gin/proto/generated/product"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	db.ConnectDatabase()

	// init layer
	orderRepo := repository.NewOrderRepository(db.DB)

	// Tạo gRPC client cho Product Service
	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	productClient := product.NewProductServiceClient(productConn)

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
