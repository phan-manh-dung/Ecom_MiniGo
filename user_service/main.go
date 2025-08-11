package main

import (
	"context"
	"fmt"
	"gin/proto/generated/user"
	"gin/user_service/db"
	"gin/user_service/email"
	"gin/user_service/handler"
	userredis "gin/user_service/redis"
	"gin/user_service/repository"
	"gin/user_service/service"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/hashicorp/consul/api"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
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
	host := "user-service"

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

	// Connect to database
	db.ConnectDatabase()

	// Get port from environment variable, default to 50051
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	servicePort := 50051 // default
	if p, err := strconv.Atoi(port); err == nil {
		servicePort = p
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-19112.c52.us-east-1-4.ec2.redns.redis-cloud.com:19112",
		Username: "default",
		Password: "pA4GVyJVQTLUeCXNBsKnauUAtKQND7Jl",
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")

	// Initialize email service
	emailService := email.NewEmailService()

	// Initialize Redis subscriber
	subscriber := userredis.NewRedisSubscriber(redisClient, emailService)

	// Start Redis subscriber in goroutine
	go func() {
		subscriber.SubscribeToOrderEvents(ctx)
	}()

	// Initialize layers
	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Đăng ký health check service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Register user service
	user.RegisterUserServiceServer(grpcServer, userHandler)

	// Listen on port 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	//
	// Đăng ký Consul trước khi serve
	log.Printf("Attempting to register user-service with Consul on port %d", servicePort)
	err = registerServiceWithConsul("user-service", servicePort)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
	log.Printf("Successfully registered user-service with Consul")

	fmt.Printf("User service gRPC server starting on :%d\n", servicePort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
