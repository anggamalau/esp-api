package handlers

import (
	"context"
	"time"

	"backend/models"
	"backend/services"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	menuService *services.MenuService
}

func NewMenuHandler(menuService *services.MenuService) *MenuHandler {
	return &MenuHandler{
		menuService: menuService,
	}
}

// Menu CRUD operations

// CreateMenu godoc
// @Summary      Create a new menu
// @Description  Create a new menu item (Admin only)
// @Tags         Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      models.MenuCreateRequest  true  "Menu creation data"
// @Success      201      {object}  models.SwaggerResponse{data=models.MenuResponse}
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      403      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /admin/menus [post]
func (h *MenuHandler) CreateMenu(c *fiber.Ctx) error {
	var req models.MenuCreateRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.menuService.CreateMenu(ctx, &req)
	if err != nil {
		if utils.IsValidationError(err) {
			return utils.ValidationErrorResponse(c, err)
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create menu", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Menu created successfully", response)
}

// GetAllMenus godoc
// @Summary      Get all menus
// @Description  Get all menu items (Admin only)
// @Tags         Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  models.SwaggerResponse{data=[]models.MenuResponse}
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      403      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /admin/menus [get]
func (h *MenuHandler) GetAllMenus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.menuService.GetAllMenus(ctx)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch menus", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Menus fetched successfully", response)
}

// GetMenuByID godoc
// @Summary      Get menu by ID
// @Description  Get a specific menu by ID (Admin only)
// @Tags         Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Menu ID"
// @Success      200  {object}  models.SwaggerResponse{data=models.MenuResponse}
// @Failure      400  {object}  models.SwaggerErrorResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /admin/menus/{id} [get]
func (h *MenuHandler) GetMenuByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Menu ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.menuService.GetMenuByID(ctx, id)
	if err != nil {
		if err == utils.ErrMenuNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Menu not found")
		}
		if err == utils.ErrInvalidID {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid menu ID format")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch menu", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Menu fetched successfully", response)
}

// UpdateMenu godoc
// @Summary      Update menu
// @Description  Update a menu item (Admin only)
// @Tags         Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                    true  "Menu ID"
// @Param        request  body      models.MenuUpdateRequest  true  "Menu update data"
// @Success      200      {object}  models.SwaggerResponse{data=models.MenuResponse}
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      403      {object}  models.SwaggerErrorResponse
// @Failure      404      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /admin/menus/{id} [put]
func (h *MenuHandler) UpdateMenu(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Menu ID is required")
	}

	var req models.MenuUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.menuService.UpdateMenu(ctx, id, &req)
	if err != nil {
		if utils.IsValidationError(err) {
			return utils.ValidationErrorResponse(c, err)
		}
		if err == utils.ErrMenuNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Menu not found")
		}
		if err == utils.ErrInvalidID {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid menu ID format")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update menu", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Menu updated successfully", response)
}

// DeleteMenu godoc
// @Summary      Delete menu
// @Description  Delete a menu item (Admin only)
// @Tags         Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Menu ID"
// @Success      200  {object}  models.SwaggerResponse
// @Failure      400  {object}  models.SwaggerErrorResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /admin/menus/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Menu ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.menuService.DeleteMenu(ctx, id)
	if err != nil {
		if err == utils.ErrMenuNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Menu not found")
		}
		if err == utils.ErrInvalidID {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid menu ID format")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete menu", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Menu deleted successfully", nil)
}

// Permission operations

// GrantPermission godoc
// @Summary      Grant menu permission to role
// @Description  Grant access to a menu for a specific role (Admin only)
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        role     path      string  true   "Role name"
// @Param        menuId   path      string  true   "Menu ID"
// @Success      200      {object}  models.SwaggerResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      403      {object}  models.SwaggerErrorResponse
// @Failure      404      {object}  models.SwaggerErrorResponse
// @Failure      409      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /admin/roles/{role}/menus/{menuId} [post]
func (h *MenuHandler) GrantPermission(c *fiber.Ctx) error {
	role := c.Params("role")
	menuID := c.Params("menuId")
	adminID := c.Locals("userID").(string)

	if role == "" || menuID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Role and menu ID are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.menuService.GrantPermission(ctx, role, menuID, adminID)
	if err != nil {
		if err == utils.ErrMenuNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Menu not found")
		}
		if err == utils.ErrPermissionAlreadyExists {
			return utils.ErrorResponse(c, fiber.StatusConflict, "Permission already exists")
		}
		if err == utils.ErrInvalidID {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to grant permission", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Permission granted successfully", nil)
}

