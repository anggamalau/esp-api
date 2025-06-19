package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Role constants for security and consistency
const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleModerator = "moderator"
)

// ValidRoles slice for easy iteration and validation
var ValidRoles = []string{RoleAdmin, RoleUser, RoleModerator}

var Validator *validator.Validate

func InitValidator() {
	Validator = validator.New()
}

func ValidateStruct(s interface{}) error {
	return Validator.Struct(s)
}

// IsValidRole checks if the provided role is valid
func IsValidRole(role string) bool {
	for _, validRole := range ValidRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// ValidateRole returns an error if the role is invalid
func ValidateRole(role string) error {
	if !IsValidRole(role) {
		return errors.New("invalid role: must be one of admin, user, or moderator")
	}
	return nil
}
