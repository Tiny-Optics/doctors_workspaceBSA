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
	ErrReferralConfigNotFound = errors.New("referral configuration not found")
)

// ReferralConfigRepository handles database operations for REDCap referral configuration
type ReferralConfigRepository struct {
	collection *mongo.Collection
}

// NewReferralConfigRepository creates a new ReferralConfigRepository
func NewReferralConfigRepository(db *mongo.Database) *ReferralConfigRepository {
	return &ReferralConfigRepository{
		collection: db.Collection("referral_config"),
	}
}

// GetConfig retrieves the singleton referral configuration
func (r *ReferralConfigRepository) GetConfig(ctx context.Context) (*models.ReferralConfig, error) {
	var config models.ReferralConfig
	err := r.collection.FindOne(ctx, bson.M{}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrReferralConfigNotFound
		}
		return nil, err
	}
	return &config, nil
}

// CreateConfig creates the initial referral configuration
// This should only be called once during initial setup
func (r *ReferralConfigRepository) CreateConfig(ctx context.Context, config *models.ReferralConfig) error {
	// Check if config already exists
	existing, err := r.GetConfig(ctx)
	if err != nil && err != ErrReferralConfigNotFound {
		return err
	}
	if existing != nil {
		return errors.New("referral configuration already exists")
	}

	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, config)
	if err != nil {
		return err
	}

	config.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// UpdateConfig updates the referral configuration
func (r *ReferralConfigRepository) UpdateConfig(ctx context.Context, config *models.ReferralConfig) error {
	config.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"redcap_url": config.RedCapURL,
			"is_enabled": config.IsEnabled,
			"updated_at": config.UpdatedAt,
			"updated_by": config.UpdatedBy,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrReferralConfigNotFound
	}

	return nil
}
