package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"payment-service/handlers"
	"payment-service/pubsub"
)

func main() {
	// Initialize Redis Pub/Sub client
	pubsub.InitRedis()

	// Start event consumer in a goroutine
	go pubsub.StartConsumer()

	// Set up Gin router
	r := gin.Default()

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Mock payment route (for direct calls if needed)
	r.POST("/payments/process", handlers.ProcessPayment)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Payment Service running on port %s", port)
	if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
