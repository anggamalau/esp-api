package database

import (
	"context"
	"log"
	"time"

	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SeedInitialData() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check if menus already exist
	menuCollection := DB.Collection("menus")
	count, err := menuCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Println("Error checking existing menus:", err)
		return
	}

	if count > 0 {
		log.Println("Menus already exist, skipping seed")
		return
	}

	// Define initial menus
	initialMenus := []models.Menu{
		{
			ID:          primitive.NewObjectID(),
			Name:        "Dashboard",
			Description: "Main dashboard with overview and analytics",
			Icon:        "dashboard",
			Path:        "/dashboard",
			Order:       1,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "User Management",
			Description: "Manage users, roles and permissions",
			Icon:        "users",
			Path:        "/users",
			Order:       2,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Reports",
			Description: "Generate and view various reports",
			Icon:        "chart-bar",
			Path:        "/reports",
			Order:       3,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Settings",
			Description: "System settings and configuration",
			Icon:        "cog",
			Path:        "/settings",
			Order:       4,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Help & Support",
			Description: "Help documentation and support resources",
			Icon:        "question-circle",
			Path:        "/help",
			Order:       5,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Insert menus
	var menuDocs []interface{}
	for _, menu := range initialMenus {
		menuDocs = append(menuDocs, menu)
	}

	_, err = menuCollection.InsertMany(ctx, menuDocs)
	if err != nil {
		log.Println("Error inserting initial menus:", err)
		return
	}

	log.Println("Successfully seeded initial menus:")
	for _, menu := range initialMenus {
		log.Printf("- %s (%s)", menu.Name, menu.Path)
	}

	log.Println("Note: Admin role has access to all menus by default (no explicit permissions needed)")
	log.Println("Use the admin endpoints to grant menu access to other roles (liaison, voice, finance)")
}
