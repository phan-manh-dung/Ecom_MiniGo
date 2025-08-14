package service_manager

import (
	"gin/proto/generated/order"
	"gin/proto/generated/product"
	"gin/proto/generated/user"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceManager quản lý tất cả gRPC clients
type ServiceManager struct {
	UserClient    user.UserServiceClient
	ProductClient product.ProductServiceClient
	OrderClient   order.OrderServiceClient
}

// NewServiceManager tạo instance mới của ServiceManager
func NewServiceManager() (*ServiceManager, error) {
	log.Printf("Initializing ServiceManager...")

	// Kết nối User Service
	userAddr := os.Getenv("USER_SERVICE_ADDR")
	if userAddr == "" {
		userAddr = "localhost:50051"
	}
	userConn, err := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: Failed to connect to User Service: %v", err)
		log.Printf("User Service will not be available")
	} else {
		log.Printf("Connected to User Service at %s", userAddr)
	}

	// Kết nối Product Service
	productAddr := os.Getenv("PRODUCT_SERVICE_ADDR")
	if productAddr == "" {
		productAddr = "localhost:60051"
	}
	productConn, err := grpc.Dial(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: Failed to connect to Product Service: %v", err)
		log.Printf("Product Service will not be available")
	} else {
		log.Printf("Connected to Product Service at %s", productAddr)
	}

	// Kết nối Order Service
	orderAddr := os.Getenv("ORDER_SERVICE_ADDR")
	if orderAddr == "" {
		orderAddr = "localhost:40051"
	}
	orderConn, err := grpc.Dial(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: Failed to connect to Order Service: %v", err)
		log.Printf("Order Service will not be available")
	} else {
		log.Printf("Connected to Order Service at %s", orderAddr)
	}

	// Tạo clients
	var userClient user.UserServiceClient
	var productClient product.ProductServiceClient
	var orderClient order.OrderServiceClient

	if userConn != nil {
		userClient = user.NewUserServiceClient(userConn)
	}
	if productConn != nil {
		productClient = product.NewProductServiceClient(productConn)
	}
	if orderConn != nil {
		orderClient = order.NewOrderServiceClient(orderConn)
	}

	return &ServiceManager{
		UserClient:    userClient,
		ProductClient: productClient,
		OrderClient:   orderClient,
	}, nil
}
