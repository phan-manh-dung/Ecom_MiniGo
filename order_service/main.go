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

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	redis.InitRedis()

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
	// Đăng ký health check service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	// register order service
	order.RegisterOrderServiceServer(grpcServer, orderHandler)

	// Listen on port 50053
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	err = registerServiceWithConsul("order-service", 50053)
	if err != nil {
		panic(err)
	}

	fmt.Println("Order service gRPC server starting on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
