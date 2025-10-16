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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	fmt.Println("Adding 30 new users to the database...")

	// 30 diverse users with realistic South African data
	sampleUsers := []struct {
		FirstName         string
		LastName          string
		Email             string
		Role              models.UserRole
		AdminLevel        models.AdminLevel
		Institution       string
		Specialty         string
		Location          string
		RegistrationNumber string
		PhoneNumber       string
		IsActive          bool
	}{
		// Haematologists (10)
		{
			FirstName:         "Dr. Thabo",
			LastName:          "Mthembu",
			Email:             "thabo.mthembu@nhls.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "National Health Laboratory Service",
			Specialty:         "Haematology",
			Location:          "Johannesburg",
			RegistrationNumber: "HPCSA234567",
			PhoneNumber:       "+27 11 242 8000",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Nomsa",
			LastName:          "Dlamini",
			Email:             "nomsa.dlamini@wits.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "University of the Witwatersrand",
			Specialty:         "Haematology",
			Location:          "Johannesburg",
			RegistrationNumber: "HPCSA234568",
			PhoneNumber:       "+27 11 717 1000",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Sipho",
			LastName:          "Nkosi",
			Email:             "sipho.nkosi@uct.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "University of Cape Town",
			Specialty:         "Haematology",
			Location:          "Cape Town",
			RegistrationNumber: "HPCSA234569",
			PhoneNumber:       "+27 21 650 9111",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Lerato",
			LastName:          "Molefe",
			Email:             "lerato.molefe@sun.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "Stellenbosch University",
			Specialty:         "Haematology",
			Location:          "Stellenbosch",
			RegistrationNumber: "HPCSA234570",
			PhoneNumber:       "+27 21 808 9111",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Mandla",
			LastName:          "Zulu",
			Email:             "mandla.zulu@up.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "University of Pretoria",
			Specialty:         "Haematology",
			Location:          "Pretoria",
			RegistrationNumber: "HPCSA234571",
			PhoneNumber:       "+27 12 420 3111",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Busisiwe",
			LastName:          "Mkhize",
			Email:             "busisiwe.mkhize@ukzn.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "University of KwaZulu-Natal",
			Specialty:         "Haematology",
			Location:          "Durban",
			RegistrationNumber: "HPCSA234572",
			PhoneNumber:       "+27 31 260 1111",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Themba",
			LastName:          "Mabena",
			Email:             "themba.mabena@ru.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "Rhodes University",
			Specialty:         "Haematology",
			Location:          "Grahamstown",
			RegistrationNumber: "HPCSA234573",
			PhoneNumber:       "+27 46 603 8111",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Ntombi",
			LastName:          "Sithole",
			Email:             "ntombi.sithole@ufs.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "University of the Free State",
			Specialty:         "Haematology",
			Location:          "Bloemfontein",
			RegistrationNumber: "HPCSA234574",
			PhoneNumber:       "+27 51 401 9111",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Khaya",
			LastName:          "Mthembu",
			Email:             "khaya.mthembu@nwu.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "North-West University",
			Specialty:         "Haematology",
			Location:          "Potchefstroom",
			RegistrationNumber: "HPCSA234575",
			PhoneNumber:       "+27 18 299 1111",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Zanele",
			LastName:          "Mthembu",
			Email:             "zanele.mthembu@unisa.ac.za",
			Role:              models.RoleHaematologist,
			Institution:       "University of South Africa",
			Specialty:         "Haematology",
			Location:          "Pretoria",
			RegistrationNumber: "HPCSA234576",
			PhoneNumber:       "+27 11 670 9000",
			IsActive:          true,
		},

		// Physicians (10)
		{
			FirstName:         "Dr. Andile",
			LastName:          "Ngcobo",
			Email:             "andile.ngcobo@health.gov.za",
			Role:              models.RolePhysician,
			Institution:       "Department of Health",
			Specialty:         "Internal Medicine",
			Location:          "Cape Town",
			RegistrationNumber: "HPCSA234577",
			PhoneNumber:       "+27 21 483 5000",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Nomfundo",
			LastName:          "Mthembu",
			Email:             "nomfundo.mthembu@medunsa.ac.za",
			Role:              models.RolePhysician,
			Institution:       "Sefako Makgatho Health Sciences University",
			Specialty:         "Internal Medicine",
			Location:          "Pretoria",
			RegistrationNumber: "HPCSA234578",
			PhoneNumber:       "+27 12 521 4000",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Sibusiso",
			LastName:          "Mthembu",
			Email:             "sibusiso.mthembu@wsu.ac.za",
			Role:              models.RolePhysician,
			Institution:       "Walter Sisulu University",
			Specialty:         "Internal Medicine",
			Location:          "Mthatha",
			RegistrationNumber: "HPCSA234579",
			PhoneNumber:       "+27 47 502 2200",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Thandeka",
			LastName:          "Mthembu",
			Email:             "thandeka.mthembu@ul.ac.za",
			Role:              models.RolePhysician,
			Institution:       "University of Limpopo",
			Specialty:         "Internal Medicine",
			Location:          "Polokwane",
			RegistrationNumber: "HPCSA234580",
			PhoneNumber:       "+27 15 268 3000",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Mthunzi",
			LastName:          "Mthembu",
			Email:             "mthunzi.mthembu@cut.ac.za",
			Role:              models.RolePhysician,
			Institution:       "Central University of Technology",
			Specialty:         "Internal Medicine",
			Location:          "Bloemfontein",
			RegistrationNumber: "HPCSA234581",
			PhoneNumber:       "+27 51 507 3911",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Nolwazi",
			LastName:          "Mthembu",
			Email:             "nolwazi.mthembu@dut.ac.za",
			Role:              models.RolePhysician,
			Institution:       "Durban University of Technology",
			Specialty:         "Internal Medicine",
			Location:          "Durban",
			RegistrationNumber: "HPCSA234582",
			PhoneNumber:       "+27 31 373 2000",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Bongani",
			LastName:          "Mthembu",
			Email:             "bongani.mthembu@cput.ac.za",
			Role:              models.RolePhysician,
			Institution:       "Cape Peninsula University of Technology",
			Specialty:         "Internal Medicine",
			Location:          "Cape Town",
			RegistrationNumber: "HPCSA234583",
			PhoneNumber:       "+27 21 460 3911",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Nomsa",
			LastName:          "Mthembu",
			Email:             "nomsa.mthembu@tut.ac.za",
			Role:              models.RolePhysician,
			Institution:       "Tshwane University of Technology",
			Specialty:         "Internal Medicine",
			Location:          "Pretoria",
			RegistrationNumber: "HPCSA234584",
			PhoneNumber:       "+27 12 382 5911",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Sipho",
			LastName:          "Mthembu",
			Email:             "sipho.mthembu@vut.ac.za",
			Role:              models.RolePhysician,
			Institution:       "Vaal University of Technology",
			Specialty:         "Internal Medicine",
			Location:          "Vanderbijlpark",
			RegistrationNumber: "HPCSA234585",
			PhoneNumber:       "+27 16 950 9000",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Thabo",
			LastName:          "Mthembu",
			Email:             "thabo.mthembu@mandela.ac.za",
			Role:              models.RolePhysician,
			Institution:       "Nelson Mandela University",
			Specialty:         "Internal Medicine",
			Location:          "Port Elizabeth",
			RegistrationNumber: "HPCSA234586",
			PhoneNumber:       "+27 41 504 1111",
			IsActive:          true,
		},

		// Data Capturers (8)
		{
			FirstName:   "Nomsa",
			LastName:    "Mthembu",
			Email:       "nomsa.mthembu@data.co.za",
			Role:        models.RoleDataCapturer,
			Institution: "Data Solutions Ltd",
			Specialty:   "Data Management",
			Location:    "Johannesburg",
			PhoneNumber: "+27 11 123 4567",
			IsActive:    true,
		},
		{
			FirstName:   "Sipho",
			LastName:    "Mthembu",
			Email:       "sipho.mthembu@research.co.za",
			Role:        models.RoleDataCapturer,
			Institution: "Medical Research Council",
			Specialty:   "Research Data",
			Location:    "Cape Town",
			PhoneNumber: "+27 21 938 0911",
			IsActive:    true,
		},
		{
			FirstName:   "Thabo",
			LastName:    "Mthembu",
			Email:       "thabo.mthembu@stats.co.za",
			Role:        models.RoleDataCapturer,
			Institution: "Statistics South Africa",
			Specialty:   "Health Statistics",
			Location:    "Pretoria",
			PhoneNumber: "+27 12 310 8911",
			IsActive:    true,
		},
		{
			FirstName:   "Lerato",
			LastName:    "Mthembu",
			Email:       "lerato.mthembu@healthdata.co.za",
			Role:        models.RoleDataCapturer,
			Institution: "Health Data Systems",
			Specialty:   "Health Informatics",
			Location:    "Durban",
			PhoneNumber: "+27 31 123 4567",
			IsActive:    true,
		},
		{
			FirstName:   "Mandla",
			LastName:    "Mthembu",
			Email:       "mandla.mthembu@clinical.co.za",
			Role:        models.RoleDataCapturer,
			Institution: "Clinical Research Institute",
			Specialty:   "Clinical Data",
			Location:    "Stellenbosch",
			PhoneNumber: "+27 21 808 9111",
			IsActive:    true,
		},
		{
			FirstName:   "Busisiwe",
			LastName:    "Mthembu",
			Email:       "busisiwe.mthembu@biomed.co.za",
			Role:        models.RoleDataCapturer,
			Institution: "Biomedical Research Centre",
			Specialty:   "Biomedical Data",
			Location:    "Bloemfontein",
			PhoneNumber: "+27 51 401 9111",
			IsActive:    true,
		},
		{
			FirstName:   "Themba",
			LastName:    "Mthembu",
			Email:       "themba.mthembu@epidemiology.co.za",
			Role:        models.RoleDataCapturer,
			Institution: "Epidemiology Research Unit",
			Specialty:   "Epidemiological Data",
			Location:    "Grahamstown",
			PhoneNumber: "+27 46 603 8111",
			IsActive:    true,
		},
		{
			FirstName:   "Ntombi",
			LastName:    "Mthembu",
			Email:       "ntombi.mthembu@publichealth.co.za",
			Role:        models.RoleDataCapturer,
			Institution: "Public Health Institute",
			Specialty:   "Public Health Data",
			Location:    "Port Elizabeth",
			PhoneNumber: "+27 41 504 1111",
			IsActive:    true,
		},

		// User Managers (2)
		{
			FirstName:         "Dr. Khaya",
			LastName:          "Mthembu",
			Email:             "khaya.mthembu@admin.co.za",
			Role:              models.RoleHaematologist,
			AdminLevel:        models.AdminLevelUserManager,
			Institution:       "Health Administration Services",
			Specialty:         "Haematology",
			Location:          "Johannesburg",
			RegistrationNumber: "HPCSA234587",
			PhoneNumber:       "+27 11 234 5678",
			IsActive:          true,
		},
		{
			FirstName:         "Dr. Zanele",
			LastName:          "Mthembu",
			Email:             "zanele.mthembu@management.co.za",
			Role:              models.RolePhysician,
			AdminLevel:        models.AdminLevelUserManager,
			Institution:       "Medical Management Group",
			Specialty:         "Internal Medicine",
			Location:          "Cape Town",
			RegistrationNumber: "HPCSA234588",
			PhoneNumber:       "+27 21 345 6789",
			IsActive:          true,
		},
	}

	// Hash password for all users
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Insert users
	successCount := 0
	for i, su := range sampleUsers {
		now := time.Now()
		user := models.User{
			Username:     su.Email,
			Email:        su.Email,
			PasswordHash: string(passwordHash),
			Role:         su.Role,
			AdminLevel:   su.AdminLevel,
			IsActive:     su.IsActive,
			Profile: models.UserProfile{
				FirstName:         su.FirstName,
				LastName:          su.LastName,
				Institution:       su.Institution,
				Specialty:         su.Specialty,
				Location:          su.Location,
				RegistrationNumber: su.RegistrationNumber,
				PhoneNumber:       su.PhoneNumber,
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		_, err = usersCollection.InsertOne(ctx, user)
		if err != nil {
			log.Printf("Failed to insert user %d (%s): %v", i+1, su.Email, err)
		} else {
			successCount++
			fmt.Printf("User %d: %s (%s) added successfully\n", i+1, su.Email, su.Role)
		}
	}

	fmt.Printf("\n✅ Successfully added %d out of %d users!\n", successCount, len(sampleUsers))
	fmt.Println("\n📊 User Distribution:")
	fmt.Printf("- Haematologists: 10\n")
	fmt.Printf("- Physicians: 10\n")
	fmt.Printf("- Data Capturers: 8\n")
	fmt.Printf("- User Managers: 2\n")
	fmt.Printf("- Total: 30 users\n")

	fmt.Println("\n🔐 Login Credentials for all users:")
	fmt.Println("Email: [user-specific email]")
	fmt.Println("Password: Password123!")

	fmt.Println("\n🏥 Institutions Represented:")
	institutions := make(map[string]bool)
	for _, user := range sampleUsers {
		institutions[user.Institution] = true
	}
	for institution := range institutions {
		fmt.Printf("- %s\n", institution)
	}

	fmt.Println("\n📍 Locations Covered:")
	locations := make(map[string]bool)
	for _, user := range sampleUsers {
		locations[user.Location] = true
	}
	for location := range locations {
		fmt.Printf("- %s\n", location)
	}
}
