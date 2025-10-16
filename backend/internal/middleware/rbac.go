package middleware

import (
	"net/http"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
)

// RequirePermission creates a middleware that checks if the user has a specific permission
func RequirePermission(permission models.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if !user.HasPermission(permission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission creates a middleware that checks if the user has any of the specified permissions
func RequireAnyPermission(permissions ...models.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		hasAny := false
		for _, permission := range permissions {
			if user.HasPermission(permission) {
				hasAny = true
				break
			}
		}

		if !hasAny {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions creates a middleware that checks if the user has all of the specified permissions
func RequireAllPermissions(permissions ...models.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		for _, permission := range permissions {
			if !user.HasPermission(permission) {
				c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireRole creates a middleware that checks if the user has a specific role
func RequireRole(role models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if user.Role != role {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole creates a middleware that checks if the user has any of the specified roles
func RequireAnyRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if user.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin creates a middleware that checks if the user is an admin
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}

// RequireAdminLevel creates a middleware that checks if the user has a specific admin level
func RequireAdminLevel(adminLevel models.AdminLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		if user.Role != models.RoleAdmin || user.AdminLevel != adminLevel {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin creates a middleware that checks if the user is a super admin
func RequireSuperAdmin() gin.HandlerFunc {
	return RequireAdminLevel(models.AdminLevelSuperAdmin)
}
