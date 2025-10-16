package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", "SecurePass123!", false},
		{"Valid password 2", "Abcdef1@", false},
		{"Too short", "Short1!", true},
		{"No uppercase", "nouppercas1!", true},
		{"No lowercase", "NOLOWERCASE1!", true},
		{"No number", "NoNumberHere!", true},
		{"No special char", "NoSpecialChar1", true},
		{"Empty password", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateUserRequest_Validate(t *testing.T) {
	validRequest := &CreateUserRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		Password:    "ValidPass123!",
		Role:        RoleHaematologist,
		FirstName:   "John",
		LastName:    "Doe",
		Institution: "Test Hospital",
		Location:    "Cape Town",
	}

	tests := []struct {
		name    string
		modify  func(*CreateUserRequest)
		wantErr bool
	}{
		{"Valid request", func(r *CreateUserRequest) {}, false},
		{"Invalid email", func(r *CreateUserRequest) { r.Email = "invalid-email" }, true},
		{"Invalid username - too short", func(r *CreateUserRequest) { r.Username = "ab" }, true},
		{"Invalid username - special chars", func(r *CreateUserRequest) { r.Username = "test@user" }, true},
		{"Weak password", func(r *CreateUserRequest) { r.Password = "weak" }, true},
		{"Invalid role", func(r *CreateUserRequest) { r.Role = UserRole("invalid") }, true},
		{"Admin without admin level", func(r *CreateUserRequest) {
			r.Role = RoleAdmin
			r.AdminLevel = AdminLevelNone
		}, true},
		{"Admin with admin level", func(r *CreateUserRequest) {
			r.Role = RoleAdmin
			r.AdminLevel = AdminLevelSuperAdmin
		}, false},
		{"Non-admin with admin level", func(r *CreateUserRequest) {
			r.Role = RolePhysician
			r.AdminLevel = AdminLevelUserManager
		}, true},
		{"First name too short", func(r *CreateUserRequest) { r.FirstName = "J" }, true},
		{"Last name too short", func(r *CreateUserRequest) { r.LastName = "D" }, true},
		{"Institution too short", func(r *CreateUserRequest) { r.Institution = "H" }, true},
		{"Location too short", func(r *CreateUserRequest) { r.Location = "C" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the valid request
			req := &CreateUserRequest{
				Username:    validRequest.Username,
				Email:       validRequest.Email,
				Password:    validRequest.Password,
				Role:        validRequest.Role,
				AdminLevel:  validRequest.AdminLevel,
				FirstName:   validRequest.FirstName,
				LastName:    validRequest.LastName,
				Institution: validRequest.Institution,
				Location:    validRequest.Location,
			}

			// Apply modification
			tt.modify(req)

			err := req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUserRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_IsLocked(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)

	tests := []struct {
		name        string
		lockedUntil *time.Time
		want        bool
	}{
		{"Not locked - nil", nil, false},
		{"Not locked - past time", &past, false},
		{"Locked - future time", &future, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				LockedUntil: tt.lockedUntil,
			}
			if got := user.IsLocked(); got != tt.want {
				t.Errorf("User.IsLocked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	tests := []struct {
		name       string
		user       *User
		permission Permission
		want       bool
	}{
		{
			name: "Clinical user has clinical permission",
			user: &User{
				Role: RoleHaematologist,
			},
			permission: PermViewSOPs,
			want:       true,
		},
		{
			name: "Clinical user lacks admin permission",
			user: &User{
				Role: RolePhysician,
			},
			permission: PermManageUsers,
			want:       false,
		},
		{
			name: "Super admin has all permissions",
			user: &User{
				Role:       RoleAdmin,
				AdminLevel: AdminLevelSuperAdmin,
			},
			permission: PermDeleteUsers,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.HasPermission(tt.permission); got != tt.want {
				t.Errorf("User.HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_CanManageUser(t *testing.T) {
	superAdminID := primitive.NewObjectID()
	userManagerID := primitive.NewObjectID()
	clinicalUserID := primitive.NewObjectID()

	superAdmin := &User{
		ID:         superAdminID,
		Role:       RoleAdmin,
		AdminLevel: AdminLevelSuperAdmin,
	}

	userManager := &User{
		ID:         userManagerID,
		Role:       RoleAdmin,
		AdminLevel: AdminLevelUserManager,
	}

	clinicalUser := &User{
		ID:   clinicalUserID,
		Role: RolePhysician,
	}

	targetAdmin := &User{
		Role: RoleAdmin,
	}

	targetClinical := &User{
		Role: RoleHaematologist,
	}

	tests := []struct {
		name       string
		manager    *User
		targetUser *User
		want       bool
	}{
		{"Super admin can manage admin", superAdmin, targetAdmin, true},
		{"Super admin can manage clinical user", superAdmin, targetClinical, true},
		{"User manager can manage clinical user", userManager, targetClinical, true},
		{"User manager cannot manage admin", userManager, targetAdmin, false},
		{"Clinical user cannot manage anyone", clinicalUser, targetClinical, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.manager.CanManageUser(tt.targetUser); got != tt.want {
				t.Errorf("User.CanManageUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_FullName(t *testing.T) {
	user := &User{
		Profile: UserProfile{
			FirstName: "John",
			LastName:  "Doe",
		},
	}

	want := "John Doe"
	if got := user.FullName(); got != want {
		t.Errorf("User.FullName() = %v, want %v", got, want)
	}
}
