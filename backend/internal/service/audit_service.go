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
	ID          string `json:"id"`
	Action      string `json:"action"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Icon        string `json:"icon"`
	IconBg      string `json:"iconBg"`
	IconColor   string `json:"iconColor"`
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

	// Get user information - try from database first, then from details
	var userName string
	if log.UserID != nil {
		user, err := s.userRepo.FindByID(ctx, *log.UserID)
		if err == nil && user != nil {
			userName = user.Profile.FirstName + " " + user.Profile.LastName
		} else {
			// User not found (deleted), try to get name from details
			userName = s.getUserNameFromDetails(log.Details)
		}
	} else {
		// No user ID, try to get name from details
		userName = s.getUserNameFromDetails(log.Details)
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
		role := s.getRoleFromDetails(log.Details)
		email := s.getEmailFromDetails(log.Details)
		if email != "" {
			activity.Description = userName + " (" + email + ") - " + role + " - was removed from the system"
		} else {
			activity.Description = userName + " (" + role + ") was removed from the system"
		}
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

	// Institution actions
	case models.AuditActionInstitutionCreated:
		institutionName := s.getInstitutionNameFromDetails(log.Details)
		institutionType := s.getInstitutionTypeFromDetails(log.Details)
		city := s.getCityFromDetails(log.Details)
		activity.Title = "New institution created"
		activity.Description = institutionName + " (" + institutionType + ") in " + city + " was added"
		activity.Icon = "user-plus"
		activity.IconBg = "bg-green-100"
		activity.IconColor = "text-green-600"

	case models.AuditActionInstitutionUpdated:
		institutionName := s.getInstitutionNameFromDetails(log.Details)
		activity.Title = "Institution updated"
		activity.Description = institutionName + " details were updated"
		activity.Icon = "settings"
		activity.IconBg = "bg-purple-100"
		activity.IconColor = "text-purple-600"

	case models.AuditActionInstitutionDeleted:
		institutionName := s.getInstitutionNameFromDetails(log.Details)
		institutionType := s.getInstitutionTypeFromDetails(log.Details)
		city := s.getCityFromDetails(log.Details)
		activity.Title = "Institution deleted"
		activity.Description = institutionName + " (" + institutionType + ") in " + city + " was removed"
		activity.Icon = "trash"
		activity.IconBg = "bg-red-100"
		activity.IconColor = "text-red-600"

	case models.AuditActionInstitutionActivated:
		institutionName := s.getInstitutionNameFromDetails(log.Details)
		activity.Title = "Institution activated"
		activity.Description = institutionName + " was activated"
		activity.Icon = "user-check"
		activity.IconBg = "bg-green-100"
		activity.IconColor = "text-green-600"

	case models.AuditActionInstitutionDeactivated:
		institutionName := s.getInstitutionNameFromDetails(log.Details)
		activity.Title = "Institution deactivated"
		activity.Description = institutionName + " was deactivated"
		activity.Icon = "user-x"
		activity.IconBg = "bg-red-100"
		activity.IconColor = "text-red-600"

	default:
		activity.Title = "System activity"
		activity.Description = string(log.Action)
		activity.Icon = "activity"
		activity.IconBg = "bg-gray-100"
		activity.IconColor = "text-gray-600"
	}

	return activity
}

// getUserNameFromDetails extracts user name from audit log details
func (s *AuditService) getUserNameFromDetails(details map[string]interface{}) string {
	if details == nil {
		return "Unknown User"
	}

	// Try to get first_name and last_name
	firstName, firstOk := details["first_name"].(string)
	lastName, lastOk := details["last_name"].(string)

	if firstOk && lastOk {
		return firstName + " " + lastName
	}

	// Try email as fallback
	if email, ok := details["email"].(string); ok {
		return email
	}

	// Try username as fallback
	if username, ok := details["username"].(string); ok {
		return username
	}

	return "Unknown User"
}

// getEmailFromDetails extracts email from audit log details
func (s *AuditService) getEmailFromDetails(details map[string]interface{}) string {
	if details == nil {
		return ""
	}
	if email, ok := details["email"].(string); ok {
		return email
	}
	return ""
}

// getRoleFromDetails extracts role from audit log details
func (s *AuditService) getRoleFromDetails(details map[string]interface{}) string {
	if role, ok := details["role"].(string); ok {
		switch role {
		case "user":
			return "User"
		case "admin":
			return "Admin"
		// Legacy roles for backward compatibility with old audit logs
		case "haematologist":
			return "User (Haematologist)"
		case "physician":
			return "User (Physician)"
		case "data_capturer":
			return "User (Data Capturer)"
		default:
			return role
		}
	}
	return "User"
}

// getInstitutionNameFromDetails extracts institution name from audit log details
func (s *AuditService) getInstitutionNameFromDetails(details map[string]interface{}) string {
	if details == nil {
		return "Unknown Institution"
	}
	if name, ok := details["institution_name"].(string); ok {
		return name
	}
	return "Unknown Institution"
}

// getInstitutionTypeFromDetails extracts institution type from audit log details
func (s *AuditService) getInstitutionTypeFromDetails(details map[string]interface{}) string {
	if details == nil {
		return "Institution"
	}
	if instType, ok := details["type"].(string); ok {
		switch instType {
		case "university":
			return "University"
		case "hospital":
			return "Hospital"
		case "laboratory":
			return "Laboratory"
		case "research_center":
			return "Research Center"
		case "government":
			return "Government"
		case "private_practice":
			return "Private Practice"
		case "ngo":
			return "NGO"
		case "other":
			return "Other"
		default:
			return instType
		}
	}
	return "Institution"
}

// getCityFromDetails extracts city from audit log details
func (s *AuditService) getCityFromDetails(details map[string]interface{}) string {
	if details == nil {
		return "Unknown City"
	}
	if city, ok := details["city"].(string); ok {
		return city
	}
	return "Unknown City"
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
