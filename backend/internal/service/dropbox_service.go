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
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Size         uint64            `json:"size"`
	ModifiedTime time.Time         `json:"modifiedTime"`
	IsFolder     bool              `json:"isFolder"`
	Children     []DropboxFileInfo `json:"children,omitempty"`
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

	// Single-flight guard so only one refresh runs at a time
	refreshMutex sync.Mutex
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
		return s.refreshAccessTokenSingleflight(ctx)
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

	// Decrypt the app secret before using it
	decryptedAppSecret, err := s.encryptionService.Decrypt(s.cachedConfig.AppSecret)
	if err != nil {
		s.handleRefreshFailure(ctx, err)
		return fmt.Errorf("failed to decrypt app secret: %w", err)
	}

	// Debug logging (remove in production)
	fmt.Printf("DEBUG: Using AppKey: %s\n", s.cachedConfig.AppKey)
	fmt.Printf("DEBUG: Using AppSecret: %s (length: %d)\n", decryptedAppSecret, len(decryptedAppSecret))
	fmt.Printf("DEBUG: Using RefreshToken: %s (length: %d)\n", s.cachedConfig.RefreshToken, len(s.cachedConfig.RefreshToken))

	// Prepare refresh request (form-urlencoded)
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", s.cachedConfig.RefreshToken)
	formData.Set("client_id", s.cachedConfig.AppKey)
	formData.Set("client_secret", decryptedAppSecret)

	// Call Dropbox token endpoint with context-aware request and timeout
	httpClient := &http.Client{Timeout: 20 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.dropbox.com/oauth2/token", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		s.handleRefreshFailure(ctx, err)
		return fmt.Errorf("%w: failed to build refresh request: %v", ErrTokenRefreshFailed, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
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
		fmt.Printf("ERROR: Dropbox token refresh failed with status %d\n", resp.StatusCode)
		fmt.Printf("ERROR: Response body: %s\n", string(body))
		fmt.Printf("DEBUG: Request form data (without secrets): grant_type=refresh_token&client_id=%s\n", s.cachedConfig.AppKey)
		fmt.Printf("DEBUG: Refresh token length: %d\n", len(s.cachedConfig.RefreshToken))

		// Check for specific error types
		if strings.Contains(string(body), "invalid_grant") {
			fmt.Printf("ERROR: Refresh token is invalid or expired. User may need to re-authorize the app.\n")
			// Mark needs reconnection explicitly for UI
			_ = s.configRepo.UpdateHealth(ctx, false, "invalid_grant during refresh - re-authorization required")
		} else if strings.Contains(string(body), "expired") {
			fmt.Printf("ERROR: Refresh token has expired. User needs to re-authorize the app.\n")
			_ = s.configRepo.UpdateHealth(ctx, false, "refresh token expired - re-authorization required")
		}

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
		fmt.Printf("Successfully reloaded Dropbox config; new expiry: %s (now: %s)\n", s.cachedConfig.TokenExpiry.Format(time.RFC3339), time.Now().Format(time.RFC3339))
	}

	// Immediately verify the refreshed token with a lightweight live check
	verifyCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := s.quickLiveCheck(verifyCtx); err != nil {
		fmt.Printf("ERROR: Post-refresh live check failed: %v\n", err)
		// Mark health degraded so status reflects reality
		_ = s.configRepo.UpdateHealth(ctx, false, "post-refresh live check failed: "+err.Error())
		// Treat this as a refresh failure from caller's perspective
		return fmt.Errorf("%w: post-refresh live check failed: %v", ErrTokenRefreshFailed, err)
	}

	fmt.Println("Successfully refreshed Dropbox access token and verified connectivity")
	return nil
}

// refreshAccessTokenSingleflight ensures only one refresh runs at a time
func (s *DropboxService) refreshAccessTokenSingleflight(ctx context.Context) error {
	s.refreshMutex.Lock()
	defer s.refreshMutex.Unlock()
	return s.refreshAccessToken(ctx)
}

