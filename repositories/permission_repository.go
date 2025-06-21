package repositories

import (
	"context"
	"time"

	"backend/database"
	"backend/models"
	"backend/repositories/interfaces"
	"backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type permissionRepository struct {
	collection *mongo.Collection
}

func NewPermissionRepository() interfaces.PermissionRepository {
	return &permissionRepository{
		collection: database.DB.Collection("role_menu_permissions"),
	}
}

func (r *permissionRepository) GrantPermission(ctx context.Context, permission *models.RoleMenuPermission) error {
	// Check if permission already exists
	exists, err := r.CheckPermission(ctx, permission.Role, permission.MenuID.Hex())
	if err != nil {
		return err
	}
	if exists {
		return utils.ErrPermissionAlreadyExists
	}

	permission.ID = primitive.NewObjectID()
	permission.CreatedAt = time.Now()

	_, err = r.collection.InsertOne(ctx, permission)
	return err
}

func (r *permissionRepository) RevokePermission(ctx context.Context, role, menuID string) error {
	menuObjectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return utils.ErrInvalidID
	}

	filter := bson.M{
		"role":    role,
		"menu_id": menuObjectID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return utils.ErrPermissionNotFound
	}

	return nil
}

func (r *permissionRepository) GetPermissionsByRole(ctx context.Context, role string) ([]*models.RoleMenuPermission, error) {
	filter := bson.M{"role": role}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*models.RoleMenuPermission
	for cursor.Next(ctx) {
		var permission models.RoleMenuPermission
		if err := cursor.Decode(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, cursor.Err()
}

func (r *permissionRepository) GetRolesByMenu(ctx context.Context, menuID string) ([]*models.RoleMenuPermission, error) {
	menuObjectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return nil, utils.ErrInvalidID
	}

	filter := bson.M{"menu_id": menuObjectID}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*models.RoleMenuPermission
	for cursor.Next(ctx) {
		var permission models.RoleMenuPermission
		if err := cursor.Decode(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, cursor.Err()
}

func (r *permissionRepository) GetAllPermissions(ctx context.Context) ([]*models.RoleMenuPermission, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var permissions []*models.RoleMenuPermission
	for cursor.Next(ctx) {
		var permission models.RoleMenuPermission
		if err := cursor.Decode(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, &permission)
	}

	return permissions, cursor.Err()
}

func (r *permissionRepository) CheckPermission(ctx context.Context, role, menuID string) (bool, error) {
	menuObjectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return false, utils.ErrInvalidID
	}

	filter := bson.M{
		"role":    role,
		"menu_id": menuObjectID,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *permissionRepository) RevokeAllPermissionsForMenu(ctx context.Context, menuID string) error {
	menuObjectID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return utils.ErrInvalidID
	}

	filter := bson.M{"menu_id": menuObjectID}
	_, err = r.collection.DeleteMany(ctx, filter)
	return err
}

func (r *permissionRepository) RevokeAllPermissionsForRole(ctx context.Context, role string) error {
	filter := bson.M{"role": role}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
