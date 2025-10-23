package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PasswordResetToken represents a password reset token for forgot password functionality
type PasswordResetToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	Token     string             `bson:"token" json:"token"`          // JWT token for API calls
	Code      string             `bson:"code" json:"code"`            // 6-digit verification code
	ExpiresAt time.Time          `bson:"expires_at" json:"expiresAt"` // Token expiration time
	Used      bool               `bson:"used" json:"used"`            // Whether token has been used
	IPAddress string             `bson:"ip_address" json:"ipAddress"` // IP address of requester
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"` // Creation timestamp
}

// ForgotPasswordRequest represents the request to initiate password reset
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ValidateResetCodeRequest represents the request to validate reset code
type ValidateResetCodeRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// ResetPasswordRequest represents the request to reset password with token
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

// ForgotPasswordResponse represents the response for forgot password request
type ForgotPasswordResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// ValidateResetCodeResponse represents the response for code validation
type ValidateResetCodeResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// ResetPasswordResponse represents the response for password reset
type ResetPasswordResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// IsExpired checks if the password reset token has expired
func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// IsValid checks if the password reset token is valid (not expired and not used)
func (p *PasswordResetToken) IsValid() bool {
	return !p.IsExpired() && !p.Used
}
