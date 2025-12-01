package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"product-service/db"
	"product-service/models"
)

// CreateProduct handles the creation of a new product
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.ProductCollection.InsertOne(ctx, product)
	if err != nil {
		log.Printf("Error inserting product: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	product.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, product)
}

// GetProduct retrieves a product by its ID
func GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var product models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.ProductCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		if err == context.DeadlineExceeded {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "Database query timed out"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListProducts retrieves all products
func ListProducts(c *gin.Context) {
	var products []models.Product
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := db.ProductCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &products); err != nil {
		log.Printf("Error decoding products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
