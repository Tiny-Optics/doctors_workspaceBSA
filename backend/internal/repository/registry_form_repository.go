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
	ErrFormSchemaNotFound = errors.New("form schema not found")
	ErrNoActiveForm       = errors.New("no active form schema found")
)

// RegistryFormRepository handles database operations for registry form schemas
type RegistryFormRepository struct {
	collection *mongo.Collection
}

// NewRegistryFormRepository creates a new RegistryFormRepository
func NewRegistryFormRepository(db *mongo.Database) *RegistryFormRepository {
	collection := db.Collection("registry_form_schemas")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create indexes
	_, _ = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	})

	return &RegistryFormRepository{
		collection: collection,
	}
}

// Create creates a new form schema
func (r *RegistryFormRepository) Create(ctx context.Context, schema *models.RegistryFormSchema) error {
	now := time.Now()
	schema.CreatedAt = now
	schema.UpdatedAt = now
	schema.IsActive = false // New forms are inactive by default

	result, err := r.collection.InsertOne(ctx, schema)
	if err != nil {
		return err
	}

	schema.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a form schema by ID
func (r *RegistryFormRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.RegistryFormSchema, error) {
	var schema models.RegistryFormSchema
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&schema)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrFormSchemaNotFound
		}
		return nil, err
	}
	return &schema, nil
}

// FindActive finds the currently active form schema
func (r *RegistryFormRepository) FindActive(ctx context.Context) (*models.RegistryFormSchema, error) {
	var schema models.RegistryFormSchema
	err := r.collection.FindOne(ctx, bson.M{"is_active": true}).Decode(&schema)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoActiveForm
		}
		return nil, err
	}
	return &schema, nil
}

// List retrieves all form schemas with pagination
func (r *RegistryFormRepository) List(ctx context.Context, page, limit int) ([]*models.RegistryFormSchema, int64, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var schemas []*models.RegistryFormSchema
	if err = cursor.All(ctx, &schemas); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return schemas, total, nil
}

// Update updates a form schema
func (r *RegistryFormRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrFormSchemaNotFound
	}
	return nil
}

// Delete deletes a form schema
func (r *RegistryFormRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrFormSchemaNotFound
	}
	return nil
}

// SetActive sets a form schema as active and deactivates all others
func (r *RegistryFormRepository) SetActive(ctx context.Context, id primitive.ObjectID) error {
	// For development without replica set, do operations without transaction
	// In production, consider using a replica set for atomic operations

	now := time.Now()

	// Deactivate all forms
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{},
		bson.M{"$set": bson.M{"is_active": false, "updated_at": now}},
	)
	if err != nil {
		return err
	}

	// Activate the specified form
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"is_active": true, "updated_at": now}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrFormSchemaNotFound
	}

	return nil
}

// DeactivateAll deactivates all form schemas
func (r *RegistryFormRepository) DeactivateAll(ctx context.Context) error {
	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{},
		bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now()}},
	)
	return err
}

// Count returns the total number of form schemas
func (r *RegistryFormRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
