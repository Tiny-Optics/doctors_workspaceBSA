package handlers

import (
	"net/http"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

// DropboxAdminHandler handles admin operations for Dropbox configuration
type DropboxAdminHandler struct {
	oauthService *service.DropboxOAuthService
}

// NewDropboxAdminHandler creates a new DropboxAdminHandler
func NewDropboxAdminHandler(oauthService *service.DropboxOAuthService) *DropboxAdminHandler {
	return &DropboxAdminHandler{
		oauthService: oauthService,
	}
}

// GetStatus godoc
// @Summary Get Dropbox connection status
// @Description Get the current status of the Dropbox connection
// @Tags admin-dropbox
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/dropbox/status [get]
// @Security BearerAuth
func (h *DropboxAdminHandler) GetStatus(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only super admins can view Dropbox status
	if !user.HasPermission(models.PermManageSystem) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	status, err := h.oauthService.GetStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// InitiateAuthRequest represents the request to initiate OAuth
type InitiateAuthRequest struct {
	AppKey       string `json:"appKey" binding:"required"`
	AppSecret    string `json:"appSecret" binding:"required"`
	ParentFolder string `json:"parentFolder"` // Optional - empty means Dropbox root
	RedirectURI  string `json:"redirectUri"`  // Optional
}

// InitiateAuth godoc
// @Summary Initiate Dropbox OAuth flow
// @Description Generate the authorization URL for admin to visit
// @Tags admin-dropbox
// @Accept json
// @Produce json
// @Param request body InitiateAuthRequest true "OAuth configuration"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/dropbox/authorize [post]
// @Security BearerAuth
func (h *DropboxAdminHandler) InitiateAuth(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only super admins can configure Dropbox
	if !user.HasPermission(models.PermManageSystem) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	var req InitiateAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store credentials temporarily in session (or you could return them to be sent in callback)
	// For simplicity, we'll include them in the state parameter
	// In production, you might want to store them in Redis or session storage

	authURL, err := h.oauthService.GenerateAuthorizationURL(req.AppKey, req.RedirectURI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authUrl":      authURL,
		"message":      "Visit this URL to authorize Dropbox access",
		"instructions": "After authorizing, you will receive a code. Use the callback endpoint to complete setup.",
	})
}

// CompleteAuthRequest represents the request to complete OAuth
type CompleteAuthRequest struct {
	Code         string `json:"code" binding:"required"`
	AppKey       string `json:"appKey" binding:"required"`
	AppSecret    string `json:"appSecret" binding:"required"`
	ParentFolder string `json:"parentFolder"` // Optional - empty means Dropbox root
	RedirectURI  string `json:"redirectUri"`  // Must match the one used in authorization
}

// CompleteAuth godoc
// @Summary Complete Dropbox OAuth flow
// @Description Exchange authorization code for tokens and save configuration
// @Tags admin-dropbox
// @Accept json
// @Produce json
// @Param request body CompleteAuthRequest true "Authorization code and configuration"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/dropbox/callback [post]
// @Security BearerAuth
func (h *DropboxAdminHandler) CompleteAuth(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only super admins can configure Dropbox
	if !user.HasPermission(models.PermManageSystem) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	var req CompleteAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	config, err := h.oauthService.ExchangeCodeForTokens(
		c.Request.Context(),
		req.Code,
		req.AppKey,
		req.AppSecret,
		req.RedirectURI,
		req.ParentFolder,
		user,
		ipAddress,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dropbox successfully authorized and configured",
		"status":  config.GetPublicStatus(),
	})
}

// ForceRefresh godoc
// @Summary Force refresh Dropbox access token
// @Description Manually trigger a token refresh
// @Tags admin-dropbox
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/dropbox/refresh [post]
// @Security BearerAuth
func (h *DropboxAdminHandler) ForceRefresh(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only super admins can refresh tokens
	if !user.HasPermission(models.PermManageSystem) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	if err := h.oauthService.ForceRefresh(c.Request.Context(), user, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refresh successful",
	})
}

// TestConnection godoc
// @Summary Test Dropbox connection
// @Description Test if the Dropbox connection is working
// @Tags admin-dropbox
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/dropbox/test [post]
// @Security BearerAuth
func (h *DropboxAdminHandler) TestConnection(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only super admins can test connection
	if !user.HasPermission(models.PermManageSystem) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	if err := h.oauthService.TestConnection(c.Request.Context(), user, ipAddress); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Dropbox connection is working correctly",
	})
}

// DeleteConfiguration godoc
// @Summary Delete Dropbox configuration
// @Description Delete the Dropbox configuration (requires re-authorization)
// @Tags admin-dropbox
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/dropbox/configuration [delete]
// @Security BearerAuth
func (h *DropboxAdminHandler) DeleteConfiguration(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only super admins can delete configuration
	if !user.HasPermission(models.PermManageSystem) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	if err := h.oauthService.DeleteConfiguration(c.Request.Context(), user, ipAddress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dropbox configuration deleted successfully",
	})
}
