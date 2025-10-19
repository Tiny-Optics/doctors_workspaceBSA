package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrOAuthConfigNotFound = errors.New("oauth configuration not found in environment")
	ErrInvalidAuthCode     = errors.New("invalid authorization code")
)

// DropboxOAuthService handles Dropbox OAuth operations for admin
type DropboxOAuthService struct {
	configRepo        *repository.DropboxConfigRepository
	auditRepo         *repository.AuditRepository
	encryptionService *EncryptionService
	dropboxService    *DropboxService
}

// NewDropboxOAuthService creates a new DropboxOAuthService
func NewDropboxOAuthService(
	configRepo *repository.DropboxConfigRepository,
	auditRepo *repository.AuditRepository,
	encryptionService *EncryptionService,
	dropboxService *DropboxService,
) *DropboxOAuthService {
	return &DropboxOAuthService{
		configRepo:        configRepo,
		auditRepo:         auditRepo,
		encryptionService: encryptionService,
		dropboxService:    dropboxService,
	}
}

// GenerateAuthorizationURL generates the Dropbox OAuth authorization URL for admin to visit
func (s *DropboxOAuthService) GenerateAuthorizationURL(appKey, redirectURI string) (string, error) {
	if appKey == "" {
		return "", ErrOAuthConfigNotFound
	}

	// Build authorization URL
	// Using code flow with offline access (to get refresh token)
	params := url.Values{}
	params.Add("client_id", appKey)
	params.Add("response_type", "code")
	params.Add("token_access_type", "offline") // Request refresh token

	if redirectURI != "" {
		params.Add("redirect_uri", redirectURI)
	}

	authURL := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?%s", params.Encode())
	return authURL, nil
}

// ExchangeCodeForTokens exchanges the authorization code for access and refresh tokens
func (s *DropboxOAuthService) ExchangeCodeForTokens(
	ctx context.Context,
	code, appKey, appSecret, redirectURI, parentFolder string,
	createdBy *models.User,
	ipAddress string,
) (*models.DropboxConfig, error) {
	if code == "" {
		return nil, ErrInvalidAuthCode
	}
	if appKey == "" || appSecret == "" {
		return nil, ErrOAuthConfigNotFound
	}

	// Prepare token request (form-urlencoded)
	formData := url.Values{}
	formData.Set("code", code)
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", appKey)
	formData.Set("client_secret", appSecret)
	
	if redirectURI != "" {
		formData.Set("redirect_uri", redirectURI)
	}

	// Call Dropbox token endpoint with form-urlencoded data
	resp, err := http.Post(
		"https://api.dropbox.com/oauth2/token",
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call token endpoint: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var tokenResp models.DropboxOAuthResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	// Validate response
	if tokenResp.AccessToken == "" || tokenResp.RefreshToken == "" {
		return nil, errors.New("invalid token response: missing tokens")
	}

	// Encrypt tokens before storing
	encryptedAccessToken, err := s.encryptionService.Encrypt(tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt access token: %w", err)
	}

	encryptedRefreshToken, err := s.encryptionService.Encrypt(tokenResp.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt refresh token: %w", err)
	}

	encryptedAppSecret, err := s.encryptionService.Encrypt(appSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt app secret: %w", err)
	}

	// Check if config already exists
	existingConfig, err := s.configRepo.GetConfig(ctx)
	if err != nil && err != repository.ErrDropboxConfigNotFound {
		return nil, fmt.Errorf("failed to check existing config: %w", err)
	}

	var config *models.DropboxConfig

	if existingConfig != nil {
		// Update existing configuration
		existingConfig.AppKey = appKey
		existingConfig.AppSecret = encryptedAppSecret
		existingConfig.RefreshToken = encryptedRefreshToken
		existingConfig.AccessToken = encryptedAccessToken
		existingConfig.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		existingConfig.ParentFolder = parentFolder
		existingConfig.IsConnected = true
		existingConfig.ConsecutiveFailures = 0
		existingConfig.LastRefreshSuccess = time.Now()
		existingConfig.LastError = ""

		if err := s.configRepo.UpdateConfig(ctx, existingConfig); err != nil {
			return nil, fmt.Errorf("failed to update config: %w", err)
		}

		config = existingConfig
	} else {
		// Create new configuration
		config = &models.DropboxConfig{
			AppKey:              appKey,
			AppSecret:           encryptedAppSecret,
			RefreshToken:        encryptedRefreshToken,
			AccessToken:         encryptedAccessToken,
			TokenExpiry:         time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
			ParentFolder:        parentFolder,
			IsConnected:         true,
			ConsecutiveFailures: 0,
			LastRefreshSuccess:  time.Now(),
			LastError:           "",
			CreatedBy:           &createdBy.ID,
		}

		if err := s.configRepo.CreateConfig(ctx, config); err != nil {
			return nil, fmt.Errorf("failed to create config: %w", err)
		}
	}

	// Reload dropbox service configuration
	if err := s.dropboxService.loadConfigFromDB(ctx); err != nil {
		fmt.Printf("Warning: Failed to reload Dropbox service: %v\n", err)
	}

	// Create parent folder if it doesn't exist
	if s.dropboxService.IsConfigured() {
		fmt.Printf("Ensuring parent folder exists: %s\n", parentFolder)
		// Use empty relative path to create just the parent folder
		if err := s.dropboxService.CreateFolder(""); err != nil {
			fmt.Printf("Warning: Failed to create parent folder: %v\n", err)
			// Don't fail the authorization, just log the warning
		} else {
			fmt.Printf("Successfully ensured parent folder exists: %s\n", parentFolder)
		}
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &createdBy.ID,
		PerformedBy: &createdBy.ID,
		Action:      "dropbox.authorize",
		Details: bson.M{
			"config_id":     config.ID.Hex(),
			"account_id":    tokenResp.AccountID,
			"parent_folder": parentFolder,
		},
		IPAddress: ipAddress,
	})

	fmt.Println("Successfully authorized and configured Dropbox")
	return config, nil
}

