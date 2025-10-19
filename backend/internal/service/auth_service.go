package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost          = 12
	tokenExpiry         = 24 * time.Hour
	refreshTokenExpiry  = 30 * 24 * time.Hour
	maxFailedAttempts   = 5
	accountLockDuration = 30 * time.Minute
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrAccountLocked      = errors.New("account is locked due to too many failed login attempts")
	ErrAccountInactive    = errors.New("account is not active")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	auditRepo   *repository.AuditRepository
	jwtSecret   string
}

// NewAuthService creates a new AuthService
func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	auditRepo *repository.AuditRepository,
) *AuthService {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production" // Should be set via env var
	}

	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		auditRepo:   auditRepo,
		jwtSecret:   jwtSecret,
	}
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a hashed password with a plain text password
func (s *AuthService) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// generateToken generates a random token
func (s *AuthService) generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateJWT generates a JWT token for a user
func (s *AuthService) generateJWT(user *models.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(tokenExpiry)

	claims := jwt.MapClaims{
		"user_id":     user.ID.Hex(),
		"email":       user.Email,
		"role":        user.Role,
		"admin_level": user.AdminLevel,
		"exp":         expiresAt.Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateJWT validates a JWT token and returns the claims
func (s *AuthService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest, ipAddress, userAgent string) (*models.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Log failed login attempt
		s.auditRepo.Create(ctx, &models.AuditLog{
			Action:    models.AuditActionLoginFailed,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			Details: map[string]interface{}{
				"email":  req.Email,
				"reason": "user not found",
			},
		})
		return nil, ErrInvalidCredentials
	}

	// Check if account is locked
	if user.IsLocked() {
		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionLoginFailed,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			Details: map[string]interface{}{
				"email":  req.Email,
				"reason": "account locked",
			},
		})
		return nil, ErrAccountLocked
	}

	// Check if account is active
	if !user.IsActive {
		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionLoginFailed,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			Details: map[string]interface{}{
				"email":  req.Email,
				"reason": "account inactive",
			},
		})
		return nil, ErrAccountInactive
	}

	// Check password
	if !s.CheckPassword(user.PasswordHash, req.Password) {
		// Increment failed login attempts
		user.FailedLoginAttempts++
		s.userRepo.IncrementFailedLoginAttempts(ctx, user.ID)

		// Lock account if too many failed attempts
		if user.FailedLoginAttempts >= maxFailedAttempts {
			lockUntil := time.Now().Add(accountLockDuration)
			s.userRepo.LockAccount(ctx, user.ID, lockUntil)

			s.auditRepo.Create(ctx, &models.AuditLog{
				UserID:    &user.ID,
				Action:    models.AuditActionAccountLocked,
				IPAddress: ipAddress,
				UserAgent: userAgent,
				Details: map[string]interface{}{
					"lock_until": lockUntil,
					"reason":     "too many failed login attempts",
				},
			})

			return nil, ErrAccountLocked
		}

		s.auditRepo.Create(ctx, &models.AuditLog{
			UserID:    &user.ID,
			Action:    models.AuditActionLoginFailed,
			IPAddress: ipAddress,
			UserAgent: userAgent,
			Details: map[string]interface{}{
				"email":           req.Email,
				"reason":          "invalid password",
				"failed_attempts": user.FailedLoginAttempts,
			},
		})

		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, expiresAt, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := s.generateToken()
	if err != nil {
		return nil, err
	}

	// Create session
	session := &models.Session{
		UserID:           user.ID,
		Token:            token,
		RefreshToken:     refreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: time.Now().Add(refreshTokenExpiry),
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	// Update last login
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, err
	}

	// Log successful login
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:    &user.ID,
		Action:    models.AuditActionLoginSuccess,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Details: map[string]interface{}{
			"email": req.Email,
		},
	})

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout logs out a user by deleting their session
func (s *AuthService) Logout(ctx context.Context, token string, userID primitive.ObjectID, ipAddress string) error {
	// Delete session
	if err := s.sessionRepo.Delete(ctx, token); err != nil {
		return err
	}

	// Log logout
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:    &userID,
		Action:    models.AuditActionLogout,
		IPAddress: ipAddress,
	})

	return nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	// Find session by refresh token
	session, err := s.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, ErrAccountInactive
	}

	// Generate new JWT token
	token, expiresAt, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := s.generateToken()
	if err != nil {
		return nil, err
	}

	// Delete old session
	s.sessionRepo.Delete(ctx, session.Token)

	// Create new session
	newSession := &models.Session{
		UserID:           user.ID,
		Token:            token,
		RefreshToken:     newRefreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: time.Now().Add(refreshTokenExpiry),
		IPAddress:        session.IPAddress,
		UserAgent:        session.UserAgent,
	}

	if err := s.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		User:         user,
		ExpiresAt:    expiresAt,
	}, nil
}

// GetUserFromToken extracts user information from a JWT token
func (s *AuthService) GetUserFromToken(ctx context.Context, tokenString string) (*models.User, error) {
	claims, err := s.ValidateJWT(tokenString)
	if err != nil {
		return nil, err
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdatePassword updates a user's password
func (s *AuthService) UpdatePassword(ctx context.Context, userID primitive.ObjectID, newPasswordHash, ipAddress string) error {
	// Update password in database
	err := s.userRepo.Update(ctx, userID, bson.M{
		"password_hash": newPasswordHash,
	})
	if err != nil {
		return err
	}

	// Log password change
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &userID,
		PerformedBy: &userID,
		Action:      "password.change",
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"changed_at": time.Now(),
		},
	})

	return nil
}
