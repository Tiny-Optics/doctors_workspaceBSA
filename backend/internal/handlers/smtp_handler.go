package handlers

import (
	"net/http"

	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

// SMTPHandler handles SMTP-related requests
type SMTPHandler struct {
	registryService *service.RegistryService
}

// NewSMTPHandler creates a new SMTPHandler
func NewSMTPHandler(registryService *service.RegistryService) *SMTPHandler {
	return &SMTPHandler{
		registryService: registryService,
	}
}

// CheckSMTPConfiguration godoc
// @Summary Check if SMTP is configured (public)
// @Description Check if SMTP configuration is complete for password reset functionality
// @Tags smtp
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /smtp/status [get]
func (h *SMTPHandler) CheckSMTPConfiguration(c *gin.Context) {
	config, err := h.registryService.GetPublicSMTPConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"isConfigured": false,
			"message":      "SMTP not configured",
		})
		return
	}

	// Check if SMTP configuration is complete
	isComplete := config.IsComplete()

	c.JSON(http.StatusOK, gin.H{
		"isConfigured": isComplete,
		"message":      "SMTP configuration status",
	})
}
