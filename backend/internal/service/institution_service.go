package service

import (
	"context"
	"errors"

	"backend/internal/models"
	"backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInstitutionHasUsers = errors.New("cannot delete institution: users are still associated with it")
)

// InstitutionService handles business logic for institutions
type InstitutionService struct {
	institutionRepo *repository.InstitutionRepository
	userRepo        *repository.UserRepository
	auditRepo       *repository.AuditRepository
}

// NewInstitutionService creates a new InstitutionService
func NewInstitutionService(institutionRepo *repository.InstitutionRepository, userRepo *repository.UserRepository, auditRepo *repository.AuditRepository) *InstitutionService {
	return &InstitutionService{
		institutionRepo: institutionRepo,
		userRepo:        userRepo,
		auditRepo:       auditRepo,
	}
}

// CreateInstitution creates a new institution
func (s *InstitutionService) CreateInstitution(ctx context.Context, req *models.CreateInstitutionRequest, createdBy *models.User, ipAddress string) (*models.Institution, error) {
	// Check if user has permission
	if !createdBy.HasPermission(models.PermManageUsers) {
		return nil, ErrUnauthorized
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	institution := &models.Institution{
		Name:       req.Name,
		ShortName:  req.ShortName,
		Type:       req.Type,
		Country:    req.Country,
		Province:   req.Province,
		City:       req.City,
		Address:    req.Address,
		PostalCode: req.PostalCode,
		Phone:      req.Phone,
		Email:      req.Email,
		Website:    req.Website,
		IsActive:   true,
		CreatedBy:  &createdBy.ID,
	}

	if err := s.institutionRepo.Create(ctx, institution); err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		PerformedBy: &createdBy.ID,
		Action:      models.AuditActionInstitutionCreated,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"institution_id":   institution.ID.Hex(),
			"institution_name": institution.Name,
			"type":             string(institution.Type),
			"city":             institution.City,
		},
	})

	return institution, nil
}

// GetInstitution retrieves an institution by ID
func (s *InstitutionService) GetInstitution(ctx context.Context, id primitive.ObjectID) (*models.Institution, error) {
	return s.institutionRepo.FindByID(ctx, id)
}

// UpdateInstitution updates an institution
func (s *InstitutionService) UpdateInstitution(ctx context.Context, id primitive.ObjectID, req *models.UpdateInstitutionRequest, updatedBy *models.User, ipAddress string) (*models.Institution, error) {
	// Check if user has permission
	if !updatedBy.HasPermission(models.PermManageUsers) {
		return nil, ErrUnauthorized
	}

	// Get existing institution
	institution, err := s.institutionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Build update document
	update := bson.M{}
	if req.Name != nil {
		// Check for duplicate name
		if *req.Name != institution.Name {
			exists, err := s.institutionRepo.NameExists(ctx, *req.Name, &id)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, repository.ErrDuplicateInstitution
			}
			update["name"] = *req.Name
		}
	}
	if req.ShortName != nil {
		update["short_name"] = *req.ShortName
	}
	if req.Type != nil {
		if !req.Type.IsValid() {
			return nil, models.ErrInvalidInstitutionType
		}
		update["type"] = *req.Type
	}
	if req.Country != nil {
		update["country"] = *req.Country
	}
	if req.Province != nil {
		update["province"] = *req.Province
	}
	if req.City != nil {
		update["city"] = *req.City
	}
	if req.Address != nil {
		update["address"] = *req.Address
	}
	if req.PostalCode != nil {
		update["postal_code"] = *req.PostalCode
	}
	if req.Phone != nil {
		update["phone"] = *req.Phone
	}
	if req.Email != nil {
		update["email"] = *req.Email
	}
	if req.Website != nil {
		update["website"] = *req.Website
	}
	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}

	if len(update) == 0 {
		return institution, nil
	}

	if err := s.institutionRepo.Update(ctx, id, update); err != nil {
		return nil, err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		PerformedBy: &updatedBy.ID,
		Action:      models.AuditActionInstitutionUpdated,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"institution_id":   id.Hex(),
			"institution_name": institution.Name,
			"updated_fields":   update,
		},
	})

	return s.institutionRepo.FindByID(ctx, id)
}

