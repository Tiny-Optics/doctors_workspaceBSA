package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Get environment variables
	host := os.Getenv("BLUEPRINT_DB_HOST")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	username := os.Getenv("BLUEPRINT_DB_USERNAME")
	password := os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")
	database := os.Getenv("BLUEPRINT_DB_DATABASE")

	// Set defaults
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

	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Get database
	db := client.Database(database)

	// Create sample users
	sampleUsers := []models.User{
		{
			Username:     "dr.sarah.smith",
			Email:        "dr.sarah.smith@grooteschuur.co.za",
			PasswordHash: hashPassword("Password123!"),
			Role:         models.RoleHaematologist,
			AdminLevel:   models.AdminLevelNone,
			IsActive:     true,
			Profile: models.UserProfile{
				FirstName:          "Dr. Sarah",
				LastName:           "Smith",
				Institution:        "Groote Schuur Hospital",
				Specialty:          "Haematology",
				Location:           "Cape Town",
				RegistrationNumber: "HPCSA123456",
				PhoneNumber:        "+27 21 404 9111",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:     "dr.michael.jones",
			Email:        "dr.michael.jones@tygerberg.co.za",
			PasswordHash: hashPassword("Password123!"),
			Role:         models.RolePhysician,
			AdminLevel:   models.AdminLevelNone,
			IsActive:     true,
			Profile: models.UserProfile{
				FirstName:          "Dr. Michael",
				LastName:           "Jones",
				Institution:        "Tygerberg Hospital",
				Specialty:          "Internal Medicine",
				Location:           "Cape Town",
				RegistrationNumber: "HPCSA789012",
				PhoneNumber:        "+27 21 938 9111",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:     "dr.emily.davis",
			Email:        "dr.emily.davis@redcross.co.za",
			PasswordHash: hashPassword("Password123!"),
			Role:         models.RoleHaematologist,
			AdminLevel:   models.AdminLevelNone,
			IsActive:     true,
			Profile: models.UserProfile{
				FirstName:          "Dr. Emily",
				LastName:           "Davis",
				Institution:        "Red Cross War Memorial Children's Hospital",
				Specialty:          "Paediatric Haematology",
				Location:           "Cape Town",
				RegistrationNumber: "HPCSA345678",
				PhoneNumber:        "+27 21 658 5111",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:     "james.wilson",
			Email:        "james.wilson@uct.ac.za",
			PasswordHash: hashPassword("Password123!"),
			Role:         models.RoleDataCapturer,
			AdminLevel:   models.AdminLevelNone,
			IsActive:     true,
			Profile: models.UserProfile{
				FirstName:   "James",
				LastName:    "Wilson",
				Institution: "University of Cape Town",
				Specialty:   "Research Data Management",
				Location:    "Cape Town",
				PhoneNumber: "+27 21 650 9111",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:     "dr.robert.brown",
			Email:        "dr.robert.brown@sun.ac.za",
			PasswordHash: hashPassword("Password123!"),
			Role:         models.RolePhysician,
			AdminLevel:   models.AdminLevelNone,
			IsActive:     true,
			Profile: models.UserProfile{
				FirstName:          "Dr. Robert",
				LastName:           "Brown",
				Institution:        "Stellenbosch University",
				Specialty:          "Oncology",
				Location:           "Stellenbosch",
				RegistrationNumber: "HPCSA901234",
				PhoneNumber:        "+27 21 938 9111",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:     "lisa.garcia",
			Email:        "lisa.garcia@wits.ac.za",
			PasswordHash: hashPassword("Password123!"),
			Role:         models.RoleDataCapturer,
			AdminLevel:   models.AdminLevelNone,
			IsActive:     true,
			Profile: models.UserProfile{
				FirstName:   "Lisa",
				LastName:    "Garcia",
				Institution: "University of the Witwatersrand",
				Specialty:   "Clinical Research",
				Location:    "Johannesburg",
				PhoneNumber: "+27 11 717 1000",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:     "dr.david.miller",
			Email:        "dr.david.miller@nhls.ac.za",
			PasswordHash: hashPassword("Password123!"),
			Role:         models.RoleHaematologist,
			AdminLevel:   models.AdminLevelUserManager,
			IsActive:     true,
			Profile: models.UserProfile{
				FirstName:          "Dr. David",
				LastName:           "Miller",
				Institution:        "National Health Laboratory Service",
				Specialty:          "Laboratory Haematology",
				Location:           "Johannesburg",
				RegistrationNumber: "HPCSA567890",
				PhoneNumber:        "+27 11 386 6000",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Username:     "dr.anna.taylor",
			Email:        "dr.anna.taylor@up.ac.za",
			PasswordHash: hashPassword("Password123!"),
			Role:         models.RolePhysician,
			AdminLevel:   models.AdminLevelNone,
			IsActive:     false, // Inactive user for testing
			Profile: models.UserProfile{
				FirstName:          "Dr. Anna",
				LastName:           "Taylor",
				Institution:        "University of Pretoria",
				Specialty:          "Haematology",
				Location:           "Pretoria",
				RegistrationNumber: "HPCSA234567",
				PhoneNumber:        "+27 12 420 9111",
			},
			CreatedAt: time.Now().AddDate(0, 0, -30), // Created 30 days ago
			UpdatedAt: time.Now().AddDate(0, 0, -5),  // Updated 5 days ago
		},
	}

	// Insert users into database
	collection := db.Collection("users")

	fmt.Println("Adding sample users to the database...")

	for i, user := range sampleUsers {
		// Check if user already exists
		var existingUser models.User
		err := collection.FindOne(ctx, map[string]interface{}{
			"email": user.Email,
		}).Decode(&existingUser)

		if err == nil {
			fmt.Printf("User %d: %s already exists, skipping...\n", i+1, user.Email)
			continue
		}

		// Insert new user
		result, err := collection.InsertOne(ctx, user)
		if err != nil {
			log.Printf("Failed to insert user %s: %v", user.Email, err)
			continue
		}

		fmt.Printf("User %d: %s (%s) added successfully with ID: %v\n",
			i+1, user.Email, user.Role, result.InsertedID)
	}

	fmt.Println("\nSample users added successfully!")
	fmt.Println("\nLogin credentials for all users:")
	fmt.Println("Email: [user-email]")
	fmt.Println("Password: Password123!")
	fmt.Println("\nUsers added:")
	for _, user := range sampleUsers {
		status := "Active"
		if !user.IsActive {
			status = "Inactive"
		}
		fmt.Printf("- %s %s (%s) - %s - %s\n",
			user.Profile.FirstName, user.Profile.LastName, user.Role, user.Profile.Institution, status)
	}
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hash)
}
