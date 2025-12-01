package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var ProductCollection *mongo.Collection

func InitDB() {
	// Set client options
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default for local dev
	}
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	MongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")

	// Set up collection
	databaseName := os.Getenv("MONGO_DB_NAME")
	if databaseName == "" {
		databaseName = "ecommerce"
	}
	ProductCollection = MongoClient.Database(databaseName).Collection("products")
}
