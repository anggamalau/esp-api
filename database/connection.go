package database

import (
	"context"
	"log"
	"time"

	"backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectMongoDB() {
	clientOptions := options.Client().ApplyURI(config.AppConfig.MongoDBURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	log.Println("Successfully connected to MongoDB")

	// Set the database
	DB = client.Database("esp_backend_db")

	// Create indexes
	createIndexes()
}

func createIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create unique index for user email
	userCollection := DB.Collection("users")
	emailIndex := mongo.IndexModel{
		Keys:    map[string]interface{}{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := userCollection.Indexes().CreateOne(ctx, emailIndex)
	if err != nil {
		log.Println("Warning: Failed to create email index:", err)
	}

	// Create index for refresh tokens
	tokenCollection := DB.Collection("refresh_tokens")
	tokenIndex := mongo.IndexModel{
		Keys: map[string]interface{}{"token": 1},
	}

	_, err = tokenCollection.Indexes().CreateOne(ctx, tokenIndex)
	if err != nil {
		log.Println("Warning: Failed to create token index:", err)
	}

	// Create TTL index for expired tokens
	expiryIndex := mongo.IndexModel{
		Keys:    map[string]interface{}{"expires_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err = tokenCollection.Indexes().CreateOne(ctx, expiryIndex)
	if err != nil {
		log.Println("Warning: Failed to create TTL index:", err)
	}

	log.Println("Database indexes created successfully")
}
