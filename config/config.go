package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	Host             string
	MongoDBURI       string
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
	BcryptRounds     int
	AppEnv           string

	// Swagger Configuration
	SwaggerEnabled  bool
	SwaggerHost     string
	SwaggerBasePath string
	SwaggerSchemes  string
	SwaggerTitle    string
	SwaggerVersion  string
	SwaggerUIPath   string
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	bcryptRounds, err := strconv.Atoi(getEnv("BCRYPT_ROUNDS", "12"))
	if err != nil {
		bcryptRounds = 12
	}

	AppConfig = &Config{
		Port:             getEnv("PORT", "3000"),
		Host:             getEnv("HOST", "localhost"),
		MongoDBURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017/backend_db"),
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "your-access-secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret"),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "168h"),
		BcryptRounds:     bcryptRounds,
		AppEnv:           getEnv("APP_ENV", "development"),

		// Swagger Configuration
		SwaggerEnabled:  getEnvBool("SWAGGER_ENABLED", true),
		SwaggerHost:     getEnv("SWAGGER_HOST", "localhost:3000"),
		SwaggerBasePath: getEnv("SWAGGER_BASE_PATH", "/api/v1"),
		SwaggerSchemes:  getEnv("SWAGGER_SCHEMES", "http"),
		SwaggerTitle:    getEnv("SWAGGER_TITLE", "Backend API"),
		SwaggerVersion:  getEnv("SWAGGER_VERSION", "1.0"),
		SwaggerUIPath:   getEnv("SWAGGER_UI_PATH", "/swagger"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return defaultValue
}

// ShouldEnableSwagger determines if Swagger should be enabled based on environment and configuration
func (c *Config) ShouldEnableSwagger() bool {
	// Rule 1: Explicit configuration override
	if !c.SwaggerEnabled {
		return false
	}

	// Rule 2: Development environment default
	if c.AppEnv == "development" || c.AppEnv == "dev" {
		return true
	}

	// Rule 3: Staging environment (optional)
	if c.AppEnv == "staging" && c.SwaggerEnabled {
		return true
	}

	// Rule 4: Production - only if explicitly enabled
	if c.AppEnv == "production" && c.SwaggerEnabled {
		return true
	}

	// Default: disabled
	return false
}
