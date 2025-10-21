package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrUnauthorizedRegistryAccess = errors.New("unauthorized access to registry configuration")
	ErrFormNotReady               = errors.New("registry form is not ready for submissions")
	ErrConfigNotComplete          = errors.New("registry configuration is not complete")
	ErrNoActiveFormSchema         = errors.New("no active form schema available")
	ErrInvalidFormData            = errors.New("invalid form data for the schema")
)

// RegistryService handles business logic for the registry
type RegistryService struct {
	configRepo     *repository.RegistryConfigRepository
	formRepo       *repository.RegistryFormRepository
	submissionRepo *repository.RegistrySubmissionRepository
	userRepo       *repository.UserRepository
	auditRepo      *repository.AuditRepository
	dropboxService *DropboxService
	emailService   *EmailService
}

// NewRegistryService creates a new RegistryService
func NewRegistryService(
	configRepo *repository.RegistryConfigRepository,
	formRepo *repository.RegistryFormRepository,
	submissionRepo *repository.RegistrySubmissionRepository,
	userRepo *repository.UserRepository,
	auditRepo *repository.AuditRepository,
	dropboxService *DropboxService,
	emailService *EmailService,
) *RegistryService {
	return &RegistryService{
		configRepo:     configRepo,
		formRepo:       formRepo,
		submissionRepo: submissionRepo,
		userRepo:       userRepo,
		auditRepo:      auditRepo,
		dropboxService: dropboxService,
		emailService:   emailService,
	}
}

// Configuration Management

// GetConfiguration retrieves the registry configuration
func (s *RegistryService) GetConfiguration(ctx context.Context) (*models.RegistryConfig, error) {
	return s.configRepo.GetConfig(ctx)
}

// GetPublicConfiguration retrieves public configuration (video and documents only, no SMTP)
func (s *RegistryService) GetPublicConfiguration(ctx context.Context) (map[string]interface{}, error) {
	config, err := s.configRepo.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"videoUrl":      config.VideoURL,
		"documentsPath": config.DocumentsPath,
	}, nil
}

// UpdateConfiguration updates the registry configuration
func (s *RegistryService) UpdateConfiguration(
	ctx context.Context,
	req *models.UpdateRegistryConfigRequest,
	user *models.User,
	encryptionService *EncryptionService,
	ipAddress string,
) (*models.RegistryConfig, error) {
	// Check admin permission
	if !user.HasPermission(models.PermManageSystem) {
		return nil, ErrUnauthorizedRegistryAccess
	}

	// Get existing config or create new
	config, err := s.configRepo.GetConfig(ctx)
	if err != nil && err != repository.ErrRegistryConfigNotFound {
		return nil, err
	}

	// If no config exists, create a new one
	if err == repository.ErrRegistryConfigNotFound {
		config = &models.RegistryConfig{
			NotificationEmails: []string{},
		}
	}

	// Update fields
	if req.VideoURL != nil {
		config.VideoURL = *req.VideoURL
	}
	if req.DocumentsPath != nil {
		config.DocumentsPath = *req.DocumentsPath
	}
	if req.NotificationEmails != nil {
		config.NotificationEmails = *req.NotificationEmails
	}

	// Update SMTP config
	if req.SMTPHost != nil {
		config.SMTPConfig.Host = *req.SMTPHost
	}
	if req.SMTPPort != nil {
		config.SMTPConfig.Port = *req.SMTPPort
	}
	if req.SMTPUsername != nil {
		config.SMTPConfig.Username = *req.SMTPUsername
	}
	if req.SMTPPassword != nil && *req.SMTPPassword != "" {
		// Encrypt password before storing
		encryptedPassword, err := encryptionService.Encrypt(*req.SMTPPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt SMTP password: %w", err)
		}
		config.SMTPConfig.Password = encryptedPassword
	}
	if req.SMTPFromEmail != nil {
		config.SMTPConfig.FromEmail = *req.SMTPFromEmail
	}
	if req.SMTPFromName != nil {
		config.SMTPConfig.FromName = *req.SMTPFromName
	}

	config.UpdatedBy = &user.ID

	// Validate before saving
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Save configuration
	if err := s.configRepo.CreateOrUpdate(ctx, config); err != nil {
		return nil, err
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &user.ID,
		PerformedBy: &user.ID,
		Action:      "registry_config_updated",
		Details: map[string]interface{}{
			"updated_by": user.Email,
		},
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	})

	return config, nil
}

