package handlers

import (
	"context"
	"time"

	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get current user's profile information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SwaggerUserResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user profile", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile retrieved successfully", user.ToResponse())
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update current user's profile information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.UserUpdateRequest  true  "Updated user data"
// @Success      200      {object}  models.SwaggerUserResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      404      {object}  models.SwaggerErrorResponse
// @Failure      409      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req models.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.userService.UpdateUser(ctx, userID, &req)
	if err != nil {
		if utils.IsValidationError(err) {
			return utils.ValidationErrorResponse(c, err)
		}
		if err == utils.ErrUserNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		if err == utils.ErrUserAlreadyExists {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Email already taken")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile updated successfully", user.ToResponse())
}

// DeleteProfile godoc
// @Summary      Delete user profile
// @Description  Delete current user's account permanently
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SwaggerResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /users/profile [delete]
func (h *UserHandler) DeleteProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.userService.DeleteUser(ctx, userID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete profile", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Profile deleted successfully", nil)
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Change current user's password by providing current and new password
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.ChangePasswordRequest  true  "Password change data"
// @Success      200      {object}  models.SwaggerChangePasswordResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      404      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /users/change-password [put]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req models.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.userService.ChangePassword(ctx, userID, &req)
	if err != nil {
		if utils.IsValidationError(err) {
			return utils.ValidationErrorResponse(c, err)
		}
		if err == utils.ErrInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Current password is incorrect")
		}
		if err == utils.ErrUserNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		if err.Error() == "password confirmation does not match" || err.Error() == "new password must be different from current password" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to change password", err.Error())
	}

	response := models.ChangePasswordResponse{
		Message: "Your password has been updated successfully",
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Password changed successfully", response)
}
