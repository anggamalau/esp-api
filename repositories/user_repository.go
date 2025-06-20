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

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() interfaces.UserRepository {
	return &userRepository{
		collection: database.DB.Collection("users"),
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, utils.ErrUserNotFound
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"name":               user.Name,
			"email":              user.Email,
			"role":               user.Role,
			"is_verified":        user.IsVerified,
			"verified_at":        user.VerifiedAt,
			"verified_by":        user.VerifiedBy,
			"verification_notes": user.VerificationNotes,
			"updated_at":         user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return utils.ErrUserAlreadyExists
		}
		return err
	}

	if result.MatchedCount == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.ErrUserNotFound
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

func (r *userRepository) GetPendingUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"is_verified": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) VerifyUser(ctx context.Context, userID, adminID string, notes string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	adminObjectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"is_verified":        true,
			"verified_at":        &now,
			"verified_by":        &adminObjectID,
			"verification_notes": notes,
			"updated_at":         now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": userObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// UpdatePassword updates only the user's password
func (r *userRepository) UpdatePassword(ctx context.Context, userID, hashedPassword string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"password":   hashedPassword,
			"updated_at": now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": userObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

// UpdatePasswordResetInfo updates password reset tracking information
func (r *userRepository) UpdatePasswordResetInfo(ctx context.Context, userID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"last_password_reset": &now,
			"updated_at":          now,
		},
		"$inc": bson.M{
			"password_reset_count": 1,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": userObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}