// Form Schema Management

// CreateFormSchema creates a new form schema
func (s *RegistryService) CreateFormSchema(
	ctx context.Context,
	req *models.CreateFormSchemaRequest,
	user *models.User,
	ipAddress string,
) (*models.RegistryFormSchema, error) {
	// Check admin permission
	if !user.HasPermission(models.PermManageUsers) {
		return nil, ErrUnauthorizedRegistryAccess
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Create schema
	schema := &models.RegistryFormSchema{
		FormName:    req.FormName,
		Description: req.Description,
		Fields:      req.Fields,
		IsActive:    false, // New forms are inactive by default
		CreatedBy:   &user.ID,
		UpdatedBy:   &user.ID,
	}

	// Validate schema
	if err := schema.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if err := s.formRepo.Create(ctx, schema); err != nil {
		return nil, err
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &user.ID,
		PerformedBy: &user.ID,
		Action:      "registry_form_created",
		Details: map[string]interface{}{
			"form_id":   schema.ID.Hex(),
			"form_name": schema.FormName,
		},
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	})

	return schema, nil
}

// UpdateFormSchema updates an existing form schema
func (s *RegistryService) UpdateFormSchema(
	ctx context.Context,
	id primitive.ObjectID,
	req *models.UpdateFormSchemaRequest,
	user *models.User,
	ipAddress string,
) (*models.RegistryFormSchema, error) {
	// Check admin permission
	if !user.HasPermission(models.PermManageUsers) {
		return nil, ErrUnauthorizedRegistryAccess
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get existing schema
	schema, err := s.formRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Build update document
	update := bson.M{
		"updated_by": user.ID,
	}

	if req.FormName != nil {
		update["form_name"] = *req.FormName
	}
	if req.Description != nil {
		update["description"] = *req.Description
	}
	if req.Fields != nil {
		update["fields"] = *req.Fields
	}
	if req.IsActive != nil {
		// If activating this form, deactivate all others
		if *req.IsActive {
			if err := s.formRepo.SetActive(ctx, id); err != nil {
				return nil, err
			}
		} else {
			update["is_active"] = false
		}
	}

	// Update
	if len(update) > 1 { // More than just updated_by
		if err := s.formRepo.Update(ctx, id, update); err != nil {
			return nil, err
		}
	}

	// Get updated schema
	schema, err = s.formRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &user.ID,
		PerformedBy: &user.ID,
		Action:      "registry_form_updated",
		Details: map[string]interface{}{
			"form_id":   schema.ID.Hex(),
			"form_name": schema.FormName,
		},
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	})

	return schema, nil
}

// GetFormSchema retrieves a form schema by ID
func (s *RegistryService) GetFormSchema(ctx context.Context, id primitive.ObjectID) (*models.RegistryFormSchema, error) {
	return s.formRepo.FindByID(ctx, id)
}

// GetActiveFormSchema retrieves the currently active form schema
func (s *RegistryService) GetActiveFormSchema(ctx context.Context) (*models.RegistryFormSchema, error) {
	return s.formRepo.FindActive(ctx)
}

// ListFormSchemas retrieves all form schemas with pagination
func (s *RegistryService) ListFormSchemas(ctx context.Context, page, limit int) ([]*models.RegistryFormSchema, int64, error) {
	return s.formRepo.List(ctx, page, limit)
}

