package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/service"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("Starting database seed...")

	// Get MongoDB connection details from environment
	host := os.Getenv("BLUEPRINT_DB_HOST")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	database := os.Getenv("BLUEPRINT_DB_DATABASE")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "27017"
	}
	if database == "" {
		database = "doctors_workspace"
	}

	// Connect to MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(ctx)

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}
	log.Println("Connected to MongoDB successfully")

	// Get database
	db := client.Database(database)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	auditRepo := repository.NewAuditRepository(db)

	// Create indexes
	log.Println("Creating indexes...")
	if err := userRepo.CreateIndexes(ctx); err != nil {
		log.Fatal("Failed to create user indexes:", err)
	}
	if err := sessionRepo.CreateIndexes(ctx); err != nil {
		log.Fatal("Failed to create session indexes:", err)
	}
	if err := auditRepo.CreateIndexes(ctx); err != nil {
		log.Fatal("Failed to create audit indexes:", err)
	}
	log.Println("Indexes created successfully")

	// Initialize auth service
	authService := service.NewAuthService(userRepo, sessionRepo, auditRepo)

	// Check if super admin already exists
	superAdminEmail := os.Getenv("SUPER_ADMIN_EMAIL")
	if superAdminEmail == "" {
		superAdminEmail = "admin@bloodsa.org.za"
	}

	existingAdmin, err := userRepo.FindByEmail(ctx, superAdminEmail)
	if err == nil && existingAdmin != nil {
		log.Printf("Super admin already exists: %s (%s)\n", existingAdmin.Email, existingAdmin.Username)
		log.Println("Seed completed (no changes needed)")
		return
	}

	// Get super admin details from environment or use defaults
	superAdminUsername := os.Getenv("SUPER_ADMIN_USERNAME")
	if superAdminUsername == "" {
		superAdminUsername = "superadmin"
	}

	superAdminPassword := os.Getenv("SUPER_ADMIN_PASSWORD")
	if superAdminPassword == "" {
		superAdminPassword = "BloodSA2025!"
		log.Println("WARNING: Using default password. Please change it after first login!")
	}

	superAdminFirstName := os.Getenv("SUPER_ADMIN_FIRST_NAME")
	if superAdminFirstName == "" {
		superAdminFirstName = "Super"
	}

	superAdminLastName := os.Getenv("SUPER_ADMIN_LAST_NAME")
	if superAdminLastName == "" {
		superAdminLastName = "Admin"
	}

	// Hash password
	passwordHash, err := authService.HashPassword(superAdminPassword)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Create super admin user
	superAdmin := &models.User{
		Username:     superAdminUsername,
		Email:        superAdminEmail,
		PasswordHash: passwordHash,
		Role:         models.RoleAdmin,
		AdminLevel:   models.AdminLevelSuperAdmin,
		IsActive:     true,
		Profile: models.UserProfile{
			FirstName:   superAdminFirstName,
			LastName:    superAdminLastName,
			Institution: "BLOODSA",
			Location:    "South Africa",
		},
	}

	// Create the user
	if err := userRepo.Create(ctx, superAdmin); err != nil {
		log.Fatal("Failed to create super admin:", err)
	}

	// Log audit entry
	auditLog := &models.AuditLog{
		UserID:      &superAdmin.ID,
		PerformedBy: &superAdmin.ID,
		Action:      models.AuditActionUserCreated,
		IPAddress:   "system",
		Details: map[string]interface{}{
			"username":    superAdmin.Username,
			"email":       superAdmin.Email,
			"role":        superAdmin.Role,
			"admin_level": superAdmin.AdminLevel,
			"created_by":  "seed_script",
		},
	}
	if err := auditRepo.Create(ctx, auditLog); err != nil {
		log.Println("Warning: Failed to create audit log:", err)
	}

	log.Println("✅ Super admin created successfully!")
	log.Println("================================================")
	log.Printf("Username: %s\n", superAdmin.Username)
	log.Printf("Email: %s\n", superAdmin.Email)
	log.Printf("Password: %s\n", superAdminPassword)
	log.Printf("Role: %s\n", superAdmin.Role)
	log.Printf("Admin Level: %s\n", superAdmin.AdminLevel)
	log.Println("================================================")
	log.Println("⚠️  IMPORTANT: Please change the password after first login!")
	log.Println("Seed completed successfully")
}
