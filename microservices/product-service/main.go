package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"product-service/db"
	"product-service/handlers"
)

func main() {
	// Initialize Database
	db.InitDB()

	// Set up Gin router
	r := gin.Default()

	// Product routes
	r.POST("/products", handlers.CreateProduct)
	r.GET("/products/:id", handlers.GetProduct)
	r.GET("/products", handlers.ListProducts)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Product Service running on port %s", port)
	if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
