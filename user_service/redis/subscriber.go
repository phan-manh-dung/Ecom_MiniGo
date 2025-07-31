package redis

import (
	"context"
	"encoding/json"
	"gin/user_service/email"
	"log"

	"github.com/redis/go-redis/v9"
)

type OrderCancelledEvent struct {
	OrderID   uint32 `json:"order_id"`
	UserID    uint32 `json:"user_id"`
	UserEmail string `json:"user_email"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type RedisSubscriber struct {
	client       *redis.Client
	emailService *email.EmailService
}

func NewRedisSubscriber(client *redis.Client, emailService *email.EmailService) *RedisSubscriber {
	return &RedisSubscriber{
		client:       client,
		emailService: emailService,
	}
}

// SubscribeToOrderEvents subscribe các events liên quan đến order
func (r *RedisSubscriber) SubscribeToOrderEvents(ctx context.Context) {
	pubsub := r.client.Subscribe(ctx, "order.cancelled")
	defer pubsub.Close()

	log.Println("User Service: Subscribed to order.cancelled events")

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			continue
		}

		log.Printf("Received event: %s", msg.Payload)
		r.handleOrderCancelledEvent(ctx, msg.Payload)
	}
}

// handleOrderCancelledEvent xử lý event order cancelled
func (r *RedisSubscriber) handleOrderCancelledEvent(ctx context.Context, payload string) {
	var event OrderCancelledEvent
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		return
	}

	log.Printf("Processing order cancelled event for order %d, user %d", event.OrderID, event.UserID)

	// Gửi email thông báo
	if err := r.emailService.SendOrderCancelledEmail(event.UserEmail, event.OrderID); err != nil {
		log.Printf("Failed to send email for order %d: %v", event.OrderID, err)
		return
	}

	log.Printf("Email sent successfully for order %d to %s", event.OrderID, event.UserEmail)
}
