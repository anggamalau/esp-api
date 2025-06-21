package services

import (
	"context"
	"errors"
	"time"

	"backend/models"
	"backend/repositories/interfaces"
	"backend/utils"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *models.UserCreateRequest) (*models.User, error) {
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
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, userID string, req *models.UserUpdateRequest) (*models.User, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if email is already taken by another user
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
		if err == nil && existingUser.ID != user.ID {
			return nil, utils.ErrUserAlreadyExists
		}
		if err != nil && err != utils.ErrUserNotFound {
			return nil, err
		}
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	// Update user
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *UserService) ValidateUserCredentials(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == utils.ErrUserNotFound {
			return nil, utils.ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, utils.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	// Check if new password matches confirm password
	if req.NewPassword != req.ConfirmPassword {
		return errors.New("password confirmation does not match")
	}

	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !utils.CheckPasswordHash(req.CurrentPassword, user.Password) {
		return utils.ErrInvalidCredentials
	}

	// Check if new password is different from current password
	if utils.CheckPasswordHash(req.NewPassword, user.Password) {
		return errors.New("new password must be different from current password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password in database
	err = s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}
