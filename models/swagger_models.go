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
	Status  string    `json:"status" example:"success"`
	Message string    `json:"message" example:"Token refreshed successfully"`
	Data    TokenPair `json:"data"`
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

// SwaggerRegisterPendingResponse represents registration pending response for Swagger documentation
type SwaggerRegisterPendingResponse struct {
	Success bool                    `json:"success" example:"true"`
	Message string                  `json:"message" example:"Registration successful. Your account is pending admin verification."`
	Data    RegisterPendingResponse `json:"data"`
	Error   string                  `json:"error,omitempty" example:""`
}

// SwaggerPendingUsersResponse represents pending users list response for Swagger documentation
type SwaggerPendingUsersResponse struct {
	Success bool                  `json:"success" example:"true"`
	Message string                `json:"message" example:"Pending users retrieved successfully"`
	Data    []PendingUserResponse `json:"data"`
	Error   string                `json:"error,omitempty" example:""`
}

type SwaggerForgotPasswordResponse struct {
	Status  string                 `json:"status" example:"success"`
	Message string                 `json:"message" example:"Request processed successfully"`
	Data    ForgotPasswordResponse `json:"data"`
}

// Menu-related Swagger models

// SwaggerMenuResponse represents menu response for Swagger documentation
type SwaggerMenuResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Menu retrieved successfully"`
	Data    MenuResponse `json:"data"`
	Error   string       `json:"error,omitempty" example:""`
}

// SwaggerMenuListResponse represents menu list response for Swagger documentation
type SwaggerMenuListResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"Menus retrieved successfully"`
	Data    []MenuResponse `json:"data"`
	Error   string         `json:"error,omitempty" example:""`
}

// SwaggerUserMenuResponse represents user menu response for Swagger documentation
type SwaggerUserMenuResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"User menus retrieved successfully"`
	Data    []UserMenuResponse `json:"data"`
	Error   string             `json:"error,omitempty" example:""`
}

// SwaggerPermissionResponse represents permission response for Swagger documentation
type SwaggerPermissionResponse struct {
	Success bool                       `json:"success" example:"true"`
	Message string                     `json:"message" example:"Permission retrieved successfully"`
	Data    RoleMenuPermissionResponse `json:"data"`
	Error   string                     `json:"error,omitempty" example:""`
}

// SwaggerPermissionListResponse represents permission list response for Swagger documentation
type SwaggerPermissionListResponse struct {
	Success bool                         `json:"success" example:"true"`
	Message string                       `json:"message" example:"Permissions retrieved successfully"`
	Data    []RoleMenuPermissionResponse `json:"data"`
	Error   string                       `json:"error,omitempty" example:""`
}

// SwaggerRolePermissionSummaryResponse represents role permission summary response for Swagger documentation
type SwaggerRolePermissionSummaryResponse struct {
	Success bool                    `json:"success" example:"true"`
	Message string                  `json:"message" example:"Role summary retrieved successfully"`
	Data    []RolePermissionSummary `json:"data"`
	Error   string                  `json:"error,omitempty" example:""`
}
