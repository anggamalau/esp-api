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
	userRepo  interfaces.UserRepository
	tokenRepo interfaces.TokenRepository
}

func NewAuthService(userRepo interfaces.UserRepository, tokenRepo interfaces.TokenRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.UserCreateRequest) (*models.LoginResponse, error) {
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

	// Create user
	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
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
