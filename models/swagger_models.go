package models

// SwaggerResponse represents a standard API response for Swagger documentation
type SwaggerResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:""`
}

// SwaggerErrorResponse represents an error response for Swagger documentation
type SwaggerErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error occurred"`
	Data    string `json:"data" example:"null"`
	Error   string `json:"error" example:"Detailed error message"`
}

// SwaggerLoginResponse represents login response for Swagger documentation
type SwaggerLoginResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Login successful"`
	Data    LoginResponse `json:"data"`
	Error   string        `json:"error,omitempty" example:""`
}

// SwaggerUserResponse represents user profile response for Swagger documentation
type SwaggerUserResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Profile retrieved successfully"`
	Data    UserResponse `json:"data"`
	Error   string       `json:"error,omitempty" example:""`
}

// SwaggerTokenResponse represents token refresh response for Swagger documentation
type SwaggerTokenResponse struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"Token refreshed successfully"`
	Data    TokenPair `json:"data"`
	Error   string    `json:"error,omitempty" example:""`
}

// SwaggerHealthResponse represents health check response for Swagger documentation
type SwaggerHealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Server is running"`
}

// SwaggerValidationError represents validation error details
type SwaggerValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email is required"`
}

// SwaggerValidationErrorResponse represents validation error response
type SwaggerValidationErrorResponse struct {
	Success bool                     `json:"success" example:"false"`
	Message string                   `json:"message" example:"Validation failed"`
	Data    string                   `json:"data" example:"null"`
	Error   []SwaggerValidationError `json:"error"`
}
