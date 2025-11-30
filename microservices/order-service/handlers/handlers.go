package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"order-service/db"
	"order-service/pubsub"
)

type Order struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id" binding:"required"`
	ProductID   int     `json:"product_id" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required"`
	Status      string  `json:"status"`
}

func CreateOrder(c *gin.Context) {
	var order Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order.Status = "PENDING_PAYMENT"
	var orderID int
	err := db.DB.QueryRow(
		"INSERT INTO orders (user_id, product_id, quantity, total_amount, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		order.UserID, order.ProductID, order.Quantity, order.TotalAmount, order.Status,
	).Scan(&orderID)

	if err != nil {
		log.Printf("Error inserting order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	order.ID = orderID

	// Publish ORDER_CREATED event
	if err := pubsub.PublishOrderCreated(order.ID, order.TotalAmount); err != nil {
		log.Printf("Warning: Failed to publish event for order %d: %v", order.ID, err)
	}

	c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID format"})
		return
	}

	var order Order
	err = db.DB.QueryRow(
		"SELECT id, user_id, product_id, quantity, total_amount, status FROM orders WHERE id = $1",
		orderID,
	).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.TotalAmount, &order.Status)

	if err != nil {
		log.Printf("Error retrieving order: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
