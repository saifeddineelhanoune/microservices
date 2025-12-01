package pubsub

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: "", 
		DB: 0,        
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis for Payment Service!")
}

type OrderEvent struct {
	OrderID int     `json:"order_id"`
	Amount  float64 `json:"amount"`
}

func StartConsumer() {
	pubsub := RedisClient.Subscribe(Ctx, "order_events")
	defer pubsub.Close()

	log.Println("Starting Redis consumer for 'order_events' channel...")

	for {
		msg, err := pubsub.ReceiveMessage(Ctx)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		log.Printf("Received message from %s: %s", msg.Channel, msg.Payload)

		var event OrderEvent
		if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
			log.Printf("Error unmarshalling event payload: %v", err)
			continue
		}

		processPayment(event)
	}
}

func processPayment(event OrderEvent) {
	log.Printf("Processing payment for Order ID: %d, Amount: %.2f", event.OrderID, event.Amount)

	// --- Payment Simulation Logic ---
	// In a real system, this would involve calling a third-party payment gateway API.
	time.Sleep(2 * time.Second) // Simulate network latency and processing time

	// Simple success/failure logic
	if event.Amount > 1000.00 {
		log.Printf("Payment for Order ID %d FAILED (Amount too high for simulation)", event.OrderID)
		// In a real system, send an event back to order-service to update status to FAILED
	} else {
		log.Printf("Payment for Order ID %d SUCCESSFUL", event.OrderID)
		// In a real system, send an event back to order-service to update status to PAID
	}
	// --- End Simulation ---
}
