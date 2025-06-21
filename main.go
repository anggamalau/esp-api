// @title           Backend API
// @version         1.0
// @description     A robust backend service built with Go Fiber framework, featuring MongoDB integration and JWT-based authentication with refresh token support.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:3000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"log"
	"strings"

	"backend/config"
	"backend/database"
	_ "backend/docs"
	"backend/handlers"
	"backend/repositories"
	"backend/routes"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize validator
	utils.InitValidator()

	// Connect to MongoDB
	database.ConnectMongoDB()

	// Configure Swagger based on environment
	configureSwagger()

	// Initialize repositories
	userRepo := repositories.NewUserRepository()
	tokenRepo := repositories.NewTokenRepository()
	permissionRepo := repositories.NewPermissionRepository()
	menuRepo := repositories.NewMenuRepository(permissionRepo)

	// Initialize services
	emailService := services.NewEmailService()
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenRepo, emailService)
	adminService := services.NewAdminService(userRepo)
	menuService := services.NewMenuService(menuRepo, permissionRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	adminHandler := handlers.NewAdminHandler(adminService)
	menuHandler := handlers.NewMenuHandler(menuService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return utils.ErrorResponse(c, code, "Internal Server Error", err.Error())
		},
	})

	// Setup routes
	routes.SetupRoutes(app, authHandler, userHandler, adminHandler, menuHandler, userRepo)

	// Log Swagger status
	logSwaggerStatus()

	// Start server
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

// configureSwagger sets up Swagger documentation based on configuration
func configureSwagger() {
	if config.AppConfig.ShouldEnableSwagger() {
		// Note: docs package will be imported after generation
		log.Println("Swagger configuration will be applied when docs are generated")
	}
}

// logSwaggerStatus logs the current Swagger configuration status
func logSwaggerStatus() {
	log.Printf("=== Swagger Configuration ===")
	log.Printf("Enabled: %v", config.AppConfig.ShouldEnableSwagger())
	log.Printf("Environment: %s", config.AppConfig.AppEnv)

	if config.AppConfig.ShouldEnableSwagger() {
		schemes := strings.Split(config.AppConfig.SwaggerSchemes, ",")
		scheme := "http"
		if len(schemes) > 0 {
			scheme = strings.TrimSpace(schemes[0])
		}

		log.Printf("Swagger UI: %s://%s%s/index.html", scheme, config.AppConfig.SwaggerHost, config.AppConfig.SwaggerUIPath)
		log.Printf("API Spec JSON: %s://%s%s/doc.json", scheme, config.AppConfig.SwaggerHost, config.AppConfig.SwaggerUIPath)
		log.Printf("Title: %s", config.AppConfig.SwaggerTitle)
		log.Printf("Version: %s", config.AppConfig.SwaggerVersion)
	} else {
		log.Printf("Swagger UI is disabled for environment: %s", config.AppConfig.AppEnv)
		log.Printf("Access /swagger/* will return 404")
	}
	log.Printf("============================")
}
