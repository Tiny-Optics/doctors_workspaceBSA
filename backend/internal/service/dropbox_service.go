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
	"path/filepath"
	"strings"
	"sync"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

var (
	ErrDropboxNotConfigured = errors.New("dropbox is not properly configured")
	ErrFolderCreationFailed = errors.New("failed to create folder in dropbox")
	ErrFolderNotFound       = errors.New("folder not found in dropbox")
	ErrFileNotFound         = errors.New("file not found in dropbox")
	ErrTokenRefreshFailed   = errors.New("failed to refresh access token")
)

// DropboxFileInfo represents metadata about a file in Dropbox
type DropboxFileInfo struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         uint64    `json:"size"`
	ModifiedTime time.Time `json:"modifiedTime"`
	IsFolder     bool      `json:"isFolder"`
}

// DropboxService handles Dropbox operations with automatic token refresh
type DropboxService struct {
	configRepo        *repository.DropboxConfigRepository
	encryptionService *EncryptionService

	// In-memory cache to avoid DB hits on every request
	cachedConfig *models.DropboxConfig
	cachedClient files.Client
	cacheMutex   sync.RWMutex

	// Configuration
	isConfigured   bool
	alertThreshold int // Number of consecutive failures before alerting
}

// NewDropboxService creates a new DropboxService with DB-backed configuration
func NewDropboxService(configRepo *repository.DropboxConfigRepository, encryptionService *EncryptionService) *DropboxService {
	service := &DropboxService{
		configRepo:        configRepo,
		encryptionService: encryptionService,
		alertThreshold:    3,
		isConfigured:      false,
	}

	// Try to load configuration from database
	ctx := context.Background()
	if err := service.loadConfigFromDB(ctx); err != nil {
		if err != repository.ErrDropboxConfigNotFound {
			fmt.Printf("Warning: Failed to load Dropbox config from DB: %v\n", err)
		}
		// Not configured yet, admin needs to set it up
		return service
	}

	return service
}

// IsConfigured returns whether Dropbox is properly configured
func (s *DropboxService) IsConfigured() bool {
	return s.isConfigured
}

// loadConfigFromDB loads the Dropbox configuration from database and caches it
func (s *DropboxService) loadConfigFromDB(ctx context.Context) error {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	config, err := s.configRepo.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Decrypt tokens
	decryptedAccessToken, err := s.encryptionService.Decrypt(config.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt access token: %w", err)
	}

	decryptedRefreshToken, err := s.encryptionService.Decrypt(config.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt refresh token: %w", err)
	}

	// Store decrypted versions in cached config
	config.AccessToken = decryptedAccessToken
	config.RefreshToken = decryptedRefreshToken

	// Create Dropbox client with current access token
	dropboxConfig := dropbox.Config{
		Token:    decryptedAccessToken,
		LogLevel: dropbox.LogOff,
	}

	s.cachedConfig = config
	s.cachedClient = files.New(dropboxConfig)
	s.isConfigured = true

	return nil
}

// ensureValidToken checks if the access token is valid and refreshes if necessary
func (s *DropboxService) ensureValidToken(ctx context.Context) error {
	if !s.isConfigured {
		return ErrDropboxNotConfigured
	}

	s.cacheMutex.RLock()
	needsRefresh := s.cachedConfig.IsTokenExpired()
	s.cacheMutex.RUnlock()

	if needsRefresh {
		return s.refreshAccessToken(ctx)
	}

	return nil
}

