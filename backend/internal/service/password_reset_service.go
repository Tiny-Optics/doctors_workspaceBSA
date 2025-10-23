package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrUserInactive              = errors.New("user account is inactive")
	ErrTooManyRequests           = errors.New("too many password reset requests")
	ErrInvalidResetCode          = errors.New("invalid or expired reset code")
	ErrResetTokenUsed            = errors.New("reset token has already been used")
	ErrSMTPNotConfigured         = errors.New("SMTP configuration is not complete")
	ErrPasswordResetTokenExpired = errors.New("password reset token has expired")
)

const (
	// Password reset token expires in 15 minutes
	passwordResetTokenExpiry = 15 * time.Minute
	// Maximum 3 password reset requests per hour per user
	maxRequestsPerHour = 3
	// Maximum 5 password reset requests per hour per IP
	maxRequestsPerHourPerIP = 5
)

// PasswordResetService handles password reset operations
type PasswordResetService struct {
	userRepo          *repository.UserRepository
	passwordResetRepo *repository.PasswordResetRepository
	auditRepo         *repository.AuditRepository
	emailService      *EmailService
	encryptionService *EncryptionService
	registryService   *RegistryService
	jwtSecret         string
}

// NewPasswordResetService creates a new PasswordResetService
func NewPasswordResetService(
	userRepo *repository.UserRepository,
	passwordResetRepo *repository.PasswordResetRepository,
	auditRepo *repository.AuditRepository,
	emailService *EmailService,
	encryptionService *EncryptionService,
	registryService *RegistryService,
	jwtSecret string,
) *PasswordResetService {
	return &PasswordResetService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		auditRepo:         auditRepo,
		emailService:      emailService,
		encryptionService: encryptionService,
		registryService:   registryService,
		jwtSecret:         jwtSecret,
	}
}

// RequestPasswordReset initiates a password reset request
func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, email, ipAddress string) (*models.ForgotPasswordResponse, error) {
	// Normalize email
	email = strings.ToLower(strings.TrimSpace(email))

	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Log failed attempt for security
		s.auditRepo.Create(ctx, &models.AuditLog{
			Action:    models.AuditActionPasswordResetRequested,
			IPAddress: ipAddress,
			Details: map[string]interface{}{
				"email":  email,
				"reason": "user not found",
			},
		})
		return nil, ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive {
		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionPasswordResetRequested,
			IPAddress: ipAddress,
			Details: map[string]interface{}{
				"email":  email,
				"reason": "account inactive",
			},
		})
		return nil, ErrUserInactive
	}

	// Check rate limiting for user
	oneHourAgo := time.Now().Add(-time.Hour)
	userRequestCount, err := s.passwordResetRepo.CountRecentRequests(ctx, user.ID, oneHourAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to check user rate limit: %w", err)
	}
	if userRequestCount >= maxRequestsPerHour {
		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionPasswordResetRequested,
			IPAddress: ipAddress,
			Details: map[string]interface{}{
				"email":  email,
				"reason": "too many requests per user",
			},
		})
		return nil, ErrTooManyRequests
	}

	// Check rate limiting for IP
	ipRequestCount, err := s.passwordResetRepo.CountRecentRequestsByIP(ctx, ipAddress, oneHourAgo)
	if err != nil {
		return nil, fmt.Errorf("failed to check IP rate limit: %w", err)
	}
	if ipRequestCount >= maxRequestsPerHourPerIP {
		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionPasswordResetRequested,
			IPAddress: ipAddress,
			Details: map[string]interface{}{
				"email":  email,
				"reason": "too many requests per IP",
			},
		})
		return nil, ErrTooManyRequests
	}

	// Generate 6-digit verification code
	code, err := s.generateVerificationCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Generate JWT token for API calls
	token, err := s.generateResetToken(user.ID, code)
	if err != nil {
		return nil, fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Create password reset token
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		Code:      code,
		ExpiresAt: time.Now().Add(passwordResetTokenExpiry),
		Used:      false,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := s.passwordResetRepo.Create(ctx, resetToken); err != nil {
		return nil, fmt.Errorf("failed to save password reset token: %w", err)
	}

	// Get SMTP configuration from registry service
	smtpConfig, err := s.registryService.GetPublicSMTPConfig(ctx)
	if err != nil {
		// Log SMTP config failure
		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionPasswordResetRequested,
			IPAddress: ipAddress,
			Details: map[string]interface{}{
				"email":  email,
				"reason": "SMTP config not available",
				"error":  err.Error(),
			},
		})
		return nil, ErrSMTPNotConfigured
	}

	// Send email with verification code
	fmt.Printf("DEBUG: Attempting to send password reset email to %s with code %s\n", user.Email, code)
	fmt.Printf("DEBUG: SMTP Config - Host: %s, Port: %d, Username: %s, FromEmail: %s\n",
		smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.FromEmail)

	if err := s.emailService.SendPasswordResetEmail(*smtpConfig, user.Email, code, user.Profile.FirstName+" "+user.Profile.LastName); err != nil {
		fmt.Printf("DEBUG: Email send failed: %v\n", err)
		// Log email failure but don't fail the request
		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionPasswordResetRequested,
			IPAddress: ipAddress,
			Details: map[string]interface{}{
				"email":  email,
				"reason": "email send failed",
				"error":  err.Error(),
			},
		})
		// Continue with success response even if email fails
	} else {
		fmt.Printf("DEBUG: Email sent successfully to %s\n", user.Email)
	}

	// Log successful request
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:    &user.ID,
		Action:    models.AuditActionPasswordResetRequested,
		IPAddress: ipAddress,
		Details: map[string]interface{}{
			"email": email,
		},
	})

	return &models.ForgotPasswordResponse{
		Message: "Password reset code has been sent to your email address",
		Success: true,
	}, nil
}

