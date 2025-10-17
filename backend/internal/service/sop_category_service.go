package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"backend/internal/models"
	"backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrDuplicateSlug    = errors.New("category with this name already exists")
)

// SOPCategoryService handles business logic for SOP categories
type SOPCategoryService struct {
	categoryRepo   *repository.SOPCategoryRepository
	dropboxService *DropboxService
	auditRepo      *repository.AuditRepository
	userRepo       *repository.UserRepository
}

// NewSOPCategoryService creates a new SOPCategoryService
func NewSOPCategoryService(
	categoryRepo *repository.SOPCategoryRepository,
	dropboxService *DropboxService,
	auditRepo *repository.AuditRepository,
	userRepo *repository.UserRepository,
) *SOPCategoryService {
	return &SOPCategoryService{
		categoryRepo:   categoryRepo,
		dropboxService: dropboxService,
		auditRepo:      auditRepo,
		userRepo:       userRepo,
	}
}

// CreateCategory creates a new SOP category
func (s *SOPCategoryService) CreateCategory(
	ctx context.Context,
	req *models.CreateSOPCategoryRequest,
	createdBy *models.User,
	ipAddress string,
) (*models.SOPCategory, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check permission - only super admins can manage categories
	if !createdBy.HasPermission(models.PermDeleteUsers) {
		return nil, ErrUnauthorized
	}

	// Generate slug from name
	slug := models.GenerateSlug(req.Name)

	// Check if slug already exists
	exists, err := s.categoryRepo.ExistsBySlug(ctx, slug, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug existence: %w", err)
	}
	if exists {
		return nil, ErrDuplicateSlug
	}

	// Create category object
	category := &models.SOPCategory{
		Name:         req.Name,
		Slug:         slug,
		Description:  req.Description,
		ImagePath:    req.ImagePath,
		DropboxPath:  "SOPS/" + req.Name,
		DisplayOrder: req.DisplayOrder,
		IsActive:     true,
		CreatedBy:    &createdBy.ID,
	}

	// Validate category
	if err := category.Validate(); err != nil {
		return nil, err
	}

	// Create folder in Dropbox
	if s.dropboxService.IsConfigured() {
		err := s.dropboxService.CreateFolder(category.DropboxPath)
		if err != nil {
			// Log the error but don't fail the category creation
			// In production, you'd want proper logging here
			fmt.Printf("Warning: Failed to create Dropbox folder for category %s: %v\n", category.Name, err)
		}
	}

	// Create category in database
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		if err == repository.ErrDuplicateSlug {
			return nil, ErrDuplicateSlug
		}
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &createdBy.ID,
		PerformedBy: &createdBy.ID,
		Action:      "sop_category.create",
		Details: bson.M{
			"category_id":   category.ID.Hex(),
			"category_name": category.Name,
			"slug":          category.Slug,
		},
		IPAddress: ipAddress,
	})

	return category, nil
}

// GetCategory retrieves a category by ID
func (s *SOPCategoryService) GetCategory(ctx context.Context, id primitive.ObjectID, user *models.User) (*models.SOPCategory, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Non-admins can only see active categories
	if !category.IsActive && !user.HasPermission(models.PermDeleteUsers) {
		return nil, ErrCategoryNotFound
	}

	return category, nil
}

// GetCategoryBySlug retrieves a category by slug
func (s *SOPCategoryService) GetCategoryBySlug(ctx context.Context, slug string, user *models.User) (*models.SOPCategory, error) {
	category, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Non-admins can only see active categories
	if !category.IsActive && !user.HasPermission(models.PermDeleteUsers) {
		return nil, ErrCategoryNotFound
	}

	return category, nil
}

// ListCategories lists categories with filters
func (s *SOPCategoryService) ListCategories(
	ctx context.Context,
	user *models.User,
	search string,
	page, limit int,
) ([]*models.SOPCategory, int64, error) {
	// Build filter
	filter := repository.SOPCategoryFilter{
		Search: search,
	}

	// Non-admins can only see active categories
	if !user.HasPermission(models.PermDeleteUsers) {
		isActive := true
		filter.IsActive = &isActive
	}

	// Get categories
	categories, err := s.categoryRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list categories: %w", err)
	}

	// Get total count
	total, err := s.categoryRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return categories, total, nil
}

