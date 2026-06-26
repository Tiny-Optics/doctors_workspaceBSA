package repository

import (
	"context"
	"time"

	"backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const workingPartyCategoriesCollection = "working_party_categories"

// WorkingPartyCategoryRepository handles working party category database operations
type WorkingPartyCategoryRepository struct {
	collection *mongo.Collection
}

// NewWorkingPartyCategoryRepository creates a new WorkingPartyCategoryRepository
func NewWorkingPartyCategoryRepository(db *mongo.Database) *WorkingPartyCategoryRepository {
	collection := db.Collection(workingPartyCategoriesCollection)

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
		// Indexes might already exist
		_ = err
	}

	return &WorkingPartyCategoryRepository{
		collection: collection,
	}
}

// Create creates a new working party category
func (r *WorkingPartyCategoryRepository) Create(ctx context.Context, category *models.WorkingPartyCategory) error {
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
func (r *WorkingPartyCategoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.WorkingPartyCategory, error) {
	var category models.WorkingPartyCategory

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
func (r *WorkingPartyCategoryRepository) FindBySlug(ctx context.Context, slug string) (*models.WorkingPartyCategory, error) {
	var category models.WorkingPartyCategory

	err := r.collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil
}

// WorkingPartyCategoryFilter represents filters for listing categories
type WorkingPartyCategoryFilter struct {
	IsActive *bool
	Search   string
}

// List returns a paginated list of categories
func (r *WorkingPartyCategoryRepository) List(ctx context.Context, filter WorkingPartyCategoryFilter, page, limit int) ([]*models.WorkingPartyCategory, error) {
	mongoFilter := bson.M{}

	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}

	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"slug": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "display_order", Value: 1}, {Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []*models.WorkingPartyCategory
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

// Count returns the total count of categories matching the filter
func (r *WorkingPartyCategoryRepository) Count(ctx context.Context, filter WorkingPartyCategoryFilter) (int64, error) {
	mongoFilter := bson.M{}

	if filter.IsActive != nil {
		mongoFilter["is_active"] = *filter.IsActive
	}

	if filter.Search != "" {
		mongoFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"slug": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	count, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Update updates a category
func (r *WorkingPartyCategoryRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
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
func (r *WorkingPartyCategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
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
func (r *WorkingPartyCategoryRepository) Activate(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{"is_active": true})
}

// Deactivate deactivates a category
func (r *WorkingPartyCategoryRepository) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{"is_active": false})
}

// ExistsBySlug checks if a category with the given slug exists
func (r *WorkingPartyCategoryRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *primitive.ObjectID) (bool, error) {
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
