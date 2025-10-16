package models

import (
	"testing"
)

func TestUserRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role UserRole
		want bool
	}{
		{"Valid haematologist", RoleHaematologist, true},
		{"Valid physician", RolePhysician, true},
		{"Valid data capturer", RoleDataCapturer, true},
		{"Valid admin", RoleAdmin, true},
		{"Invalid role", UserRole("invalid"), false},
		{"Empty role", UserRole(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("UserRole.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRole_IsClinical(t *testing.T) {
	tests := []struct {
		name string
		role UserRole
		want bool
	}{
		{"Haematologist is clinical", RoleHaematologist, true},
		{"Physician is clinical", RolePhysician, true},
		{"Data capturer is clinical", RoleDataCapturer, true},
		{"Admin is not clinical", RoleAdmin, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsClinical(); got != tt.want {
				t.Errorf("UserRole.IsClinical() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminLevel_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		level AdminLevel
		want  bool
	}{
		{"Valid none", AdminLevelNone, true},
		{"Valid user manager", AdminLevelUserManager, true},
		{"Valid super admin", AdminLevelSuperAdmin, true},
		{"Invalid level", AdminLevel("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.IsValid(); got != tt.want {
				t.Errorf("AdminLevel.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPermissionsForRole(t *testing.T) {
	clinicalPermissions := []Permission{
		PermViewSOPs,
		PermDownloadSOPs,
		PermAccessReferrals,
		PermViewRegistry,
		PermUploadEthicsApproval,
	}

	tests := []struct {
		name       string
		role       UserRole
		adminLevel AdminLevel
		wantCount  int
		mustHave   []Permission
	}{
		{
			name:      "Haematologist permissions",
			role:      RoleHaematologist,
			wantCount: 5,
			mustHave:  clinicalPermissions,
		},
		{
			name:      "Physician permissions",
			role:      RolePhysician,
			wantCount: 5,
			mustHave:  clinicalPermissions,
		},
		{
			name:      "Data capturer permissions",
			role:      RoleDataCapturer,
			wantCount: 5,
			mustHave:  clinicalPermissions,
		},
		{
			name:       "User manager permissions",
			role:       RoleAdmin,
			adminLevel: AdminLevelUserManager,
			wantCount:  7,
			mustHave:   append(clinicalPermissions, PermManageUsers, PermAssignRoles),
		},
		{
			name:       "Super admin permissions",
			role:       RoleAdmin,
			adminLevel: AdminLevelSuperAdmin,
			wantCount:  10,
			mustHave: append(clinicalPermissions, PermManageUsers, PermAssignRoles,
				PermViewAuditLogs, PermManageSystem, PermDeleteUsers),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPermissionsForRole(tt.role, tt.adminLevel)

			if len(got) != tt.wantCount {
				t.Errorf("GetPermissionsForRole() returned %d permissions, want %d", len(got), tt.wantCount)
			}

			for _, perm := range tt.mustHave {
				found := false
				for _, p := range got {
					if p == perm {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetPermissionsForRole() missing permission %s", perm)
				}
			}
		})
	}
}

func TestHasPermission(t *testing.T) {
	tests := []struct {
		name       string
		role       UserRole
		adminLevel AdminLevel
		permission Permission
		want       bool
	}{
		{
			name:       "Clinical user has view SOPs",
			role:       RoleHaematologist,
			permission: PermViewSOPs,
			want:       true,
		},
		{
			name:       "Clinical user does not have manage users",
			role:       RolePhysician,
			permission: PermManageUsers,
			want:       false,
		},
		{
			name:       "User manager has manage users",
			role:       RoleAdmin,
			adminLevel: AdminLevelUserManager,
			permission: PermManageUsers,
			want:       true,
		},
		{
			name:       "User manager does not have delete users",
			role:       RoleAdmin,
			adminLevel: AdminLevelUserManager,
			permission: PermDeleteUsers,
			want:       false,
		},
		{
			name:       "Super admin has delete users",
			role:       RoleAdmin,
			adminLevel: AdminLevelSuperAdmin,
			permission: PermDeleteUsers,
			want:       true,
		},
		{
			name:       "Super admin has all permissions",
			role:       RoleAdmin,
			adminLevel: AdminLevelSuperAdmin,
			permission: PermViewAuditLogs,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPermission(tt.role, tt.adminLevel, tt.permission); got != tt.want {
				t.Errorf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}
