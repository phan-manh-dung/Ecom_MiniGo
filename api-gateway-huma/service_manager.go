package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceManager quản lý tất cả gRPC clients
type ServiceManager struct {
	UserClient    interface{} // Placeholder cho user.UserServiceClient
	ProductClient interface{} // Placeholder cho product.ProductServiceClient
	OrderClient   interface{} // Placeholder cho order.OrderServiceClient
}

// getServiceAddressFromConsul lấy địa chỉ service từ Consul
func getServiceAddressFromConsul(serviceName string) (string, error) {
	config := api.DefaultConfig()
	consulAddr := os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		consulAddr = "localhost:8500"
	}
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("failed to create consul client: %v", err)
	}

	services, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil || len(services) == 0 {
		return "", fmt.Errorf("service %s not found in consul", serviceName)
	}

	svc := services[0].Service
	return fmt.Sprintf("%s:%d", svc.Address, svc.Port), nil
}

// waitForService chờ service có sẵn trong Consul
func waitForService(serviceName string, maxRetries int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to find service %s (attempt %d/%d)", serviceName, i+1, maxRetries)
		addr, err := getServiceAddressFromConsul(serviceName)
		if err == nil {
			log.Printf("Service %s found at %s", serviceName, addr)
			return addr, nil
		}
		log.Printf("Service %s not found yet: %v", serviceName, err)
		if i < maxRetries-1 {
			log.Printf("Waiting for service %s... (attempt %d/%d)", serviceName, i+1, maxRetries)
			time.Sleep(5 * time.Second)
		}
	}
	return "", fmt.Errorf("service %s not available after %d retries", serviceName, maxRetries)
}

// NewServiceManager tạo và kết nối tất cả services
func NewServiceManager() (*ServiceManager, error) {
	log.Printf("Initializing ServiceManager...")

	// Connect to User Service
	userAddr, err := waitForService("user-service", 12)
	if err != nil {
		return nil, fmt.Errorf("failed to get user service address: %v", err)
	}
	log.Printf("Connecting to user service at: %s", userAddr)
	userConn, err := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %v", err)
	}
	log.Printf("Successfully connected to user service")

	// Connect to Product Service
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

	// Connect to Order Service
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
		UserClient:    userConn,
		ProductClient: productConn,
		OrderClient:   orderConn,
	}, nil
}

// Close đóng tất cả connections
func (sm *ServiceManager) Close() {
	log.Printf("Closing all service connections...")
	// Note: gRPC connections are managed by the client and will be closed when the program exits
}
