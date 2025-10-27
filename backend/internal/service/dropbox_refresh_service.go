package service

import (
	"context"
	"fmt"
	"time"
)

// DropboxRefreshService handles automatic background refresh of Dropbox tokens
type DropboxRefreshService struct {
	dropboxService *DropboxService
	ticker         *time.Ticker
	done           chan bool
	isRunning      bool
}

// NewDropboxRefreshService creates a new DropboxRefreshService
func NewDropboxRefreshService(dropboxService *DropboxService) *DropboxRefreshService {
	return &DropboxRefreshService{
		dropboxService: dropboxService,
		done:           make(chan bool),
		isRunning:      false,
	}
}

// Start begins the background token refresh service
// Refreshes every 3 hours (before 4-hour Dropbox token expiry)
func (s *DropboxRefreshService) Start() {
	if s.isRunning {
		fmt.Println("Dropbox refresh service is already running")
		return
	}

	// Refresh every 3 hours (before 4-hour expiry)
	s.ticker = time.NewTicker(3 * time.Hour)
	s.isRunning = true

	fmt.Println("Starting Dropbox background refresh service (every 3 hours)")

	go func() {
		// Do an initial refresh check on startup
		s.refreshIfNeeded()

		for {
			select {
			case <-s.ticker.C:
				s.refreshIfNeeded()
			case <-s.done:
				fmt.Println("Dropbox refresh service stopped")
				return
			}
		}
	}()
}

// Stop stops the background token refresh service
func (s *DropboxRefreshService) Stop() {
	if !s.isRunning {
		return
	}

	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
	s.isRunning = false
	fmt.Println("Stopping Dropbox background refresh service")
}

// IsRunning returns whether the service is currently running
func (s *DropboxRefreshService) IsRunning() bool {
	return s.isRunning
}

// refreshIfNeeded checks if Dropbox is configured and refreshes the token if needed
func (s *DropboxRefreshService) refreshIfNeeded() {
	if !s.dropboxService.IsConfigured() {
		fmt.Println("Dropbox not configured, skipping background refresh")
		return
	}

	ctx := context.Background()
	fmt.Println("Performing background Dropbox token refresh...")

	if err := s.dropboxService.ensureValidToken(ctx); err != nil {
		fmt.Printf("Background token refresh failed: %v\n", err)
	} else {
		fmt.Println("Background token refresh successful")
	}
}

// ForceRefreshNow manually triggers a refresh (useful for testing)
func (s *DropboxRefreshService) ForceRefreshNow() error {
	if !s.dropboxService.IsConfigured() {
		return ErrDropboxNotConfigured
	}

	ctx := context.Background()
	return s.dropboxService.ensureValidToken(ctx)
}