// GetStatus returns the current Dropbox connection status
func (s *DropboxOAuthService) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	config, err := s.configRepo.GetConfig(ctx)
	if err != nil {
		if err == repository.ErrDropboxConfigNotFound {
			return map[string]interface{}{
				"configured": false,
				"message":    "Dropbox not configured. Please authorize the app.",
			}, nil
		}
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return config.GetPublicStatus(), nil
}

// ForceRefresh manually triggers a token refresh
func (s *DropboxOAuthService) ForceRefresh(
	ctx context.Context,
	performedBy *models.User,
	ipAddress string,
) error {
	// This will trigger a refresh through the dropbox service
	if err := s.dropboxService.ensureValidToken(ctx); err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &performedBy.ID,
		PerformedBy: &performedBy.ID,
		Action:      "dropbox.force_refresh",
		Details:     bson.M{},
		IPAddress:   ipAddress,
	})

	return nil
}

// TestConnection tests if the Dropbox connection is working
func (s *DropboxOAuthService) TestConnection(
	ctx context.Context,
	performedBy *models.User,
	ipAddress string,
) error {
	if err := s.dropboxService.TestConnection(ctx); err != nil {
		// Update health status
		s.configRepo.UpdateHealth(ctx, false, err.Error())
		return fmt.Errorf("connection test failed: %w", err)
	}

	// Update health status
	s.configRepo.UpdateHealth(ctx, true, "")

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &performedBy.ID,
		PerformedBy: &performedBy.ID,
		Action:      "dropbox.test_connection",
		Details: bson.M{
			"success": true,
		},
		IPAddress: ipAddress,
	})

	return nil
}

// DeleteConfiguration deletes the Dropbox configuration
// This will require re-authorization
func (s *DropboxOAuthService) DeleteConfiguration(
	ctx context.Context,
	performedBy *models.User,
	ipAddress string,
) error {
	config, err := s.configRepo.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if err := s.configRepo.DeleteConfig(ctx); err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	// Mark dropbox service as not configured
	s.dropboxService.isConfigured = false

	// Audit log
	s.auditRepo.Create(ctx, &models.AuditLog{
		UserID:      &performedBy.ID,
		PerformedBy: &performedBy.ID,
		Action:      "dropbox.delete_configuration",
		Details: bson.M{
			"config_id": config.ID.Hex(),
		},
		IPAddress: ipAddress,
	})

	return nil
}
