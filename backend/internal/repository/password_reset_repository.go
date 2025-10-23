package repository

import (
	"context"
	"errors"
	"time"

	"backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrPasswordResetTokenNotFound = errors.New("password reset token not found")
	ErrPasswordResetTokenExpired  = errors.New("password reset token expired")
	ErrPasswordResetTokenUsed     = errors.New("password reset token already used")
)

// PasswordResetRepository handles password reset token operations
type PasswordResetRepository struct {
	collection *mongo.Collection
}

// NewPasswordResetRepository creates a new PasswordResetRepository
func NewPasswordResetRepository(db *mongo.Database) *PasswordResetRepository {
	return &PasswordResetRepository{
		collection: db.Collection("password_reset_tokens"),
	}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(ctx context.Context, token *models.PasswordResetToken) error {
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

// FindByToken finds a password reset token by token string
func (r *PasswordResetRepository) FindByToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&resetToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPasswordResetTokenNotFound
		}
		return nil, err
	}
	return &resetToken, nil
}

// FindByCode finds a password reset token by verification code
func (r *PasswordResetRepository) FindByCode(ctx context.Context, code string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&resetToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPasswordResetTokenNotFound
		}
		return nil, err
	}
	return &resetToken, nil
}

// FindByUserID finds the most recent password reset token for a user
func (r *PasswordResetRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	opts := options.FindOne().SetSort(bson.M{"created_at": -1})
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}, opts).Decode(&resetToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrPasswordResetTokenNotFound
		}
		return nil, err
	}
	return &resetToken, nil
}

// MarkAsUsed marks a password reset token as used
func (r *PasswordResetRepository) MarkAsUsed(ctx context.Context, tokenID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": tokenID},
		bson.M{"$set": bson.M{"used": true}},
	)
	return err
}

// CleanupExpiredTokens removes expired password reset tokens
func (r *PasswordResetRepository) CleanupExpiredTokens(ctx context.Context) error {
	_, err := r.collection.DeleteMany(
		ctx,
		bson.M{
			"expires_at": bson.M{"$lt": time.Now()},
		},
	)
	return err
}

// CountRecentRequests counts password reset requests for a user in the last hour
func (r *PasswordResetRepository) CountRecentRequests(ctx context.Context, userID primitive.ObjectID, since time.Time) (int64, error) {
	count, err := r.collection.CountDocuments(
		ctx,
		bson.M{
			"user_id":    userID,
			"created_at": bson.M{"$gte": since},
		},
	)
	return count, err
}

// CountRecentRequestsByIP counts password reset requests from an IP in the last hour
func (r *PasswordResetRepository) CountRecentRequestsByIP(ctx context.Context, ipAddress string, since time.Time) (int64, error) {
	count, err := r.collection.CountDocuments(
		ctx,
		bson.M{
			"ip_address": ipAddress,
			"created_at": bson.M{"$gte": since},
		},
	)
	return count, err
}