// refreshAccessToken refreshes the Dropbox access token using the refresh token
func (s *DropboxService) refreshAccessToken(ctx context.Context) error {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	if s.cachedConfig == nil || s.cachedConfig.RefreshToken == "" {
		return ErrDropboxNotConfigured
	}

	fmt.Println("Refreshing Dropbox access token...")

	// Prepare refresh request (form-urlencoded)
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", s.cachedConfig.RefreshToken)
	formData.Set("client_id", s.cachedConfig.AppKey)
	formData.Set("client_secret", s.cachedConfig.AppSecret)

	// Call Dropbox token endpoint with form-urlencoded data
	resp, err := http.Post(
		"https://api.dropbox.com/oauth2/token",
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(formData.Encode()),
	)
	if err != nil {
		s.handleRefreshFailure(ctx, err)
		return fmt.Errorf("%w: %v", ErrTokenRefreshFailed, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.handleRefreshFailure(ctx, err)
		return fmt.Errorf("%w: failed to read response: %v", ErrTokenRefreshFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("status %d: %s", resp.StatusCode, string(body))
		s.handleRefreshFailure(ctx, errors.New(errMsg))
		return fmt.Errorf("%w: %s", ErrTokenRefreshFailed, errMsg)
	}

	// Parse response
	var tokenResp models.DropboxOAuthResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		s.handleRefreshFailure(ctx, err)
		return fmt.Errorf("%w: failed to parse response: %v", ErrTokenRefreshFailed, err)
	}

	// Encrypt the new access token
	encryptedAccessToken, err := s.encryptionService.Encrypt(tokenResp.AccessToken)
	if err != nil {
		s.handleRefreshFailure(ctx, err)
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	// Update database with new token (refresh token not included in response for refresh grant)
	if err := s.configRepo.UpdateTokens(ctx, encryptedAccessToken, "", tokenResp.ExpiresIn); err != nil {
		s.handleRefreshFailure(ctx, err)
		return fmt.Errorf("failed to update tokens in database: %w", err)
	}

	// Update cached config
	s.cachedConfig.AccessToken = tokenResp.AccessToken // Store decrypted in memory
	s.cachedConfig.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	// Create new Dropbox client with refreshed token
	dropboxConfig := dropbox.Config{
		Token:    tokenResp.AccessToken,
		LogLevel: dropbox.LogOff,
	}
	s.cachedClient = files.New(dropboxConfig)

	// Reset failure count on success
	if err := s.configRepo.ResetFailures(ctx); err != nil {
		fmt.Printf("Warning: Failed to reset failure count: %v\n", err)
	}

	// Reload configuration from database to ensure cached config is in sync
	// This is critical to prevent refresh token issues after 24+ hours
	if err := s.loadConfigFromDB(ctx); err != nil {
		fmt.Printf("Warning: Failed to reload config after refresh: %v\n", err)
		// Don't return error here as the refresh was successful
		// The cached config will be updated on next operation
	} else {
		fmt.Println("Successfully reloaded Dropbox config from database after token refresh")
	}

	fmt.Println("Successfully refreshed Dropbox access token")
	return nil
}

// handleRefreshFailure handles token refresh failures
func (s *DropboxService) handleRefreshFailure(ctx context.Context, err error) {
	fmt.Printf("ERROR: Dropbox token refresh failed: %v\n", err)

	// Increment failure count in database
	if dbErr := s.configRepo.IncrementFailures(ctx); dbErr != nil {
		fmt.Printf("ERROR: Failed to increment failure count: %v\n", dbErr)
	}

	// Update health status
	if dbErr := s.configRepo.UpdateHealth(ctx, false, err.Error()); dbErr != nil {
		fmt.Printf("ERROR: Failed to update health status: %v\n", dbErr)
	}

	// Check if we should send alert
	if s.cachedConfig != nil && s.cachedConfig.ConsecutiveFailures >= s.alertThreshold {
		s.sendAlert(err)
	}
}

// sendAlert sends an alert when Dropbox connection fails
// TODO: Implement email notifications
func (s *DropboxService) sendAlert(err error) {
	fmt.Printf("ALERT: Dropbox has failed %d times. Last error: %v\n", s.cachedConfig.ConsecutiveFailures, err)
	fmt.Println("TODO: Send email notification to admin")
	// TODO: Implement email service integration
	// emailService.SendAlert(adminEmail, "Dropbox Token Refresh Failed", message)
}

// CreateFolder creates a folder in Dropbox
func (s *DropboxService) CreateFolder(relativePath string) error {
	ctx := context.Background()
	if err := s.ensureValidToken(ctx); err != nil {
		return err
	}

	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	fullPath := s.getFullPath(relativePath, parentFolder)

	// Check if folder already exists
	_, err := client.GetMetadata(files.NewGetMetadataArg(fullPath))
	if err == nil {
		// Folder already exists
		return nil
	}

	// Create the folder
	createArg := files.NewCreateFolderArg(fullPath)
	_, err = client.CreateFolderV2(createArg)
	if err != nil {
		// Check if error is because folder already exists
		if strings.Contains(err.Error(), "path/conflict/folder") {
			return nil
		}
		return fmt.Errorf("%w: %v", ErrFolderCreationFailed, err)
	}

	return nil
}

// UploadFile uploads a file to Dropbox
func (s *DropboxService) UploadFile(ctx context.Context, file io.Reader, remotePath string) error {
	if err := s.ensureValidToken(ctx); err != nil {
		return err
	}

	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	fullPath := s.getFullPath(remotePath, parentFolder)

	// Create upload argument
	uploadArg := files.NewUploadArg(fullPath)
	uploadArg.Mode = &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}}

	// Upload the file
	_, err := client.Upload(uploadArg, file)
	if err != nil {
		return fmt.Errorf("failed to upload file to Dropbox: %w", err)
	}

	return nil
}

