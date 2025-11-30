package pubsub

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: "",
		DB: 0,        
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis!")
}

func PublishOrderCreated(orderID int, totalAmount float64) error {
	message := map[string]interface{}{
		"order_id": orderID,
		"amount":   totalAmount,
	}
	msg := fmt.Sprintf(`{"order_id": %d, "amount": %.2f}`, orderID, totalAmount)

	err := RedisClient.Publish(Ctx, "order_events", msg).Err()
	if err != nil {
		log.Printf("Failed to publish ORDER_CREATED event: %v", err)
		return err
	}
	log.Printf("Published ORDER_CREATED event for order %d", orderID)
	return nil
}
