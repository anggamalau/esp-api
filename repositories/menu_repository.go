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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type menuRepository struct {
	collection *mongo.Collection
	permRepo   interfaces.PermissionRepository
}

func NewMenuRepository(permRepo interfaces.PermissionRepository) interfaces.MenuRepository {
	return &menuRepository{
		collection: database.DB.Collection("menus"),
		permRepo:   permRepo,
	}
}

func (r *menuRepository) Create(ctx context.Context, menu *models.Menu) error {
	menu.ID = primitive.NewObjectID()
	menu.CreatedAt = time.Now()
	menu.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, menu)
	return err
}

func (r *menuRepository) GetAll(ctx context.Context) ([]*models.Menu, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []*models.Menu
	for cursor.Next(ctx) {
		var menu models.Menu
		if err := cursor.Decode(&menu); err != nil {
			return nil, err
		}
		menus = append(menus, &menu)
	}

	return menus, cursor.Err()
}

func (r *menuRepository) GetByID(ctx context.Context, id string) (*models.Menu, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.ErrInvalidID
	}

	var menu models.Menu
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&menu)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrMenuNotFound
		}
		return nil, err
	}

	return &menu, nil
}

func (r *menuRepository) GetActiveMenus(ctx context.Context) ([]*models.Menu, error) {
	filter := bson.M{"is_active": true}
	opts := options.Find().SetSort(bson.D{{"order", 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []*models.Menu
	for cursor.Next(ctx) {
		var menu models.Menu
		if err := cursor.Decode(&menu); err != nil {
			return nil, err
		}
		menus = append(menus, &menu)
	}

	return menus, cursor.Err()
}

func (r *menuRepository) Update(ctx context.Context, id string, menu *models.Menu) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrInvalidID
	}

	menu.UpdatedAt = time.Now()
	update := bson.M{"$set": menu}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return utils.ErrMenuNotFound
	}

	return nil
}

func (r *menuRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrInvalidID
	}

	// First, revoke all permissions for this menu
	if err := r.permRepo.RevokeAllPermissionsForMenu(ctx, id); err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return utils.ErrMenuNotFound
	}

	return nil
}

func (r *menuRepository) GetMenusOrderedByOrder(ctx context.Context) ([]*models.Menu, error) {
	opts := options.Find().SetSort(bson.D{{"order", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []*models.Menu
	for cursor.Next(ctx) {
		var menu models.Menu
		if err := cursor.Decode(&menu); err != nil {
			return nil, err
		}
		menus = append(menus, &menu)
	}

	return menus, cursor.Err()
}

func (r *menuRepository) GetMenusByRole(ctx context.Context, role string) ([]*models.Menu, error) {
	// Admin has access to all active menus
	if role == "admin" {
		return r.GetActiveMenus(ctx)
	}

	// Get permissions for the role
	permissions, err := r.permRepo.GetPermissionsByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	if len(permissions) == 0 {
		return []*models.Menu{}, nil
	}

	// Extract menu IDs from permissions
	var menuIDs []primitive.ObjectID
	for _, perm := range permissions {
		menuIDs = append(menuIDs, perm.MenuID)
	}

	// Get menus by IDs that are also active
	filter := bson.M{
		"_id":       bson.M{"$in": menuIDs},
		"is_active": true,
	}
	opts := options.Find().SetSort(bson.D{{"order", 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var menus []*models.Menu
	for cursor.Next(ctx) {
		var menu models.Menu
		if err := cursor.Decode(&menu); err != nil {
			return nil, err
		}
		menus = append(menus, &menu)
	}

	return menus, cursor.Err()
}