// DeleteFormSchema deletes a form schema
func (s *RegistryService) DeleteFormSchema(
	ctx context.Context,
	id primitive.ObjectID,
	user *models.User,
	ipAddress string,
) error {
	// Check admin permission
	if !user.HasPermission(models.PermManageUsers) {
		return ErrUnauthorizedRegistryAccess
	}

	// Get schema for audit
	schema, err := s.formRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete
	if err := s.formRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &user.ID,
		PerformedBy: &user.ID,
		Action:      "registry_form_deleted",
		Details: map[string]interface{}{
			"form_id":   schema.ID.Hex(),
			"form_name": schema.FormName,
		},
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	})

	return nil
}

// Submission Management

// ValidatePreSubmission checks if the system is ready to accept submissions
func (s *RegistryService) ValidatePreSubmission(ctx context.Context) error {
	// Check if there's an active form
	_, err := s.formRepo.FindActive(ctx)
	if err != nil {
		if err == repository.ErrNoActiveForm {
			return ErrNoActiveFormSchema
		}
		return err
	}

	// Check if configuration is complete
	config, err := s.configRepo.GetConfig(ctx)
	if err != nil {
		return ErrConfigNotComplete
	}

	// Check notification emails
	if len(config.NotificationEmails) == 0 {
		return errors.New("no notification emails configured")
	}

	// Check SMTP configuration
	if !config.SMTPConfig.IsComplete() {
		return errors.New("SMTP configuration is incomplete")
	}

	return nil
}

// SubmitForm handles form submission with document uploads
func (s *RegistryService) SubmitForm(
	ctx context.Context,
	req *models.CreateSubmissionRequest,
	files []*multipart.FileHeader,
	user *models.User,
	ipAddress string,
) (*models.RegistrySubmission, error) {
	// Pre-submission validation
	if err := s.ValidatePreSubmission(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormNotReady, err)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get form schema
	formSchemaID, _ := primitive.ObjectIDFromHex(req.FormSchemaID)
	schema, err := s.formRepo.FindByID(ctx, formSchemaID)
	if err != nil {
		return nil, err
	}

	// Validate form data against schema
	if err := s.validateFormData(req.FormData, schema); err != nil {
		return nil, err
	}

	// Validate files (at least one required)
	if len(files) == 0 {
		return nil, models.ErrNoDocumentsUploaded
	}

	// Create submission record
	submission := &models.RegistrySubmission{
		UserID:       user.ID,
		FormSchemaID: formSchemaID,
		FormData:     req.FormData,
		Status:       models.SubmissionStatusSubmitted,
	}

	// Create submission to get ID
	if err := s.submissionRepo.Create(ctx, submission); err != nil {
		return nil, err
	}

	// Upload documents to Dropbox
	dropboxPath := fmt.Sprintf("BLDS_approvals/%s/%s", user.Username, submission.ID.Hex())
	submission.DocumentsPath = dropboxPath

	uploadedFiles := []string{}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		// Upload to Dropbox
		remotePath := filepath.Join(dropboxPath, fileHeader.Filename)
		if err := s.dropboxService.UploadFile(ctx, file, remotePath); err != nil {
			return nil, fmt.Errorf("failed to upload file to Dropbox: %w", err)
		}

		uploadedFiles = append(uploadedFiles, fileHeader.Filename)
	}

	submission.UploadedDocuments = uploadedFiles

	// Update submission with documents info
	if err := s.submissionRepo.Update(ctx, submission.ID, bson.M{
		"documents_path":     submission.DocumentsPath,
		"uploaded_documents": submission.UploadedDocuments,
	}); err != nil {
		return nil, err
	}

	// Send email notifications
	if err := s.sendSubmissionNotification(ctx, submission, user, schema); err != nil {
		// Log error but don't fail the submission
		fmt.Printf("Warning: Failed to send notification email: %v\n", err)
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &user.ID,
		PerformedBy: &user.ID,
		Action:      "registry_submission_created",
		Details: map[string]interface{}{
			"submission_id": submission.ID.Hex(),
			"form_id":       schema.ID.Hex(),
			"files_count":   len(uploadedFiles),
		},
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	})

	return submission, nil
}

