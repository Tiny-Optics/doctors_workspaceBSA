package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidCategoryName = errors.New("category name is required and must be between 1 and 100 characters")
	ErrInvalidSlug         = errors.New("invalid slug format")
	ErrInvalidImageFormat  = errors.New("image must be jpg, jpeg, png, or webp format")
	ErrInvalidDescription  = errors.New("description must be less than 1000 characters")
	ErrInvalidDisplayOrder = errors.New("display order must be a positive number")
)

// SOPCategory represents a category for Standard Operating Procedures
type SOPCategory struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name         string              `bson:"name" json:"name"`
	Slug         string              `bson:"slug" json:"slug"`
	Description  string              `bson:"description,omitempty" json:"description,omitempty"`
	ImagePath    string              `bson:"image_path,omitempty" json:"imagePath,omitempty"`
	DropboxPath  string              `bson:"dropbox_path" json:"dropboxPath"`
	DisplayOrder int                 `bson:"display_order" json:"displayOrder"`
	IsActive     bool                `bson:"is_active" json:"isActive"`
	CreatedAt    time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updated_at" json:"updatedAt"`
	CreatedBy    *primitive.ObjectID `bson:"created_by,omitempty" json:"createdBy,omitempty"`
}

// CreateSOPCategoryRequest represents the request to create a new SOP category
type CreateSOPCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ImagePath    string `json:"imagePath"`
	DisplayOrder int    `json:"displayOrder"`
}

// UpdateSOPCategoryRequest represents the request to update an SOP category
type UpdateSOPCategoryRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	ImagePath    *string `json:"imagePath"`
	DisplayOrder *int    `json:"displayOrder"`
	IsActive     *bool   `json:"isActive"`
}

// Validate validates the SOPCategory fields
func (c *SOPCategory) Validate() error {
	// Validate name
	if err := ValidateCategoryName(c.Name); err != nil {
		return err
	}

	// Validate slug
	if c.Slug == "" {
		return ErrInvalidSlug
	}

	// Validate description
	if len(c.Description) > 1000 {
		return ErrInvalidDescription
	}

	// Validate image path format if provided
	if c.ImagePath != "" {
		if err := ValidateImagePath(c.ImagePath); err != nil {
			return err
		}
	}

	// Validate display order
	if c.DisplayOrder < 0 {
		return ErrInvalidDisplayOrder
	}

	return nil
}

// ValidateCategoryName validates the category name
func ValidateCategoryName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 100 {
		return ErrInvalidCategoryName
	}
	return nil
}

// ValidateImagePath validates the image file path
func ValidateImagePath(imagePath string) error {
	if imagePath == "" {
		return nil
	}

	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}
	imagePath = strings.ToLower(imagePath)

	isValid := false
	for _, ext := range validExtensions {
		if strings.HasSuffix(imagePath, ext) {
			isValid = true
			break
		}
	}

	if !isValid {
		return ErrInvalidImageFormat
	}

	return nil
}

// GenerateSlug generates a URL-friendly slug from a name
func GenerateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove all non-alphanumeric characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// GetDropboxPath generates the Dropbox path for the category
func (c *SOPCategory) GetDropboxPath() string {
	if c.DropboxPath != "" {
		return c.DropboxPath
	}
	return "SOPS/" + c.Name
}

// Validate validates the CreateSOPCategoryRequest
func (req *CreateSOPCategoryRequest) Validate() error {
	if err := ValidateCategoryName(req.Name); err != nil {
		return err
	}

	if len(req.Description) > 1000 {
		return ErrInvalidDescription
	}

	if req.ImagePath != "" {
		if err := ValidateImagePath(req.ImagePath); err != nil {
			return err
		}
	}

	if req.DisplayOrder < 0 {
		return ErrInvalidDisplayOrder
	}

	return nil
}

// Validate validates the UpdateSOPCategoryRequest
func (req *UpdateSOPCategoryRequest) Validate() error {
	if req.Name != nil {
		if err := ValidateCategoryName(*req.Name); err != nil {
			return err
		}
	}

	if req.Description != nil && len(*req.Description) > 1000 {
		return ErrInvalidDescription
	}

	if req.ImagePath != nil && *req.ImagePath != "" {
		if err := ValidateImagePath(*req.ImagePath); err != nil {
			return err
		}
	}

	if req.DisplayOrder != nil && *req.DisplayOrder < 0 {
		return ErrInvalidDisplayOrder
	}

	return nil
}