// ListFiles lists all files in a Dropbox folder
func (s *DropboxService) ListFiles(relativePath string) ([]DropboxFileInfo, error) {
	ctx := context.Background()
	if err := s.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	fullPath := s.getFullPath(relativePath, parentFolder)

	listArg := files.NewListFolderArg(fullPath)
	result, err := client.ListFolder(listArg)
	if err != nil {
		if strings.Contains(err.Error(), "path/not_found") {
			return nil, ErrFolderNotFound
		}
		return nil, fmt.Errorf("failed to list folder: %w", err)
	}

	var fileInfos []DropboxFileInfo

	// Process initial batch
	for _, entry := range result.Entries {
		fileInfo := s.entryToFileInfo(entry)
		if fileInfo != nil {
			fileInfos = append(fileInfos, *fileInfo)
		}
	}

	// Handle pagination if there are more results
	for result.HasMore {
		continueArg := files.NewListFolderContinueArg(result.Cursor)
		result, err = client.ListFolderContinue(continueArg)
		if err != nil {
			return fileInfos, fmt.Errorf("failed to continue listing folder: %w", err)
		}

		for _, entry := range result.Entries {
			fileInfo := s.entryToFileInfo(entry)
			if fileInfo != nil {
				fileInfos = append(fileInfos, *fileInfo)
			}
		}
	}

	return fileInfos, nil
}

// GetFileDownloadLink generates a temporary download link for a file
func (s *DropboxService) GetFileDownloadLink(relativePath string) (string, error) {
	ctx := context.Background()
	if err := s.ensureValidToken(ctx); err != nil {
		return "", err
	}

	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	fullPath := s.getFullPath(relativePath, parentFolder)

	// Create a temporary link (valid for 4 hours)
	arg := files.NewGetTemporaryLinkArg(fullPath)
	result, err := client.GetTemporaryLink(arg)
	if err != nil {
		if strings.Contains(err.Error(), "path/not_found") {
			return "", ErrFileNotFound
		}
		return "", fmt.Errorf("failed to get download link: %w", err)
	}

	return result.Link, nil
}

