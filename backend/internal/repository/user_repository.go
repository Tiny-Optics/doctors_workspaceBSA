package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

// UserRepository handles database operations for users
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// CreateIndexes creates necessary indexes for the users collection
func (r *UserRepository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "role", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// Check which field caused the duplicate
			if err.Error() != "" {
				return ErrUserAlreadyExists
			}
		}
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email (case-insensitive)
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	// Use case-insensitive regex for email comparison
	emailRegex := bson.M{"$regex": "^" + strings.ReplaceAll(email, ".", "\\.") + "$", "$options": "i"}
	err := r.collection.FindOne(ctx, bson.M{"email": emailRegex}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername finds a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
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
		return ErrUserNotFound
	}
	return nil
}

// Delete deletes a user (hard delete)
func (r *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}
	return nil
}

// Deactivate deactivates a user (soft delete)
func (r *UserRepository) Deactivate(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{"is_active": false})
}

// Activate activates a user
func (r *UserRepository) Activate(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{"is_active": true})
}

// List retrieves users with pagination and filtering
func (r *UserRepository) List(ctx context.Context, filter bson.M, limit, skip int64) ([]*models.User, error) {
	opts := options.Find().
		SetLimit(limit).
		SetSkip(skip).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// Count counts users matching a filter
func (r *UserRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.collection.CountDocuments(ctx, filter)
}

// UpdateLastLogin updates the last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	return r.Update(ctx, id, bson.M{
		"last_login_at":         now,
		"failed_login_attempts": 0,
		"locked_until":          nil,
	})
}

// IncrementFailedLoginAttempts increments the failed login attempts counter
func (r *UserRepository) IncrementFailedLoginAttempts(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$inc": bson.M{"failed_login_attempts": 1},
			"$set": bson.M{"updated_at": time.Now()},
		},
	)
	return err
}

// LockAccount locks a user account until the specified time
func (r *UserRepository) LockAccount(ctx context.Context, id primitive.ObjectID, until time.Time) error {
	return r.Update(ctx, id, bson.M{"locked_until": until})
}

// UnlockAccount unlocks a user account
func (r *UserRepository) UnlockAccount(ctx context.Context, id primitive.ObjectID) error {
	return r.Update(ctx, id, bson.M{
		"locked_until":          nil,
		"failed_login_attempts": 0,
	})
}

// EmailExists checks if an email already exists (case-insensitive)
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	// Use case-insensitive regex for email comparison
	emailRegex := bson.M{"$regex": "^" + strings.ReplaceAll(email, ".", "\\.") + "$", "$options": "i"}
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": emailRegex})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UsernameExists checks if a username already exists
func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
