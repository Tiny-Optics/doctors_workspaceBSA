package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token and get user
		user, err := authService.GetUserFromToken(context.Background(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "account is not active"})
			c.Abort()
			return
		}

		// Store user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("token", token)

		c.Next()
	}
}

// GetUserFromContext retrieves the user from the Gin context
func GetUserFromContext(c *gin.Context) (*models.User, error) {
	userInterface, exists := c.Get("user")
	if !exists {
		return nil, http.ErrNotSupported
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		return nil, http.ErrNotSupported
	}

	return user, nil
}

// GetUserIDFromContext retrieves the user ID from the Gin context
func GetUserIDFromContext(c *gin.Context) (primitive.ObjectID, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return primitive.NilObjectID, http.ErrNotSupported
	}

	userID, ok := userIDInterface.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, http.ErrNotSupported
	}

	return userID, nil
}

// GetTokenFromContext retrieves the token from the Gin context
func GetTokenFromContext(c *gin.Context) (string, error) {
	tokenInterface, exists := c.Get("token")
	if !exists {
		return "", http.ErrNotSupported
	}

	token, ok := tokenInterface.(string)
	if !ok {
		return "", http.ErrNotSupported
	}

	return token, nil
}

// GetIPAddress extracts the client IP address from the request
func GetIPAddress(c *gin.Context) string {
	// Check for X-Forwarded-For header first (proxy/load balancer)
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}

	// Check for X-Real-IP header
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to remote address
	return c.ClientIP()
}