// withClientRetry executes an operation with the Dropbox client and, on
// authentication failure, refreshes the token once and retries the operation.
func (s *DropboxService) withClientRetry(ctx context.Context, op func(files.Client, string) error) error {
	if err := s.ensureValidToken(ctx); err != nil {
		return err
	}

	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	err := op(client, parentFolder)
	if err == nil {
		return nil
	}

	// Detect auth-related errors that require refresh
	if isAuthError(err) {
		if rErr := s.refreshAccessTokenSingleflight(ctx); rErr != nil {
			return err
		}
		s.cacheMutex.RLock()
		client = s.cachedClient
		parentFolder = s.cachedConfig.ParentFolder
		s.cacheMutex.RUnlock()
		return op(client, parentFolder)
	}

	return err
}

func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	// Dropbox SDK surfaces these tokens in error strings
	return strings.Contains(msg, "invalid_access_token") ||
		strings.Contains(msg, "expired_access_token") ||
		strings.Contains(msg, "invalid grant") ||
		strings.Contains(msg, "401")
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
	return s.withClientRetry(ctx, func(client files.Client, parentFolder string) error {
		fullPath := s.getFullPath(relativePath, parentFolder)
		// Check if folder already exists
		_, err := client.GetMetadata(files.NewGetMetadataArg(fullPath))
		if err == nil {
			return nil
		}
		// Create the folder
		createArg := files.NewCreateFolderArg(fullPath)
		_, err = client.CreateFolderV2(createArg)
		if err != nil {
			if strings.Contains(err.Error(), "path/conflict/folder") {
				return nil
			}
			return fmt.Errorf("%w: %v", ErrFolderCreationFailed, err)
		}
		return nil
	})
}

// UploadFile uploads a file to Dropbox
func (s *DropboxService) UploadFile(ctx context.Context, file io.Reader, remotePath string) error {
	return s.withClientRetry(ctx, func(client files.Client, parentFolder string) error {
		fullPath := s.getFullPath(remotePath, parentFolder)
		uploadArg := files.NewUploadArg(fullPath)
		uploadArg.Mode = &files.WriteMode{Tagged: dropbox.Tagged{Tag: "overwrite"}}
		_, err := client.Upload(uploadArg, file)
		if err != nil {
			return fmt.Errorf("failed to upload file to Dropbox: %w", err)
		}
		return nil
	})
}

// ListFiles lists all files in a Dropbox folder
func (s *DropboxService) ListFiles(relativePath string) ([]DropboxFileInfo, error) {
	ctx := context.Background()
	var out []DropboxFileInfo
	err := s.withClientRetry(ctx, func(client files.Client, parentFolder string) error {
		fullPath := s.getFullPath(relativePath, parentFolder)
		listArg := files.NewListFolderArg(fullPath)
		result, err := client.ListFolder(listArg)
		if err != nil {
			if strings.Contains(err.Error(), "path/not_found") {
				return ErrFolderNotFound
			}
			return fmt.Errorf("failed to list folder: %w", err)
		}
		var fileInfos []DropboxFileInfo
		for _, entry := range result.Entries {
			fileInfo := s.entryToFileInfo(entry)
			if fileInfo != nil {
				fileInfos = append(fileInfos, *fileInfo)
			}
		}
		for result.HasMore {
			continueArg := files.NewListFolderContinueArg(result.Cursor)
			result, err = client.ListFolderContinue(continueArg)
			if err != nil {
				return fmt.Errorf("failed to continue listing folder: %w", err)
			}
			for _, entry := range result.Entries {
				fileInfo := s.entryToFileInfo(entry)
				if fileInfo != nil {
					fileInfos = append(fileInfos, *fileInfo)
				}
			}
		}
		out = fileInfos
		return nil
	})
	return out, err
}