// GetFileMetadata retrieves metadata for a specific file
func (s *DropboxService) GetFileMetadata(relativePath string) (*DropboxFileInfo, error) {
	ctx := context.Background()
	if err := s.ensureValidToken(ctx); err != nil {
		return nil, err
	}

	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	fullPath := s.getFullPath(relativePath, parentFolder)

	metadataArg := files.NewGetMetadataArg(fullPath)
	metadata, err := client.GetMetadata(metadataArg)
	if err != nil {
		if strings.Contains(err.Error(), "path/not_found") {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get file metadata: %w", err)
	}

	fileInfo := s.metadataToFileInfo(metadata)
	return fileInfo, nil
}

// RenameFolder renames a folder in Dropbox
func (s *DropboxService) RenameFolder(oldRelativePath, newRelativePath string) error {
	ctx := context.Background()
	if err := s.ensureValidToken(ctx); err != nil {
		return err
	}

	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	oldFullPath := s.getFullPath(oldRelativePath, parentFolder)
	newFullPath := s.getFullPath(newRelativePath, parentFolder)

	fmt.Printf("Dropbox RenameFolder - Parent: '%s', Old relative: '%s', New relative: '%s'\n", parentFolder, oldRelativePath, newRelativePath)
	fmt.Printf("Dropbox RenameFolder - Old full path: '%s', New full path: '%s'\n", oldFullPath, newFullPath)

	moveArg := files.NewRelocationArg(oldFullPath, newFullPath)
	_, err := client.MoveV2(moveArg)
	if err != nil {
		fmt.Printf("Dropbox RenameFolder ERROR: %v\n", err)
		if strings.Contains(err.Error(), "from_lookup/not_found") {
			return fmt.Errorf("%w: folder '%s' not found in Dropbox", ErrFolderNotFound, oldFullPath)
		}
		return fmt.Errorf("failed to rename folder from '%s' to '%s': %w", oldFullPath, newFullPath, err)
	}

	fmt.Printf("Dropbox RenameFolder SUCCESS: '%s' → '%s'\n", oldFullPath, newFullPath)
	return nil
}

// TestConnection tests the Dropbox connection
func (s *DropboxService) TestConnection(ctx context.Context) error {
	if err := s.ensureValidToken(ctx); err != nil {
		return err
	}

	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	// Normalize parent folder (empty = root)
	testPath := parentFolder
	if testPath == "" {
		testPath = "" // Dropbox API uses empty string for root
	}

	// Try to list the parent folder
	listArg := files.NewListFolderArg(testPath)
	_, err := client.ListFolder(listArg)
	if err != nil {
		return fmt.Errorf("dropbox connection test failed: %w", err)
	}

	return nil
}

// Helper methods

func (s *DropboxService) getFullPath(relativePath string, parentFolder string) string {
	// Clean the relative path
	relativePath = strings.TrimPrefix(relativePath, "/")
	relativePath = strings.TrimSuffix(relativePath, "/")

	// Decode URL-encoded path (for categories with spaces)
	if decoded, err := url.PathUnescape(relativePath); err == nil {
		relativePath = decoded
	}

	// If parent folder is empty, use root
	if parentFolder == "" {
		parentFolder = "/"
	}

	// Join with parent folder
	if relativePath == "" {
		return parentFolder
	}

	// Ensure result starts with / for Dropbox
	result := filepath.Join(parentFolder, relativePath)
	if !strings.HasPrefix(result, "/") {
		result = "/" + result
	}

	return result
}

func (s *DropboxService) entryToFileInfo(entry files.IsMetadata) *DropboxFileInfo {
	switch meta := entry.(type) {
	case *files.FileMetadata:
		return &DropboxFileInfo{
			Name:         meta.Name,
			Path:         meta.PathDisplay,
			Size:         meta.Size,
			ModifiedTime: meta.ServerModified,
			IsFolder:     false,
		}
	case *files.FolderMetadata:
		return &DropboxFileInfo{
			Name:     meta.Name,
			Path:     meta.PathDisplay,
			IsFolder: true,
		}
	default:
		return nil
	}
}

func (s *DropboxService) metadataToFileInfo(metadata files.IsMetadata) *DropboxFileInfo {
	return s.entryToFileInfo(metadata)
}
