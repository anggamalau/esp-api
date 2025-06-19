package routes

import (
	"log"

	"backend/config"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// SetupSwaggerRoutes configures Swagger routes based on configuration
func SetupSwaggerRoutes(app *fiber.App) {
	if !config.AppConfig.ShouldEnableSwagger() {
		log.Printf("Swagger documentation is disabled (Environment: %s)", config.AppConfig.AppEnv)

		// Optional: Serve a disabled message
		app.Get("/swagger/*", func(c *fiber.Ctx) error {
			return c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "API documentation is not available in this environment",
				"error":   "Swagger is disabled",
			})
		})
		return
	}

	log.Printf("Swagger documentation enabled at %s/*", config.AppConfig.SwaggerUIPath)
	log.Printf("Environment: %s", config.AppConfig.AppEnv)
	log.Printf("Swagger UI: http://%s%s/index.html", config.AppConfig.SwaggerHost, config.AppConfig.SwaggerUIPath)

	// Enable Swagger UI
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Redirect /swagger to /swagger/index.html
	app.Get("/swagger", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})

	// API documentation endpoint
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})
}

// GetSwaggerStatus returns current Swagger configuration status
func GetSwaggerStatus() fiber.Map {
	return fiber.Map{
		"enabled":     config.AppConfig.ShouldEnableSwagger(),
		"environment": config.AppConfig.AppEnv,
		"ui_url":      config.AppConfig.SwaggerHost + config.AppConfig.SwaggerUIPath,
		"config": fiber.Map{
			"host":      config.AppConfig.SwaggerHost,
			"base_path": config.AppConfig.SwaggerBasePath,
			"schemes":   config.AppConfig.SwaggerSchemes,
			"title":     config.AppConfig.SwaggerTitle,
			"version":   config.AppConfig.SwaggerVersion,
		},
	}
}
