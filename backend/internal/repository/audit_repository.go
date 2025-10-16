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

// AuditRepository handles database operations for audit logs
type AuditRepository struct {
	collection *mongo.Collection
}

// NewAuditRepository creates a new AuditRepository
func NewAuditRepository(db *mongo.Database) *AuditRepository {
	return &AuditRepository{
		collection: db.Collection("audit_logs"),
	}
}

// CreateIndexes creates necessary indexes for the audit_logs collection
func (r *AuditRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "performed_by", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "action", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "timestamp", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "ip_address", Value: 1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create creates a new audit log entry
func (r *AuditRepository) Create(ctx context.Context, log *models.AuditLog) error {
	log.Timestamp = time.Now()

	result, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return err
	}

	log.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// List retrieves audit logs with pagination and filtering
func (r *AuditRepository) List(ctx context.Context, filter bson.M, limit, skip int64) ([]*models.AuditLog, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*models.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// Count counts audit logs matching a filter
func (r *AuditRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

// FindByUserID retrieves audit logs for a specific user
func (r *AuditRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, limit, skip int64) ([]*models.AuditLog, error) {
	return r.List(ctx, bson.M{"user_id": userID}, limit, skip)
}

// FindByAction retrieves audit logs for a specific action
func (r *AuditRepository) FindByAction(ctx context.Context, action models.AuditAction, limit, skip int64) ([]*models.AuditLog, error) {
	return r.List(ctx, bson.M{"action": action}, limit, skip)
}

// FindByDateRange retrieves audit logs within a date range
func (r *AuditRepository) FindByDateRange(ctx context.Context, start, end time.Time, limit, skip int64) ([]*models.AuditLog, error) {
	filter := bson.M{
		"timestamp": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}
	return r.List(ctx, filter, limit, skip)
}
