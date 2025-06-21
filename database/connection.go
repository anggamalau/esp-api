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

	// Create indexes for menus collection
	menuCollection := DB.Collection("menus")
	menuNameIndex := mongo.IndexModel{
		Keys:    map[string]interface{}{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = menuCollection.Indexes().CreateOne(ctx, menuNameIndex)
	if err != nil {
		log.Println("Warning: Failed to create menu name index:", err)
	}

	menuPathIndex := mongo.IndexModel{
		Keys:    map[string]interface{}{"path": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = menuCollection.Indexes().CreateOne(ctx, menuPathIndex)
	if err != nil {
		log.Println("Warning: Failed to create menu path index:", err)
	}

	menuOrderIndex := mongo.IndexModel{
		Keys: map[string]interface{}{"order": 1},
	}

	_, err = menuCollection.Indexes().CreateOne(ctx, menuOrderIndex)
	if err != nil {
		log.Println("Warning: Failed to create menu order index:", err)
	}

	// Create indexes for role_menu_permissions collection
	permissionCollection := DB.Collection("role_menu_permissions")
	roleMenuIndex := mongo.IndexModel{
		Keys:    map[string]interface{}{"role": 1, "menu_id": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = permissionCollection.Indexes().CreateOne(ctx, roleMenuIndex)
	if err != nil {
		log.Println("Warning: Failed to create role-menu permission index:", err)
	}

	roleIndex := mongo.IndexModel{
		Keys: map[string]interface{}{"role": 1},
	}

	_, err = permissionCollection.Indexes().CreateOne(ctx, roleIndex)
	if err != nil {
		log.Println("Warning: Failed to create role index:", err)
	}

	menuIdIndex := mongo.IndexModel{
		Keys: map[string]interface{}{"menu_id": 1},
	}

	_, err = permissionCollection.Indexes().CreateOne(ctx, menuIdIndex)
	if err != nil {
		log.Println("Warning: Failed to create menu_id index:", err)
	}

	log.Println("Database indexes created successfully")

	// Seed initial data
	SeedInitialData()
}
