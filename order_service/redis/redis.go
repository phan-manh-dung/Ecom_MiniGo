package redis

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// InitRedis khởi tạo Redis client
func InitRedis() {

	// Get Redis address from environment variable, default to localhost
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // No password for local Redis
		DB:       0,
	})

	// Test connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}

// [Tạo PublishOrderCancelled publish event] khi order bị hủy
func PublishOrderCancelled(ctx context.Context, orderID uint32, userID uint32, userEmail string) error {
	event := map[string]interface{}{
		"order_id":   orderID,
		"user_id":    userID,
		"user_email": userEmail,
		"status":     "CANCELLED",
		"timestamp":  "2024-01-01T00:00:00Z",
	}

	// Convert to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Publish event
	err = RedisClient.Publish(ctx, "order.cancelled", eventJSON).Err()
	if err != nil {
		return err
	}

	log.Printf("Published order.cancelled event for order %d", orderID)
	return nil
}
