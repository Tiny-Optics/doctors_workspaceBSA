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

var (
	host     = os.Getenv("BLUEPRINT_DB_HOST")
	port     = os.Getenv("BLUEPRINT_DB_PORT")
	username = os.Getenv("BLUEPRINT_DB_USERNAME")
	password = os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")
	database = os.Getenv("BLUEPRINT_DB_DATABASE")
)

func New() Service {
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "27017"
	}
	if database == "" {
		database = "doctors_workspace"
	}

	// Build connection URI
	var uri string
	if username != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", host, port)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	return &service{
		db:     client,
		dbName: database,
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
