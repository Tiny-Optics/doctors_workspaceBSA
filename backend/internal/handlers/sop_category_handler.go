package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SOPCategoryHandler handles SOP category requests
type SOPCategoryHandler struct {
	categoryService *service.SOPCategoryService
}

// NewSOPCategoryHandler creates a new SOPCategoryHandler
func NewSOPCategoryHandler(categoryService *service.SOPCategoryService) *SOPCategoryHandler {
	return &SOPCategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create a new SOP category
// @Description Create a new SOP category (requires super admin permission)
// @Tags sop-categories
// @Accept json
// @Produce json
// @Param request body models.CreateSOPCategoryRequest true "Category information"
// @Success 201 {object} models.SOPCategory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /sops/categories [post]
// @Security BearerAuth
func (h *SOPCategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateSOPCategoryRequest
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

	category, err := h.categoryService.CreateCategory(c.Request.Context(), &req, user, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == service.ErrDuplicateSlug || err == repository.ErrDuplicateSlug {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategory godoc
// @Summary Get a category by ID
// @Description Get SOP category information by ID
// @Tags sop-categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.SOPCategory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /sops/categories/:id [get]
// @Security BearerAuth
func (h *SOPCategoryHandler) GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	category, err := h.categoryService.GetCategory(c.Request.Context(), id, user)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrCategoryNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// ListCategories godoc
// @Summary List all categories
// @Description List SOP categories with pagination and filters
// @Tags sop-categories
// @Produce json
// @Param search query string false "Search term"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /sops/categories [get]
// @Security BearerAuth
func (h *SOPCategoryHandler) ListCategories(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Parse query parameters
	search := c.Query("search")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}

	categories, total, err := h.categoryService.ListCategories(
		c.Request.Context(),
		user,
		search,
		page,
		limit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an SOP category (requires super admin permission)
// @Tags sop-categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body models.UpdateSOPCategoryRequest true "Update information"
// @Success 200 {object} models.SOPCategory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /sops/categories/:id [put]
// @Security BearerAuth
func (h *SOPCategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req models.UpdateSOPCategoryRequest
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

	category, err := h.categoryService.UpdateCategory(c.Request.Context(), id, &req, user, ipAddress)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == service.ErrCategoryNotFound {
			statusCode = http.StatusNotFound
		} else if err == service.ErrDuplicateSlug || err == repository.ErrDuplicateSlug {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete an SOP category from database (Dropbox folder remains) (requires super admin permission)
// @Tags sop-categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /sops/categories/:id [delete]
// @Security BearerAuth
func (h *SOPCategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	err = h.categoryService.DeleteCategory(c.Request.Context(), id, user, ipAddress)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrUnauthorized {
			statusCode = http.StatusForbidden
		} else if err == service.ErrCategoryNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

// GetCategoryFiles godoc
// @Summary List files in a category
// @Description List all files in a category's Dropbox folder
// @Tags sop-categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /sops/categories/:id/files [get]
// @Security BearerAuth
func (h *SOPCategoryHandler) GetCategoryFiles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	files, err := h.categoryService.GetCategoryFiles(c.Request.Context(), id, user)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrCategoryNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

// DownloadFile godoc
// @Summary Get download link for a file
// @Description Get a temporary download link for a specific file in a category
// @Tags sop-categories
// @Produce json
// @Param id path string true "Category ID"
// @Param path query string true "File path within the category folder"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /sops/categories/:id/files/download [get]
// @Security BearerAuth
func (h *SOPCategoryHandler) DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	filePath := c.Query("path")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file path is required"})
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	downloadLink, err := h.categoryService.GetFileDownloadLink(c.Request.Context(), id, filePath, user)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrCategoryNotFound {
			statusCode = http.StatusNotFound
		} else if err.Error() == "file not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"downloadLink": downloadLink,
	})
}

// UploadImage godoc
// @Summary Upload category image
// @Description Upload an image for SOP category (requires super admin permission)
// @Tags sop-categories
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file (jpg, jpeg, png, webp, max 5MB)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /sops/images/upload [post]
// @Security BearerAuth
func (h *SOPCategoryHandler) UploadImage(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check permission - only super admins can upload images
	if !user.HasPermission(models.PermDeleteUsers) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	// Get file from form
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image size must be less than 5MB"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !validExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image format. Allowed: jpg, jpeg, png, webp"})
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads/sops"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), ext)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	// Return relative path for database storage
	imagePath := "/uploads/sops/" + filename

	c.JSON(http.StatusOK, gin.H{
		"imagePath": imagePath,
		"message":   "image uploaded successfully",
	})
}

// SeedCategories godoc
// @Summary Seed initial SOP categories
// @Description Create initial SOP categories if none exist (requires super admin permission)
// @Tags sop-categories
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /sops/seed [post]
// @Security BearerAuth
func (h *SOPCategoryHandler) SeedCategories(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Check permission - only super admins can seed
	if !user.HasPermission(models.PermDeleteUsers) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	// Check if categories already exist
	categories, total, err := h.categoryService.ListCategories(c.Request.Context(), user, "", 1, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing categories"})
		return
	}

	if total > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":    "categories already exist",
			"count":      total,
			"categories": categories,
		})
		return
	}

	ipAddress := middleware.GetIPAddress(c)

	// Seed initial categories
	seedData := []struct {
		name         string
		description  string
		displayOrder int
	}{
		{
			name:         "Anemia",
			description:  "Standard operating procedures for anemia diagnosis and treatment",
			displayOrder: 1,
		},
		{
			name:         "Lymphoma",
			description:  "Standard operating procedures for lymphoma management",
			displayOrder: 2,
		},
		{
			name:         "Myeloma",
			description:  "Standard operating procedures for multiple myeloma treatment",
			displayOrder: 3,
		},
		{
			name:         "General Business",
			description:  "General business procedures and administrative guidelines",
			displayOrder: 4,
		},
	}

	createdCategories := make([]interface{}, 0)

	for _, data := range seedData {
		req := &models.CreateSOPCategoryRequest{
			Name:         data.name,
			Description:  data.description,
			DisplayOrder: data.displayOrder,
		}

		category, err := h.categoryService.CreateCategory(c.Request.Context(), req, user, ipAddress)
		if err != nil {
			// Log error but continue with other categories
			continue
		}
		createdCategories = append(createdCategories, category)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "categories seeded successfully",
		"count":      len(createdCategories),
		"categories": createdCategories,
	})
}

// Helper function to generate random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
