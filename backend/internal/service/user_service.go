package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUnauthorized              = errors.New("unauthorized to perform this action")
	ErrCannotModifyOwnAdminLevel = errors.New("cannot modify your own admin level")
)

// UserService handles user management operations
type UserService struct {
	userRepo    *repository.UserRepository
	auditRepo   *repository.AuditRepository
	authService *AuthService
}

// NewUserService creates a new UserService
func NewUserService(
	userRepo *repository.UserRepository,
	auditRepo *repository.AuditRepository,
	authService *AuthService,
) *UserService {
	return &UserService{
		userRepo:    userRepo,
		auditRepo:   auditRepo,
		authService: authService,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest, createdBy *models.User, ipAddress string) (*models.User, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if creator has permission to create users
	if !createdBy.HasPermission(models.PermManageUsers) {
		return nil, ErrUnauthorized
	}

	// Check if user manager is trying to create an admin
	if req.Role == models.RoleAdmin && createdBy.AdminLevel == models.AdminLevelUserManager {
		return nil, errors.New("user managers cannot create admin accounts")
	}

	// Normalize email to lowercase for case-insensitive comparison
	normalizedEmail := strings.ToLower(strings.TrimSpace(req.Email))
	
	// Check if email already exists
	emailExists, err := s.userRepo.EmailExists(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, repository.ErrDuplicateEmail
	}

	// Check if username already exists
	usernameExists, err := s.userRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, repository.ErrDuplicateUsername
	}

	// Hash password
	passwordHash, err := s.authService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Parse institution ID
	institutionID, err := primitive.ObjectIDFromHex(req.InstitutionID)
	if err != nil {
		return nil, errors.New("invalid institution ID format")
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        normalizedEmail,
		PasswordHash: passwordHash,
		Role:         req.Role,
		AdminLevel:   req.AdminLevel,
		IsActive:     true,
		Profile: models.UserProfile{
			FirstName:          req.FirstName,
			LastName:           req.LastName,
			InstitutionID:      &institutionID,
			Specialty:          req.Specialty,
			RegistrationNumber: req.RegistrationNumber,
			PhoneNumber:        req.PhoneNumber,
		},
		CreatedBy: &createdBy.ID,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &user.ID,
		PerformedBy: &createdBy.ID,
		Action:      models.AuditActionUserCreated,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"username":    user.Username,
			"email":       user.Email,
			"role":        user.Role,
			"admin_level": user.AdminLevel,
		},
	})

	return user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, userID primitive.ObjectID, req *models.UpdateUserRequest, updatedBy *models.User, ipAddress string) (*models.User, error) {
	// Get target user
	targetUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if updater can manage this user
	if !updatedBy.CanManageUser(targetUser) && updatedBy.ID != targetUser.ID {
		return nil, ErrUnauthorized
	}

	// Prevent users from changing their own admin level
	if req.AdminLevel != nil && updatedBy.ID == targetUser.ID {
		return nil, ErrCannotModifyOwnAdminLevel
	}

	// Build update document
	update := bson.M{}
	details := make(map[string]interface{})

	if req.FirstName != nil {
		update["profile.first_name"] = *req.FirstName
		details["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		update["profile.last_name"] = *req.LastName
		details["last_name"] = *req.LastName
	}
	if req.InstitutionID != nil {
		institutionID, err := primitive.ObjectIDFromHex(*req.InstitutionID)
		if err != nil {
			return nil, errors.New("invalid institution ID format")
		}
		update["profile.institution_id"] = institutionID
		details["institution_id"] = *req.InstitutionID
	}
	if req.Specialty != nil {
		update["profile.specialty"] = *req.Specialty
		details["specialty"] = *req.Specialty
	}
	if req.RegistrationNumber != nil {
		update["profile.registration_number"] = *req.RegistrationNumber
		details["registration_number"] = *req.RegistrationNumber
	}
	if req.PhoneNumber != nil {
		update["profile.phone_number"] = *req.PhoneNumber
		details["phone_number"] = *req.PhoneNumber
	}

	// Only admins can change role, admin level, and active status
	if updatedBy.HasPermission(models.PermManageUsers) {
		if req.Role != nil {
			// User managers cannot change someone to admin
			if *req.Role == models.RoleAdmin && updatedBy.AdminLevel == models.AdminLevelUserManager {
				return nil, errors.New("user managers cannot assign admin role")
			}
			update["role"] = *req.Role
			details["role"] = *req.Role
		}
		if req.AdminLevel != nil {
			// Only super admins can change admin level
			if updatedBy.AdminLevel != models.AdminLevelSuperAdmin {
				return nil, errors.New("only super admins can change admin level")
			}
			update["admin_level"] = *req.AdminLevel
			details["admin_level"] = *req.AdminLevel
		}
		if req.IsActive != nil {
			update["is_active"] = *req.IsActive
			details["is_active"] = *req.IsActive
		}
	}

	// Update user
	if len(update) > 0 {
		if err := s.userRepo.Update(ctx, userID, update); err != nil {
			return nil, err
		}
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &userID,
		PerformedBy: &updatedBy.ID,
		Action:      models.AuditActionUserUpdated,
		IPAddress:   ipAddress,
		Details:     details,
	})

	// Get updated user
	return s.userRepo.FindByID(ctx, userID)
}

