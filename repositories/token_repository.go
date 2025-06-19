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

type tokenRepository struct {
	collection *mongo.Collection
}

func NewTokenRepository() interfaces.TokenRepository {
	return &tokenRepository{
		collection: database.DB.Collection("refresh_tokens"),
	}
}

func (r *tokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	token.ID = primitive.NewObjectID()
	token.CreatedAt = time.Now()
	token.IsRevoked = false

	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *tokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&refreshToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.ErrTokenNotFound
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *tokenRepository) GetByUserID(ctx context.Context, userID string) ([]models.RefreshToken, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":    objectID,
		"is_revoked": false,
		"expires_at": bson.M{"$gt": time.Now()},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tokens []models.RefreshToken
	err = cursor.All(ctx, &tokens)
	return tokens, err
}

func (r *tokenRepository) RevokeToken(ctx context.Context, token string) error {
	filter := bson.M{"token": token}
	update := bson.M{"$set": bson.M{"is_revoked": true}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return utils.ErrTokenNotFound
	}

	return nil
}

func (r *tokenRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"user_id": objectID, "is_revoked": false}
	update := bson.M{"$set": bson.M{"is_revoked": true}}

	_, err = r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	filter := bson.M{
		"$or": []bson.M{
			{"expires_at": bson.M{"$lt": time.Now()}},
			{"is_revoked": true},
		},
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
