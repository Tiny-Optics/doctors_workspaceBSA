package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

var (
	ErrDropboxNotConfigured = errors.New("dropbox is not properly configured")
	ErrFolderCreationFailed = errors.New("failed to create folder in dropbox")
	ErrFolderNotFound       = errors.New("folder not found in dropbox")
	ErrFileNotFound         = errors.New("file not found in dropbox")
)

// DropboxFileInfo represents metadata about a file in Dropbox
type DropboxFileInfo struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         uint64    `json:"size"`
	ModifiedTime time.Time `json:"modifiedTime"`
	IsFolder     bool      `json:"isFolder"`
}

// DropboxService handles Dropbox operations
type DropboxService struct {
	client       files.Client
	parentFolder string
	isConfigured bool
}

// NewDropboxService creates a new DropboxService
func NewDropboxService() *DropboxService {
	accessToken := os.Getenv("DROPBOX_APP_API_ACCESS_TOKEN")
	parentFolder := os.Getenv("DROPBOX_APP_PARENT_FOLDER")

	if accessToken == "" || parentFolder == "" {
		return &DropboxService{
			isConfigured: false,
		}
	}

	config := dropbox.Config{
		Token:    accessToken,
		LogLevel: dropbox.LogOff,
	}

	client := files.New(config)

	return &DropboxService{
		client:       client,
		parentFolder: parentFolder,
		isConfigured: true,
	}
}

// IsConfigured returns whether Dropbox is properly configured
func (s *DropboxService) IsConfigured() bool {
	return s.isConfigured
}

// CreateFolder creates a folder in Dropbox
func (s *DropboxService) CreateFolder(relativePath string) error {
	if !s.isConfigured {
		return ErrDropboxNotConfigured
	}

	fullPath := s.getFullPath(relativePath)

	// Check if folder already exists
	_, err := s.client.GetMetadata(files.NewGetMetadataArg(fullPath))
	if err == nil {
		// Folder already exists
		return nil
	}

	// Create the folder
	createArg := files.NewCreateFolderArg(fullPath)
	_, err = s.client.CreateFolderV2(createArg)
	if err != nil {
		// Check if error is because folder already exists
		if strings.Contains(err.Error(), "path/conflict/folder") {
			return nil
		}
		return fmt.Errorf("%w: %v", ErrFolderCreationFailed, err)
	}

	return nil
}

// ListFiles lists all files in a Dropbox folder
func (s *DropboxService) ListFiles(relativePath string) ([]DropboxFileInfo, error) {
	if !s.isConfigured {
		return nil, ErrDropboxNotConfigured
	}

	fullPath := s.getFullPath(relativePath)

	listArg := files.NewListFolderArg(fullPath)
	result, err := s.client.ListFolder(listArg)
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
		result, err = s.client.ListFolderContinue(continueArg)
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
	if !s.isConfigured {
		return "", ErrDropboxNotConfigured
	}

	fullPath := s.getFullPath(relativePath)

	// Create a temporary link (valid for 4 hours)
	arg := files.NewGetTemporaryLinkArg(fullPath)
	result, err := s.client.GetTemporaryLink(arg)
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
	if !s.isConfigured {
		return nil, ErrDropboxNotConfigured
	}

	fullPath := s.getFullPath(relativePath)

	metadataArg := files.NewGetMetadataArg(fullPath)
	metadata, err := s.client.GetMetadata(metadataArg)
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
	if !s.isConfigured {
		return ErrDropboxNotConfigured
	}

	oldFullPath := s.getFullPath(oldRelativePath)
	newFullPath := s.getFullPath(newRelativePath)

	moveArg := files.NewRelocationArg(oldFullPath, newFullPath)
	_, err := s.client.MoveV2(moveArg)
	if err != nil {
		if strings.Contains(err.Error(), "from_lookup/not_found") {
			return ErrFolderNotFound
		}
		return fmt.Errorf("failed to rename folder: %w", err)
	}

	return nil
}

// TestConnection tests the Dropbox connection
func (s *DropboxService) TestConnection(ctx context.Context) error {
	if !s.isConfigured {
		return ErrDropboxNotConfigured
	}

	// Try to list the parent folder
	listArg := files.NewListFolderArg(s.parentFolder)
	_, err := s.client.ListFolder(listArg)
	if err != nil {
		return fmt.Errorf("dropbox connection test failed: %w", err)
	}

	return nil
}

// Helper methods

func (s *DropboxService) getFullPath(relativePath string) string {
	// Clean the relative path
	relativePath = strings.TrimPrefix(relativePath, "/")
	relativePath = strings.TrimSuffix(relativePath, "/")

	// Join with parent folder
	if relativePath == "" {
		return s.parentFolder
	}

	return filepath.Join(s.parentFolder, relativePath)
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