// validateFormData validates submitted data against the form schema
func (s *RegistryService) validateFormData(data map[string]interface{}, schema *models.RegistryFormSchema) error {
	// Check required fields
	for _, field := range schema.Fields {
		if field.Required && field.Type != models.FieldTypeFile {
			value, exists := data[field.ID]
			if !exists || value == nil || value == "" {
				return fmt.Errorf("required field '%s' is missing", field.Label)
			}
		}
	}

	// Additional validation can be added here based on field types and validation rules
	return nil
}

// sendSubmissionNotification sends email notification for new submission
func (s *RegistryService) sendSubmissionNotification(
	ctx context.Context,
	submission *models.RegistrySubmission,
	user *models.User,
	schema *models.RegistryFormSchema,
) error {
	// Get config for SMTP and notification emails
	config, err := s.configRepo.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Get Dropbox shared link for the folder
	dropboxLink := fmt.Sprintf("https://www.dropbox.com/home/%s", submission.DocumentsPath)

	// Prepare email data
	emailData := SubmissionNotificationData{
		UserName:     user.Profile.FirstName + " " + user.Profile.LastName,
		UserEmail:    user.Email,
		SubmissionID: submission.ID.Hex(),
		DropboxLink:  dropboxLink,
		FormName:     schema.FormName,
		SubmittedAt:  submission.CreatedAt.Format("January 2, 2006 at 3:04 PM"),
	}

	// Send email
	return s.emailService.SendSubmissionNotification(
		config.SMTPConfig,
		config.NotificationEmails,
		emailData,
	)
}

// GetUserSubmissions retrieves submissions for a specific user
func (s *RegistryService) GetUserSubmissions(
	ctx context.Context,
	userID primitive.ObjectID,
	page, limit int,
) ([]*models.RegistrySubmission, int64, error) {
	return s.submissionRepo.FindByUser(ctx, userID, page, limit)
}

// GetAllSubmissions retrieves all submissions (admin only)
func (s *RegistryService) GetAllSubmissions(
	ctx context.Context,
	user *models.User,
	page, limit int,
	filter bson.M,
) ([]*models.RegistrySubmission, int64, error) {
	// Check admin permission
	if !user.HasPermission(models.PermManageUsers) {
		return nil, 0, ErrUnauthorizedRegistryAccess
	}

	return s.submissionRepo.List(ctx, page, limit, filter)
}

// GetSubmission retrieves a specific submission
func (s *RegistryService) GetSubmission(
	ctx context.Context,
	id primitive.ObjectID,
	user *models.User,
) (*models.RegistrySubmission, error) {
	submission, err := s.submissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user has permission (owner or admin)
	if submission.UserID != user.ID && !user.HasPermission(models.PermManageUsers) {
		return nil, errors.New("unauthorized to view this submission")
	}

	return submission, nil
}

// UpdateSubmissionStatus updates the status of a submission (admin only)
func (s *RegistryService) UpdateSubmissionStatus(
	ctx context.Context,
	id primitive.ObjectID,
	req *models.UpdateSubmissionStatusRequest,
	user *models.User,
	ipAddress string,
) (*models.RegistrySubmission, error) {
	// Check admin permission
	if !user.HasPermission(models.PermManageUsers) {
		return nil, ErrUnauthorizedRegistryAccess
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Update status
	if err := s.submissionRepo.UpdateStatus(ctx, id, req.Status, &user.ID, req.ReviewNotes); err != nil {
		return nil, err
	}

	// Get updated submission
	submission, err := s.submissionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &user.ID,
		PerformedBy: &user.ID,
		Action:      "registry_submission_status_updated",
		Details: map[string]interface{}{
			"submission_id": submission.ID.Hex(),
			"new_status":    req.Status,
		},
		IPAddress: ipAddress,
		Timestamp: time.Now(),
	})

	return submission, nil
}
