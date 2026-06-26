package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WorkingPartyCategory represents a category for Working Parties documents
type WorkingPartyCategory struct {
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

// CreateWorkingPartyCategoryRequest represents the request to create a new working party category
type CreateWorkingPartyCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ImagePath    string `json:"imagePath"`
	DisplayOrder int    `json:"displayOrder"`
}

// UpdateWorkingPartyCategoryRequest represents the request to update a working party category
type UpdateWorkingPartyCategoryRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	ImagePath    *string `json:"imagePath"`
	DisplayOrder *int    `json:"displayOrder"`
	IsActive     *bool   `json:"isActive"`
}

// Validate validates the WorkingPartyCategory fields
func (c *WorkingPartyCategory) Validate() error {
	if err := ValidateCategoryName(c.Name); err != nil {
		return err
	}

	if c.Slug == "" {
		return ErrInvalidSlug
	}

	if len(c.Description) > 1000 {
		return ErrInvalidDescription
	}

	if c.ImagePath != "" {
		if err := ValidateImagePath(c.ImagePath); err != nil {
			return err
		}
	}

	if c.DisplayOrder < 0 {
		return ErrInvalidDisplayOrder
	}

	return nil
}

// GetDropboxPath generates the Dropbox path for the category
func (c *WorkingPartyCategory) GetDropboxPath() string {
	if c.DropboxPath != "" {
		return c.DropboxPath
	}
	return "WORKING_PARTIES/" + c.Name
}

// Validate validates the CreateWorkingPartyCategoryRequest
func (req *CreateWorkingPartyCategoryRequest) Validate() error {
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

// Validate validates the UpdateWorkingPartyCategoryRequest
func (req *UpdateWorkingPartyCategoryRequest) Validate() error {
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
