package server

import (
	"net/http"

	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Public routes
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	// Get MongoDB database
	db := s.db.GetDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	institutionRepo := repository.NewInstitutionRepository(db)
	sopCategoryRepo := repository.NewSOPCategoryRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo, auditRepo)
	institutionService := service.NewInstitutionService(institutionRepo, userRepo, auditRepo)
	userService := service.NewUserService(userRepo, auditRepo, authService)
	auditService := service.NewAuditService(auditRepo, userRepo)
	dropboxService := service.NewDropboxService()
	sopCategoryService := service.NewSOPCategoryService(sopCategoryRepo, dropboxService, auditRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	institutionHandler := handlers.NewInstitutionHandler(institutionService)
	statsHandler := handlers.NewStatsHandler(userService, institutionService, auditService)
	sopCategoryHandler := handlers.NewSOPCategoryHandler(sopCategoryService)

	// API routes group
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(authService))
			{
				authProtected.GET("/me", authHandler.Me)
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.POST("/change-password", authHandler.ChangePassword)
			}
		}

		// User routes (all protected)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(authService))
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", middleware.RequirePermission(models.PermManageUsers), userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.POST("/:id/activate", middleware.RequirePermission(models.PermManageUsers), userHandler.ActivateUser)
			users.POST("/:id/deactivate", middleware.RequirePermission(models.PermManageUsers), userHandler.DeactivateUser)
			users.DELETE("/:id", middleware.RequirePermission(models.PermDeleteUsers), userHandler.DeleteUser)
		}

		// Institution routes (all protected)
		institutions := api.Group("/institutions")
		institutions.Use(middleware.AuthMiddleware(authService))
		{
			institutions.GET("", institutionHandler.ListInstitutions)
			institutions.GET("/:id", institutionHandler.GetInstitution)
			institutions.POST("", middleware.RequirePermission(models.PermManageUsers), institutionHandler.CreateInstitution)
			institutions.PUT("/:id", middleware.RequirePermission(models.PermManageUsers), institutionHandler.UpdateInstitution)
			institutions.DELETE("/:id", middleware.RequirePermission(models.PermDeleteUsers), institutionHandler.DeleteInstitution)
			institutions.POST("/:id/activate", middleware.RequirePermission(models.PermManageUsers), institutionHandler.ActivateInstitution)
			institutions.POST("/:id/deactivate", middleware.RequirePermission(models.PermManageUsers), institutionHandler.DeactivateInstitution)
		}

		// Stats routes (all protected)
		stats := api.Group("/stats")
		stats.Use(middleware.AuthMiddleware(authService))
		{
			stats.GET("/admin", middleware.RequirePermission(models.PermManageUsers), statsHandler.GetAdminStats)
			stats.GET("/recent-activity", middleware.RequirePermission(models.PermManageUsers), statsHandler.GetRecentActivity)
		}

		// SOP routes
		sops := api.Group("/sops")
		sops.Use(middleware.AuthMiddleware(authService))
		{
			// Categories (read for all authenticated users, write for super admins)
			categories := sops.Group("/categories")
			{
				categories.GET("", sopCategoryHandler.ListCategories)
				categories.GET("/:id", sopCategoryHandler.GetCategory)
				categories.GET("/:id/files", sopCategoryHandler.GetCategoryFiles)
				categories.GET("/:id/files/download", sopCategoryHandler.DownloadFile)

				categories.POST("", middleware.RequirePermission(models.PermDeleteUsers), sopCategoryHandler.CreateCategory)
				categories.PUT("/:id", middleware.RequirePermission(models.PermDeleteUsers), sopCategoryHandler.UpdateCategory)
				categories.DELETE("/:id", middleware.RequirePermission(models.PermDeleteUsers), sopCategoryHandler.DeleteCategory)
			}
		}
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
