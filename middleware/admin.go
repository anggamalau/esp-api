package middleware

import (
	"context"
	"time"

	"backend/repositories/interfaces"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(userRepo interfaces.UserRepository) fiber.Handler {
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

		// Check if user has admin role
		if user.Role != "admin" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Admin access required")
		}

		// Check if admin is verified
		if !user.IsVerified {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Admin account not verified")
		}

		return c.Next()
	}
}
