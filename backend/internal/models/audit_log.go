package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionUserCreated            AuditAction = "user_created"
	AuditActionUserRegistered         AuditAction = "user_registered"
	AuditActionUserUpdated            AuditAction = "user_updated"
	AuditActionUserDeleted            AuditAction = "user_deleted"
	AuditActionUserDeactivated        AuditAction = "user_deactivated"
	AuditActionUserActivated          AuditAction = "user_activated"
	AuditActionRoleChanged            AuditAction = "role_changed"
	AuditActionAdminLevelChanged      AuditAction = "admin_level_changed"
	AuditActionLoginSuccess           AuditAction = "login_success"
	AuditActionLoginFailed            AuditAction = "login_failed"
	AuditActionLogout                 AuditAction = "logout"
	AuditActionPasswordChanged        AuditAction = "password_changed"
	AuditActionAccountLocked          AuditAction = "account_locked"
	AuditActionAccountUnlocked        AuditAction = "account_unlocked"
	AuditActionInstitutionCreated     AuditAction = "institution_created"
	AuditActionInstitutionUpdated     AuditAction = "institution_updated"
	AuditActionInstitutionDeleted     AuditAction = "institution_deleted"
	AuditActionInstitutionActivated   AuditAction = "institution_activated"
	AuditActionInstitutionDeactivated AuditAction = "institution_deactivated"
	AuditActionReferralConfigUpdated  AuditAction = "referral_config_updated"
	AuditActionReferralAccessed       AuditAction = "referral_accessed"
	AuditActionSMTPConfigUpdated      AuditAction = "smtp_config_updated"
	AuditActionPasswordResetRequested AuditAction = "password_reset_requested"
	AuditActionPasswordResetCompleted AuditAction = "password_reset_completed"
)

// AuditLog represents a log entry for audit trail
type AuditLog struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID      *primitive.ObjectID    `bson:"user_id,omitempty" json:"userId,omitempty"`
	PerformedBy *primitive.ObjectID    `bson:"performed_by,omitempty" json:"performedBy,omitempty"`
	Action      AuditAction            `bson:"action" json:"action"`
	Details     map[string]interface{} `bson:"details,omitempty" json:"details,omitempty"`
	IPAddress   string                 `bson:"ip_address" json:"ipAddress"`
	UserAgent   string                 `bson:"user_agent,omitempty" json:"userAgent,omitempty"`
	Timestamp   time.Time              `bson:"timestamp" json:"timestamp"`
}
