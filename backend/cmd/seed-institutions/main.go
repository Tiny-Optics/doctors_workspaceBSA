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
	institutionsCollection := db.Collection("institutions")

	fmt.Println("Seeding institutions...")

	now := time.Now()
	institutions := []models.Institution{
		// Major Universities
		{
			Name:       "University of Cape Town",
			ShortName:  "UCT",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Western Cape",
			City:       "Cape Town",
			Address:    "Private Bag X3, Rondebosch",
			PostalCode: "7701",
			Phone:      "+27 21 650 9111",
			Email:      "info@uct.ac.za",
			Website:    "https://www.uct.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "University of the Witwatersrand",
			ShortName:  "Wits",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Address:    "1 Jan Smuts Avenue, Braamfontein",
			PostalCode: "2000",
			Phone:      "+27 11 717 1000",
			Email:      "info@wits.ac.za",
			Website:    "https://www.wits.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Stellenbosch University",
			ShortName:  "SU",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Western Cape",
			City:       "Stellenbosch",
			Address:    "Private Bag X1, Matieland",
			PostalCode: "7602",
			Phone:      "+27 21 808 9111",
			Email:      "info@sun.ac.za",
			Website:    "https://www.sun.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "University of Pretoria",
			ShortName:  "UP",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Pretoria",
			Address:    "Lynnwood Road, Hatfield",
			PostalCode: "0002",
			Phone:      "+27 12 420 3111",
			Email:      "info@up.ac.za",
			Website:    "https://www.up.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "University of KwaZulu-Natal",
			ShortName:  "UKZN",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "KwaZulu-Natal",
			City:       "Durban",
			Address:    "King George V Avenue",
			PostalCode: "4041",
			Phone:      "+27 31 260 1111",
			Email:      "info@ukzn.ac.za",
			Website:    "https://www.ukzn.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "University of the Free State",
			ShortName:  "UFS",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Free State",
			City:       "Bloemfontein",
			Address:    "205 Nelson Mandela Drive, Park West",
			PostalCode: "9301",
			Phone:      "+27 51 401 9111",
			Email:      "info@ufs.ac.za",
			Website:    "https://www.ufs.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Rhodes University",
			ShortName:  "RU",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Eastern Cape",
			City:       "Grahamstown",
			Phone:      "+27 46 603 8111",
			Email:      "info@ru.ac.za",
			Website:    "https://www.ru.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "North-West University",
			ShortName:  "NWU",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "North West",
			City:       "Potchefstroom",
			Phone:      "+27 18 299 1111",
			Email:      "info@nwu.ac.za",
			Website:    "https://www.nwu.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "University of South Africa",
			ShortName:  "UNISA",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Pretoria",
			Phone:      "+27 12 429 3111",
			Email:      "info@unisa.ac.za",
			Website:    "https://www.unisa.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Nelson Mandela University",
			ShortName:  "Mandela",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Eastern Cape",
			City:       "Port Elizabeth",
			Phone:      "+27 41 504 1111",
			Email:      "info@mandela.ac.za",
			Website:    "https://www.mandela.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "University of Limpopo",
			ShortName:  "UL",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Limpopo",
			City:       "Polokwane",
			Phone:      "+27 15 268 2111",
			Email:      "info@ul.ac.za",
			Website:    "https://www.ul.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Walter Sisulu University",
			ShortName:  "WSU",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Eastern Cape",
			City:       "Mthatha",
			Phone:      "+27 47 502 2111",
			Email:      "info@wsu.ac.za",
			Website:    "https://www.wsu.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Cape Peninsula University of Technology",
			ShortName:  "CPUT",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Western Cape",
			City:       "Cape Town",
			Phone:      "+27 21 460 3911",
			Email:      "info@cput.ac.za",
			Website:    "https://www.cput.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Durban University of Technology",
			ShortName:  "DUT",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "KwaZulu-Natal",
			City:       "Durban",
			Phone:      "+27 31 373 2000",
			Email:      "info@dut.ac.za",
			Website:    "https://www.dut.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Tshwane University of Technology",
			ShortName:  "TUT",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Pretoria",
			Phone:      "+27 12 382 5911",
			Email:      "info@tut.ac.za",
			Website:    "https://www.tut.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Vaal University of Technology",
			ShortName:  "VUT",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Vanderbijlpark",
			Phone:      "+27 16 950 9000",
			Email:      "info@vut.ac.za",
			Website:    "https://www.vut.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Central University of Technology",
			ShortName:  "CUT",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Free State",
			City:       "Bloemfontein",
			Phone:      "+27 51 507 3911",
			Email:      "info@cut.ac.za",
			Website:    "https://www.cut.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Sefako Makgatho Health Sciences University",
			ShortName:  "SMU",
			Type:       models.InstitutionTypeUniversity,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Pretoria",
			Phone:      "+27 12 521 4111",
			Email:      "info@smu.ac.za",
			Website:    "https://www.smu.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},

		// Major Hospitals
		{
			Name:       "Groote Schuur Hospital",
			ShortName:  "GSH",
			Type:       models.InstitutionTypeHospital,
			Country:    "South Africa",
			Province:   "Western Cape",
			City:       "Cape Town",
			Address:    "Main Road, Observatory",
			PostalCode: "7925",
			Phone:      "+27 21 404 9111",
			Email:      "info@gsh.gov.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Chris Hani Baragwanath Academic Hospital",
			ShortName:  "Bara",
			Type:       models.InstitutionTypeHospital,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Address:    "26 Chris Hani Road, Diepkloof",
			PostalCode: "1864",
			Phone:      "+27 11 933 0111",
			Email:      "info@bara.gov.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Charlotte Maxeke Johannesburg Academic Hospital",
			ShortName:  "CMJAH",
			Type:       models.InstitutionTypeHospital,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Address:    "17 Jubilee Road, Parktown",
			PostalCode: "2193",
			Phone:      "+27 11 488 4911",
			Email:      "info@cmjah.gov.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Inkosi Albert Luthuli Central Hospital",
			ShortName:  "IALCH",
			Type:       models.InstitutionTypeHospital,
			Country:    "South Africa",
			Province:   "KwaZulu-Natal",
			City:       "Durban",
			Address:    "800 Vusi Mzimela Road, Cato Manor",
			PostalCode: "4091",
			Phone:      "+27 31 240 1111",
			Email:      "info@ialch.gov.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Steve Biko Academic Hospital",
			ShortName:  "SBAH",
			Type:       models.InstitutionTypeHospital,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Pretoria",
			Address:    "Corner Malherbe and Steve Biko Road, Pretoria Central",
			PostalCode: "0001",
			Phone:      "+27 12 354 1000",
			Email:      "info@sbah.gov.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},

		// Laboratories & Research
		{
			Name:       "National Health Laboratory Service",
			ShortName:  "NHLS",
			Type:       models.InstitutionTypeLaboratory,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Address:    "1 Modderfontein Road, Sandringham",
			PostalCode: "2192",
			Phone:      "+27 11 386 6000",
			Email:      "info@nhls.ac.za",
			Website:    "https://www.nhls.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "South African Medical Research Council",
			ShortName:  "SAMRC",
			Type:       models.InstitutionTypeResearchCenter,
			Country:    "South Africa",
			Province:   "Western Cape",
			City:       "Cape Town",
			Address:    "Francie van Zijl Drive, Parow Valley",
			PostalCode: "7501",
			Phone:      "+27 21 938 0911",
			Email:      "info@mrc.ac.za",
			Website:    "https://www.samrc.ac.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Biomedical Research Centre",
			ShortName:  "BRC",
			Type:       models.InstitutionTypeResearchCenter,
			Country:    "South Africa",
			Province:   "Western Cape",
			City:       "Cape Town",
			Phone:      "+27 21 456 7890",
			Email:      "info@biomed.co.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Epidemiology Research Unit",
			ShortName:  "ERU",
			Type:       models.InstitutionTypeResearchCenter,
			Country:    "South Africa",
			Province:   "KwaZulu-Natal",
			City:       "Durban",
			Phone:      "+27 31 567 8901",
			Email:      "info@epidemiology.co.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Clinical Research Institute",
			ShortName:  "CRI",
			Type:       models.InstitutionTypeResearchCenter,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Phone:      "+27 11 678 9012",
			Email:      "info@clinical.co.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Public Health Institute",
			ShortName:  "PHI",
			Type:       models.InstitutionTypeResearchCenter,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Phone:      "+27 11 123 4567",
			Email:      "info@publichealth.co.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},

		// Government
		{
			Name:       "National Department of Health",
			ShortName:  "NDoH",
			Type:       models.InstitutionTypeGovernment,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Pretoria",
			Address:    "Civitas Building, corner of Thabo Sehume and Struben Streets",
			PostalCode: "0001",
			Phone:      "+27 12 395 8000",
			Email:      "info@health.gov.za",
			Website:    "https://www.health.gov.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Statistics South Africa",
			ShortName:  "StatsSA",
			Type:       models.InstitutionTypeGovernment,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Pretoria",
			Phone:      "+27 12 310 8911",
			Email:      "info@statssa.gov.za",
			Website:    "https://www.statssa.gov.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},

		// BloodSA Organization
		{
			Name:       "South African National Blood Service",
			ShortName:  "BloodSA",
			Type:       models.InstitutionTypeNGO,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Address:    "1 Constantia Boulevard, Constantia Kloof",
			PostalCode: "1709",
			Phone:      "+27 11 761 9000",
			Email:      "info@bloodsa.org.za",
			Website:    "https://www.bloodsa.org.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},

		// Other Organizations
		{
			Name:       "Health Data Systems",
			ShortName:  "HDS",
			Type:       models.InstitutionTypeOther,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Phone:      "+27 11 234 5678",
			Email:      "info@healthdata.co.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Data Solutions Ltd",
			ShortName:  "DSL",
			Type:       models.InstitutionTypeOther,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Phone:      "+27 11 345 6789",
			Email:      "info@data.co.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Health Administration Services",
			ShortName:  "HAS",
			Type:       models.InstitutionTypeOther,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Pretoria",
			Phone:      "+27 12 789 0123",
			Email:      "info@admin.co.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			Name:       "Medical Management Group",
			ShortName:  "MMG",
			Type:       models.InstitutionTypeOther,
			Country:    "South Africa",
			Province:   "Gauteng",
			City:       "Johannesburg",
			Phone:      "+27 11 890 1234",
			Email:      "info@management.co.za",
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	// Insert institutions
	var interfaceInstitutions []interface{}
	for i := range institutions {
		interfaceInstitutions = append(interfaceInstitutions, institutions[i])
	}

	result, err := institutionsCollection.InsertMany(ctx, interfaceInstitutions)
	if err != nil {
		log.Fatalf("Failed to seed institutions: %v", err)
	}

	fmt.Printf("✅ Successfully seeded %d institutions!\n\n", len(result.InsertedIDs))

	fmt.Println("📋 Seeded Institutions:")
	for _, institution := range institutions {
		fmt.Printf("  - %s (%s) - %s, %s\n", institution.Name, institution.ShortName, institution.City, institution.Province)
	}

	fmt.Println("\n📊 Institution Types:")
	typeCount := make(map[models.InstitutionType]int)
	for _, inst := range institutions {
		typeCount[inst.Type]++
	}
	for instType, count := range typeCount {
		fmt.Printf("  - %s: %d\n", instType, count)
	}

	fmt.Println("\nInstitutions can now be used when creating users!")
}
