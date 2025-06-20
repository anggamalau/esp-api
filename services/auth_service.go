package services

import (
	"context"
	"time"

	"backend/config"
	"backend/models"
	"backend/repositories/interfaces"
	"backend/utils"
)

type AuthService struct {
	userRepo     interfaces.UserRepository
	tokenRepo    interfaces.TokenRepository
	emailService *EmailService
}

func NewAuthService(userRepo interfaces.UserRepository, tokenRepo interfaces.TokenRepository, emailService *EmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
		emailService: emailService,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.UserCreateRequest) (*models.RegisterPendingResponse, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, utils.ErrUserAlreadyExists
	}
	if err != utils.ErrUserNotFound {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user with verification disabled by default
	user := &models.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       req.Role,
		IsVerified: false, // New users need admin verification
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Return pending response (no tokens for unverified users)
	return &models.RegisterPendingResponse{
		Message: "Registration successful. Your account is pending admin verification.",
		User:    user.ToResponse(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.UserLoginRequest) (*models.LoginResponse, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return nil, utils.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, utils.ErrInvalidCredentials
	}

	// Check if user is verified
	if !user.IsVerified {
		return nil, utils.ErrUserNotVerified
	}

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		User:   user.ToResponse(),
		Tokens: *tokens,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*models.TokenPair, error) {
	// Get refresh token from database
	refreshToken, err := s.tokenRepo.GetByToken(ctx, refreshTokenString)
	if err != nil {
		return nil, utils.ErrInvalidToken
	}

	// Check if token is revoked
	if refreshToken.IsRevoked {
		return nil, utils.ErrTokenRevoked
	}

	// Check if token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, utils.ErrTokenExpired
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, refreshToken.UserID.Hex())
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	err = s.tokenRepo.RevokeToken(ctx, refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	return s.generateTokenPair(ctx, user)
}

func (s *AuthService) Logout(ctx context.Context, refreshTokenString string) error {
	return s.tokenRepo.RevokeToken(ctx, refreshTokenString)
}

func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	return s.tokenRepo.RevokeAllUserTokens(ctx, userID)
}

func (s *AuthService) generateTokenPair(ctx context.Context, user *models.User) (*models.TokenPair, error) {
	// Generate access token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenString := utils.GenerateRefreshToken()

	// Calculate expiry
	expiryDuration, err := time.ParseDuration(config.AppConfig.JWTRefreshExpiry)
	if err != nil {
		expiryDuration = 168 * time.Hour // 7 days fallback
	}

	// Save refresh token to database
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(expiryDuration),
	}

	err = s.tokenRepo.Create(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Calculate access token expiry
	accessExpiryDuration, err := time.ParseDuration(config.AppConfig.JWTAccessExpiry)
	if err != nil {
		accessExpiryDuration = 15 * time.Minute // fallback
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessExpiryDuration.Seconds()),
	}, nil
}

func (s *AuthService) CleanupExpiredTokens(ctx context.Context) error {
	return s.tokenRepo.DeleteExpiredTokens(ctx)
}

// ForgotPassword generates a new password and sends it via email
func (s *AuthService) ForgotPassword(ctx context.Context, req *models.ForgotPasswordRequest) (*models.ForgotPasswordResponse, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == utils.ErrUserNotFound {
			// For security reasons, don't reveal if email exists or not
			return &models.ForgotPasswordResponse{
				Message: "If the email address exists in our system, a new password has been sent.",
			}, nil
		}
		return nil, err
	}

	// Check if user is verified
	if !user.IsVerified {
		return nil, utils.ErrUserNotEligibleForReset
	}

	// Check rate limiting
	if err := s.checkPasswordResetRateLimit(user); err != nil {
		return nil, err
	}

	// Generate new secure password
	newPassword, err := utils.GenerateSecurePassword(config.AppConfig.PasswordResetLength)
	if err != nil {
		return nil, utils.ErrPasswordGenerationFailed
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	// Update user's password in database
	if err := s.userRepo.UpdatePassword(ctx, user.ID.Hex(), hashedPassword); err != nil {
		return nil, err
	}

	// Update password reset tracking
	if err := s.userRepo.UpdatePasswordResetInfo(ctx, user.ID.Hex()); err != nil {
		// Log error but don't fail the request
		// The password has already been updated successfully
	}

	// Send email with new password
	if err := s.emailService.SendPasswordResetEmail(user.Email, user.Name, newPassword); err != nil {
		// Log the error but still return success to user for security
		// In production, you might want to queue this for retry
		return nil, utils.ErrEmailDeliveryFailed
	}

	// Revoke all existing tokens for security
	_ = s.tokenRepo.RevokeAllUserTokens(ctx, user.ID.Hex())

	return &models.ForgotPasswordResponse{
		Message: "A new password has been sent to your email address.",
	}, nil
}

// checkPasswordResetRateLimit checks if user has exceeded password reset attempts
func (s *AuthService) checkPasswordResetRateLimit(user *models.User) error {
	maxAttempts := config.AppConfig.PasswordResetAttempts
	if maxAttempts <= 0 {
		return nil // Rate limiting disabled
	}

	// Check if user has exceeded attempts within the last hour
	if user.LastPasswordReset != nil {
		oneHourAgo := time.Now().Add(-1 * time.Hour)
		if user.LastPasswordReset.After(oneHourAgo) && user.PasswordResetCount >= maxAttempts {
			return utils.ErrPasswordResetLimitExceeded
		}
	}

	return nil
}
