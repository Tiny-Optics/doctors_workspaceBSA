package models

import (
	"errors"
	"regexp"
	"time"
	"unicode"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the Doctor's Workspace system
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`

	// Role & Permissions
	Role       UserRole   `bson:"role" json:"role"`
	AdminLevel AdminLevel `bson:"admin_level,omitempty" json:"adminLevel,omitempty"`
	IsActive   bool       `bson:"is_active" json:"isActive"`

	// Extended Profile
	Profile UserProfile `bson:"profile" json:"profile"`

	// Metadata
	CreatedAt   time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updatedAt"`
	LastLoginAt *time.Time          `bson:"last_login_at,omitempty" json:"lastLoginAt,omitempty"`
	CreatedBy   *primitive.ObjectID `bson:"created_by,omitempty" json:"createdBy,omitempty"`

	// Security
	FailedLoginAttempts int        `bson:"failed_login_attempts" json:"-"`
	LockedUntil         *time.Time `bson:"locked_until,omitempty" json:"-"`
}

// UserProfile contains extended profile information for a user
type UserProfile struct {
	FirstName          string              `bson:"first_name" json:"firstName"`
	LastName           string              `bson:"last_name" json:"lastName"`
	InstitutionID      *primitive.ObjectID `bson:"institution_id,omitempty" json:"institutionId,omitempty"`
	Specialty          string              `bson:"specialty,omitempty" json:"specialty,omitempty"`
	RegistrationNumber string              `bson:"registration_number,omitempty" json:"registrationNumber,omitempty"`
	PhoneNumber        string              `bson:"phone_number,omitempty" json:"phoneNumber,omitempty"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	Username           string     `json:"username" binding:"required"`
	Email              string     `json:"email" binding:"required,email"`
	Password           string     `json:"password" binding:"required"`
	Role               UserRole   `json:"role" binding:"required"`
	AdminLevel         AdminLevel `json:"adminLevel,omitempty"`
	FirstName          string     `json:"firstName" binding:"required"`
	LastName           string     `json:"lastName" binding:"required"`
	InstitutionID      string     `json:"institutionId" binding:"required"`
	Specialty          string     `json:"specialty,omitempty"`
	RegistrationNumber string     `json:"registrationNumber,omitempty"`
	PhoneNumber        string     `json:"phoneNumber,omitempty"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FirstName          *string     `json:"firstName,omitempty"`
	LastName           *string     `json:"lastName,omitempty"`
	InstitutionID      *string     `json:"institutionId,omitempty"`
	Specialty          *string     `json:"specialty,omitempty"`
	RegistrationNumber *string     `json:"registrationNumber,omitempty"`
	PhoneNumber        *string     `json:"phoneNumber,omitempty"`
	Role               *UserRole   `json:"role,omitempty"`
	AdminLevel         *AdminLevel `json:"adminLevel,omitempty"`
	IsActive           *bool       `json:"isActive,omitempty"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
	User         *User     `json:"user"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

// Validation errors
var (
	ErrInvalidEmail              = errors.New("invalid email format")
	ErrInvalidUsername           = errors.New("username must be 3-50 characters, alphanumeric and underscores only")
	ErrWeakPassword              = errors.New("password must be at least 8 characters and include uppercase, lowercase, number, and special character")
	ErrInvalidRole               = errors.New("invalid user role")
	ErrInvalidAdminLevel         = errors.New("invalid admin level")
	ErrProfileFieldTooShort      = errors.New("profile field is too short (minimum 2 characters)")
	ErrProfileFieldTooLong       = errors.New("profile field is too long (maximum 100 characters)")
	ErrFieldRequired             = errors.New("required field is missing")
	ErrInvalidRegistrationNumber = errors.New("invalid registration number format")
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)
)

// Validate validates the CreateUserRequest
func (req *CreateUserRequest) Validate() error {
	// Validate email
	if !emailRegex.MatchString(req.Email) {
		return ErrInvalidEmail
	}

	// Validate username
	if !usernameRegex.MatchString(req.Username) {
		return ErrInvalidUsername
	}

	// Validate password
	if err := ValidatePassword(req.Password); err != nil {
		return err
	}

	// Validate role
	if !req.Role.IsValid() {
		return ErrInvalidRole
	}

	// Validate admin level if role is admin
	if req.Role == RoleAdmin {
		if !req.AdminLevel.IsValid() {
			return ErrInvalidAdminLevel
		}
		if req.AdminLevel == AdminLevelNone {
			return errors.New("admin role requires an admin level")
		}
	} else if req.AdminLevel != AdminLevelNone {
		return errors.New("non-admin roles cannot have an admin level")
	}

	// Validate profile fields
	if len(req.FirstName) < 2 || len(req.FirstName) > 100 {
		return ErrProfileFieldTooShort
	}
	if len(req.LastName) < 2 || len(req.LastName) > 100 {
		return ErrProfileFieldTooShort
	}

	// Validate institution ID
	if req.InstitutionID == "" {
		return errors.New("institution ID is required")
	}
	if _, err := primitive.ObjectIDFromHex(req.InstitutionID); err != nil {
		return errors.New("invalid institution ID format")
	}

	return nil
}

// ValidatePassword checks if a password meets security requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrWeakPassword
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return ErrWeakPassword
	}

	return nil
}

// IsLocked returns true if the user account is currently locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permission Permission) bool {
	return HasPermission(u.Role, u.AdminLevel, permission)
}

// GetPermissions returns all permissions for the user
func (u *User) GetPermissions() []Permission {
	return GetPermissionsForRole(u.Role, u.AdminLevel)
}

// CanManageUser checks if the user can manage another user based on roles
func (u *User) CanManageUser(targetUser *User) bool {
	// Must have manage users permission
	if !u.HasPermission(PermManageUsers) {
		return false
	}

	// Super admin can manage anyone
	if u.AdminLevel == AdminLevelSuperAdmin {
		return true
	}

	// User managers can only manage non-admin users
	if u.AdminLevel == AdminLevelUserManager {
		return targetUser.Role != RoleAdmin
	}

	return false
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.Profile.FirstName + " " + u.Profile.LastName
}
