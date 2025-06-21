package services

import (
	"context"
	"time"

	"backend/models"
	"backend/repositories/interfaces"
	"backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuService struct {
	menuRepo       interfaces.MenuRepository
	permissionRepo interfaces.PermissionRepository
	userRepo       interfaces.UserRepository
}

func NewMenuService(menuRepo interfaces.MenuRepository, permissionRepo interfaces.PermissionRepository, userRepo interfaces.UserRepository) *MenuService {
	return &MenuService{
		menuRepo:       menuRepo,
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

// Menu CRUD operations

func (s *MenuService) CreateMenu(ctx context.Context, req *models.MenuCreateRequest) (*models.MenuResponse, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	menu := &models.Menu{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Path:        req.Path,
		Order:       req.Order,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.menuRepo.Create(ctx, menu); err != nil {
		return nil, err
	}

	response := menu.ToResponse()
	return &response, nil
}

func (s *MenuService) GetAllMenus(ctx context.Context) ([]*models.MenuResponse, error) {
	menus, err := s.menuRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*models.MenuResponse
	for _, menu := range menus {
		response := menu.ToResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *MenuService) GetMenuByID(ctx context.Context, id string) (*models.MenuResponse, error) {
	menu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := menu.ToResponse()
	return &response, nil
}

func (s *MenuService) UpdateMenu(ctx context.Context, id string, req *models.MenuUpdateRequest) (*models.MenuResponse, error) {
	// Validate input
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Get existing menu
	existingMenu, err := s.menuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		existingMenu.Name = req.Name
	}
	if req.Description != "" {
		existingMenu.Description = req.Description
	}
	if req.Icon != "" {
		existingMenu.Icon = req.Icon
	}
	if req.Path != "" {
		existingMenu.Path = req.Path
	}
	if req.Order != 0 {
		existingMenu.Order = req.Order
	}
	if req.IsActive != nil {
		existingMenu.IsActive = *req.IsActive
	}

	existingMenu.UpdatedAt = time.Now()

	if err := s.menuRepo.Update(ctx, id, existingMenu); err != nil {
		return nil, err
	}

	response := existingMenu.ToResponse()
	return &response, nil
}

func (s *MenuService) DeleteMenu(ctx context.Context, id string) error {
	return s.menuRepo.Delete(ctx, id)
}

// Permission operations

func (s *MenuService) GrantPermission(ctx context.Context, role, menuID, adminID string) error {
	// Validate that menu exists
	_, err := s.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return err
	}

	// Get admin info
	admin, err := s.userRepo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}

	menuObjectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return utils.ErrInvalidID
	}

	adminObjectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return utils.ErrInvalidID
	}

	permission := &models.RoleMenuPermission{
		Role:          role,
		MenuID:        menuObjectID,
		GrantedByID:   adminObjectID,
		GrantedByName: admin.Name,
		CreatedAt:     time.Now(),
	}

	return s.permissionRepo.GrantPermission(ctx, permission)
}

func (s *MenuService) RevokePermission(ctx context.Context, role, menuID string) error {
	return s.permissionRepo.RevokePermission(ctx, role, menuID)
}

func (s *MenuService) GetPermissionsByRole(ctx context.Context, role string) ([]*models.RoleMenuPermissionResponse, error) {
	permissions, err := s.permissionRepo.GetPermissionsByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	var responses []*models.RoleMenuPermissionResponse
	for _, perm := range permissions {
		// Get menu name
		menu, err := s.menuRepo.GetByID(ctx, perm.MenuID.Hex())
		if err != nil {
			continue // Skip if menu not found
		}

		response := perm.ToResponse(menu.Name)
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *MenuService) GetRolesByMenu(ctx context.Context, menuID string) ([]*models.RoleMenuPermissionResponse, error) {
	permissions, err := s.permissionRepo.GetRolesByMenu(ctx, menuID)
	if err != nil {
		return nil, err
	}

	// Get menu name
	menu, err := s.menuRepo.GetByID(ctx, menuID)
	if err != nil {
		return nil, err
	}

	var responses []*models.RoleMenuPermissionResponse
	for _, perm := range permissions {
		response := perm.ToResponse(menu.Name)
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *MenuService) GetAllPermissions(ctx context.Context) ([]*models.RoleMenuPermissionResponse, error) {
	permissions, err := s.permissionRepo.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	var responses []*models.RoleMenuPermissionResponse
	for _, perm := range permissions {
		// Get menu name
		menu, err := s.menuRepo.GetByID(ctx, perm.MenuID.Hex())
		if err != nil {
			continue // Skip if menu not found
		}

		response := perm.ToResponse(menu.Name)
		responses = append(responses, &response)
	}

	return responses, nil
}

// User menu access

func (s *MenuService) GetUserMenus(ctx context.Context, userRole string) ([]*models.UserMenuResponse, error) {
	menus, err := s.menuRepo.GetMenusByRole(ctx, userRole)
	if err != nil {
		return nil, err
	}

	var responses []*models.UserMenuResponse
	for _, menu := range menus {
		response := menu.ToUserMenuResponse()
		responses = append(responses, &response)
	}

	return responses, nil
}

// Permission summary

func (s *MenuService) GetRolePermissionSummary(ctx context.Context) ([]*models.RolePermissionSummary, error) {
	roles := []string{"admin", "liaison", "voice", "finance"}
	var summaries []*models.RolePermissionSummary

	for _, role := range roles {
		menus, err := s.menuRepo.GetMenusByRole(ctx, role)
		if err != nil {
			return nil, err
		}

		var menuResponses []models.MenuResponse
		for _, menu := range menus {
			menuResponses = append(menuResponses, menu.ToResponse())
		}

		summary := &models.RolePermissionSummary{
			Role:      role,
			MenuCount: len(menuResponses),
			Menus:     menuResponses,
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// Helper method to get user by ID
func (s *MenuService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
