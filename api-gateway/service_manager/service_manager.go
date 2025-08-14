package service_manager

import (
	"fmt"
	"gin/proto/generated/order"
	"gin/proto/generated/product"
	"gin/proto/generated/user"
	"log"
	"os"
	"time"

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
	if err != nil {
		return "", fmt.Errorf("failed to query consul: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("service %s not found in consul", serviceName)
	}

	service := services[0]
	return fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port), nil
}

// waitForService chờ service có sẵn trong Consul
func waitForService(serviceName string, maxRetries int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		addr, err := getServiceAddressFromConsul(serviceName)
		if err == nil {
			log.Printf("Service %s found at %s", serviceName, addr)
			return addr, nil
		}
		log.Printf("Waiting for service %s... (attempt %d/%d)", serviceName, i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	return "", fmt.Errorf("service %s not available after %d retries", serviceName, maxRetries)
}

// NewServiceManager tạo instance mới của ServiceManager
func NewServiceManager() (*ServiceManager, error) {
	log.Printf("Initializing ServiceManager with Consul discovery...")

	// Kết nối User Service trực tiếp
	userAddr := "localhost:50051"
	log.Printf("Connecting to User Service at %s", userAddr)

	userConn, err := grpc.Dial(userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: Failed to connect to User Service at %s: %v", userAddr, err)
		log.Printf("User Service will not be available")
	} else {
		log.Printf("Connected to User Service at %s", userAddr)
	}

	// Kết nối Product Service trực tiếp
	productAddr := "localhost:60051"
	log.Printf("Connecting to Product Service at %s", productAddr)

	productConn, err := grpc.Dial(productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: Failed to connect to Product Service at %s: %v", productAddr, err)
		log.Printf("Product Service will not be available")
	} else {
		log.Printf("Connected to Product Service at %s", productAddr)
	}

	// Kết nối Order Service trực tiếp
	orderAddr := "localhost:40051"
	log.Printf("Connecting to Order Service at %s", orderAddr)

	orderConn, err := grpc.Dial(orderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: Failed to connect to Order Service at %s: %v", orderAddr, err)
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
