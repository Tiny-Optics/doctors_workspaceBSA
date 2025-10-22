package handlers

import (
	"net/http"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

// ReferralHandler handles referral-related requests
type ReferralHandler struct {
	referralService *service.ReferralService
}

// NewReferralHandler creates a new ReferralHandler
func NewReferralHandler(referralService *service.ReferralService) *ReferralHandler {
	return &ReferralHandler{
		referralService: referralService,
	}
}

// GetConfig godoc
// @Summary Get referral configuration (user view)
// @Description Get referral configuration for authenticated users (URL and enabled status only)
// @Tags referrals
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /referrals/config [get]
// @Security BearerAuth
func (h *ReferralHandler) GetConfig(c *gin.Context) {
	config, err := h.referralService.GetReferralConfig(c.Request.Context())
	if err != nil {
		if err == service.ErrReferralNotConfigured {
			c.JSON(http.StatusOK, gin.H{
				"isConfigured": false,
				"isEnabled":    false,
				"redcapUrl":    "",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return limited info for regular users
	c.JSON(http.StatusOK, gin.H{
		"isConfigured": config.IsConfigured,
		"isEnabled":    config.IsEnabled,
		"redcapUrl":    config.RedCapURL,
	})
}

// GetAdminConfig godoc
// @Summary Get full referral configuration (admin)
// @Description Get complete referral configuration including metadata (admin only)
// @Tags referrals
// @Produce json
// @Success 200 {object} models.ReferralConfig
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/referrals/config [get]
// @Security BearerAuth
func (h *ReferralHandler) GetAdminConfig(c *gin.Context) {
	config, err := h.referralService.GetReferralConfig(c.Request.Context())
	if err != nil {
		if err == service.ErrReferralNotConfigured {
			c.JSON(http.StatusNotFound, gin.H{"error": "referral system not configured"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateConfig godoc
// @Summary Update referral configuration
// @Description Update REDCap link and enabled status (admin only)
// @Tags referrals
// @Accept json
// @Produce json
// @Param request body models.UpdateReferralConfigRequest true "Configuration update"
// @Success 200 {object} models.ReferralConfig
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/referrals/config [put]
// @Security BearerAuth
func (h *ReferralHandler) UpdateConfig(c *gin.Context) {
	var req models.UpdateReferralConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get IP address
	ipAddress := c.ClientIP()

	// Update configuration
	config, err := h.referralService.UpdateReferralConfig(c.Request.Context(), &req, userID, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// LogAccess godoc
// @Summary Log referral access
// @Description Log when a user accesses the referral link and return the URL
// @Tags referrals
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /referrals/access [post]
// @Security BearerAuth
func (h *ReferralHandler) LogAccess(c *gin.Context) {
	// Get user ID from context
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get IP address and user agent
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	// Log the access
	if err := h.referralService.LogReferralAccess(c.Request.Context(), userID, ipAddress, userAgent); err != nil {
		if err == service.ErrReferralNotConfigured {
			c.JSON(http.StatusNotFound, gin.H{"error": "referral system not configured"})
			return
		}
		if err == service.ErrReferralDisabled {
			c.JSON(http.StatusForbidden, gin.H{"error": "referral system is currently disabled"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the URL
	url, err := h.referralService.GetReferralURL(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"redirectUrl": url})
}
