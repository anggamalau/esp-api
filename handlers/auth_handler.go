package handlers

import (
	"context"
	"time"

	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      models.UserCreateRequest  true  "User registration data"
// @Success      201      {object}  models.SwaggerLoginResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      409      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.UserCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.authService.Register(ctx, &req)
	if err != nil {
		if utils.IsValidationError(err) {
			return utils.ValidationErrorResponse(c, err)
		}
		if err == utils.ErrUserAlreadyExists {
			return utils.ErrorResponse(c, fiber.StatusConflict, "User already exists")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to register user", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "User registered successfully", response)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      models.UserLoginRequest  true  "User login credentials"
// @Success      200      {object}  models.SwaggerLoginResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.authService.Login(ctx, &req)
	if err != nil {
		if utils.IsValidationError(err) {
			return utils.ValidationErrorResponse(c, err)
		}
		if err == utils.ErrInvalidCredentials {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to login", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Login successful", response)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get a new access token using refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      models.RefreshTokenRequest  true  "Refresh token"
// @Success      200      {object}  models.SwaggerTokenResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tokens, err := h.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if utils.IsAuthError(err) {
			return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired refresh token")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh token", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Token refreshed successfully", tokens)
}

// Logout godoc
// @Summary      User logout
// @Description  Revoke refresh token to logout user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      models.RefreshTokenRequest  true  "Refresh token to revoke"
// @Success      200      {object}  models.SwaggerResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.authService.Logout(ctx, req.RefreshToken)
	if err != nil {
		if err == utils.ErrTokenNotFound {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid refresh token")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Logout successful", nil)
}

// LogoutAll godoc
// @Summary      Logout from all devices
// @Description  Revoke all refresh tokens for the user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SwaggerResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /users/logout-all [post]
func (h *AuthHandler) LogoutAll(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.authService.LogoutAll(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to logout from all devices", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Logged out from all devices successfully", nil)
}
