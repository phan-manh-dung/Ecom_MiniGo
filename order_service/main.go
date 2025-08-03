package main

import (
	"fmt"
	"gin/order_service/db"
	"gin/order_service/handler"
	"gin/order_service/redis"
	"gin/order_service/repository"
	"gin/order_service/service"
	"gin/proto/generated/order"
	"gin/proto/generated/product"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func registerServiceWithConsul(serviceName string, servicePort int) error {
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
		return err
	}

	// Sử dụng container name thay vì localhost
	host := "order-service"

	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%d", serviceName, servicePort),
		Name:    serviceName,
		Address: host,
		Port:    servicePort,
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", host, servicePort),
			Interval:                       "10s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	return client.Agent().ServiceRegister(registration)
}

func main() {
	db.ConnectDatabase()
	redis.InitRedis()

	port := os.Getenv("PORT")
	if port == "" {
		port = "40051"
	}

	servicePort := 40051 // default
	if p, err := strconv.Atoi(port); err == nil {
		servicePort = p
	}

	// init layer
	orderRepo := repository.NewOrderRepository(db.DB)

	// Tạo gRPC client cho Product Service
	productServiceAddr := os.Getenv("PRODUCT_SERVICE_ADDR")
	if productServiceAddr == "" {
		productServiceAddr = "product-service:60051"
	}
	productConn, err := grpc.Dial(productServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	productClient := product.NewProductServiceClient(productConn)

	orderService := service.NewOrderService(orderRepo, productClient)
	orderHandler := handler.NewOrderHandler(orderService)

	// create gRPC server
	grpcServer := grpc.NewServer()
	// Đăng ký health check service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	// register order service
	order.RegisterOrderServiceServer(grpcServer, orderHandler)

	// Listen on port 40051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Attempting to register order-service with Consul on port %d", servicePort)
	err = registerServiceWithConsul("order-service", servicePort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
	log.Printf("Successfully registered order-service with Consul")

	fmt.Println("Order service gRPC server starting on :40051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
