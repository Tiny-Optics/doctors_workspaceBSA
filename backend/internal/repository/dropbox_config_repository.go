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
	ErrDropboxConfigNotFound = errors.New("dropbox configuration not found")
)

// DropboxConfigRepository handles database operations for Dropbox configuration
type DropboxConfigRepository struct {
	collection *mongo.Collection
}

// NewDropboxConfigRepository creates a new DropboxConfigRepository
func NewDropboxConfigRepository(db *mongo.Database) *DropboxConfigRepository {
	return &DropboxConfigRepository{
		collection: db.Collection("dropbox_config"),
	}
}

// GetConfig retrieves the singleton Dropbox configuration
func (r *DropboxConfigRepository) GetConfig(ctx context.Context) (*models.DropboxConfig, error) {
	var config models.DropboxConfig
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrDropboxConfigNotFound
		}
		return nil, err
	}
	return &config, nil
}

// CreateConfig creates the initial Dropbox configuration
// This should only be called once during initial setup
func (r *DropboxConfigRepository) CreateConfig(ctx context.Context, config *models.DropboxConfig) error {
	// Check if config already exists
	existing, err := r.GetConfig(ctx)
	if err != nil && err != ErrDropboxConfigNotFound {
		return err
	}
	if existing != nil {
		return errors.New("dropbox configuration already exists")
	}

	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now
	config.LastRefreshAttempt = now

	result, err := r.collection.InsertOne(ctx, config)
	if err != nil {
		return err
	}

	config.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// UpdateTokens updates the access token, refresh token (if provided), and expiry
func (r *DropboxConfigRepository) UpdateTokens(ctx context.Context, accessToken string, refreshToken string, expiresIn int) error {
	update := bson.M{
		"access_token": accessToken,
		"token_expiry": time.Now().Add(time.Duration(expiresIn) * time.Second),
		"updated_at":   time.Now(),
	}

	// Only update refresh token if provided (it's only returned on initial auth)
	if refreshToken != "" {
		update["refresh_token"] = refreshToken
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{"$set": update},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrDropboxConfigNotFound
	}
	return nil
}

// UpdateConfig updates the entire configuration
func (r *DropboxConfigRepository) UpdateConfig(ctx context.Context, config *models.DropboxConfig) error {
	config.UpdatedAt = time.Now()

	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": config.ID},
		config,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrDropboxConfigNotFound
	}
	return nil
}

// UpdateHealth updates the health monitoring fields
func (r *DropboxConfigRepository) UpdateHealth(ctx context.Context, isConnected bool, lastError string) error {
	update := bson.M{
		"is_connected":         isConnected,
		"last_refresh_attempt": time.Now(),
		"last_error":           lastError,
		"updated_at":           time.Now(),
	}

	if isConnected {
		update["last_refresh_success"] = time.Now()
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{"$set": update},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrDropboxConfigNotFound
	}
	return nil
}

// IncrementFailures increments the consecutive failures counter
func (r *DropboxConfigRepository) IncrementFailures(ctx context.Context) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{
			"$inc": bson.M{"consecutive_failures": 1},
			"$set": bson.M{
				"is_connected": false,
				"updated_at":   time.Now(),
			},
		},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrDropboxConfigNotFound
	}
	return nil
}

// ResetFailures resets the consecutive failures counter
func (r *DropboxConfigRepository) ResetFailures(ctx context.Context) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{},
		bson.M{
			"$set": bson.M{
				"consecutive_failures": 0,
				"is_connected":         true,
				"last_refresh_success": time.Now(),
				"updated_at":           time.Now(),
			},
		},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrDropboxConfigNotFound
	}
	return nil
}

// DeleteConfig deletes the Dropbox configuration
// Use with caution - this will require re-authorization
func (r *DropboxConfigRepository) DeleteConfig(ctx context.Context) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrDropboxConfigNotFound
	}
	return nil
}

// EnsureIndexes creates necessary indexes for the collection
func (r *DropboxConfigRepository) EnsureIndexes(ctx context.Context) error {
	// We don't need indexes since this is a singleton collection
	// But we can add one for good measure
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetUnique(false),
	}

	_, err := r.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}
