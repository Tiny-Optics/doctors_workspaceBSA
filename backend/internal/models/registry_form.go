package models

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidFormName      = errors.New("form name is required and must be between 1 and 200 characters")
	ErrInvalidFieldID       = errors.New("field ID is required and must be alphanumeric with underscores")
	ErrInvalidFieldLabel    = errors.New("field label is required")
	ErrInvalidFieldType     = errors.New("invalid field type")
	ErrDuplicateFieldID     = errors.New("duplicate field ID found")
	ErrNoFileField          = errors.New("at least one file field is required")
	ErrInvalidSelectOptions = errors.New("select and radio fields must have at least one option")
	ErrInvalidMaxLength     = errors.New("max length must be greater than 0")
	ErrInvalidMinValue      = errors.New("min value must be less than or equal to max value")
)

// FormFieldType represents the type of form field
type FormFieldType string

const (
	FieldTypeText     FormFieldType = "text"
	FieldTypeTextarea FormFieldType = "textarea"
	FieldTypeSelect   FormFieldType = "select"
	FieldTypeRadio    FormFieldType = "radio"
	FieldTypeDate     FormFieldType = "date"
	FieldTypeNumber   FormFieldType = "number"
	FieldTypeEmail    FormFieldType = "email"
	FieldTypeFile     FormFieldType = "file"
)

// IsValid checks if the field type is valid
func (ft FormFieldType) IsValid() bool {
	switch ft {
	case FieldTypeText, FieldTypeTextarea, FieldTypeSelect, FieldTypeRadio,
		FieldTypeDate, FieldTypeNumber, FieldTypeEmail, FieldTypeFile:
		return true
	}
	return false
}

// ValidationRules represents validation rules for a form field
type ValidationRules struct {
	MinLength *int     `bson:"min_length,omitempty" json:"minLength,omitempty"`
	MaxLength *int     `bson:"max_length,omitempty" json:"maxLength,omitempty"`
	MinValue  *float64 `bson:"min_value,omitempty" json:"minValue,omitempty"`
	MaxValue  *float64 `bson:"max_value,omitempty" json:"maxValue,omitempty"`
	Pattern   *string  `bson:"pattern,omitempty" json:"pattern,omitempty"`
}

// FormField represents a single field in the form schema
type FormField struct {
	ID              string          `bson:"id" json:"id"`
	Label           string          `bson:"label" json:"label"`
	Type            FormFieldType   `bson:"type" json:"type"`
	Required        bool            `bson:"required" json:"required"`
	Placeholder     string          `bson:"placeholder,omitempty" json:"placeholder,omitempty"`
	HelpText        string          `bson:"help_text,omitempty" json:"helpText,omitempty"`
	Options         []string        `bson:"options,omitempty" json:"options,omitempty"`              // For select/radio
	AllowMultiple   bool            `bson:"allow_multiple,omitempty" json:"allowMultiple,omitempty"` // For file fields
	ValidationRules ValidationRules `bson:"validation_rules,omitempty" json:"validationRules,omitempty"`
	DisplayOrder    int             `bson:"display_order" json:"displayOrder"`
}

// RegistryFormSchema represents the schema for the registry submission form
type RegistryFormSchema struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	FormName    string              `bson:"form_name" json:"formName"`
	Description string              `bson:"description,omitempty" json:"description,omitempty"`
	Fields      []FormField         `bson:"fields" json:"fields"`
	IsActive    bool                `bson:"is_active" json:"isActive"`
	CreatedAt   time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updatedAt"`
	CreatedBy   *primitive.ObjectID `bson:"created_by,omitempty" json:"createdBy,omitempty"`
	UpdatedBy   *primitive.ObjectID `bson:"updated_by,omitempty" json:"updatedBy,omitempty"`
}

// CreateFormSchemaRequest represents the request to create a form schema
type CreateFormSchemaRequest struct {
	FormName    string      `json:"formName" binding:"required"`
	Description string      `json:"description"`
	Fields      []FormField `json:"fields" binding:"required"`
}

