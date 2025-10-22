package models

import (
	"errors"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidRedCapURL  = errors.New("invalid REDCap URL format")
	ErrRedCapURLRequired = errors.New("REDCap URL is required")
)

// ReferralConfig represents the configuration for REDCap referral integration
type ReferralConfig struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	RedCapURL    string              `bson:"redcap_url" json:"redcapUrl"`
	IsEnabled    bool                `bson:"is_enabled" json:"isEnabled"`
	IsConfigured bool                `bson:"is_configured" json:"isConfigured"`
	CreatedAt    time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updated_at" json:"updatedAt"`
	UpdatedBy    *primitive.ObjectID `bson:"updated_by,omitempty" json:"updatedBy,omitempty"`
}

// UpdateReferralConfigRequest represents the request to update referral configuration
type UpdateReferralConfigRequest struct {
	RedCapURL *string `json:"redcapUrl,omitempty"`
	IsEnabled *bool   `json:"isEnabled,omitempty"`
}

// Validate validates the ReferralConfig
func (c *ReferralConfig) Validate() error {
	// If enabled, URL must be present and valid
	if c.IsEnabled {
		if c.RedCapURL == "" {
			return ErrRedCapURLRequired
		}

		if err := validateURL(c.RedCapURL); err != nil {
			return ErrInvalidRedCapURL
		}
	}

	return nil
}

// validateURL checks if the URL is valid
func validateURL(urlStr string) error {
	if urlStr == "" {
		return errors.New("URL cannot be empty")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	// Ensure it has a scheme (http or https)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must start with http:// or https://")
	}

	// Ensure it has a host
	if parsedURL.Host == "" {
		return errors.New("URL must have a valid host")
	}

	return nil
}
