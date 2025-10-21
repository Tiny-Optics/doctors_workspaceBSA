package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidVideoURL          = errors.New("video URL is required")
	ErrInvalidDocumentsPath     = errors.New("documents path is required")
	ErrInvalidNotificationEmail = errors.New("invalid email address in notification list")
	ErrInvalidSMTPHost          = errors.New("SMTP host is required")
	ErrInvalidSMTPPort          = errors.New("SMTP port must be between 1 and 65535")
	ErrInvalidSMTPUsername      = errors.New("SMTP username is required")
	ErrInvalidSMTPPassword      = errors.New("SMTP password is required")
	ErrInvalidFromEmail         = errors.New("from email is required")
	ErrNoNotificationEmails     = errors.New("at least one notification email is required")
)

// SMTPConfig represents SMTP configuration for email notifications
type SMTPConfig struct {
	Host      string `bson:"host" json:"host"`
	Port      int    `bson:"port" json:"port"`
	Username  string `bson:"username" json:"username"`
	Password  string `bson:"password" json:"-"` // Encrypted, not sent in JSON
	FromEmail string `bson:"from_email" json:"fromEmail"`
	FromName  string `bson:"from_name,omitempty" json:"fromName,omitempty"`
}

// RegistryConfig represents the configuration for the African HOPeR Registry
type RegistryConfig struct {
	ID                 primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	VideoURL           string              `bson:"video_url" json:"videoUrl"`
	DocumentsPath      string              `bson:"documents_path" json:"documentsPath"`
	NotificationEmails []string            `bson:"notification_emails" json:"notificationEmails"`
	SMTPConfig         SMTPConfig          `bson:"smtp_config" json:"smtpConfig"`
	CreatedAt          time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt          time.Time           `bson:"updated_at" json:"updatedAt"`
	UpdatedBy          *primitive.ObjectID `bson:"updated_by,omitempty" json:"updatedBy,omitempty"`
}

// UpdateRegistryConfigRequest represents the request to update registry configuration
type UpdateRegistryConfigRequest struct {
	VideoURL           *string   `json:"videoUrl,omitempty"`
	DocumentsPath      *string   `json:"documentsPath,omitempty"`
	NotificationEmails *[]string `json:"notificationEmails,omitempty"`
	SMTPHost           *string   `json:"smtpHost,omitempty"`
	SMTPPort           *int      `json:"smtpPort,omitempty"`
	SMTPUsername       *string   `json:"smtpUsername,omitempty"`
	SMTPPassword       *string   `json:"smtpPassword,omitempty"`
	SMTPFromEmail      *string   `json:"smtpFromEmail,omitempty"`
	SMTPFromName       *string   `json:"smtpFromName,omitempty"`
}

// Validate validates the RegistryConfig
func (c *RegistryConfig) Validate() error {
	if c.VideoURL == "" {
		return ErrInvalidVideoURL
	}

	if c.DocumentsPath == "" {
		return ErrInvalidDocumentsPath
	}

	if len(c.NotificationEmails) == 0 {
		return ErrNoNotificationEmails
	}

	// Validate email format for notification emails
	for _, email := range c.NotificationEmails {
		if err := ValidateEmail(email); err != nil {
			return ErrInvalidNotificationEmail
		}
	}

	// Validate SMTP config
	if err := c.SMTPConfig.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate validates the SMTPConfig
func (s *SMTPConfig) Validate() error {
	if s.Host == "" {
		return ErrInvalidSMTPHost
	}

	if s.Port < 1 || s.Port > 65535 {
		return ErrInvalidSMTPPort
	}

	if s.Username == "" {
		return ErrInvalidSMTPUsername
	}

	if s.Password == "" {
		return ErrInvalidSMTPPassword
	}

	if s.FromEmail == "" {
		return ErrInvalidFromEmail
	}

	if err := ValidateEmail(s.FromEmail); err != nil {
		return ErrInvalidFromEmail
	}

	return nil
}

// IsComplete checks if SMTP configuration is complete (for pre-submission validation)
func (s *SMTPConfig) IsComplete() bool {
	return s.Host != "" && s.Port > 0 && s.Username != "" && s.Password != "" && s.FromEmail != ""
}
