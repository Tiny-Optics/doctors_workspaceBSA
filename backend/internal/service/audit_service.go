package service

import (
	"context"
	"fmt"
	"time"

	"backend/internal/models"
	"backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson"
)

// AuditService handles audit log operations
type AuditService struct {
	auditRepo *repository.AuditRepository
	userRepo  *repository.UserRepository
}

// NewAuditService creates a new AuditService
func NewAuditService(auditRepo *repository.AuditRepository, userRepo *repository.UserRepository) *AuditService {
	return &AuditService{
		auditRepo: auditRepo,
		userRepo:  userRepo,
	}
}

// RecentActivityItem represents a recent activity entry with enriched data
type RecentActivityItem struct {
	ID          string                 `json:"id"`
	Action      string                 `json:"action"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Time        string                 `json:"time"`
	Icon        string                 `json:"icon"`
	IconBg      string                 `json:"iconBg"`
	IconColor   string                 `json:"iconColor"`
}

// GetRecentActivity retrieves recent activity with user information
func (s *AuditService) GetRecentActivity(ctx context.Context, limit int64) ([]RecentActivityItem, error) {
	// Fetch recent audit logs
	logs, err := s.auditRepo.List(ctx, bson.M{}, limit, 0)
	if err != nil {
		return nil, err
	}

	activities := make([]RecentActivityItem, 0, len(logs))
	for _, log := range logs {
		activity := s.enrichAuditLog(ctx, log)
		activities = append(activities, activity)
	}

	return activities, nil
}

// enrichAuditLog converts an audit log to a user-friendly activity item
func (s *AuditService) enrichAuditLog(ctx context.Context, log *models.AuditLog) RecentActivityItem {
	activity := RecentActivityItem{
		ID:     log.ID.Hex(),
		Action: string(log.Action),
		Time:   formatTimeAgo(log.Timestamp),
	}

	// Get user information if available
	var userName string
	if log.UserID != nil {
		user, err := s.userRepo.FindByID(ctx, *log.UserID)
		if err == nil && user != nil {
			userName = user.Profile.FirstName + " " + user.Profile.LastName
		}
	}

	// Get performer information if available
	var performerName string
	if log.PerformedBy != nil {
		performer, err := s.userRepo.FindByID(ctx, *log.PerformedBy)
		if err == nil && performer != nil {
			performerName = performer.Profile.FirstName + " " + performer.Profile.LastName
		}
	}

	// Set title, description, and icon based on action type
	switch log.Action {
	case models.AuditActionUserCreated:
		activity.Title = "New user registered"
		activity.Description = userName + " joined as " + s.getRoleFromDetails(log.Details)
		activity.Icon = "user-plus"
		activity.IconBg = "bg-green-100"
		activity.IconColor = "text-green-600"

	case models.AuditActionLoginSuccess:
		activity.Title = "User login"
		activity.Description = userName + " logged in"
		activity.Icon = "login"
		activity.IconBg = "bg-blue-100"
		activity.IconColor = "text-blue-600"

	case models.AuditActionUserUpdated:
		activity.Title = "Profile updated"
		if performerName != "" && userName != "" && performerName != userName {
			activity.Description = performerName + " updated " + userName + "'s profile"
		} else {
			activity.Description = userName + " updated their profile"
		}
		activity.Icon = "settings"
		activity.IconBg = "bg-purple-100"
		activity.IconColor = "text-purple-600"

	case models.AuditActionUserDeactivated:
		activity.Title = "User deactivated"
		activity.Description = userName + " was deactivated"
		activity.Icon = "user-x"
		activity.IconBg = "bg-red-100"
		activity.IconColor = "text-red-600"

	case models.AuditActionUserActivated:
		activity.Title = "User activated"
		activity.Description = userName + " was activated"
		activity.Icon = "user-check"
		activity.IconBg = "bg-green-100"
		activity.IconColor = "text-green-600"

	case models.AuditActionUserDeleted:
		activity.Title = "User deleted"
		activity.Description = userName + " was removed from the system"
		activity.Icon = "trash"
		activity.IconBg = "bg-red-100"
		activity.IconColor = "text-red-600"

	case models.AuditActionPasswordChanged:
		activity.Title = "Password changed"
		activity.Description = userName + " changed their password"
		activity.Icon = "key"
		activity.IconBg = "bg-yellow-100"
		activity.IconColor = "text-yellow-600"

	case models.AuditActionRoleChanged:
		activity.Title = "Role changed"
		activity.Description = userName + "'s role was updated"
		activity.Icon = "shield"
		activity.IconBg = "bg-orange-100"
		activity.IconColor = "text-orange-600"

	default:
		activity.Title = "System activity"
		activity.Description = string(log.Action)
		activity.Icon = "activity"
		activity.IconBg = "bg-gray-100"
		activity.IconColor = "text-gray-600"
	}

	return activity
}

// getRoleFromDetails extracts role from audit log details
func (s *AuditService) getRoleFromDetails(details map[string]interface{}) string {
	if role, ok := details["role"].(string); ok {
		switch role {
		case "haematologist":
			return "Haematologist"
		case "physician":
			return "Physician"
		case "data_capturer":
			return "Data Capturer"
		case "admin":
			return "Admin"
		default:
			return role
		}
	}
	return "User"
}

// formatTimeAgo formats a time.Time as a relative time string
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	} else if diff < 24*time.Hour {
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if diff < 30*24*time.Hour {
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else {
		months := int(diff.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}
}