// DeleteUser deletes a user (hard delete)
func (s *UserService) DeleteUser(ctx context.Context, userID primitive.ObjectID, deletedBy *models.User, ipAddress string) error {
	// Check if deleter has permission
	if !deletedBy.HasPermission(models.PermDeleteUsers) {
		return ErrUnauthorized
	}

	// Get target user
	targetUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Prevent deleting yourself
	if deletedBy.ID == targetUser.ID {
		return errors.New("cannot delete your own account")
	}

	// Delete user
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &userID,
		PerformedBy: &deletedBy.ID,
		Action:      models.AuditActionUserDeleted,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"username":   targetUser.Username,
			"email":      targetUser.Email,
			"first_name": targetUser.Profile.FirstName,
			"last_name":  targetUser.Profile.LastName,
			"role":       string(targetUser.Role),
		},
	})

	return nil
}

// DeactivateUser deactivates a user (soft delete)
func (s *UserService) DeactivateUser(ctx context.Context, userID primitive.ObjectID, deactivatedBy *models.User, ipAddress string) error {
	// Check if deactivator has permission
	if !deactivatedBy.HasPermission(models.PermManageUsers) {
		return ErrUnauthorized
	}

	// Get target user
	targetUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if deactivator can manage this user
	if !deactivatedBy.CanManageUser(targetUser) {
		return ErrUnauthorized
	}

	// Prevent deactivating yourself
	if deactivatedBy.ID == targetUser.ID {
		return errors.New("cannot deactivate your own account")
	}

	// Deactivate user
	if err := s.userRepo.Deactivate(ctx, userID); err != nil {
		return err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &userID,
		PerformedBy: &deactivatedBy.ID,
		Action:      models.AuditActionUserDeactivated,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"username": targetUser.Username,
			"email":    targetUser.Email,
		},
	})

	return nil
}

// ActivateUser activates a user
func (s *UserService) ActivateUser(ctx context.Context, userID primitive.ObjectID, activatedBy *models.User, ipAddress string) error {
	// Check if activator has permission
	if !activatedBy.HasPermission(models.PermManageUsers) {
		return ErrUnauthorized
	}

	// Get target user
	targetUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if activator can manage this user
	if !activatedBy.CanManageUser(targetUser) {
		return ErrUnauthorized
	}

	// Activate user
	if err := s.userRepo.Activate(ctx, userID); err != nil {
		return err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &userID,
		PerformedBy: &activatedBy.ID,
		Action:      models.AuditActionUserActivated,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"username": targetUser.Username,
			"email":    targetUser.Email,
		},
	})

	return nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// ListUsers retrieves users with pagination, filtering, and searching
func (s *UserService) ListUsers(ctx context.Context, role *models.UserRole, isActive *bool, search string, limit, skip int64) ([]*models.User, int64, error) {
	filter := bson.M{}

	if role != nil {
		filter["role"] = *role
	}
	if isActive != nil {
		filter["is_active"] = *isActive
	}

	// Add search filter if search query is provided
	if search != "" {
		// Search across multiple fields: firstName, lastName, email, role
		// Note: Institution search requires a join/lookup with institutions collection
		filter["$or"] = []bson.M{
			{"profile.first_name": bson.M{"$regex": search, "$options": "i"}},
			{"profile.last_name": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
			{"role": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	users, err := s.userRepo.List(ctx, filter, limit, skip)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.userRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// CountUsers counts users with optional filtering
func (s *UserService) CountUsers(ctx context.Context, isActive *bool) (int64, error) {
	filter := bson.M{}
	if isActive != nil {
		filter["is_active"] = *isActive
	}
	return s.userRepo.Count(ctx, filter)
}

// CountUsersCreatedAfter counts users created after a specific time
func (s *UserService) CountUsersCreatedAfter(ctx context.Context, after time.Time) (int64, error) {
	filter := bson.M{
		"created_at": bson.M{"$gte": after},
	}
	return s.userRepo.Count(ctx, filter)
}

// GetRoleDistribution returns the count of users per role
func (s *UserService) GetRoleDistribution(ctx context.Context) (map[models.UserRole]int64, error) {
	roles := []models.UserRole{
		models.RoleHaematologist,
		models.RolePhysician,
		models.RoleDataCapturer,
		models.UserRole("admin"), // Admin role
	}

	distribution := make(map[models.UserRole]int64)
	for _, role := range roles {
		count, err := s.userRepo.Count(ctx, bson.M{"role": role})
		if err != nil {
			return nil, err
		}
		distribution[role] = count
	}

	return distribution, nil
}