// ListFilesRecursive lists all files and folders recursively, building a tree structure
func (s *DropboxService) ListFilesRecursive(relativePath string) ([]DropboxFileInfo, error) {
	ctx := context.Background()
	var out []DropboxFileInfo
	err := s.withClientRetry(ctx, func(client files.Client, parentFolder string) error {
		fullPath := s.getFullPath(relativePath, parentFolder)
		listArg := files.NewListFolderArg(fullPath)
		result, err := client.ListFolder(listArg)
		if err != nil {
			if strings.Contains(err.Error(), "path/not_found") {
				return ErrFolderNotFound
			}
			return fmt.Errorf("failed to list folder: %w", err)
		}
		var fileInfos []DropboxFileInfo
		for _, entry := range result.Entries {
			fileInfo := s.entryToFileInfo(entry)
			if fileInfo != nil {
				if fileInfo.IsFolder {
					children, err := s.listFolderRecursive(client, fileInfo.Path, parentFolder)
					if err != nil {
						fmt.Printf("Warning: Could not access folder %s: %v\n", fileInfo.Path, err)
						fileInfo.Children = []DropboxFileInfo{}
					} else {
						fileInfo.Children = children
					}
				}
				fileInfos = append(fileInfos, *fileInfo)
			}
		}
		for result.HasMore {
			continueArg := files.NewListFolderContinueArg(result.Cursor)
			result, err = client.ListFolderContinue(continueArg)
			if err != nil {
				return fmt.Errorf("failed to continue listing folder: %w", err)
			}
			for _, entry := range result.Entries {
				fileInfo := s.entryToFileInfo(entry)
				if fileInfo != nil {
					if fileInfo.IsFolder {
						children, err := s.listFolderRecursive(client, fileInfo.Path, parentFolder)
						if err != nil {
							fmt.Printf("Warning: Could not access folder %s: %v\n", fileInfo.Path, err)
							fileInfo.Children = []DropboxFileInfo{}
						} else {
							fileInfo.Children = children
						}
					}
					fileInfos = append(fileInfos, *fileInfo)
				}
			}
		}
		out = s.filterEmptyFolders(fileInfos)
		return nil
	})
	return out, err
}

// listFolderRecursive is a helper method to recursively list folder contents
func (s *DropboxService) listFolderRecursive(client files.Client, folderPath, parentFolder string) ([]DropboxFileInfo, error) {
	listArg := files.NewListFolderArg(folderPath)
	result, err := client.ListFolder(listArg)
	if err != nil {
		if strings.Contains(err.Error(), "path/not_found") {
			return []DropboxFileInfo{}, nil
		}
		return nil, err
	}

	var fileInfos []DropboxFileInfo

	// Process initial batch
	for _, entry := range result.Entries {
		fileInfo := s.entryToFileInfo(entry)
		if fileInfo != nil {
			// If it's a folder, recursively get its contents
			if fileInfo.IsFolder {
				children, err := s.listFolderRecursive(client, fileInfo.Path, parentFolder)
				if err != nil {
					// If we can't access the folder, skip it but continue
					fmt.Printf("Warning: Could not access folder %s: %v\n", fileInfo.Path, err)
					fileInfo.Children = []DropboxFileInfo{}
				} else {
					fileInfo.Children = children
				}
			}
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
				// If it's a folder, recursively get its contents
				if fileInfo.IsFolder {
					children, err := s.listFolderRecursive(client, fileInfo.Path, parentFolder)
					if err != nil {
						// If we can't access the folder, skip it but continue
						fmt.Printf("Warning: Could not access folder %s: %v\n", fileInfo.Path, err)
						fileInfo.Children = []DropboxFileInfo{}
					} else {
						fileInfo.Children = children
					}
				}
				fileInfos = append(fileInfos, *fileInfo)
			}
		}
	}

	// Filter out empty folders
	return s.filterEmptyFolders(fileInfos), nil
}

