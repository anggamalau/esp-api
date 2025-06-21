package routes

import (
	"backend/handlers"
	"backend/middleware"
	"backend/repositories/interfaces"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, adminHandler *handlers.AdminHandler, menuHandler *handlers.MenuHandler, userRepo interfaces.UserRepository) {
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
	auth.Post("/forgot-password", authHandler.ForgotPassword)

	// Protected routes
	protected := api.Group("/users", middleware.AuthMiddleware())
	protected.Get("/profile", userHandler.GetProfile)
	protected.Put("/profile", userHandler.UpdateProfile)
	protected.Delete("/profile", userHandler.DeleteProfile)
	protected.Put("/change-password", userHandler.ChangePassword)
	protected.Post("/logout-all", authHandler.LogoutAll)
	protected.Get("/menus", menuHandler.GetUserMenus)

	// Admin-only routes
	admin := api.Group("/admin", middleware.AuthMiddleware(), middleware.AdminMiddleware(userRepo))
	admin.Get("/users/pending", adminHandler.GetPendingUsers)
	admin.Post("/users/:id/verify", adminHandler.VerifyUser)
	admin.Get("/users/:id", adminHandler.GetUserDetails)

	// Menu management routes (Admin only)
	admin.Post("/menus", menuHandler.CreateMenu)
	admin.Get("/menus", menuHandler.GetAllMenus)
	admin.Get("/menus/:id", menuHandler.GetMenuByID)
	admin.Put("/menus/:id", menuHandler.UpdateMenu)
	admin.Delete("/menus/:id", menuHandler.DeleteMenu)
	admin.Get("/menus/:id/roles", menuHandler.GetRolesByMenu)

	// Permission management routes (Admin only)
	admin.Post("/roles/:role/menus/:menuId", menuHandler.GrantPermission)
	admin.Delete("/roles/:role/menus/:menuId", menuHandler.RevokePermission)
	admin.Get("/roles/:role/menus", menuHandler.GetPermissionsByRole)
	admin.Get("/roles/permissions", menuHandler.GetAllPermissions)
	admin.Get("/roles/summary", menuHandler.GetRolePermissionSummary)
}
