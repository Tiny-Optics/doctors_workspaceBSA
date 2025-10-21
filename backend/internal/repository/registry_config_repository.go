package repository

import (
	"context"
	"errors"
	"time"

	"backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrRegistryConfigNotFound = errors.New("registry configuration not found")
)

// RegistryConfigRepository handles database operations for registry configuration
type RegistryConfigRepository struct {
	collection *mongo.Collection
}

// NewRegistryConfigRepository creates a new RegistryConfigRepository
func NewRegistryConfigRepository(db *mongo.Database) *RegistryConfigRepository {
	collection := db.Collection("registry_config")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create index for updated_at
	_, _ = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "updated_at", Value: -1}},
	})

	return &RegistryConfigRepository{
		collection: collection,
	}
}

// GetConfig retrieves the singleton registry configuration
func (r *RegistryConfigRepository) GetConfig(ctx context.Context) (*models.RegistryConfig, error) {
	var config models.RegistryConfig
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrRegistryConfigNotFound
		}
		return nil, err
	}
	return &config, nil
}

// CreateOrUpdate creates or updates the registry configuration (upsert)
func (r *RegistryConfigRepository) CreateOrUpdate(ctx context.Context, config *models.RegistryConfig) error {
	now := time.Now()

	// Check if config exists
	existingConfig, err := r.GetConfig(ctx)

	if err == ErrRegistryConfigNotFound {
		// Create new config
		config.CreatedAt = now
		config.UpdatedAt = now

		result, err := r.collection.InsertOne(ctx, config)
		if err != nil {
			return err
		}
		config.ID = result.InsertedID.(primitive.ObjectID)
		return nil
	}

	if err != nil {
		return err
	}

	// Update existing config
	config.ID = existingConfig.ID
	config.CreatedAt = existingConfig.CreatedAt
	config.UpdatedAt = now

	_, err = r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": existingConfig.ID},
		config,
	)

	return err
}

// Update updates specific fields of the configuration
func (r *RegistryConfigRepository) Update(ctx context.Context, update bson.M) error {
	config, err := r.GetConfig(ctx)
	if err != nil {
		return err
	}

	update["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": update},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrRegistryConfigNotFound
	}
	return nil
}

// AddNotificationEmail adds an email to the notification list
func (r *RegistryConfigRepository) AddNotificationEmail(ctx context.Context, email string) error {
	config, err := r.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Check if email already exists
	for _, existingEmail := range config.NotificationEmails {
		if existingEmail == email {
			return nil // Already exists, no need to add
		}
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{
			"$push": bson.M{"notification_emails": email},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)

	return err
}

// RemoveNotificationEmail removes an email from the notification list
func (r *RegistryConfigRepository) RemoveNotificationEmail(ctx context.Context, email string) error {
	config, err := r.GetConfig(ctx)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{
			"$pull": bson.M{"notification_emails": email},
			"$set":  bson.M{"updated_at": time.Now()},
		},
	)

	return err
}

// ConfigExists checks if a configuration exists
func (r *RegistryConfigRepository) ConfigExists(ctx context.Context) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteConfig deletes the configuration (for testing purposes)
func (r *RegistryConfigRepository) DeleteConfig(ctx context.Context) error {
	config, err := r.GetConfig(ctx)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": config.ID})
	return err
}
