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
	ErrSubmissionNotFound = errors.New("submission not found")
)

// RegistrySubmissionRepository handles database operations for registry submissions
type RegistrySubmissionRepository struct {
	collection *mongo.Collection
}

// NewRegistrySubmissionRepository creates a new RegistrySubmissionRepository
func NewRegistrySubmissionRepository(db *mongo.Database) *RegistrySubmissionRepository {
	collection := db.Collection("registry_submissions")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create indexes
	_, _ = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "form_schema_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
	})

	return &RegistrySubmissionRepository{
		collection: collection,
	}
}

// Create creates a new submission
func (r *RegistrySubmissionRepository) Create(ctx context.Context, submission *models.RegistrySubmission) error {
	now := time.Now()
	submission.CreatedAt = now
	submission.UpdatedAt = now
	submission.Status = models.SubmissionStatusSubmitted

	result, err := r.collection.InsertOne(ctx, submission)
	if err != nil {
		return err
	}

	submission.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a submission by ID
func (r *RegistrySubmissionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.RegistrySubmission, error) {
	var submission models.RegistrySubmission
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&submission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSubmissionNotFound
		}
		return nil, err
	}
	return &submission, nil
}

// FindByUser finds all submissions for a specific user with pagination
func (r *RegistrySubmissionRepository) FindByUser(ctx context.Context, userID primitive.ObjectID, page, limit int) ([]*models.RegistrySubmission, int64, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	filter := bson.M{"user_id": userID}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var submissions []*models.RegistrySubmission
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, 0, err
	}

	// Get total count for this user
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

// List retrieves all submissions with pagination and optional filters
func (r *RegistrySubmissionRepository) List(ctx context.Context, page, limit int, filter bson.M) ([]*models.RegistrySubmission, int64, error) {
	skip := (page - 1) * limit

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var submissions []*models.RegistrySubmission
	if err = cursor.All(ctx, &submissions); err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

// Update updates a submission
func (r *RegistrySubmissionRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
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
		return ErrSubmissionNotFound
	}
	return nil
}

// UpdateStatus updates the status of a submission
func (r *RegistrySubmissionRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.SubmissionStatus, reviewedBy *primitive.ObjectID, notes string) error {
	now := time.Now()
	update := bson.M{
		"status":      status,
		"updated_at":  now,
		"reviewed_by": reviewedBy,
		"reviewed_at": now,
	}

	if notes != "" {
		update["review_notes"] = notes
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrSubmissionNotFound
	}
	return nil
}

// Delete deletes a submission
func (r *RegistrySubmissionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrSubmissionNotFound
	}
	return nil
}

// Count returns the total number of submissions
func (r *RegistrySubmissionRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

// CountByStatus returns the count of submissions by status
func (r *RegistrySubmissionRepository) CountByStatus(ctx context.Context, status models.SubmissionStatus) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

// CountByUser returns the total number of submissions for a user
func (r *RegistrySubmissionRepository) CountByUser(ctx context.Context, userID primitive.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
}