// RevokePermission godoc
// @Summary      Revoke menu permission from role
// @Description  Revoke access to a menu for a specific role (Admin only)
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        role     path      string  true   "Role name"
// @Param        menuId   path      string  true   "Menu ID"
// @Success      200      {object}  models.SwaggerResponse
// @Failure      400      {object}  models.SwaggerErrorResponse
// @Failure      401      {object}  models.SwaggerErrorResponse
// @Failure      403      {object}  models.SwaggerErrorResponse
// @Failure      404      {object}  models.SwaggerErrorResponse
// @Failure      500      {object}  models.SwaggerErrorResponse
// @Router       /admin/roles/{role}/menus/{menuId} [delete]
func (h *MenuHandler) RevokePermission(c *fiber.Ctx) error {
	role := c.Params("role")
	menuID := c.Params("menuId")

	if role == "" || menuID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Role and menu ID are required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.menuService.RevokePermission(ctx, role, menuID)
	if err != nil {
		if err == utils.ErrPermissionNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Permission not found")
		}
		if err == utils.ErrInvalidID {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID format")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to revoke permission", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Permission revoked successfully", nil)
}

// GetPermissionsByRole godoc
// @Summary      Get permissions by role
// @Description  Get all menu permissions for a specific role (Admin only)
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        role   path      string  true  "Role name"
// @Success      200    {object}  models.SwaggerResponse{data=[]models.RoleMenuPermissionResponse}
// @Failure      400    {object}  models.SwaggerErrorResponse
// @Failure      401    {object}  models.SwaggerErrorResponse
// @Failure      403    {object}  models.SwaggerErrorResponse
// @Failure      500    {object}  models.SwaggerErrorResponse
// @Router       /admin/roles/{role}/menus [get]
func (h *MenuHandler) GetPermissionsByRole(c *fiber.Ctx) error {
	role := c.Params("role")
	if role == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Role is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.menuService.GetPermissionsByRole(ctx, role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch permissions", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Permissions fetched successfully", response)
}

// GetRolesByMenu godoc
// @Summary      Get roles by menu
// @Description  Get all roles that have access to a specific menu (Admin only)
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Menu ID"
// @Success      200  {object}  models.SwaggerResponse{data=[]models.RoleMenuPermissionResponse}
// @Failure      400  {object}  models.SwaggerErrorResponse
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Failure      404  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /admin/menus/{id}/roles [get]
func (h *MenuHandler) GetRolesByMenu(c *fiber.Ctx) error {
	menuID := c.Params("id")
	if menuID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Menu ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.menuService.GetRolesByMenu(ctx, menuID)
	if err != nil {
		if err == utils.ErrMenuNotFound {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Menu not found")
		}
		if err == utils.ErrInvalidID {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid menu ID format")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch roles", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Roles fetched successfully", response)
}

// GetAllPermissions godoc
// @Summary      Get all permissions
// @Description  Get all role-menu permissions (Admin only)
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SwaggerResponse{data=[]models.RoleMenuPermissionResponse}
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /admin/roles/permissions [get]
func (h *MenuHandler) GetAllPermissions(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.menuService.GetAllPermissions(ctx)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch permissions", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Permissions fetched successfully", response)
}

// GetRolePermissionSummary godoc
// @Summary      Get role permission summary
// @Description  Get summary of permissions for all roles (Admin only)
// @Tags         Permission Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SwaggerResponse{data=[]models.RolePermissionSummary}
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      403  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /admin/roles/summary [get]
func (h *MenuHandler) GetRolePermissionSummary(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := h.menuService.GetRolePermissionSummary(ctx)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch role summary", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Role summary fetched successfully", response)
}

// User menu access

// GetUserMenus godoc
// @Summary      Get user accessible menus
// @Description  Get menus accessible by the current user based on their role
// @Tags         User Menu Access
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.SwaggerResponse{data=[]models.UserMenuResponse}
// @Failure      401  {object}  models.SwaggerErrorResponse
// @Failure      500  {object}  models.SwaggerErrorResponse
// @Router       /users/menus [get]
func (h *MenuHandler) GetUserMenus(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user to determine role
	user, err := h.menuService.GetUserByID(ctx, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "User not found")
	}

	response, err := h.menuService.GetUserMenus(ctx, user.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch user menus", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "User menus fetched successfully", response)
}
