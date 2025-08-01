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
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func registerServiceWithConsul(serviceName string, servicePort int) error {
	config := api.DefaultConfig()
	config.Address = "localhost:8500"

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	host := "localhost"

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

	// Get port from environment variable, default to 50051
	port := os.Getenv("PORT")
	if port == "" {
		port = "60051"
	}

	servicePort := 60051 // default
	if p, err := strconv.Atoi(port); err == nil {
		servicePort = p
	}

	// init layer
	productRepo := repository.NewUserRepository(db.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// create gRPC server
	grpcServer := grpc.NewServer()
	// Đăng ký health check service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	// register product service
	product.RegisterProductServiceServer(grpcServer, productHandler)

	// Listen on port 60051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	err = registerServiceWithConsul("product-service", servicePort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	fmt.Println("Product service gRPC server starting on :60051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
