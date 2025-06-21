package interfaces

import (
	"context"

	"backend/models"
)

type PermissionRepository interface {
	// Permission CRUD operations
	GrantPermission(ctx context.Context, permission *models.RoleMenuPermission) error
	RevokePermission(ctx context.Context, role, menuID string) error
	GetPermissionsByRole(ctx context.Context, role string) ([]*models.RoleMenuPermission, error)
	GetRolesByMenu(ctx context.Context, menuID string) ([]*models.RoleMenuPermission, error)
	GetAllPermissions(ctx context.Context) ([]*models.RoleMenuPermission, error)
	CheckPermission(ctx context.Context, role, menuID string) (bool, error)

	// Bulk operations
	RevokeAllPermissionsForMenu(ctx context.Context, menuID string) error
	RevokeAllPermissionsForRole(ctx context.Context, role string) error
}
