package main

import (
	"fmt"
	"log"
	"os"
	"time"

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
	consulAddr := os.Getenv("CONSUL_ADDR")
	log.Printf("CONSUL_ADDR environment variable: '%s'", consulAddr)
	if consulAddr == "" {
		consulAddr = "localhost:8500"
		log.Printf("CONSUL_ADDR is empty, using default: %s", consulAddr)
	}
	config.Address = consulAddr
	log.Printf("Connecting to Consul at: %s", config.Address)

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

// waitForService waits for a service to be available in Consul
func waitForService(serviceName string, maxRetries int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to find service %s (attempt %d/%d)", serviceName, i+1, maxRetries)
		addr, err := getServiceAddressFromConsul(serviceName)
		if err == nil {
			log.Printf("Service %s found at %s", serviceName, addr)
			return addr, nil
		}
		log.Printf("Service %s not found yet: %v", serviceName, err)
		log.Printf("Waiting for service %s... (attempt %d/%d)", serviceName, i+1, maxRetries)
		time.Sleep(5 * time.Second)
	}
	return "", fmt.Errorf("service %s not available after %d retries", serviceName, maxRetries)
}

// NewServiceManager tạo và kết nối tất cả services
func NewServiceManager() (*ServiceManager, error) {
	log.Printf("Initializing ServiceManager...")

	userAddr, err := waitForService("user-service", 12) // Wait up to 1 minute
	if err != nil {
		return nil, fmt.Errorf("failed to get user service address: %v", err)
	}
	log.Printf("Connecting to user service at: %s", userAddr)
	userConn, err := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}
	log.Printf("Successfully connected to user service")

	productAddr, err := waitForService("product-service", 12)
	if err != nil {
		return nil, fmt.Errorf("failed to get product service address: %v", err)
	}
	log.Printf("Connecting to product service at: %s", productAddr)
	productConn, err := grpc.Dial(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %v", err)
	}
	log.Printf("Successfully connected to product service")

	orderAddr, err := waitForService("order-service", 12)
	if err != nil {
		return nil, fmt.Errorf("failed to get order service address: %v", err)
	}
	log.Printf("Connecting to order service at: %s", orderAddr)
	orderConn, err := grpc.Dial(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %v", err)
	}
	log.Printf("Successfully connected to order service")

	log.Printf("All services connected successfully")
	return &ServiceManager{
		UserClient:    user.NewUserServiceClient(userConn),
		ProductClient: product.NewProductServiceClient(productConn),
		OrderClient:   order.NewOrderServiceClient(orderConn),
	}, nil
}

func main() {
	log.Printf("Starting API Gateway...")

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
