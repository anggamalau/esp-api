package interfaces

import (
	"context"

	"backend/models"
)

type MenuRepository interface {
	// Menu CRUD operations
	Create(ctx context.Context, menu *models.Menu) error
	GetAll(ctx context.Context) ([]*models.Menu, error)
	GetByID(ctx context.Context, id string) (*models.Menu, error)
	GetActiveMenus(ctx context.Context) ([]*models.Menu, error)
	Update(ctx context.Context, id string, menu *models.Menu) error
	Delete(ctx context.Context, id string) error

	// Menu ordering
	GetMenusOrderedByOrder(ctx context.Context) ([]*models.Menu, error)

	// Menu by role access
	GetMenusByRole(ctx context.Context, role string) ([]*models.Menu, error)
}
