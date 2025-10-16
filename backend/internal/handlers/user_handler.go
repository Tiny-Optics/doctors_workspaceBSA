package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler handles user management requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user (requires manage users permission)
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User information"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /users [post]
// @Security BearerAuth
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBy, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	user, err := h.userService.CreateUser(c.Request.Context(), &req, createdBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrDuplicateEmail || err == repository.ErrDuplicateUsername {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body models.UpdateUserRequest true "Updated user information"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
// @Security BearerAuth
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedBy, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, &req, updatedBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized || err == service.ErrCannotModifyOwnAdminLevel {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Permanently delete a user (requires delete users permission)
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
// @Security BearerAuth
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	deletedBy, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	err = h.userService.DeleteUser(c.Request.Context(), userID, deletedBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// DeactivateUser godoc
// @Summary Deactivate a user
// @Description Deactivate a user account (soft delete)
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id}/deactivate [post]
// @Security BearerAuth
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	deactivatedBy, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	err = h.userService.DeactivateUser(c.Request.Context(), userID, deactivatedBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deactivated successfully"})
}

// ActivateUser godoc
// @Summary Activate a user
// @Description Activate a user account
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id}/activate [post]
// @Security BearerAuth
func (h *UserHandler) ActivateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	activatedBy, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	err = h.userService.ActivateUser(c.Request.Context(), userID, activatedBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrUserNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user activated successfully"})
}

// ListUsers godoc
// @Summary List users
// @Description Get a list of users with optional filtering, searching and pagination
// @Tags users
// @Produce json
// @Param role query string false "Filter by role"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name, email, role, or institution"
// @Param limit query int false "Limit number of results" default(20)
// @Param skip query int false "Skip number of results" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users [get]
// @Security BearerAuth
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse query parameters
	var role *models.UserRole
	if roleParam := c.Query("role"); roleParam != "" {
		r := models.UserRole(roleParam)
		if r.IsValid() {
			role = &r
		}
	}

	var isActive *bool
	if isActiveParam := c.Query("is_active"); isActiveParam != "" {
		active := isActiveParam == "true"
		isActive = &active
	}

	// Get search query parameter
	search := c.Query("search")

	limit := int64(20)
	if limitParam := c.Query("limit"); limitParam != "" {
		if l, err := strconv.ParseInt(limitParam, 10, 64); err == nil && l > 0 {
			limit = l
		}
	}

	skip := int64(0)
	if skipParam := c.Query("skip"); skipParam != "" {
		if s, err := strconv.ParseInt(skipParam, 10, 64); err == nil && s >= 0 {
			skip = s
		}
	}

	users, count, err := h.userService.ListUsers(c.Request.Context(), role, isActive, search, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": count,
		"limit": limit,
		"skip":  skip,
	})
}
