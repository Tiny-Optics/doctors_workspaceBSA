package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DropboxConfig stores Dropbox OAuth configuration and tokens
// This is a singleton - only one document should exist in the collection
type DropboxConfig struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// OAuth Configuration
	AppKey    string `bson:"app_key" json:"appKey"`
	AppSecret string `bson:"app_secret" json:"-"` // Never expose in JSON, encrypted in DB

	// Tokens (all encrypted in database)
	RefreshToken string    `bson:"refresh_token" json:"-"` // Never expose in JSON
	AccessToken  string    `bson:"access_token" json:"-"`  // Never expose in JSON
	TokenExpiry  time.Time `bson:"token_expiry" json:"tokenExpiry"`

	// Configuration
	ParentFolder string `bson:"parent_folder" json:"parentFolder"` // e.g., "/SOPS"

	// Health Monitoring
	IsConnected         bool      `bson:"is_connected" json:"isConnected"`
	LastRefreshSuccess  time.Time `bson:"last_refresh_success" json:"lastRefreshSuccess"`
	LastRefreshAttempt  time.Time `bson:"last_refresh_attempt" json:"lastRefreshAttempt"`
	ConsecutiveFailures int       `bson:"consecutive_failures" json:"consecutiveFailures"`
	LastError           string    `bson:"last_error" json:"lastError"`

	// Metadata
	CreatedAt time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time           `bson:"updated_at" json:"updatedAt"`
	CreatedBy *primitive.ObjectID `bson:"created_by,omitempty" json:"createdBy,omitempty"`
}

// IsTokenExpired checks if the access token has expired or will expire soon (within 5 minutes)
func (d *DropboxConfig) IsTokenExpired() bool {
	if d.TokenExpiry.IsZero() {
		return true
	}
	// Consider expired if less than 5 minutes remaining
	return time.Now().Add(5 * time.Minute).After(d.TokenExpiry)
}

// NeedsReconnection checks if the configuration needs manual reconnection
func (d *DropboxConfig) NeedsReconnection() bool {
	// If we have 3+ consecutive failures or no refresh token, need reconnection
	return d.ConsecutiveFailures >= 3 || d.RefreshToken == ""
}

// GetPublicStatus returns a safe view of the config status (no tokens)
func (d *DropboxConfig) GetPublicStatus() map[string]interface{} {
	return map[string]interface{}{
		"isConnected":         d.IsConnected,
		"tokenExpiry":         d.TokenExpiry,
		"lastRefreshSuccess":  d.LastRefreshSuccess,
		"lastRefreshAttempt":  d.LastRefreshAttempt,
		"consecutiveFailures": d.ConsecutiveFailures,
		"lastError":           d.LastError,
		"needsReconnection":   d.NeedsReconnection(),
		"parentFolder":        d.ParentFolder,
	}
}

// DropboxOAuthResponse represents the response from Dropbox OAuth token endpoint
type DropboxOAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`    // Seconds until expiry
	RefreshToken string `json:"refresh_token"` // Only on initial auth
	Scope        string `json:"scope"`
	UID          string `json:"uid"`
	AccountID    string `json:"account_id"`
}
