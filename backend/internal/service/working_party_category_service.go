package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"backend/internal/models"
	"backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const workingPartiesDropboxRoot = "WORKING_PARTIES"

var (
	ErrWorkingPartyCategoryNotFound = errors.New("working party category not found")
	ErrDuplicateWorkingPartySlug    = errors.New("working party category with this name already exists")
)

// WorkingPartyCategoryService handles business logic for working party categories
type WorkingPartyCategoryService struct {
	categoryRepo   *repository.WorkingPartyCategoryRepository
	dropboxService *DropboxService
	auditRepo      *repository.AuditRepository
	userRepo       *repository.UserRepository
}

// NewWorkingPartyCategoryService creates a new WorkingPartyCategoryService
func NewWorkingPartyCategoryService(
	categoryRepo *repository.WorkingPartyCategoryRepository,
	dropboxService *DropboxService,
	auditRepo *repository.AuditRepository,
	userRepo *repository.UserRepository,
) *WorkingPartyCategoryService {
	return &WorkingPartyCategoryService{
		categoryRepo:   categoryRepo,
		dropboxService: dropboxService,
		auditRepo:      auditRepo,
		userRepo:       userRepo,
	}
}

// CreateCategory creates a new working party category
func (s *WorkingPartyCategoryService) CreateCategory(
	ctx context.Context,
	req *models.CreateWorkingPartyCategoryRequest,
	createdBy *models.User,
	ipAddress string,
) (*models.WorkingPartyCategory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !createdBy.HasPermission(models.PermDeleteUsers) {
		return nil, ErrUnauthorized
	}

	slug := models.GenerateSlug(req.Name)

	exists, err := s.categoryRepo.ExistsBySlug(ctx, slug, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug existence: %w", err)
	}
	if exists {
		return nil, ErrDuplicateWorkingPartySlug
	}

	category := &models.WorkingPartyCategory{
		Name:         req.Name,
		Slug:         slug,
		Description:  req.Description,
		ImagePath:    req.ImagePath,
		DropboxPath:  workingPartiesDropboxRoot + "/" + url.PathEscape(req.Name),
		DisplayOrder: req.DisplayOrder,
		IsActive:     true,
		CreatedBy:    &createdBy.ID,
	}

	if err := category.Validate(); err != nil {
		return nil, err
	}

	if s.dropboxService.IsConfigured() {
		if err := s.dropboxService.CreateFolder(workingPartiesDropboxRoot); err != nil {
			fmt.Printf("ERROR: Failed to create Dropbox parent folder '%s': %v\n", workingPartiesDropboxRoot, err)
		}
		fmt.Printf("Attempting to create Dropbox folder at path: %s\n", category.DropboxPath)
		err := s.dropboxService.CreateFolder(category.DropboxPath)
		if err != nil {
			fmt.Printf("ERROR: Failed to create Dropbox folder for category '%s' at path '%s': %v\n", category.Name, category.DropboxPath, err)
		} else {
			fmt.Printf("SUCCESS: Created Dropbox folder for category '%s' at path '%s'\n", category.Name, category.DropboxPath)
		}
	} else {
		fmt.Printf("WARNING: Dropbox is not configured - skipping folder creation for category '%s'\n", category.Name)
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		if err == repository.ErrDuplicateSlug {
			return nil, ErrDuplicateWorkingPartySlug
		}
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &createdBy.ID,
		PerformedBy: &createdBy.ID,
		Action:      "working_party_category.create",
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
func (s *WorkingPartyCategoryService) GetCategory(ctx context.Context, id primitive.ObjectID, user *models.User) (*models.WorkingPartyCategory, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			return nil, ErrWorkingPartyCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if !category.IsActive && !user.HasPermission(models.PermDeleteUsers) {
		return nil, ErrWorkingPartyCategoryNotFound
	}

	return category, nil
}

// GetCategoryBySlug retrieves a category by slug
func (s *WorkingPartyCategoryService) GetCategoryBySlug(ctx context.Context, slug string, user *models.User) (*models.WorkingPartyCategory, error) {
	category, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			return nil, ErrWorkingPartyCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if !category.IsActive && !user.HasPermission(models.PermDeleteUsers) {
		return nil, ErrWorkingPartyCategoryNotFound
	}

	return category, nil
}

// ListCategories lists categories with filters
func (s *WorkingPartyCategoryService) ListCategories(
	ctx context.Context,
	user *models.User,
	search string,
	page, limit int,
) ([]*models.WorkingPartyCategory, int64, error) {
	filter := repository.WorkingPartyCategoryFilter{
		Search: search,
	}

	if !user.HasPermission(models.PermDeleteUsers) {
		isActive := true
		filter.IsActive = &isActive
	}

	categories, err := s.categoryRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list categories: %w", err)
	}

	total, err := s.categoryRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return categories, total, nil
}

// UpdateCategory updates a category
func (s *WorkingPartyCategoryService) UpdateCategory(
	ctx context.Context,
	id primitive.ObjectID,
	req *models.UpdateWorkingPartyCategoryRequest,
	updatedBy *models.User,
	ipAddress string,
) (*models.WorkingPartyCategory, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	if !updatedBy.HasPermission(models.PermDeleteUsers) {
		return nil, ErrUnauthorized
	}

	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			return nil, ErrWorkingPartyCategoryNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	update := bson.M{}
	oldDropboxPath := category.DropboxPath
	nameChanged := false

	if req.Name != nil && *req.Name != category.Name {
		update["name"] = *req.Name
		slug := models.GenerateSlug(*req.Name)

		exists, err := s.categoryRepo.ExistsBySlug(ctx, slug, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check slug existence: %w", err)
		}
		if exists {
			return nil, ErrDuplicateWorkingPartySlug
		}

		update["slug"] = slug
		update["dropbox_path"] = workingPartiesDropboxRoot + "/" + url.PathEscape(*req.Name)
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

	if len(update) == 0 {
		return category, nil
	}

	if err := s.categoryRepo.Update(ctx, id, update); err != nil {
		if err == repository.ErrDuplicateSlug {
			return nil, ErrDuplicateWorkingPartySlug
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	if nameChanged && s.dropboxService.IsConfigured() {
		newDropboxPath := update["dropbox_path"].(string)
		fmt.Printf("Attempting to rename Dropbox folder from '%s' to '%s'\n", oldDropboxPath, newDropboxPath)
		err := s.dropboxService.RenameFolder(oldDropboxPath, newDropboxPath)
		if err != nil {
			fmt.Printf("ERROR: Failed to rename Dropbox folder from '%s' to '%s': %v\n", oldDropboxPath, newDropboxPath, err)
			return nil, fmt.Errorf("category updated in database but failed to rename Dropbox folder: %w", err)
		}
		fmt.Printf("SUCCESS: Renamed Dropbox folder from '%s' to '%s'\n", oldDropboxPath, newDropboxPath)
	}

	updatedCategory, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated category: %w", err)
	}

	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &updatedBy.ID,
		PerformedBy: &updatedBy.ID,
		Action:      "working_party_category.update",
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
func (s *WorkingPartyCategoryService) DeleteCategory(
	ctx context.Context,
	id primitive.ObjectID,
	deletedBy *models.User,
	ipAddress string,
) error {
	if !deletedBy.HasPermission(models.PermDeleteUsers) {
		return ErrUnauthorized
	}

	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		if err == repository.ErrCategoryNotFound {
			return ErrWorkingPartyCategoryNotFound
		}
		return fmt.Errorf("failed to get category: %w", err)
	}

	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		if err == repository.ErrCategoryNotFound {
			return ErrWorkingPartyCategoryNotFound
		}
		return fmt.Errorf("failed to delete category: %w", err)
	}

	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &deletedBy.ID,
		PerformedBy: &deletedBy.ID,
		Action:      "working_party_category.delete",
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
func (s *WorkingPartyCategoryService) GetCategoryFiles(
	ctx context.Context,
	id primitive.ObjectID,
	user *models.User,
) ([]DropboxFileInfo, error) {
	category, err := s.GetCategory(ctx, id, user)
	if err != nil {
		return nil, err
	}

	if !s.dropboxService.IsConfigured() {
		return nil, errors.New("dropbox is not configured")
	}

	files, err := s.dropboxService.ListFilesRecursive(category.DropboxPath)
	if err != nil {
		if err == ErrFolderNotFound {
			return []DropboxFileInfo{}, nil
		}
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	categoryPath := strings.TrimPrefix(category.DropboxPath, "/")
	categoryPath = strings.TrimSuffix(categoryPath, "/")

	s.makePathsRelative(files, categoryPath)

	return files, nil
}

// GetFileDownloadLink generates a download link for a specific file
func (s *WorkingPartyCategoryService) GetFileDownloadLink(
	ctx context.Context,
	categoryID primitive.ObjectID,
	filePath string,
	user *models.User,
) (string, error) {
	category, err := s.GetCategory(ctx, categoryID, user)
	if err != nil {
		return "", err
	}

	if !s.dropboxService.IsConfigured() {
		return "", errors.New("dropbox is not configured")
	}

	dropboxPath := category.DropboxPath
	if decoded, err := url.PathUnescape(dropboxPath); err == nil {
		dropboxPath = decoded
	}

	dropboxPath = strings.TrimPrefix(dropboxPath, "/")
	dropboxPath = strings.TrimSuffix(dropboxPath, "/")
	filePath = strings.TrimPrefix(filePath, "/")

	fullPath := "/" + filepath.Join(dropboxPath, filePath)
	fullPath = strings.ReplaceAll(fullPath, "\\", "/")
	fullPath = strings.ReplaceAll(fullPath, "//", "/")

	link, err := s.dropboxService.GetFileDownloadLink(fullPath)
	if err != nil {
		if err == ErrFileNotFound {
			return "", errors.New("file not found")
		}
		return "", fmt.Errorf("failed to get download link: %w", err)
	}

	return link, nil
}

func (s *WorkingPartyCategoryService) makePathsRelative(files []DropboxFileInfo, categoryPath string) {
	for i := range files {
		fullPath := strings.TrimPrefix(files[i].Path, "/")
		if strings.HasPrefix(fullPath, categoryPath+"/") {
			files[i].Path = strings.TrimPrefix(fullPath, categoryPath+"/")
		} else if fullPath == categoryPath {
			files[i].Path = ""
		}

		if files[i].IsFolder && len(files[i].Children) > 0 {
			s.makePathsRelative(files[i].Children, categoryPath)
		}
	}
}

// CountCategories returns the total count of working party categories
func (s *WorkingPartyCategoryService) CountCategories(ctx context.Context, activeOnly bool) (int64, error) {
	filter := repository.WorkingPartyCategoryFilter{}

	if activeOnly {
		isActive := true
		filter.IsActive = &isActive
	}

	count, err := s.categoryRepo.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return count, nil
}
