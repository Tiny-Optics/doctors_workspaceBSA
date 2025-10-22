package service

import (
	"context"
	"errors"

	"backend/internal/models"
	"backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrReferralNotConfigured = errors.New("referral system is not configured")
	ErrReferralDisabled      = errors.New("referral system is currently disabled")
)

// ReferralService handles REDCap referral operations
type ReferralService struct {
	referralRepo *repository.ReferralConfigRepository
	auditRepo    *repository.AuditRepository
}

// NewReferralService creates a new ReferralService
func NewReferralService(
	referralRepo *repository.ReferralConfigRepository,
	auditRepo *repository.AuditRepository,
) *ReferralService {
	return &ReferralService{
		referralRepo: referralRepo,
		auditRepo:    auditRepo,
	}
}

// GetReferralConfig retrieves the referral configuration
func (s *ReferralService) GetReferralConfig(ctx context.Context) (*models.ReferralConfig, error) {
	config, err := s.referralRepo.GetConfig(ctx)
	if err != nil {
		if err == repository.ErrReferralConfigNotFound {
			return nil, ErrReferralNotConfigured
		}
		return nil, err
	}

	// Set IsConfigured based on whether a REDCap URL is provided
	config.IsConfigured = config.RedCapURL != ""

	return config, nil
}

// UpdateReferralConfig updates the referral configuration
func (s *ReferralService) UpdateReferralConfig(
	ctx context.Context,
	req *models.UpdateReferralConfigRequest,
	userID primitive.ObjectID,
	ipAddress string,
) (*models.ReferralConfig, error) {
	// Get existing config or create new one
	config, err := s.referralRepo.GetConfig(ctx)
	if err != nil {
		if err == repository.ErrReferralConfigNotFound {
			// Create new config
			config = &models.ReferralConfig{
				RedCapURL: "",
				IsEnabled: false,
			}
		} else {
			return nil, err
		}
	}

	// Track old values for audit log
	oldURL := config.RedCapURL
	oldEnabled := config.IsEnabled

	// Update fields if provided
	if req.RedCapURL != nil {
		config.RedCapURL = *req.RedCapURL
	}
	if req.IsEnabled != nil {
		config.IsEnabled = *req.IsEnabled
	}

	// Set IsConfigured based on whether a REDCap URL is provided
	config.IsConfigured = config.RedCapURL != ""

	config.UpdatedBy = &userID

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Save to database
	if config.ID.IsZero() {
		err = s.referralRepo.CreateConfig(ctx, config)
	} else {
		err = s.referralRepo.UpdateConfig(ctx, config)
	}
	if err != nil {
		return nil, err
	}

	// Create audit log
	details := map[string]interface{}{
		"old_url":     oldURL,
		"new_url":     config.RedCapURL,
		"old_enabled": oldEnabled,
		"new_enabled": config.IsEnabled,
	}

	auditLog := &models.AuditLog{
		UserID:      &userID,
		PerformedBy: &userID,
		Action:      models.AuditActionReferralConfigUpdated,
		Details:     details,
		IPAddress:   ipAddress,
	}

	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return config, nil
}

// LogReferralAccess logs when a user accesses the referral link
func (s *ReferralService) LogReferralAccess(
	ctx context.Context,
	userID primitive.ObjectID,
	ipAddress string,
	userAgent string,
) error {
	// Get config to verify it's enabled
	config, err := s.GetReferralConfig(ctx)
	if err != nil {
		return err
	}

	if !config.IsEnabled {
		return ErrReferralDisabled
	}

	// Create audit log
	auditLog := &models.AuditLog{
		UserID:    &userID,
		Action:    models.AuditActionReferralAccessed,
		Details:   map[string]interface{}{"redcap_url": config.RedCapURL},
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	return s.auditRepo.Create(ctx, auditLog)
}

// GetReferralURL gets the REDCap URL if referrals are enabled
func (s *ReferralService) GetReferralURL(ctx context.Context) (string, error) {
	config, err := s.GetReferralConfig(ctx)
	if err != nil {
		return "", err
	}

	if !config.IsEnabled {
		return "", ErrReferralDisabled
	}

	return config.RedCapURL, nil
}
