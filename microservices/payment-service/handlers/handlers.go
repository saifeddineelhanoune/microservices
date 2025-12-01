package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentRequest struct {
	OrderID int     `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
}

func ProcessPayment(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received direct payment request for Order ID: %d, Amount: %.2f", req.OrderID, req.Amount)

	time.Sleep(1 * time.Second)

	if req.Amount > 1000.00 {
		c.JSON(http.StatusOK, gin.H{"status": "failed", "message": "Payment failed due to high amount"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Payment processed successfully"})
	}
}
