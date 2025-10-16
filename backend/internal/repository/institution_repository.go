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
	ErrInstitutionNotFound    = errors.New("institution not found")
	ErrDuplicateInstitution   = errors.New("institution with this name already exists")
)

// InstitutionRepository handles institution database operations
type InstitutionRepository struct {
	collection *mongo.Collection
}

// NewInstitutionRepository creates a new InstitutionRepository
func NewInstitutionRepository(db *mongo.Database) *InstitutionRepository {
	return &InstitutionRepository{
		collection: db.Collection("institutions"),
	}
}

// Create creates a new institution
func (r *InstitutionRepository) Create(ctx context.Context, institution *models.Institution) error {
	// Check if institution with same name already exists
	existing, _ := r.FindByName(ctx, institution.Name)
	if existing != nil {
		return ErrDuplicateInstitution
	}

	institution.ID = primitive.NewObjectID()
	institution.CreatedAt = time.Now()
	institution.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, institution)
	return err
}

// FindByID finds an institution by ID
func (r *InstitutionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Institution, error) {
	var institution models.Institution
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&institution)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInstitutionNotFound
		}
		return nil, err
	}
	return &institution, nil
}

// FindByName finds an institution by name
func (r *InstitutionRepository) FindByName(ctx context.Context, name string) (*models.Institution, error) {
	var institution models.Institution
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&institution)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrInstitutionNotFound
		}
		return nil, err
	}
	return &institution, nil
}

// Update updates an institution
func (r *InstitutionRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
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
		return ErrInstitutionNotFound
	}
	return nil
}

// Delete deletes an institution
func (r *InstitutionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrInstitutionNotFound
	}
	return nil
}

// List retrieves institutions with pagination and filtering
func (r *InstitutionRepository) List(ctx context.Context, filter bson.M, limit, skip int64) ([]*models.Institution, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "name", Value: 1}}) // Sort by name alphabetically

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var institutions []*models.Institution
	if err := cursor.All(ctx, &institutions); err != nil {
		return nil, err
	}

	return institutions, nil
}

// Count counts institutions matching a filter
func (r *InstitutionRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

// Activate activates an institution
func (r *InstitutionRepository) Activate(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{"is_active": true})
}

// Deactivate deactivates an institution
func (r *InstitutionRepository) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{"is_active": false})
}

// NameExists checks if an institution name already exists
func (r *InstitutionRepository) NameExists(ctx context.Context, name string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{"name": name}
	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

