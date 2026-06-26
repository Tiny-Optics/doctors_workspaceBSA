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

// WorkingPartyCategoryHandler handles working party category requests
type WorkingPartyCategoryHandler struct {
	categoryService *service.WorkingPartyCategoryService
}

// NewWorkingPartyCategoryHandler creates a new WorkingPartyCategoryHandler
func NewWorkingPartyCategoryHandler(categoryService *service.WorkingPartyCategoryService) *WorkingPartyCategoryHandler {
	return &WorkingPartyCategoryHandler{
		categoryService: categoryService,
	}
}

// CreateCategory godoc
// @Summary Create a new working party category
// @Description Create a new working party category (requires super admin permission)
// @Tags working-party-categories
// @Accept json
// @Produce json
// @Param request body models.CreateWorkingPartyCategoryRequest true "Category information"
// @Success 201 {object} models.WorkingPartyCategory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /working-parties/categories [post]
// @Security BearerAuth
func (h *WorkingPartyCategoryHandler) CreateCategory(c *gin.Context) {
	var req models.CreateWorkingPartyCategoryRequest
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
		} else if err == service.ErrDuplicateWorkingPartySlug || err == repository.ErrDuplicateSlug {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategory godoc
// @Summary Get a working party category by ID
// @Description Get working party category information by ID
// @Tags working-party-categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.WorkingPartyCategory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /working-parties/categories/:id [get]
// @Security BearerAuth
func (h *WorkingPartyCategoryHandler) GetCategory(c *gin.Context) {
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
		if err == service.ErrWorkingPartyCategoryNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// ListCategories godoc
// @Summary List all working party categories
// @Description List working party categories with pagination and filters
// @Tags working-party-categories
// @Produce json
// @Param search query string false "Search term"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /working-parties/categories [get]
// @Security BearerAuth
func (h *WorkingPartyCategoryHandler) ListCategories(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

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
// @Summary Update a working party category
// @Description Update a working party category (requires super admin permission)
// @Tags working-party-categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param request body models.UpdateWorkingPartyCategoryRequest true "Update information"
// @Success 200 {object} models.WorkingPartyCategory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /working-parties/categories/:id [put]
// @Security BearerAuth
func (h *WorkingPartyCategoryHandler) UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req models.UpdateWorkingPartyCategoryRequest
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
		} else if err == service.ErrWorkingPartyCategoryNotFound {
			statusCode = http.StatusNotFound
		} else if err == service.ErrDuplicateWorkingPartySlug || err == repository.ErrDuplicateSlug {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete a working party category
// @Description Delete a working party category from database (Dropbox folder remains) (requires super admin permission)
// @Tags working-party-categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /working-parties/categories/:id [delete]
// @Security BearerAuth
func (h *WorkingPartyCategoryHandler) DeleteCategory(c *gin.Context) {
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
		} else if err == service.ErrWorkingPartyCategoryNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

// GetCategoryFiles godoc
// @Summary List files in a working party category
// @Description List all files in a category's Dropbox folder
// @Tags working-party-categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /working-parties/categories/:id/files [get]
// @Security BearerAuth
func (h *WorkingPartyCategoryHandler) GetCategoryFiles(c *gin.Context) {
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
		if err == service.ErrWorkingPartyCategoryNotFound {
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
// @Summary Get download link for a working party file
// @Description Get a temporary download link for a specific file in a category
// @Tags working-party-categories
// @Produce json
// @Param id path string true "Category ID"
// @Param path query string true "File path within the category folder"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /working-parties/categories/:id/files/download [get]
// @Security BearerAuth
func (h *WorkingPartyCategoryHandler) DownloadFile(c *gin.Context) {
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
		if err == service.ErrWorkingPartyCategoryNotFound {
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
// @Summary Upload working party category image
// @Description Upload an image for working party category (requires super admin permission)
// @Tags working-party-categories
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image file (jpg, jpeg, png, webp, max 5MB)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /working-parties/images/upload [post]
// @Security BearerAuth
func (h *WorkingPartyCategoryHandler) UploadImage(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if !user.HasPermission(models.PermDeleteUsers) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image size must be less than 5MB"})
		return
	}

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

	uploadDir := "./uploads/working-parties"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), workingPartyRandomString(8), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
		return
	}

	imagePath := "/uploads/working-parties/" + filename

	c.JSON(http.StatusOK, gin.H{
		"imagePath": imagePath,
		"message":   "image uploaded successfully",
	})
}

func workingPartyRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
