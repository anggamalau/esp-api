package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                 primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name               string              `json:"name" bson:"name" validate:"required,min=2,max=50"`
	Email              string              `json:"email" bson:"email" validate:"required,email"`
	Password           string              `json:"-" bson:"password" validate:"required,min=6"`
	Role               string              `json:"role" bson:"role" validate:"required,oneof=admin liaison voice finance"`
	IsVerified         bool                `json:"is_verified" bson:"is_verified"`
	VerifiedAt         *time.Time          `json:"verified_at,omitempty" bson:"verified_at,omitempty"`
	VerifiedBy         *primitive.ObjectID `json:"verified_by,omitempty" bson:"verified_by,omitempty"`
	VerificationNotes  string              `json:"verification_notes,omitempty" bson:"verification_notes,omitempty"`
	LastPasswordReset  *time.Time          `json:"last_password_reset,omitempty" bson:"last_password_reset,omitempty"`
	PasswordResetCount int                 `json:"password_reset_count" bson:"password_reset_count"`
	CreatedAt          time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at" bson:"updated_at"`
}

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
	Role     string `json:"role" validate:"required,oneof=admin liaison voice finance" example:"user"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

type UserResponse struct {
	ID                string     `json:"id" example:"507f1f77bcf86cd799439011"`
	Name              string     `json:"name" example:"John Doe"`
	Email             string     `json:"email" example:"john@example.com"`
	Role              string     `json:"role" example:"user"`
	IsVerified        bool       `json:"is_verified" example:"true"`
	VerifiedAt        *time.Time `json:"verified_at,omitempty" example:"2024-01-01T00:00:00Z"`
	VerificationNotes string     `json:"verification_notes,omitempty" example:"Verified by admin"`
	CreatedAt         time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt         time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type UserUpdateRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=50" example:"Jane Doe"`
	Email string `json:"email" validate:"omitempty,email" example:"jane@example.com"`
}

// Admin role management models
type AdminUserRoleUpdateRequest struct {
	Role string `json:"role" validate:"required,oneof=admin liaison voice finance" example:"liaison"`
}

type AdminUserRoleUpdateResponse struct {
	Message string       `json:"message" example:"User role updated successfully"`
	User    UserResponse `json:"user"`
}

// Admin verification models
type VerificationRequest struct {
	Notes string `json:"notes" validate:"omitempty,max=500" example:"Identity verified through company records"`
}

type PendingUserResponse struct {
	ID        string    `json:"id" example:"507f1f77bcf86cd799439011"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john@example.com"`
	Role      string    `json:"role" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type RegisterPendingResponse struct {
	Message string       `json:"message" example:"Registration successful. Your account is pending admin verification."`
	User    UserResponse `json:"user"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message" example:"New password has been sent to your email address"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"oldpassword123"`
	NewPassword     string `json:"new_password" validate:"required,min=6" example:"newpassword456"`
	ConfirmPassword string `json:"confirm_password" validate:"required" example:"newpassword456"`
}

type ChangePasswordResponse struct {
	Message string `json:"message" example:"Your password has been updated successfully"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:                u.ID.Hex(),
		Name:              u.Name,
		Email:             u.Email,
		Role:              u.Role,
		IsVerified:        u.IsVerified,
		VerifiedAt:        u.VerifiedAt,
		VerificationNotes: u.VerificationNotes,
		CreatedAt:         u.CreatedAt,
		UpdatedAt:         u.UpdatedAt,
	}
}

func (u *User) ToPendingResponse() PendingUserResponse {
	return PendingUserResponse{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