// UpdateCategory updates a category
func (s *SOPCategoryService) UpdateCategory(
	ctx context.Context,
	id primitive.ObjectID,
	req *models.UpdateSOPCategoryRequest,
	updatedBy *models.User,
	ipAddress string,
) (*models.SOPCategory, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check permission
	if !updatedBy.HasPermission(models.PermDeleteUsers) {
		return nil, ErrUnauthorized
	}

	// Get existing category
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Build update map
	update := bson.M{}
	oldDropboxPath := category.DropboxPath
	nameChanged := false

	if req.Name != nil && *req.Name != category.Name {
		update["name"] = *req.Name
		slug := models.GenerateSlug(*req.Name)

		// Check if new slug already exists
		exists, err := s.categoryRepo.ExistsBySlug(ctx, slug, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check slug existence: %w", err)
		}
		if exists {
			return nil, ErrDuplicateSlug
		}

		update["slug"] = slug
		update["dropbox_path"] = "SOPS/" + *req.Name
		nameChanged = true
	}

	if req.Description != nil {
		update["description"] = *req.Description
	}

	if req.ImagePath != nil {
		update["image_path"] = *req.ImagePath
	}

	if req.DisplayOrder != nil {
		update["display_order"] = *req.DisplayOrder
	}

	if req.IsActive != nil {
		update["is_active"] = *req.IsActive
	}

	// If nothing to update, return existing category
	if len(update) == 0 {
		return category, nil
	}

	// Update in database
	if err := s.categoryRepo.Update(ctx, id, update); err != nil {
		if err == repository.ErrDuplicateSlug {
			return nil, ErrDuplicateSlug
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// If name changed, rename Dropbox folder
	if nameChanged && s.dropboxService.IsConfigured() {
		newDropboxPath := update["dropbox_path"].(string)
		err := s.dropboxService.RenameFolder(oldDropboxPath, newDropboxPath)
		if err != nil {
			// Log the error but don't fail the update
			fmt.Printf("Warning: Failed to rename Dropbox folder from %s to %s: %v\n", oldDropboxPath, newDropboxPath, err)
		}
	}

	// Get updated category
	updatedCategory, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated category: %w", err)
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &updatedBy.ID,
		PerformedBy: &updatedBy.ID,
		Action:      "sop_category.update",
		Details: bson.M{
			"category_id":   id.Hex(),
			"category_name": updatedCategory.Name,
			"changes":       update,
		},
		IPAddress: ipAddress,
	})

	return updatedCategory, nil
}

// DeleteCategory deletes a category from the database (Dropbox folder remains)
func (s *SOPCategoryService) DeleteCategory(
	ctx context.Context,
	id primitive.ObjectID,
	deletedBy *models.User,
	ipAddress string,
) error {
	// Check permission
	if !deletedBy.HasPermission(models.PermDeleteUsers) {
		return ErrUnauthorized
	}

	// Get category for audit log
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			return ErrCategoryNotFound
		}
		return fmt.Errorf("failed to get category: %w", err)
	}

	// Delete from database (NOT from Dropbox)
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		if err == repository.ErrCategoryNotFound {
			return ErrCategoryNotFound
		}
		return fmt.Errorf("failed to delete category: %w", err)
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &deletedBy.ID,
		PerformedBy: &deletedBy.ID,
		Action:      "sop_category.delete",
		Details: bson.M{
			"category_id":   id.Hex(),
			"category_name": category.Name,
			"slug":          category.Slug,
			"dropbox_path":  category.DropboxPath,
		},
		IPAddress: ipAddress,
	})

	return nil
}

// GetCategoryFiles lists all files in a category's Dropbox folder
func (s *SOPCategoryService) GetCategoryFiles(
	ctx context.Context,
	id primitive.ObjectID,
	user *models.User,
) ([]DropboxFileInfo, error) {
	// Get category
	category, err := s.GetCategory(ctx, id, user)
	if err != nil {
		return nil, err
	}

	// Check if Dropbox is configured
	if !s.dropboxService.IsConfigured() {
		return nil, errors.New("dropbox is not configured")
	}

	// List files from Dropbox
	files, err := s.dropboxService.ListFiles(category.DropboxPath)
	if err != nil {
		if err == ErrFolderNotFound {
			// Return empty list if folder doesn't exist yet
			return []DropboxFileInfo{}, nil
		}
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// GetFileDownloadLink generates a download link for a specific file
func (s *SOPCategoryService) GetFileDownloadLink(
	ctx context.Context,
	categoryID primitive.ObjectID,
	filePath string,
	user *models.User,
) (string, error) {
	// Get category to verify it exists and user has access
	category, err := s.GetCategory(ctx, categoryID, user)
	if err != nil {
		return "", err
	}

	// Check if Dropbox is configured
	if !s.dropboxService.IsConfigured() {
		return "", errors.New("dropbox is not configured")
	}

	// Construct full path
	fullPath := filepath.Join(category.DropboxPath, filePath)

	// Get download link
	link, err := s.dropboxService.GetFileDownloadLink(fullPath)
	if err != nil {
		if err == ErrFileNotFound {
			return "", errors.New("file not found")
		}
		return "", fmt.Errorf("failed to get download link: %w", err)
	}

	return link, nil
}
