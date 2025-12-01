package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"user-service/db"
	"user-service/handlers"
)

func main() {
	// Initialize Database
	db.InitDB()

	// Set up Gin router
	r := gin.Default()

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected route example
	// r.GET("/profile", middleware.AuthMiddleware(), handlers.GetProfile)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("User Service running on port %s", port)
	if err := r.Run(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
