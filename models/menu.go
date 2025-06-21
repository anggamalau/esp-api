package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Menu represents a menu item in the system
type Menu struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Description string             `json:"description" bson:"description" validate:"omitempty,max=200"`
	Icon        string             `json:"icon" bson:"icon" validate:"omitempty,max=50"`
	Path        string             `json:"path" bson:"path" validate:"required,max=100"`
	Order       int                `json:"order" bson:"order" validate:"min=0"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// RoleMenuPermission represents the junction table for role-menu access
type RoleMenuPermission struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Role          string             `json:"role" bson:"role" validate:"required,oneof=admin liaison voice finance"`
	MenuID        primitive.ObjectID `json:"menu_id" bson:"menu_id" validate:"required"`
	GrantedByID   primitive.ObjectID `json:"granted_by_id" bson:"granted_by_id" validate:"required"`
	GrantedByName string             `json:"granted_by_name" bson:"granted_by_name"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
}

// Request/Response models for API

type MenuCreateRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=50" example:"Dashboard"`
	Description string `json:"description" validate:"omitempty,max=200" example:"Main dashboard view"`
	Icon        string `json:"icon" validate:"omitempty,max=50" example:"dashboard-icon"`
	Path        string `json:"path" validate:"required,max=100" example:"/dashboard"`
	Order       int    `json:"order" validate:"min=0" example:"1"`
}

type MenuUpdateRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=50" example:"Dashboard"`
	Description string `json:"description" validate:"omitempty,max=200" example:"Main dashboard view"`
	Icon        string `json:"icon" validate:"omitempty,max=50" example:"dashboard-icon"`
	Path        string `json:"path" validate:"omitempty,max=100" example:"/dashboard"`
	Order       int    `json:"order" validate:"omitempty,min=0" example:"1"`
	IsActive    *bool  `json:"is_active" validate:"omitempty" example:"true"`
}

type MenuResponse struct {
	ID          string    `json:"id" example:"507f1f77bcf86cd799439011"`
	Name        string    `json:"name" example:"Dashboard"`
	Description string    `json:"description" example:"Main dashboard view"`
	Icon        string    `json:"icon" example:"dashboard-icon"`
	Path        string    `json:"path" example:"/dashboard"`
	Order       int       `json:"order" example:"1"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type RoleMenuPermissionRequest struct {
	Role   string `json:"role" validate:"required,oneof=admin liaison voice finance" example:"liaison"`
	MenuID string `json:"menu_id" validate:"required" example:"507f1f77bcf86cd799439011"`
}

type RoleMenuPermissionResponse struct {
	ID            string    `json:"id" example:"507f1f77bcf86cd799439011"`
	Role          string    `json:"role" example:"liaison"`
	MenuID        string    `json:"menu_id" example:"507f1f77bcf86cd799439011"`
	MenuName      string    `json:"menu_name" example:"Dashboard"`
	GrantedByID   string    `json:"granted_by_id" example:"507f1f77bcf86cd799439011"`
	GrantedByName string    `json:"granted_by_name" example:"Admin User"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type UserMenuResponse struct {
	ID          string `json:"id" example:"507f1f77bcf86cd799439011"`
	Name        string `json:"name" example:"Dashboard"`
	Description string `json:"description" example:"Main dashboard view"`
	Icon        string `json:"icon" example:"dashboard-icon"`
	Path        string `json:"path" example:"/dashboard"`
	Order       int    `json:"order" example:"1"`
}

type RolePermissionSummary struct {
	Role      string         `json:"role" example:"liaison"`
	MenuCount int            `json:"menu_count" example:"3"`
	Menus     []MenuResponse `json:"menus"`
}

// Helper methods

func (m *Menu) ToResponse() MenuResponse {
	return MenuResponse{
		ID:          m.ID.Hex(),
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
		Path:        m.Path,
		Order:       m.Order,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (m *Menu) ToUserMenuResponse() UserMenuResponse {
	return UserMenuResponse{
		ID:          m.ID.Hex(),
		Name:        m.Name,
		Description: m.Description,
		Icon:        m.Icon,
		Path:        m.Path,
		Order:       m.Order,
	}
}

func (rmp *RoleMenuPermission) ToResponse(menuName string) RoleMenuPermissionResponse {
	return RoleMenuPermissionResponse{
		ID:            rmp.ID.Hex(),
		Role:          rmp.Role,
		MenuID:        rmp.MenuID.Hex(),
		MenuName:      menuName,
		GrantedByID:   rmp.GrantedByID.Hex(),
		GrantedByName: rmp.GrantedByName,
		CreatedAt:     rmp.CreatedAt,
	}
}