// UpdateFormSchemaRequest represents the request to update a form schema
type UpdateFormSchemaRequest struct {
	FormName    *string      `json:"formName,omitempty"`
	Description *string      `json:"description,omitempty"`
	Fields      *[]FormField `json:"fields,omitempty"`
	IsActive    *bool        `json:"isActive,omitempty"`
}

// Validate validates the RegistryFormSchema
func (s *RegistryFormSchema) Validate() error {
	if err := ValidateFormName(s.FormName); err != nil {
		return err
	}

	if len(s.Fields) == 0 {
		return errors.New("form must have at least one field")
	}

	// Check for at least one file field
	hasFileField := false
	fieldIDs := make(map[string]bool)

	for _, field := range s.Fields {
		// Validate field
		if err := field.Validate(); err != nil {
			return err
		}

		// Check for duplicate field IDs
		if fieldIDs[field.ID] {
			return ErrDuplicateFieldID
		}
		fieldIDs[field.ID] = true

		// Check for file field
		if field.Type == FieldTypeFile {
			hasFileField = true
		}
	}

	if !hasFileField {
		return ErrNoFileField
	}

	return nil
}

// Validate validates a FormField
func (f *FormField) Validate() error {
	// Validate ID (alphanumeric and underscores only)
	if f.ID == "" {
		return ErrInvalidFieldID
	}
	validID := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validID.MatchString(f.ID) {
		return ErrInvalidFieldID
	}

	// Validate label
	if strings.TrimSpace(f.Label) == "" {
		return ErrInvalidFieldLabel
	}

	// Validate type
	if !f.Type.IsValid() {
		return ErrInvalidFieldType
	}

	// Validate options for select/radio fields
	if (f.Type == FieldTypeSelect || f.Type == FieldTypeRadio) && len(f.Options) == 0 {
		return ErrInvalidSelectOptions
	}

	// Validate validation rules
	if f.ValidationRules.MaxLength != nil && *f.ValidationRules.MaxLength <= 0 {
		return ErrInvalidMaxLength
	}

	if f.ValidationRules.MinValue != nil && f.ValidationRules.MaxValue != nil {
		if *f.ValidationRules.MinValue > *f.ValidationRules.MaxValue {
			return ErrInvalidMinValue
		}
	}

	// Validate pattern if provided
	if f.ValidationRules.Pattern != nil && *f.ValidationRules.Pattern != "" {
		_, err := regexp.Compile(*f.ValidationRules.Pattern)
		if err != nil {
			return errors.New("invalid validation pattern: " + err.Error())
		}
	}

	return nil
}

// ValidateFormName validates the form name
func ValidateFormName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 200 {
		return ErrInvalidFormName
	}
	return nil
}

// Validate validates the CreateFormSchemaRequest
func (req *CreateFormSchemaRequest) Validate() error {
	if err := ValidateFormName(req.FormName); err != nil {
		return err
	}

	if len(req.Fields) == 0 {
		return errors.New("form must have at least one field")
	}

	// Check for at least one file field
	hasFileField := false
	fieldIDs := make(map[string]bool)

	for _, field := range req.Fields {
		if err := field.Validate(); err != nil {
			return err
		}

		// Check for duplicate field IDs
		if fieldIDs[field.ID] {
			return ErrDuplicateFieldID
		}
		fieldIDs[field.ID] = true

		if field.Type == FieldTypeFile {
			hasFileField = true
		}
	}

	if !hasFileField {
		return ErrNoFileField
	}

	return nil
}

// Validate validates the UpdateFormSchemaRequest
func (req *UpdateFormSchemaRequest) Validate() error {
	if req.FormName != nil {
		if err := ValidateFormName(*req.FormName); err != nil {
			return err
		}
	}

	if req.Fields != nil {
		if len(*req.Fields) == 0 {
			return errors.New("form must have at least one field")
		}

		hasFileField := false
		fieldIDs := make(map[string]bool)

		for _, field := range *req.Fields {
			if err := field.Validate(); err != nil {
				return err
			}

			if fieldIDs[field.ID] {
				return ErrDuplicateFieldID
			}
			fieldIDs[field.ID] = true

			if field.Type == FieldTypeFile {
				hasFileField = true
			}
		}

		if !hasFileField {
			return ErrNoFileField
		}
	}

	return nil
}
