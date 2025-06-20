package handlers

import (
	"context"
	"time"

	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	adminService *services.AdminService
}

func NewAdminHandler(adminService *services.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// GetPendingUsers godoc
// @Summary      Get pending users
// @Description  Get list of users awaiting admin verification
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SwaggerPendingUsersResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /admin/users/pending [get]
func (h *AdminHandler) GetPendingUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pendingUsers, err := h.adminService.GetPendingUsers(ctx)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get pending users", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Pending users retrieved successfully", pendingUsers)
}

// VerifyUser godoc
// @Summary      Verify user
// @Description  Verify a user account to allow them to login
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "User ID"
// @Param        request  body      models.VerificationRequest  true  "Verification data"
// @Success      200      {object}  models.SwaggerResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      403      {object}  models.SwaggerErrorResponse
// @Failure      404      {object}  models.SwaggerErrorResponse
// @Failure      409      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /admin/users/{id}/verify [post]
func (h *AdminHandler) VerifyUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	adminID := c.Locals("userID").(string)

	var req models.VerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.adminService.VerifyUser(ctx, userID, adminID, &req)
	if err != nil {
		if utils.IsValidationError(err) {
			return utils.ValidationErrorResponse(c, err)
		}
		if err == utils.ErrUserNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		if err == utils.ErrUserAlreadyVerified {
			return utils.ErrorResponse(c, fiber.StatusConflict, "User already verified")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to verify user", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User verified successfully", nil)
}

// GetUserDetails godoc
// @Summary      Get user details
// @Description  Get detailed information about a specific user (admin only)
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.SwaggerUserResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /admin/users/{id} [get]
func (h *AdminHandler) GetUserDetails(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := h.adminService.GetUserByID(ctx, userID)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get user details", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User details retrieved successfully", user.ToResponse())
}
