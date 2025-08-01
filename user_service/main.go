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

	"github.com/hashicorp/consul/api"
	"github.com/redis/go-redis/v9"
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
	// Connect to database
	db.ConnectDatabase()

	addr := os.Getenv("Addr")
	username := os.Getenv("Username")
	password := os.Getenv("Password")

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
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
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Đăng ký Consul trước khi serve
	err = registerServiceWithConsul("user-service", 50051)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	fmt.Println("User service gRPC server starting on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
