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

// InstitutionHandler handles institution management requests
type InstitutionHandler struct {
	institutionService *service.InstitutionService
}

// NewInstitutionHandler creates a new InstitutionHandler
func NewInstitutionHandler(institutionService *service.InstitutionService) *InstitutionHandler {
	return &InstitutionHandler{
		institutionService: institutionService,
	}
}

// CreateInstitution godoc
// @Summary Create a new institution
// @Description Create a new institution (requires manage users permission)
// @Tags institutions
// @Accept json
// @Produce json
// @Param request body models.CreateInstitutionRequest true "Institution information"
// @Success 201 {object} models.Institution
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /institutions [post]
// @Security BearerAuth
func (h *InstitutionHandler) CreateInstitution(c *gin.Context) {
	var req models.CreateInstitutionRequest
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

	institution, err := h.institutionService.CreateInstitution(c.Request.Context(), &req, createdBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrDuplicateInstitution {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, institution)
}

// GetInstitution godoc
// @Summary Get an institution by ID
// @Description Get institution information by ID
// @Tags institutions
// @Produce json
// @Param id path string true "Institution ID"
// @Success 200 {object} models.Institution
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /institutions/{id} [get]
// @Security BearerAuth
func (h *InstitutionHandler) GetInstitution(c *gin.Context) {
	idParam := c.Param("id")
	institutionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid institution ID"})
		return
	}

	institution, err := h.institutionService.GetInstitution(c.Request.Context(), institutionID)
	if err != nil {
		if err == repository.ErrInstitutionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "institution not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, institution)
}

// UpdateInstitution godoc
// @Summary Update an institution
// @Description Update institution information (requires manage users permission)
// @Tags institutions
// @Accept json
// @Produce json
// @Param id path string true "Institution ID"
// @Param request body models.UpdateInstitutionRequest true "Updated institution information"
// @Success 200 {object} models.Institution
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /institutions/{id} [put]
// @Security BearerAuth
func (h *InstitutionHandler) UpdateInstitution(c *gin.Context) {
	idParam := c.Param("id")
	institutionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid institution ID"})
		return
	}

	var req models.UpdateInstitutionRequest
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

	institution, err := h.institutionService.UpdateInstitution(c.Request.Context(), institutionID, &req, updatedBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrInstitutionNotFound {
			statusCode = http.StatusNotFound
		} else if err == repository.ErrDuplicateInstitution {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, institution)
}

// DeleteInstitution godoc
// @Summary Delete an institution
// @Description Delete an institution (requires delete users permission)
// @Tags institutions
// @Produce json
// @Param id path string true "Institution ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /institutions/{id} [delete]
// @Security BearerAuth
func (h *InstitutionHandler) DeleteInstitution(c *gin.Context) {
	idParam := c.Param("id")
	institutionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid institution ID"})
		return
	}

	deletedBy, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	err = h.institutionService.DeleteInstitution(c.Request.Context(), institutionID, deletedBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrInstitutionNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "institution deleted successfully"})
}

// ActivateInstitution godoc
// @Summary Activate an institution
// @Description Activate an institution (requires manage users permission)
// @Tags institutions
// @Produce json
// @Param id path string true "Institution ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /institutions/{id}/activate [post]
// @Security BearerAuth
func (h *InstitutionHandler) ActivateInstitution(c *gin.Context) {
	idParam := c.Param("id")
	institutionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid institution ID"})
		return
	}

	activatedBy, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	err = h.institutionService.ActivateInstitution(c.Request.Context(), institutionID, activatedBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrInstitutionNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "institution activated successfully"})
}

// DeactivateInstitution godoc
// @Summary Deactivate an institution
// @Description Deactivate an institution (requires manage users permission)
// @Tags institutions
// @Produce json
// @Param id path string true "Institution ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /institutions/{id}/deactivate [post]
// @Security BearerAuth
func (h *InstitutionHandler) DeactivateInstitution(c *gin.Context) {
	idParam := c.Param("id")
	institutionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid institution ID"})
		return
	}

	deactivatedBy, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	err = h.institutionService.DeactivateInstitution(c.Request.Context(), institutionID, deactivatedBy, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == repository.ErrInstitutionNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "institution deactivated successfully"})
}

// ListInstitutions godoc
// @Summary List institutions
// @Description Get a list of institutions with optional filtering and pagination
// @Tags institutions
// @Produce json
// @Param type query string false "Filter by institution type"
// @Param is_active query bool false "Filter by active status"
// @Param search query string false "Search by name, city, province, or country"
// @Param limit query int false "Limit number of results" default(100)
// @Param skip query int false "Skip number of results" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /institutions [get]
// @Security BearerAuth
func (h *InstitutionHandler) ListInstitutions(c *gin.Context) {
	// Parse query parameters
	var institutionType *models.InstitutionType
	if typeParam := c.Query("type"); typeParam != "" {
		t := models.InstitutionType(typeParam)
		if t.IsValid() {
			institutionType = &t
		}
	}

	var isActive *bool
	if isActiveParam := c.Query("is_active"); isActiveParam != "" {
		active := isActiveParam == "true"
		isActive = &active
	}

	// Get search query parameter
	search := c.Query("search")

	limit := int64(100)
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

	institutions, count, err := h.institutionService.ListInstitutions(c.Request.Context(), institutionType, isActive, search, limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"institutions": institutions,
		"total":        count,
		"limit":        limit,
		"skip":         skip,
	})
}

