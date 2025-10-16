package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Get environment variables
	host := os.Getenv("BLUEPRINT_DB_HOST")
	port := os.Getenv("BLUEPRINT_DB_PORT")
	username := os.Getenv("BLUEPRINT_DB_USERNAME")
	password := os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")
	databaseName := os.Getenv("BLUEPRINT_DB_DATABASE")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "27017"
	}
	if databaseName == "" {
		databaseName = "doctors_workspace"
	}

	// Build connection URI
	var uri string
	if username != "" && password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", host, port)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database(databaseName)
	usersCollection := db.Collection("users")
	institutionsCollection := db.Collection("institutions")

	fmt.Println("🔄 Starting institution migration...")
	fmt.Println()

	// Step 1: Load all institutions
	cursor, err := institutionsCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to load institutions: %v", err)
	}

	var institutions []models.Institution
	if err = cursor.All(ctx, &institutions); err != nil {
		log.Fatalf("Failed to decode institutions: %v", err)
	}
	cursor.Close(ctx)

	fmt.Printf("📋 Loaded %d institutions\n\n", len(institutions))

	// Create a map for quick lookup
	institutionMap := make(map[string]primitive.ObjectID)
	for _, inst := range institutions {
		// Map both full name and short name
		institutionMap[strings.ToLower(inst.Name)] = inst.ID
		if inst.ShortName != "" {
			institutionMap[strings.ToLower(inst.ShortName)] = inst.ID
		}
		// Also map by domain-like keys
		switch inst.ShortName {
		case "UCT":
			institutionMap["uct"] = inst.ID
		case "Wits":
			institutionMap["wits"] = inst.ID
		case "SU":
			institutionMap["sun"] = inst.ID
		case "UP":
			institutionMap["up"] = inst.ID
		case "UKZN":
			institutionMap["ukzn"] = inst.ID
		case "UFS":
			institutionMap["ufs"] = inst.ID
		case "NHLS":
			institutionMap["nhls"] = inst.ID
		case "BloodSA":
			institutionMap["bloodsa.org.za"] = inst.ID
			institutionMap["blood sa"] = inst.ID
		}
	}

	// Add institution mappings for common variations
	additionalMappings := map[string]string{
		"stats.co.za":                 "Statistics South Africa",
		"management.co.za":            "Medical Management Group",
		"epidemiology.co.za":          "Epidemiology Research Unit",
		"publichealth.co.za":          "Public Health Institute",
		"healthdata.co.za":            "Health Data Systems",
		"data.co.za":                  "Data Solutions Ltd",
		"biomed.co.za":                "Biomedical Research Centre",
		"clinical.co.za":              "Clinical Research Institute",
		"admin.co.za":                 "Health Administration Services",
		"research.co.za":              "South African Medical Research Council",
		"health.gov.za":               "National Department of Health",
		"mandela":                     "Nelson Mandela University",
		"medunsa":                     "Sefako Makgatho Health Sciences University",
	}

	// Default institution for unmapped ones (BloodSA)
	var defaultInstitutionID primitive.ObjectID
	for _, inst := range institutions {
		if inst.ShortName == "BloodSA" {
			defaultInstitutionID = inst.ID
			break
		}
	}

	// Step 2: Load all users with old institution format
	type OldUser struct {
		ID                  primitive.ObjectID `bson:"_id"`
		Profile             bson.M             `bson:"profile"`
	}

	userCursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	var oldUsers []OldUser
	if err = userCursor.All(ctx, &oldUsers); err != nil {
		log.Fatalf("Failed to decode users: %v", err)
	}
	userCursor.Close(ctx)

	fmt.Printf("👥 Found %d users to migrate\n\n", len(oldUsers))

	// Step 3: Migrate each user
	migratedCount := 0
	skippedCount := 0

	for _, user := range oldUsers {
		// Check if already migrated (has institution_id)
		if institutionID, ok := user.Profile["institution_id"]; ok && institutionID != nil {
			skippedCount++
			continue
		}

		// Get old institution string
		oldInstitution, ok := user.Profile["institution"].(string)
		if !ok || oldInstitution == "" {
			// No institution set, use default (BloodSA)
			if defaultInstitutionID.IsZero() {
				fmt.Printf("⚠️  User %s has no institution and no default available\n", user.ID.Hex())
				skippedCount++
				continue
			}
			
			// Update user with default institution
			_, err := usersCollection.UpdateOne(
				ctx,
				bson.M{"_id": user.ID},
				bson.M{
					"$set": bson.M{
						"profile.institution_id": defaultInstitutionID,
					},
					"$unset": bson.M{
						"profile.institution": "",
						"profile.location":    "",
					},
				},
			)
			if err != nil {
				fmt.Printf("❌ Failed to migrate user %s: %v\n", user.ID.Hex(), err)
				continue
			}
			migratedCount++
			fmt.Printf("✅ Migrated user %s to default institution (BloodSA)\n", user.ID.Hex())
			continue
		}

		// Try to find matching institution
		var matchedInstitutionID primitive.ObjectID
		found := false

		// Check direct match
		if instID, ok := institutionMap[strings.ToLower(oldInstitution)]; ok {
			matchedInstitutionID = instID
			found = true
		} else {
			// Try to find by partial match
			for key, mappedName := range additionalMappings {
				if strings.Contains(strings.ToLower(oldInstitution), key) {
					if instID, ok := institutionMap[strings.ToLower(mappedName)]; ok {
						matchedInstitutionID = instID
						found = true
						break
					}
				}
			}
		}

		if !found {
			// Use default institution for unknown institutions
			if !defaultInstitutionID.IsZero() {
				matchedInstitutionID = defaultInstitutionID
				fmt.Printf("⚠️  Unknown institution '%s' for user %s, using default\n", oldInstitution, user.ID.Hex())
			} else {
				fmt.Printf("❌ Could not map institution '%s' for user %s\n", oldInstitution, user.ID.Hex())
				skippedCount++
				continue
			}
		}

		// Update user with institution ID
		_, err := usersCollection.UpdateOne(
			ctx,
			bson.M{"_id": user.ID},
			bson.M{
				"$set": bson.M{
					"profile.institution_id": matchedInstitutionID,
				},
				"$unset": bson.M{
					"profile.institution": "",
					"profile.location":    "",
				},
			},
		)
		if err != nil {
			fmt.Printf("❌ Failed to migrate user %s: %v\n", user.ID.Hex(), err)
			continue
		}

		migratedCount++
	}

	fmt.Println()
	fmt.Printf("✅ Migration complete!\n")
	fmt.Printf("   - Migrated: %d users\n", migratedCount)
	fmt.Printf("   - Already migrated: %d users\n", skippedCount)
	fmt.Printf("   - Total: %d users\n", len(oldUsers))
	fmt.Println()
	fmt.Println("🎉 Users now reference institutions by ID!")
}

