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
	ErrCategoryNotFound = errors.New("category not found")
	ErrDuplicateSlug    = errors.New("category with this slug already exists")
)

const sopCategoriesCollection = "sop_categories"

// SOPCategoryRepository handles SOP category database operations
type SOPCategoryRepository struct {
	collection *mongo.Collection
}

// NewSOPCategoryRepository creates a new SOPCategoryRepository
func NewSOPCategoryRepository(db *mongo.Database) *SOPCategoryRepository {
	collection := db.Collection(sopCategoriesCollection)

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "display_order", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		// Log error but don't fail - indexes might already exist
		// In production, you'd want proper logging here
	}

	return &SOPCategoryRepository{
		collection: collection,
	}
}

// Create creates a new SOP category
func (r *SOPCategoryRepository) Create(ctx context.Context, category *models.SOPCategory) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	if category.ID.IsZero() {
		category.ID = primitive.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateSlug
		}
		return err
	}

	return nil
}

// FindByID finds a category by ID
func (r *SOPCategoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.SOPCategory, error) {
	var category models.SOPCategory

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil
}

// FindBySlug finds a category by slug
func (r *SOPCategoryRepository) FindBySlug(ctx context.Context, slug string) (*models.SOPCategory, error) {
	var category models.SOPCategory

	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil
}

// SOPCategoryFilter represents filters for listing categories
type SOPCategoryFilter struct {
	IsActive *bool
	Search   string
}

// List returns a paginated list of categories
func (r *SOPCategoryRepository) List(ctx context.Context, filter SOPCategoryFilter, page, limit int) ([]*models.SOPCategory, error) {
	// Build filter
	mongoFilter := bson.M{}

	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}

	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Calculate skip
	skip := (page - 1) * limit

	// Set options - sort by display_order ascending, then by name
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "display_order", Value: 1}, {Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.SOPCategory
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// Count returns the total count of categories matching the filter
func (r *SOPCategoryRepository) Count(ctx context.Context, filter SOPCategoryFilter) (int64, error) {
	mongoFilter := bson.M{}

	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}

	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Update updates a category
func (r *SOPCategoryRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	// Add updated_at timestamp
	update["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrDuplicateSlug
		}
		return err
	}

	if result.MatchedCount == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

// Delete hard deletes a category from the database
func (r *SOPCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

// Activate activates a category
func (r *SOPCategoryRepository) Activate(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{"is_active": true})
}

// Deactivate deactivates a category
func (r *SOPCategoryRepository) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{"is_active": false})
}

// ExistsBySlug checks if a category with the given slug exists
func (r *SOPCategoryRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *primitive.ObjectID) (bool, error) {
	filter := bson.M{"slug": slug}

	if excludeID != nil {
		filter["_id"] = bson.M{"$ne": *excludeID}
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
