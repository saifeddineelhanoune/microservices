package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"order-service/db"
	"order-service/handlers"
	"order-service/pubsub"
)

func main() {
	db.InitDB()

	pubsub.InitRedis()

	r := gin.Default()

	r.POST("/orders", handlers.CreateOrder)
	r.GET("/orders/:id", handlers.GetOrder)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Order Service running on port %s", port)
	if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
