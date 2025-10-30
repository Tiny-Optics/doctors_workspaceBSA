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

	// Include redirect_uri if provided
	// If redirect_uri is provided, it must match exactly in token exchange
	if redirectURI != "" {
		params.Add("redirect_uri", redirectURI)
		fmt.Printf("DEBUG: Authorization URL includes redirect_uri: %s\n", redirectURI)
	} else {
		fmt.Println("DEBUG: Authorization URL without redirect_uri - user will copy code manually")
	}

	authURL := fmt.Sprintf("https://www.dropbox.com/oauth2/authorize?%s", params.Encode())
	fmt.Printf("DEBUG: Generated authorization URL (length: %d)\n", len(authURL))
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

	// Include redirect_uri if it was provided (must match authorization URL)
	// According to Dropbox docs: if redirect_uri was used in authorization, it must be included here
	if redirectURI != "" {
		formData.Set("redirect_uri", redirectURI)
		fmt.Printf("DEBUG: Including redirect_uri in token exchange: %s\n", redirectURI)
	} else {
		fmt.Println("DEBUG: No redirect_uri provided - using manual code flow")
	}

	fmt.Printf("DEBUG: Exchanging authorization code (length: %d)\n", len(code))
	fmt.Printf("DEBUG: Request form data (without secrets): grant_type=authorization_code&client_id=%s&redirect_uri=%s\n", appKey, redirectURI)

	// Call Dropbox token endpoint with form-urlencoded data
	resp, err := http.Post(
		"https://api.dropbox.com/oauth2/token",
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		fmt.Printf("ERROR: Failed to call Dropbox token endpoint: %v\n", err)
		return nil, fmt.Errorf("failed to call token endpoint: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: Failed to read Dropbox response body: %v\n", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
		fmt.Printf("ERROR: Dropbox token exchange failed:\n")
		fmt.Printf("  Status: %d\n", resp.StatusCode)
		fmt.Printf("  Response: %s\n", string(body))
		fmt.Printf("  Code length: %d\n", len(code))
		fmt.Printf("  Redirect URI: %s\n", redirectURI)
		return nil, errors.New(errMsg)
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

	// Create parent folder if it doesn't exist (skip if using root)
	if s.dropboxService.IsConfigured() && parentFolder != "" {
		fmt.Printf("Ensuring parent folder exists: %s\n", parentFolder)
		// Use empty relative path to create just the parent folder
		if err := s.dropboxService.CreateFolder(""); err != nil {
			fmt.Printf("Warning: Failed to create parent folder: %v\n", err)
			// Don't fail the authorization, just log the warning
		} else {
			fmt.Printf("Successfully ensured parent folder exists: %s\n", parentFolder)
		}
	} else if parentFolder == "" {
		fmt.Println("Using Dropbox root as parent folder (no folder creation needed)")
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

	status := config.GetPublicStatus()

	// Add refresh service status if available
	if s.dropboxService != nil {
		// We can't directly access the refresh service from here, but we can indicate
		// that background refresh is available
		status["backgroundRefreshAvailable"] = true
	}

	return status, nil
}

// ForceRefresh manually triggers a token refresh
func (s *DropboxOAuthService) ForceRefresh(
	ctx context.Context,
	performedBy *models.User,
	ipAddress string,
) error {
	// Apply a timeout to avoid proxy timeouts/hangs
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// This will trigger a refresh through the dropbox service
	if err := s.dropboxService.ensureValidToken(ctxWithTimeout); err != nil {
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
    // Apply a timeout to avoid proxy timeouts/hangs
    ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
    defer cancel()

    if err := s.dropboxService.TestConnection(ctxWithTimeout); err != nil {
		// Update health status
        s.configRepo.UpdateHealth(ctxWithTimeout, false, err.Error())
		return fmt.Errorf("connection test failed: %w", err)
	}

	// Update health status
    s.configRepo.UpdateHealth(ctxWithTimeout, true, "")

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
