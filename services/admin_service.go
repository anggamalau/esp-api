package services

import (
	"context"

	"backend/models"
	"backend/repositories/interfaces"
	"backend/utils"
)

type AdminService struct {
	userRepo interfaces.UserRepository
}

func NewAdminService(userRepo interfaces.UserRepository) *AdminService {
	return &AdminService{
		userRepo: userRepo,
	}
}

func (s *AdminService) GetPendingUsers(ctx context.Context) ([]*models.PendingUserResponse, error) {
	users, err := s.userRepo.GetPendingUsers(ctx)
	if err != nil {
		return nil, err
	}

	var pendingUsers []*models.PendingUserResponse
	for _, user := range users {
		pendingUser := user.ToPendingResponse()
		pendingUsers = append(pendingUsers, &pendingUser)
	}

	return pendingUsers, nil
}

func (s *AdminService) VerifyUser(ctx context.Context, userID, adminID string, req *models.VerificationRequest) error {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return err
	}

	// Check if user exists and is not already verified
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return utils.ErrUserAlreadyVerified
	}

	// Verify the user
	return s.userRepo.VerifyUser(ctx, userID, adminID, req.Notes)
}

func (s *AdminService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
