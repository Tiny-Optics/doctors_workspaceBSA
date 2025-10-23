package handlers

import (
	"net/http"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

// PasswordResetHandler handles password reset requests
type PasswordResetHandler struct {
	passwordResetService *service.PasswordResetService
}

// NewPasswordResetHandler creates a new PasswordResetHandler
func NewPasswordResetHandler(passwordResetService *service.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetService: passwordResetService,
	}
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send a password reset code to user's email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ForgotPasswordRequest true "Email address"
// @Success 200 {object} models.ForgotPasswordResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /auth/forgot-password [post]
func (h *PasswordResetHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	response, err := h.passwordResetService.RequestPasswordReset(c.Request.Context(), req.Email, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUserNotFound {
			statusCode = http.StatusNotFound
		} else if err == service.ErrTooManyRequests {
			statusCode = http.StatusTooManyRequests
		} else if err == service.ErrUserInactive {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ValidateResetCode godoc
// @Summary Validate password reset code
// @Description Validate the 6-digit verification code and return a token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ValidateResetCodeRequest true "Verification code"
// @Success 200 {object} models.ValidateResetCodeResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /auth/validate-reset-code [post]
func (h *PasswordResetHandler) ValidateResetCode(c *gin.Context) {
	var req models.ValidateResetCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.passwordResetService.ValidateResetCode(c.Request.Context(), req.Code)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrInvalidResetCode || err == service.ErrPasswordResetTokenExpired || err == service.ErrResetTokenUsed {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ResetPassword godoc
// @Summary Reset password with token
// @Description Set new password using the token from code validation
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.ResetPasswordRequest true "Token and new password"
// @Success 200 {object} models.ResetPasswordResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.passwordResetService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrInvalidResetCode || err == service.ErrPasswordResetTokenExpired || err == service.ErrResetTokenUsed {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
