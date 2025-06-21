package middleware

import (
	"context"
	"time"

	"backend/repositories/interfaces"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func MenuAccessMiddleware(userRepo interfaces.UserRepository, permRepo interfaces.PermissionRepository, menuID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from auth middleware
		userID, ok := c.Locals("userID").(string)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in context")
		}

		// Get user from database to check role
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := userRepo.GetByID(ctx, userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
		}

		// Admin has access to all menus
		if user.Role == "admin" {
			return c.Next()
		}

		// Check if user's role has permission to access this menu
		hasPermission, err := permRepo.CheckPermission(ctx, user.Role, menuID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check permissions")
		}

		if !hasPermission {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Menu access denied for your role")
		}

		return c.Next()
	}
}

// RoleMiddleware checks if user has one of the specified roles
func RoleMiddleware(userRepo interfaces.UserRepository, allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from auth middleware
		userID, ok := c.Locals("userID").(string)
		if !ok {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User ID not found in context")
		}

		// Get user from database to check role
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := userRepo.GetByID(ctx, userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
		}

		// Check if user has any of the allowed roles
		for _, role := range allowedRoles {
			if user.Role == role {
				return c.Next()
			}
		}

		return utils.ErrorResponse(c, fiber.StatusForbidden, "Insufficient role permissions")
	}
}
