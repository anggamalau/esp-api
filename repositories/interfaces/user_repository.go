package interfaces

import (
	"backend/models"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
	GetPendingUsers(ctx context.Context) ([]*models.User, error)
	VerifyUser(ctx context.Context, userID, adminID string, notes string) error
	UpdatePassword(ctx context.Context, userID, hashedPassword string) error
	UpdatePasswordResetInfo(ctx context.Context, userID string) error
}
