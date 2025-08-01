package main

import (
	"fmt"
	"log"

	"gin/api-gateway/handler"
	"gin/api-gateway/middleware"
	"gin/api-gateway/router"
	"gin/proto/generated/order"
	"gin/proto/generated/product"
	"gin/proto/generated/user"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceManager quản lý tất cả gRPC clients
type ServiceManager struct {
	UserClient    user.UserServiceClient
	ProductClient product.ProductServiceClient
	OrderClient   order.OrderServiceClient
}

func getServiceAddressFromConsul(serviceName string) (string, error) {
	config := api.DefaultConfig()
	config.Address = "localhost:8500"

	client, err := api.NewClient(config)
	if err != nil {
		return "", err
	}

	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil || len(services) == 0 {
		return "", fmt.Errorf("service not found")
	}

	svc := services[0].Service
	return fmt.Sprintf("%s:%d", svc.Address, svc.Port), nil
}

// NewServiceManager tạo và kết nối tất cả services
func NewServiceManager() (*ServiceManager, error) {
	userAddr, err := getServiceAddressFromConsul("user-service")
	if err != nil {
		return nil, fmt.Errorf("failed to get user service address: %v", err)
	}
	userConn, err := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}

	productAddr, err := getServiceAddressFromConsul("product-service")
	if err != nil {
		return nil, fmt.Errorf("failed to get product service address: %v", err)
	}
	productConn, err := grpc.Dial(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %v", err)
	}

	orderAddr, err := getServiceAddressFromConsul("order-service")
	if err != nil {
		return nil, fmt.Errorf("failed to get order service address: %v", err)
	}
	orderConn, err := grpc.Dial(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %v", err)
	}

	return &ServiceManager{
		UserClient:    user.NewUserServiceClient(userConn),
		ProductClient: product.NewProductServiceClient(productConn),
		OrderClient:   order.NewOrderServiceClient(orderConn),
	}, nil
}

func main() {
	// Khởi tạo service manager
	serviceManager, err := NewServiceManager()
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}

	// Tạo handlers
	userHandler := handler.NewUserServiceClient(serviceManager.UserClient)
	productHandler := handler.NewProductServiceClient(serviceManager.ProductClient)
	orderHandler := handler.NewOrderServiceClient(serviceManager.OrderClient)

	// Khởi tạo Gin router
	r := gin.New()

	// Áp dụng middleware theo thứ tự
	r.Use(middleware.CORSMiddleware())      // CORS middleware
	r.Use(middleware.RequestIDMiddleware()) // Request ID middleware
	r.Use(middleware.LoggingMiddleware())   // Logging middleware
	r.Use(middleware.AuthMiddleware())      // Authentication middleware

	// Đăng ký routes
	router.RegisterUserRoutes(r, userHandler)
	router.RegisterProductRoutes(r, productHandler)
	router.RegisterOrderRoutes(r, orderHandler)

	// Khởi động HTTP server
	fmt.Println("API Gateway starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
