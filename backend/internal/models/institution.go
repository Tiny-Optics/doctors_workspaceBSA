package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Institution represents a medical institution or organization
type Institution struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name       string              `bson:"name" json:"name"`
	ShortName  string              `bson:"short_name,omitempty" json:"shortName,omitempty"`
	Type       InstitutionType     `bson:"type" json:"type"`
	Country    string              `bson:"country" json:"country"`
	Province   string              `bson:"province,omitempty" json:"province,omitempty"`
	City       string              `bson:"city" json:"city"`
	Address    string              `bson:"address,omitempty" json:"address,omitempty"`
	PostalCode string              `bson:"postal_code,omitempty" json:"postalCode,omitempty"`
	Phone      string              `bson:"phone,omitempty" json:"phone,omitempty"`
	Email      string              `bson:"email,omitempty" json:"email,omitempty"`
	Website    string              `bson:"website,omitempty" json:"website,omitempty"`
	ImagePath  string              `bson:"image_path,omitempty" json:"imagePath,omitempty"`
	IsActive   bool                `bson:"is_active" json:"isActive"`
	CreatedAt  time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time           `bson:"updated_at" json:"updatedAt"`
	CreatedBy  *primitive.ObjectID `bson:"created_by,omitempty" json:"createdBy,omitempty"`
}

// InstitutionType represents the type of institution
type InstitutionType string

const (
	InstitutionTypeUniversity      InstitutionType = "university"
	InstitutionTypeHospital        InstitutionType = "hospital"
	InstitutionTypeLaboratory      InstitutionType = "laboratory"
	InstitutionTypeResearchCenter  InstitutionType = "research_center"
	InstitutionTypeGovernment      InstitutionType = "government"
	InstitutionTypePrivatePractice InstitutionType = "private_practice"
	InstitutionTypeNGO             InstitutionType = "ngo"
	InstitutionTypeOther           InstitutionType = "other"
)

// IsValid checks if the institution type is valid
func (t InstitutionType) IsValid() bool {
	switch t {
	case InstitutionTypeUniversity,
		InstitutionTypeHospital,
		InstitutionTypeLaboratory,
		InstitutionTypeResearchCenter,
		InstitutionTypeGovernment,
		InstitutionTypePrivatePractice,
		InstitutionTypeNGO,
		InstitutionTypeOther:
		return true
	}
	return false
}

// CreateInstitutionRequest represents the request to create a new institution
type CreateInstitutionRequest struct {
	Name       string          `json:"name" binding:"required"`
	ShortName  string          `json:"shortName,omitempty"`
	Type       InstitutionType `json:"type" binding:"required"`
	Country    string          `json:"country" binding:"required"`
	Province   string          `json:"province,omitempty"`
	City       string          `json:"city" binding:"required"`
	Address    string          `json:"address,omitempty"`
	PostalCode string          `json:"postalCode,omitempty"`
	Phone      string          `json:"phone,omitempty"`
	Email      string          `json:"email,omitempty"`
	Website    string          `json:"website,omitempty"`
	ImagePath  string          `json:"imagePath,omitempty"`
}

// UpdateInstitutionRequest represents the request to update an institution
type UpdateInstitutionRequest struct {
	Name       *string          `json:"name,omitempty"`
	ShortName  *string          `json:"shortName,omitempty"`
	Type       *InstitutionType `json:"type,omitempty"`
	Country    *string          `json:"country,omitempty"`
	Province   *string          `json:"province,omitempty"`
	City       *string          `json:"city,omitempty"`
	Address    *string          `json:"address,omitempty"`
	PostalCode *string          `json:"postalCode,omitempty"`
	Phone      *string          `json:"phone,omitempty"`
	Email      *string          `json:"email,omitempty"`
	Website    *string          `json:"website,omitempty"`
	ImagePath  *string          `json:"imagePath,omitempty"`
	IsActive   *bool            `json:"isActive,omitempty"`
}

// Validation errors
var (
	ErrInstitutionNameRequired    = errors.New("institution name is required")
	ErrInstitutionNameTooShort    = errors.New("institution name must be at least 2 characters")
	ErrInstitutionNameTooLong     = errors.New("institution name must be at most 200 characters")
	ErrInvalidInstitutionType     = errors.New("invalid institution type")
	ErrInstitutionCountryRequired = errors.New("country is required")
	ErrInstitutionCityRequired    = errors.New("city is required")
)

// Validate validates the CreateInstitutionRequest
func (req *CreateInstitutionRequest) Validate() error {
	// Validate name
	if req.Name == "" {
		return ErrInstitutionNameRequired
	}
	if len(req.Name) < 2 {
		return ErrInstitutionNameTooShort
	}
	if len(req.Name) > 200 {
		return ErrInstitutionNameTooLong
	}

	// Validate type
	if !req.Type.IsValid() {
		return ErrInvalidInstitutionType
	}

	// Validate country
	if req.Country == "" {
		return ErrInstitutionCountryRequired
	}

	// Validate city
	if req.City == "" {
		return ErrInstitutionCityRequired
	}

	return nil
}

// GetFullLocation returns the full location string
func (i *Institution) GetFullLocation() string {
	location := i.City
	if i.Province != "" {
		location += ", " + i.Province
	}
	if i.Country != "" {
		location += ", " + i.Country
	}
	return location
}

// GetDisplayName returns the display name (short name if available, otherwise full name)
func (i *Institution) GetDisplayName() string {
	if i.ShortName != "" {
		return i.ShortName
	}
	return i.Name
}