// ValidateResetCode validates a password reset code and returns a token
func (s *PasswordResetService) ValidateResetCode(ctx context.Context, code string) (*models.ValidateResetCodeResponse, error) {
	// Find token by code
	resetToken, err := s.passwordResetRepo.FindByCode(ctx, code)
	if err != nil {
		if err == repository.ErrPasswordResetTokenNotFound {
			return nil, ErrInvalidResetCode
		}
		return nil, fmt.Errorf("failed to find reset token: %w", err)
	}

	// Check if token is valid
	if !resetToken.IsValid() {
		if resetToken.IsExpired() {
			return nil, ErrPasswordResetTokenExpired
		}
		if resetToken.Used {
			return nil, ErrResetTokenUsed
		}
		return nil, ErrInvalidResetCode
	}

	// Return the JWT token for the next step
	return &models.ValidateResetCodeResponse{
		Token:   resetToken.Token,
		Message: "Code validated successfully",
		Success: true,
	}, nil
}

// ResetPassword resets the user's password using a valid token
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) (*models.ResetPasswordResponse, error) {
	// Find token
	resetToken, err := s.passwordResetRepo.FindByToken(ctx, token)
	if err != nil {
		if err == repository.ErrPasswordResetTokenNotFound {
			return nil, ErrInvalidResetCode
		}
		return nil, fmt.Errorf("failed to find reset token: %w", err)
	}

	// Check if token is valid
	if !resetToken.IsValid() {
		if resetToken.IsExpired() {
			return nil, ErrPasswordResetTokenExpired
		}
		if resetToken.Used {
			return nil, ErrResetTokenUsed
		}
		return nil, ErrInvalidResetCode
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Hash new password
	hashedPassword, err := s.hashPassword(newPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	if err := s.userRepo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(ctx, resetToken.ID); err != nil {
		// Log error but don't fail the request
		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionPasswordResetCompleted,
			IPAddress: resetToken.IPAddress,
			Details: map[string]interface{}{
				"error": "failed to mark token as used",
			},
		})
	}

	// Log successful password reset
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:    &user.ID,
		Action:    models.AuditActionPasswordResetCompleted,
		IPAddress: resetToken.IPAddress,
		Details: map[string]interface{}{
			"user_id": user.ID.Hex(),
			"email":   user.Email,
		},
	})

	return &models.ResetPasswordResponse{
		Message: "Password has been reset successfully",
		Success: true,
	}, nil
}

// generateVerificationCode generates a 6-digit verification code
func (s *PasswordResetService) generateVerificationCode() (string, error) {
	// Generate 6 random digits
	code := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += n.String()
	}
	return code, nil
}

// generateResetToken generates a JWT token for password reset
func (s *PasswordResetService) generateResetToken(userID primitive.ObjectID, code string) (string, error) {
	expiresAt := time.Now().Add(passwordResetTokenExpiry)

	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"code":    code,
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
		"type":    "password_reset",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// hashPassword hashes a password using bcrypt
func (s *PasswordResetService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CleanupExpiredTokens removes expired password reset tokens
func (s *PasswordResetService) CleanupExpiredTokens(ctx context.Context) error {
	return s.passwordResetRepo.CleanupExpiredTokens(ctx)
}
