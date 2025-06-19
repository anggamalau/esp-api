package routes

import (
	"backend/handlers"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) {
	// Middleware
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CorsMiddleware())

	// Setup Swagger routes (conditional based on configuration)
	SetupSwaggerRoutes(app)

	// API v1 routes
	api := app.Group("/api/v1")

	// Health check
	// @Summary      Health check
	// @Description  Check if the server is running
	// @Tags         System
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  models.SwaggerHealthResponse
	// @Router       /health [get]
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Swagger status endpoint
	api.Get("/swagger-status", func(c *fiber.Ctx) error {
		return c.JSON(GetSwaggerStatus())
	})

	// Authentication routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Post("/logout", authHandler.Logout)

	// Protected routes
	protected := api.Group("/users", middleware.AuthMiddleware())
	protected.Get("/profile", userHandler.GetProfile)
	protected.Put("/profile", userHandler.UpdateProfile)
	protected.Delete("/profile", userHandler.DeleteProfile)
	protected.Post("/logout-all", authHandler.LogoutAll)
}
