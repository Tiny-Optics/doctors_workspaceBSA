package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	GetDB() *mongo.Database
}

type service struct {
	db     *mongo.Client
	dbName string
}

func connectionURI() string {
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return uri
	}

	host := os.Getenv("BLUEPRINT_DB_HOST")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	username := os.Getenv("BLUEPRINT_DB_USERNAME")
	password := os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "27017"
	}
	if username != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}
	return fmt.Sprintf("mongodb://%s:%s", host, port)
}

func resolvedDatabaseName() string {
	if name := os.Getenv("BLUEPRINT_DB_DATABASE"); name != "" {
		return name
	}
	return "doctors_workspace"
}

// ConnectionLabel returns a safe description of the active connection (no credentials).
func ConnectionLabel() string {
	if os.Getenv("MONGO_URI") != "" {
		return "MongoDB Atlas (MONGO_URI)"
	}
	host := os.Getenv("BLUEPRINT_DB_HOST")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "27017"
	}
	return fmt.Sprintf("%s:%s", host, port)
}

// DatabaseName returns the configured database name.
func DatabaseName() string {
	return resolvedDatabaseName()
}

// Connect opens a MongoDB client and returns the configured database. Caller must Disconnect the client.
func Connect(ctx context.Context) (*mongo.Client, *mongo.Database, error) {
	uri := connectionURI()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, fmt.Errorf("connect: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("ping: %w", err)
	}
	dbName := resolvedDatabaseName()
	return client, client.Database(dbName), nil
}

func New() Service {
	client, _, err := Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return &service{
		db:     client,
		dbName: resolvedDatabaseName(),
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetDB() *mongo.Database {
	return s.db.Database(s.dbName)
}
