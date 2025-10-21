package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegistryHandler handles registry-related requests
type RegistryHandler struct {
	registryService   *service.RegistryService
	encryptionService *service.EncryptionService
}

// NewRegistryHandler creates a new RegistryHandler
func NewRegistryHandler(registryService *service.RegistryService, encryptionService *service.EncryptionService) *RegistryHandler {
	return &RegistryHandler{
		registryService:   registryService,
		encryptionService: encryptionService,
	}
}

// Configuration Endpoints

// GetConfiguration godoc
// @Summary Get registry configuration
// @Description Get full registry configuration (admin only)
// @Tags registry
// @Produce json
// @Success 200 {object} models.RegistryConfig
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/registry/config [get]
// @Security BearerAuth
func (h *RegistryHandler) GetConfiguration(c *gin.Context) {
	config, err := h.registryService.GetConfiguration(c.Request.Context())
	if err != nil {
		if err == repository.ErrRegistryConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// GetPublicConfiguration godoc
// @Summary Get public registry configuration
// @Description Get public configuration (video URL and documents path only)
// @Tags registry
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /registry/config [get]
// @Security BearerAuth
func (h *RegistryHandler) GetPublicConfiguration(c *gin.Context) {
	config, err := h.registryService.GetPublicConfiguration(c.Request.Context())
	if err != nil {
		if err == repository.ErrRegistryConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateConfiguration godoc
// @Summary Update registry configuration
// @Description Update registry configuration (admin only)
// @Tags registry
// @Accept json
// @Produce json
// @Param request body models.UpdateRegistryConfigRequest true "Configuration updates"
// @Success 200 {object} models.RegistryConfig
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/registry/config [put]
// @Security BearerAuth
func (h *RegistryHandler) UpdateConfiguration(c *gin.Context) {
	var req models.UpdateRegistryConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	config, err := h.registryService.UpdateConfiguration(c.Request.Context(), &req, user, h.encryptionService, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorizedRegistryAccess {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// SendTestEmail godoc
// @Summary Send test email
// @Description Send a test email to verify SMTP configuration
// @Tags registry
// @Accept json
// @Produce json
// @Param request body map[string]string true "Test email request with 'email' field"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/registry/test-email [post]
// @Security BearerAuth
func (h *RegistryHandler) SendTestEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email address"})
		return
	}

	err := h.registryService.SendTestEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "test email sent successfully"})
}

// Form Schema Endpoints

// CreateFormSchema godoc
// @Summary Create form schema
// @Description Create a new form schema (admin only)
// @Tags registry
// @Accept json
// @Produce json
// @Param request body models.CreateFormSchemaRequest true "Form schema"
// @Success 201 {object} models.RegistryFormSchema
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/registry/form-schema [post]
// @Security BearerAuth
func (h *RegistryHandler) CreateFormSchema(c *gin.Context) {
	var req models.CreateFormSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	schema, err := h.registryService.CreateFormSchema(c.Request.Context(), &req, user, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorizedRegistryAccess {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schema)
}

// GetFormSchema godoc
// @Summary Get form schema by ID
// @Description Get a specific form schema (admin only)
// @Tags registry
// @Produce json
// @Param id path string true "Form Schema ID"
// @Success 200 {object} models.RegistryFormSchema
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/registry/form-schema/{id} [get]
// @Security BearerAuth
func (h *RegistryHandler) GetFormSchema(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form schema ID"})
		return
	}

	schema, err := h.registryService.GetFormSchema(c.Request.Context(), id)
	if err != nil {
		if err == repository.ErrFormSchemaNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "form schema not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schema)
}

// GetActiveFormSchema godoc
// @Summary Get active form schema
// @Description Get the currently active form schema
// @Tags registry
// @Produce json
// @Success 200 {object} models.RegistryFormSchema
// @Failure 404 {object} map[string]string
// @Router /registry/form-schema [get]
// @Security BearerAuth
func (h *RegistryHandler) GetActiveFormSchema(c *gin.Context) {
	schema, err := h.registryService.GetActiveFormSchema(c.Request.Context())
	if err != nil {
		if err == repository.ErrNoActiveForm {
			c.JSON(http.StatusNotFound, gin.H{"error": "no active form schema"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schema)
}

// ListFormSchemas godoc
// @Summary List all form schemas
// @Description List all form schemas with pagination (admin only)
// @Tags registry
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /admin/registry/form-schemas [get]
// @Security BearerAuth
func (h *RegistryHandler) ListFormSchemas(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	schemas, total, err := h.registryService.ListFormSchemas(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schemas": schemas,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// UpdateFormSchema godoc
// @Summary Update form schema
// @Description Update an existing form schema (admin only)
// @Tags registry
// @Accept json
// @Produce json
// @Param id path string true "Form Schema ID"
// @Param request body models.UpdateFormSchemaRequest true "Form schema updates"
// @Success 200 {object} models.RegistryFormSchema
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/registry/form-schema/{id} [put]
// @Security BearerAuth
func (h *RegistryHandler) UpdateFormSchema(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form schema ID"})
		return
	}

	var req models.UpdateFormSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	schema, err := h.registryService.UpdateFormSchema(c.Request.Context(), id, &req, user, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorizedRegistryAccess {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrFormSchemaNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schema)
}

// DeleteFormSchema godoc
// @Summary Delete form schema
// @Description Delete a form schema (admin only)
// @Tags registry
// @Produce json
// @Param id path string true "Form Schema ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/registry/form-schema/{id} [delete]
// @Security BearerAuth
func (h *RegistryHandler) DeleteFormSchema(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form schema ID"})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	if err := h.registryService.DeleteFormSchema(c.Request.Context(), id, user, ipAddress); err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorizedRegistryAccess {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrFormSchemaNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Submission Endpoints

// SubmitForm godoc
// @Summary Submit form
// @Description Submit a registry form with document uploads
// @Tags registry
// @Accept multipart/form-data
// @Produce json
// @Param formData formData string true "Form data as JSON string"
// @Param formSchemaId formData string true "Form Schema ID"
// @Param documents formData file true "Documents to upload" collectionFormat(multi)
// @Success 201 {object} models.RegistrySubmission
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /registry/submit [post]
// @Security BearerAuth
func (h *RegistryHandler) SubmitForm(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	// Get form data
	formDataStr := c.PostForm("formData")
	if formDataStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formData is required"})
		return
	}

	formSchemaID := c.PostForm("formSchemaId")
	if formSchemaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formSchemaId is required"})
		return
	}

	// Parse form data JSON from the string
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(formDataStr), &formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse form data: %v", err)})
		return
	}

	// Get files (accept both "files" and "documents" field names)
	form := c.Request.MultipartForm
	files := form.File["files"]
	if len(files) == 0 {
		// Try alternative field name for backwards compatibility
		files = form.File["documents"]
	}
	
	// Note: We don't require files here - the service layer will validate
	// based on the form schema's required file fields

	// Create submission request
	req := &models.CreateSubmissionRequest{
		FormSchemaID: formSchemaID,
		FormData:     formData,
	}

	ipAddress := middleware.GetIPAddress(c)

	submission, err := h.registryService.SubmitForm(c.Request.Context(), req, files, user, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrFormNotReady || err == service.ErrNoActiveFormSchema {
			statusCode = http.StatusPreconditionFailed
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, submission)
}

// GetUserSubmissions godoc
// @Summary Get user's submissions
// @Description Get all submissions for the authenticated user
// @Tags registry
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /registry/submissions [get]
// @Security BearerAuth
func (h *RegistryHandler) GetUserSubmissions(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	submissions, total, err := h.registryService.GetUserSubmissions(c.Request.Context(), user.ID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"submissions": submissions,
		"total":       total,
		"page":        page,
		"limit":       limit,
	})
}

// GetAllSubmissions godoc
// @Summary Get all submissions
// @Description Get all submissions with filters (admin only)
// @Tags registry
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/registry/submissions [get]
// @Security BearerAuth
func (h *RegistryHandler) GetAllSubmissions(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Build filter
	filter := bson.M{}
	if status := c.Query("status"); status != "" {
		filter["status"] = status
	}

	submissions, total, err := h.registryService.GetAllSubmissions(c.Request.Context(), user, page, limit, filter)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrUnauthorizedRegistryAccess {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"submissions": submissions,
		"total":       total,
		"page":        page,
		"limit":       limit,
	})
}

// GetSubmission godoc
// @Summary Get submission by ID
// @Description Get a specific submission
// @Tags registry
// @Produce json
// @Param id path string true "Submission ID"
// @Success 200 {object} models.RegistrySubmission
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /registry/submissions/{id} [get]
// @Security BearerAuth
func (h *RegistryHandler) GetSubmission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid submission ID"})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	submission, err := h.registryService.GetSubmission(c.Request.Context(), id, user)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == repository.ErrSubmissionNotFound {
			statusCode = http.StatusNotFound
		} else if err.Error() == "unauthorized to view this submission" {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, submission)
}

// UpdateSubmissionStatus godoc
// @Summary Update submission status
// @Description Update the status of a submission (admin only)
// @Tags registry
// @Accept json
// @Produce json
// @Param id path string true "Submission ID"
// @Param request body models.UpdateSubmissionStatusRequest true "Status update"
// @Success 200 {object} models.RegistrySubmission
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/registry/submissions/{id}/status [patch]
// @Security BearerAuth
func (h *RegistryHandler) UpdateSubmissionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid submission ID"})
		return
	}

	var req models.UpdateSubmissionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	submission, err := h.registryService.UpdateSubmissionStatus(c.Request.Context(), id, &req, user, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorizedRegistryAccess {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrSubmissionNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, submission)
}

// GetExampleDocuments godoc
// @Summary List example documents
// @Description List all files in the example documents directory
// @Tags registry
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /registry/example-documents [get]
// @Security BearerAuth
func (h *RegistryHandler) GetExampleDocuments(c *gin.Context) {
	files, err := h.registryService.GetExampleDocuments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

// GetExampleDocumentDownloadLink godoc
// @Summary Get download link for example document
// @Description Get a temporary download link for a specific example document
// @Tags registry
// @Produce json
// @Param path query string true "File path"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /registry/example-documents/download [get]
// @Security BearerAuth
func (h *RegistryHandler) GetExampleDocumentDownloadLink(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file path is required"})
		return
	}

	link, err := h.registryService.GetExampleDocumentDownloadLink(c.Request.Context(), filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"link": link,
	})
}

// GetSubmissionDocumentDownloadLink godoc
// @Summary Get download link for submission document
// @Description Get a temporary download link for a specific submission document
// @Tags registry
// @Produce json
// @Param path query string true "Document path"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /registry/document-download [get]
// @Security BearerAuth
func (h *RegistryHandler) GetSubmissionDocumentDownloadLink(c *gin.Context) {
	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file path is required"})
		return
	}

	link, err := h.registryService.GetExampleDocumentDownloadLink(c.Request.Context(), filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"link": link,
	})
}
