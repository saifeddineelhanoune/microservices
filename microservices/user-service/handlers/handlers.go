package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"user-service/db"
)

type User struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process request"})
		return
	}

	// Insert user into database
	_, err = db.DB.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, string(hashedPassword))
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var storedPasswordHash string
	err := db.DB.QueryRow("SELECT password_hash FROM users WHERE email = $1", user.Email).Scan(&storedPasswordHash)
	if err != nil {
		// User not found or other DB error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(user.Password))
	if err != nil {
		// Password mismatch
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// In a real application, a JWT token would be generated and returned here
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": "mock-jwt-token-12345"})
}