// filterEmptyFolders removes folders that contain no files (only empty subfolders)
func (s *DropboxService) filterEmptyFolders(files []DropboxFileInfo) []DropboxFileInfo {
	var result []DropboxFileInfo

	for _, file := range files {
		if file.IsFolder {
			// Recursively filter children first
			filteredChildren := s.filterEmptyFolders(file.Children)

			// Check if this folder has any files (not just empty subfolders)
			hasFiles := s.folderHasFiles(file.Children)

			// Only include this folder if it has files or non-empty subfolders
			if hasFiles || len(filteredChildren) > 0 {
				file.Children = filteredChildren
				result = append(result, file)
			}
		} else {
			// Always include files
			result = append(result, file)
		}
	}

	return result
}

// folderHasFiles checks if a folder contains any files (not just empty subfolders)
func (s *DropboxService) folderHasFiles(children []DropboxFileInfo) bool {
	for _, child := range children {
		if !child.IsFolder {
			// Found a file
			return true
		}
		// Recursively check subfolders
		if s.folderHasFiles(child.Children) {
			return true
		}
	}
	return false
}

// GetFileDownloadLink generates a temporary download link for a file
func (s *DropboxService) GetFileDownloadLink(relativePath string) (string, error) {
	ctx := context.Background()
	var link string
	err := s.withClientRetry(ctx, func(client files.Client, parentFolder string) error {
		fullPath := s.getFullPath(relativePath, parentFolder)
		arg := files.NewGetTemporaryLinkArg(fullPath)
		result, err := client.GetTemporaryLink(arg)
		if err != nil {
			if strings.Contains(err.Error(), "path/not_found") {
				return ErrFileNotFound
			}
			return fmt.Errorf("failed to get download link: %w", err)
		}
		link = result.Link
		return nil
	})
	return link, err
}

// GetFileMetadata retrieves metadata for a specific file
func (s *DropboxService) GetFileMetadata(relativePath string) (*DropboxFileInfo, error) {
	ctx := context.Background()
	var info *DropboxFileInfo
	err := s.withClientRetry(ctx, func(client files.Client, parentFolder string) error {
		fullPath := s.getFullPath(relativePath, parentFolder)
		metadataArg := files.NewGetMetadataArg(fullPath)
		metadata, err := client.GetMetadata(metadataArg)
		if err != nil {
			if strings.Contains(err.Error(), "path/not_found") {
				return ErrFileNotFound
			}
			return fmt.Errorf("failed to get file metadata: %w", err)
		}
		info = s.metadataToFileInfo(metadata)
		return nil
	})
	return info, err
}

// RenameFolder renames a folder in Dropbox
func (s *DropboxService) RenameFolder(oldRelativePath, newRelativePath string) error {
	ctx := context.Background()
	return s.withClientRetry(ctx, func(client files.Client, parentFolder string) error {
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
	})
}

// TestConnection tests the Dropbox connection
func (s *DropboxService) TestConnection(ctx context.Context) error {
	return s.withClientRetry(ctx, func(client files.Client, parentFolder string) error {
		testPath := parentFolder
		if testPath == "" {
			testPath = ""
		}
		listArg := files.NewListFolderArg(testPath)
		_, err := client.ListFolder(listArg)
		if err != nil {
			return fmt.Errorf("dropbox connection test failed: %w", err)
		}
		return nil
	})
}

// quickLiveCheck performs a minimal Dropbox API call without attempting refresh
// Used to verify connectivity right after refresh to avoid false positives.
func (s *DropboxService) quickLiveCheck(ctx context.Context) error {
	s.cacheMutex.RLock()
	client := s.cachedClient
	parentFolder := s.cachedConfig.ParentFolder
	s.cacheMutex.RUnlock()

	testPath := parentFolder
	if testPath == "" {
		testPath = ""
	}
	listArg := files.NewListFolderArg(testPath)
	_, err := client.ListFolder(listArg)
	if err != nil {
		return err
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
