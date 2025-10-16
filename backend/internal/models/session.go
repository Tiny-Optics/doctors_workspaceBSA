package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session represents an authentication session
type Session struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"userId"`
	Token            string             `bson:"token" json:"token"`
	RefreshToken     string             `bson:"refresh_token" json:"refreshToken"`
	ExpiresAt        time.Time          `bson:"expires_at" json:"expiresAt"`
	RefreshExpiresAt time.Time          `bson:"refresh_expires_at" json:"refreshExpiresAt"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
	IPAddress        string             `bson:"ip_address" json:"ipAddress"`
	UserAgent        string             `bson:"user_agent" json:"userAgent"`
}

// IsExpired checks if the session token has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsRefreshExpired checks if the refresh token has expired
func (s *Session) IsRefreshExpired() bool {
	return time.Now().After(s.RefreshExpiresAt)
}