// DeleteInstitution deletes an institution
func (s *InstitutionService) DeleteInstitution(ctx context.Context, id primitive.ObjectID, deletedBy *models.User, ipAddress string) error {
	// Check if user has permission
	if !deletedBy.HasPermission(models.PermDeleteUsers) {
		return ErrUnauthorized
	}

	// Get institution for logging
	institution, err := s.institutionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if any users are still associated with this institution
	userCount, err := s.userRepo.Count(ctx, bson.M{"profile.institution_id": id})
	if err != nil {
		return err
	}
	if userCount > 0 {
		return ErrInstitutionHasUsers
	}

	if err := s.institutionRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		PerformedBy: &deletedBy.ID,
		Action:      models.AuditActionInstitutionDeleted,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"institution_id":   id.Hex(),
			"institution_name": institution.Name,
			"type":             string(institution.Type),
			"city":             institution.City,
		},
	})

	return nil
}

// ActivateInstitution activates an institution
func (s *InstitutionService) ActivateInstitution(ctx context.Context, id primitive.ObjectID, activatedBy *models.User, ipAddress string) error {
	// Check if user has permission
	if !activatedBy.HasPermission(models.PermManageUsers) {
		return ErrUnauthorized
	}

	institution, err := s.institutionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.institutionRepo.Activate(ctx, id); err != nil {
		return err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		PerformedBy: &activatedBy.ID,
		Action:      models.AuditActionInstitutionActivated,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"institution_id":   id.Hex(),
			"institution_name": institution.Name,
		},
	})

	return nil
}

// DeactivateInstitution deactivates an institution
func (s *InstitutionService) DeactivateInstitution(ctx context.Context, id primitive.ObjectID, deactivatedBy *models.User, ipAddress string) error {
	// Check if user has permission
	if !deactivatedBy.HasPermission(models.PermManageUsers) {
		return ErrUnauthorized
	}

	institution, err := s.institutionRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.institutionRepo.Deactivate(ctx, id); err != nil {
		return err
	}

	// Log audit
	s.auditRepo.Create(ctx, &models.AuditLog{
		PerformedBy: &deactivatedBy.ID,
		Action:      models.AuditActionInstitutionDeactivated,
		IPAddress:   ipAddress,
		Details: map[string]interface{}{
			"institution_id":   id.Hex(),
			"institution_name": institution.Name,
		},
	})

	return nil
}

// ListInstitutions retrieves institutions with pagination and filtering
func (s *InstitutionService) ListInstitutions(ctx context.Context, institutionType *models.InstitutionType, isActive *bool, search string, limit, skip int64) ([]*models.Institution, int64, error) {
	filter := bson.M{}

	if institutionType != nil {
		filter["type"] = *institutionType
	}
	if isActive != nil {
		filter["is_active"] = *isActive
	}

	// Add search filter if search query is provided
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"short_name": bson.M{"$regex": search, "$options": "i"}},
			{"city": bson.M{"$regex": search, "$options": "i"}},
			{"province": bson.M{"$regex": search, "$options": "i"}},
			{"country": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	institutions, err := s.institutionRepo.List(ctx, filter, limit, skip)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.institutionRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return institutions, count, nil
}

// ValidateInstitutionID checks if an institution ID exists and is active
func (s *InstitutionService) ValidateInstitutionID(ctx context.Context, id primitive.ObjectID) error {
	institution, err := s.institutionRepo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrInstitutionNotFound {
			return errors.New("institution not found")
		}
		return err
	}

	if !institution.IsActive {
		return errors.New("institution is not active")
	}

	return nil
}

// CountInstitutions counts institutions with optional filtering
func (s *InstitutionService) CountInstitutions(ctx context.Context, isActive *bool) (int64, error) {
	filter := bson.M{}
	if isActive != nil {
		filter["is_active"] = *isActive
	}
	return s.institutionRepo.Count(ctx, filter)
}

