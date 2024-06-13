// db.go
package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// DB is the MongoDB database instance
	DB *mongo.Database
)

func init() {
	// Load environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DB")

	// Set default values if not provided
	if mongoURI == "" {
		mongoURI = "mongodb+srv://sharanrohit7:QN100nYHfX9CMdyt@masterchatmodule.d97i0ny.mongodb.net"
	}

	if dbName == "" {
		dbName = "Log" // Change this to your desired database name
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB using the context
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	// Set the MongoDB database instance
	DB = client.Database(dbName)

	// Log connection details
	log.Printf("MongoDB Connection Details:\n")
	log.Printf(" - URI: %s\n", mongoURI)
	log.Printf(" - Database: %s\n", dbName)
	log.Println("Connection established successfully!")
}
