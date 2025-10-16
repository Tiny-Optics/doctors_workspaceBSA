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
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

// SessionRepository handles database operations for sessions
type SessionRepository struct {
	collection *mongo.Collection
}

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{
		collection: db.Collection("sessions"),
	}
}

// CreateIndexes creates necessary indexes for the sessions collection
func (r *SessionRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys:    bson.D{{Key: "token", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "refresh_token", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0), // TTL index
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create creates a new session
func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	session.CreatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}

	session.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByToken finds a session by token
func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if session.IsExpired() {
		return nil, ErrSessionExpired
	}

	return &session, nil
}

// FindByRefreshToken finds a session by refresh token
func (r *SessionRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	var session models.Session
	err := r.collection.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	if session.IsRefreshExpired() {
		return nil, ErrSessionExpired
	}

	return &session, nil
}

// Delete deletes a session
func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrSessionNotFound
	}
	return nil
}

// DeleteAllByUserID deletes all sessions for a user
func (r *SessionRepository) DeleteAllByUserID(ctx context.Context, userID primitive.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}

// DeleteExpired deletes all expired sessions
func (r *SessionRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
