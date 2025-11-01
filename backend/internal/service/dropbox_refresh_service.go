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
	// Log current expiry state before attempting
	s.dropboxService.cacheMutex.RLock()
	beforeExpiry := s.dropboxService.cachedConfig.TokenExpiry
	expiredBefore := s.dropboxService.cachedConfig.IsTokenExpired()
	s.dropboxService.cacheMutex.RUnlock()
	fmt.Printf("Performing background Dropbox token refresh... (expiredBefore=%t, expiry=%s)\n", expiredBefore, beforeExpiry.Format(time.RFC3339))

	if err := s.dropboxService.ensureValidToken(ctx); err != nil {
		fmt.Printf("Background token refresh failed: %v\n", err)
		return
	}

	// After ensureValidToken, log new expiry and verify quick connectivity
	s.dropboxService.cacheMutex.RLock()
	afterExpiry := s.dropboxService.cachedConfig.TokenExpiry
	expiredAfter := s.dropboxService.cachedConfig.IsTokenExpired()
	s.dropboxService.cacheMutex.RUnlock()
	fmt.Printf("Background refresh completed (expiredAfter=%t, newExpiry=%s)\n", expiredAfter, afterExpiry.Format(time.RFC3339))

	verifyCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := s.dropboxService.quickLiveCheck(verifyCtx); err != nil {
		fmt.Printf("Background refresh live check FAILED: %v\n", err)
		_ = s.dropboxService.configRepo.UpdateHealth(ctx, false, "background live check failed: "+err.Error())
		return
	}

	fmt.Println("Background token refresh verified successfully")
}

// ForceRefreshNow manually triggers a refresh (useful for testing)
func (s *DropboxRefreshService) ForceRefreshNow() error {
	if !s.dropboxService.IsConfigured() {
		return ErrDropboxNotConfigured
	}

	ctx := context.Background()
	return s.dropboxService.ensureValidToken(ctx)
}
