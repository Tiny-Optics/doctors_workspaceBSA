package handlers

import (
	"net/http"
	"time"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

// StatsHandler handles admin statistics requests
type StatsHandler struct {
	userService        *service.UserService
	institutionService *service.InstitutionService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(userService *service.UserService, institutionService *service.InstitutionService) *StatsHandler {
	return &StatsHandler{
		userService:        userService,
		institutionService: institutionService,
	}
}

// AdminStatsResponse represents the admin dashboard statistics
type AdminStatsResponse struct {
	TotalUsers         int64              `json:"totalUsers"`
	ActiveUsers        int64              `json:"activeUsers"`
	NewUsersThisMonth  int64              `json:"newUsersThisMonth"`
	NewUsersThisWeek   int64              `json:"newUsersThisWeek"`
	NewUsersToday      int64              `json:"newUsersToday"`
	TotalInstitutions  int64              `json:"totalInstitutions"`
	RoleDistribution   []RoleDistribution `json:"roleDistribution"`
}

// RoleDistribution represents user count by role
type RoleDistribution struct {
	Role  string `json:"role"`
	Count int64  `json:"count"`
}

// GetAdminStats godoc
// @Summary Get admin dashboard statistics
// @Description Get comprehensive statistics for admin dashboard
// @Tags stats
// @Produce json
// @Success 200 {object} AdminStatsResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /stats/admin [get]
// @Security BearerAuth
func (h *StatsHandler) GetAdminStats(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only admins can access this endpoint
	if !user.HasPermission(models.PermManageUsers) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	ctx := c.Request.Context()

	// Get total users
	totalUsers, err := h.userService.CountUsers(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get total users"})
		return
	}

	// Get active users
	isActive := true
	activeUsers, err := h.userService.CountUsers(ctx, &isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get active users"})
		return
	}

	// Get total institutions
	totalInstitutions, err := h.institutionService.CountInstitutions(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get total institutions"})
		return
	}

	// Get new users this month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	newUsersThisMonth, err := h.userService.CountUsersCreatedAfter(ctx, startOfMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get new users this month"})
		return
	}

	// Get new users this week
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	newUsersThisWeek, err := h.userService.CountUsersCreatedAfter(ctx, startOfWeek)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get new users this week"})
		return
	}

	// Get new users today
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	newUsersToday, err := h.userService.CountUsersCreatedAfter(ctx, startOfDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get new users today"})
		return
	}

	// Get role distribution
	roleDistribution, err := h.userService.GetRoleDistribution(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get role distribution"})
		return
	}

	// Convert role distribution to response format
	roleDistResp := make([]RoleDistribution, 0, len(roleDistribution))
	for role, count := range roleDistribution {
		roleDistResp = append(roleDistResp, RoleDistribution{
			Role:  string(role),
			Count: count,
		})
	}

	response := AdminStatsResponse{
		TotalUsers:         totalUsers,
		ActiveUsers:        activeUsers,
		NewUsersThisMonth:  newUsersThisMonth,
		NewUsersThisWeek:   newUsersThisWeek,
		NewUsersToday:      newUsersToday,
		TotalInstitutions:  totalInstitutions,
		RoleDistribution:   roleDistResp,
	}

	c.JSON(http.StatusOK, response)
}

