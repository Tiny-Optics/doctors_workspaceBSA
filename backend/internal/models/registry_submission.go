package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidSubmissionStatus = errors.New("invalid submission status")
	ErrEmptyFormData           = errors.New("form data cannot be empty")
	ErrNoDocumentsUploaded     = errors.New("at least one document must be uploaded")
)

// SubmissionStatus represents the status of a registry submission
type SubmissionStatus string

const (
	SubmissionStatusSubmitted SubmissionStatus = "submitted"
	SubmissionStatusPending   SubmissionStatus = "pending"
	SubmissionStatusApproved  SubmissionStatus = "approved"
	SubmissionStatusRejected  SubmissionStatus = "rejected"
)

// IsValid checks if the submission status is valid
func (s SubmissionStatus) IsValid() bool {
	switch s {
	case SubmissionStatusSubmitted, SubmissionStatusPending, SubmissionStatusApproved, SubmissionStatusRejected:
		return true
	}
	return false
}

// RegistrySubmission represents a user's submission to the registry
type RegistrySubmission struct {
	ID                primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID            primitive.ObjectID     `bson:"user_id" json:"userId"`
	FormSchemaID      primitive.ObjectID     `bson:"form_schema_id" json:"formSchemaId"`
	FormData          map[string]interface{} `bson:"form_data" json:"formData"`
	DocumentsPath     string                 `bson:"documents_path" json:"documentsPath"`
	UploadedDocuments []string               `bson:"uploaded_documents" json:"uploadedDocuments"`
	Status            SubmissionStatus       `bson:"status" json:"status"`
	CreatedAt         time.Time              `bson:"created_at" json:"createdAt"`
	UpdatedAt         time.Time              `bson:"updated_at" json:"updatedAt"`
	ReviewedBy        *primitive.ObjectID    `bson:"reviewed_by,omitempty" json:"reviewedBy,omitempty"`
	ReviewedAt        *time.Time             `bson:"reviewed_at,omitempty" json:"reviewedAt,omitempty"`
	ReviewNotes       string                 `bson:"review_notes,omitempty" json:"reviewNotes,omitempty"`
	
	// Populated fields (not stored in DB, only for API responses)
	UserName string `bson:"-" json:"userName,omitempty"`
	UserEmail string `bson:"-" json:"userEmail,omitempty"`
	FormName string `bson:"-" json:"formName,omitempty"`
}

// CreateSubmissionRequest represents the request to create a submission
type CreateSubmissionRequest struct {
	FormSchemaID string                 `json:"formSchemaId" binding:"required"`
	FormData     map[string]interface{} `json:"formData" binding:"required"`
}

// UpdateSubmissionStatusRequest represents the request to update submission status
type UpdateSubmissionStatusRequest struct {
	Status      SubmissionStatus `json:"status" binding:"required"`
	ReviewNotes string           `json:"reviewNotes,omitempty"`
}

// Validate validates the RegistrySubmission
func (s *RegistrySubmission) Validate() error {
	if s.UserID.IsZero() {
		return errors.New("user ID is required")
	}

	if s.FormSchemaID.IsZero() {
		return errors.New("form schema ID is required")
	}

	// Note: FormData and UploadedDocuments can be empty depending on the form schema
	// The service layer validates these based on the actual schema requirements

	if !s.Status.IsValid() {
		return ErrInvalidSubmissionStatus
	}

	return nil
}

// Validate validates the CreateSubmissionRequest
func (req *CreateSubmissionRequest) Validate() error {
	if req.FormSchemaID == "" {
		return errors.New("form schema ID is required")
	}

	if _, err := primitive.ObjectIDFromHex(req.FormSchemaID); err != nil {
		return errors.New("invalid form schema ID format")
	}

	// Note: FormData can be empty if the form only has file fields
	// The service layer will validate based on the actual form schema

	return nil
}

// Validate validates the UpdateSubmissionStatusRequest
func (req *UpdateSubmissionStatusRequest) Validate() error {
	if !req.Status.IsValid() {
		return ErrInvalidSubmissionStatus
	}

	return nil
}
