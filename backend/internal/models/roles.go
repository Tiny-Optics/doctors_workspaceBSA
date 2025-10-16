package models

// UserRole represents the type of user in the system
type UserRole string

const (
	RoleHaematologist UserRole = "haematologist"
	RolePhysician     UserRole = "physician"
	RoleDataCapturer  UserRole = "data_capturer"
	RoleAdmin         UserRole = "admin"
)

// IsValid checks if the role is a valid UserRole
func (r UserRole) IsValid() bool {
	switch r {
	case RoleHaematologist, RolePhysician, RoleDataCapturer, RoleAdmin:
		return true
	}
	return false
}

// IsClinical returns true if the role is a clinical role (non-admin)
func (r UserRole) IsClinical() bool {
	switch r {
	case RoleHaematologist, RolePhysician, RoleDataCapturer:
		return true
	}
	return false
}

// AdminLevel represents the level of administrative access
type AdminLevel string

const (
	AdminLevelNone        AdminLevel = ""
	AdminLevelUserManager AdminLevel = "user_manager"
	AdminLevelSuperAdmin  AdminLevel = "super_admin"
)

// IsValid checks if the admin level is valid
func (a AdminLevel) IsValid() bool {
	switch a {
	case AdminLevelNone, AdminLevelUserManager, AdminLevelSuperAdmin:
		return true
	}
	return false
}

// Permission represents a specific permission in the system
type Permission string

const (
	// SOP Permissions
	PermViewSOPs     Permission = "view_sops"
	PermDownloadSOPs Permission = "download_sops"

	// Referral Permissions
	PermAccessReferrals Permission = "access_referrals"

	// Registry Permissions
	PermViewRegistry         Permission = "view_registry"
	PermUploadEthicsApproval Permission = "upload_ethics_approval"

	// Admin Permissions
	PermManageUsers   Permission = "manage_users"
	PermAssignRoles   Permission = "assign_roles"
	PermViewAuditLogs Permission = "view_audit_logs"
	PermManageSystem  Permission = "manage_system"
	PermDeleteUsers   Permission = "delete_users"
)

// GetPermissionsForRole returns all permissions for a given role and admin level
func GetPermissionsForRole(role UserRole, adminLevel AdminLevel) []Permission {
	// Base permissions for all clinical users
	clinicalPermissions := []Permission{
		PermViewSOPs,
		PermDownloadSOPs,
		PermAccessReferrals,
		PermViewRegistry,
		PermUploadEthicsApproval,
	}

	// If not an admin, return clinical permissions
	if role != RoleAdmin {
		return clinicalPermissions
	}

	// Admin permissions based on admin level
	switch adminLevel {
	case AdminLevelUserManager:
		return append(clinicalPermissions, []Permission{
			PermManageUsers,
			PermAssignRoles,
		}...)
	case AdminLevelSuperAdmin:
		return append(clinicalPermissions, []Permission{
			PermManageUsers,
			PermAssignRoles,
			PermViewAuditLogs,
			PermManageSystem,
			PermDeleteUsers,
		}...)
	default:
		// Admin with no level defaults to clinical permissions
		return clinicalPermissions
	}
}

// HasPermission checks if a user with given role and admin level has a specific permission
func HasPermission(role UserRole, adminLevel AdminLevel, permission Permission) bool {
	permissions := GetPermissionsForRole(role, adminLevel)
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

