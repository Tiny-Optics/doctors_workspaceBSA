package server

import (
	"net/http"
	"os"

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
		AllowOrigins:     []string{"http://localhost:5173", "http://169.255.58.102"}, // Dev and production URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Public routes
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	// Serve static files for uploads
	r.Static("/uploads", "./uploads")

	// Get MongoDB database
	db := s.db.GetDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	institutionRepo := repository.NewInstitutionRepository(db)
	sopCategoryRepo := repository.NewSOPCategoryRepository(db)
	dropboxConfigRepo := repository.NewDropboxConfigRepository(db)
	registryConfigRepo := repository.NewRegistryConfigRepository(db)
	registryFormRepo := repository.NewRegistryFormRepository(db)
	registrySubmissionRepo := repository.NewRegistrySubmissionRepository(db)
	referralConfigRepo := repository.NewReferralConfigRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo, auditRepo)
	institutionService := service.NewInstitutionService(institutionRepo, userRepo, auditRepo)
	userService := service.NewUserService(userRepo, auditRepo, authService)
	auditService := service.NewAuditService(auditRepo, userRepo)

	// Initialize encryption service
	encryptionService, err := service.NewEncryptionService()
	if err != nil {
		panic("Failed to initialize encryption service: " + err.Error())
	}

	// Initialize Dropbox services
	dropboxService := service.NewDropboxService(dropboxConfigRepo, encryptionService)
	dropboxOAuthService := service.NewDropboxOAuthService(dropboxConfigRepo, auditRepo, encryptionService, dropboxService)
	sopCategoryService := service.NewSOPCategoryService(sopCategoryRepo, dropboxService, auditRepo, userRepo)

	// Initialize Dropbox background refresh service
	dropboxRefreshService := service.NewDropboxRefreshService(dropboxService)
	dropboxRefreshService.Start()

	// Store reference for graceful shutdown
	s.dropboxRefreshService = dropboxRefreshService

	// Initialize email and registry services
	emailService := service.NewEmailService(encryptionService)
	registryService := service.NewRegistryService(
		registryConfigRepo,
		registryFormRepo,
		registrySubmissionRepo,
		userRepo,
		institutionRepo,
		auditRepo,
		dropboxService,
		emailService,
	)
	referralService := service.NewReferralService(referralConfigRepo, auditRepo)

	// Initialize password reset service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-change-in-production"
	}
	passwordResetService := service.NewPasswordResetService(
		userRepo,
		passwordResetRepo,
		auditRepo,
		emailService,
		encryptionService,
		registryService,
		jwtSecret,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	institutionHandler := handlers.NewInstitutionHandler(institutionService)
	statsHandler := handlers.NewStatsHandler(userService, institutionService, auditService, sopCategoryService)
	sopCategoryHandler := handlers.NewSOPCategoryHandler(sopCategoryService)
	dropboxAdminHandler := handlers.NewDropboxAdminHandler(dropboxOAuthService)
	registryHandler := handlers.NewRegistryHandler(registryService, encryptionService)
	referralHandler := handlers.NewReferralHandler(referralService)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService)
	smtpHandler := handlers.NewSMTPHandler(registryService)

	// API routes group
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", userHandler.RegisterUser)
			auth.POST("/refresh", authHandler.RefreshToken)

			// Password reset routes (public)
			auth.POST("/forgot-password", passwordResetHandler.ForgotPassword)
			auth.POST("/validate-reset-code", passwordResetHandler.ValidateResetCode)
			auth.POST("/reset-password", passwordResetHandler.ResetPassword)

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

		// Public institution routes (for registration)
		api.GET("/institutions/public", institutionHandler.ListPublicInstitutions)

		// Public SMTP status route (for forgot password)
		api.GET("/smtp/status", smtpHandler.CheckSMTPConfiguration)

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

			// Image upload (super admin only)
			sops.POST("/images/upload", middleware.RequirePermission(models.PermDeleteUsers), sopCategoryHandler.UploadImage)

			// Seeding (super admin only)
			sops.POST("/seed", middleware.RequirePermission(models.PermDeleteUsers), sopCategoryHandler.SeedCategories)
		}

		// Admin routes (super admin only)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService))
		admin.Use(middleware.RequirePermission(models.PermManageSystem))
		{
			// Dropbox configuration
			dropbox := admin.Group("/dropbox")
			{
				dropbox.GET("/status", dropboxAdminHandler.GetStatus)
				dropbox.POST("/authorize", dropboxAdminHandler.InitiateAuth)
				dropbox.POST("/callback", dropboxAdminHandler.CompleteAuth)
				dropbox.POST("/refresh", dropboxAdminHandler.ForceRefresh)
				dropbox.POST("/test", dropboxAdminHandler.TestConnection)
				dropbox.DELETE("/configuration", dropboxAdminHandler.DeleteConfiguration)
			}

			// Registry configuration (super admin only)
			registry := admin.Group("/registry")
			{
				registry.GET("/config", registryHandler.GetConfiguration)
				registry.PUT("/config", registryHandler.UpdateConfiguration)
				registry.POST("/test-email", registryHandler.SendTestEmail)
				registry.GET("/submissions", registryHandler.GetAllSubmissions)
				registry.PATCH("/submissions/:id/status", registryHandler.UpdateSubmissionStatus)

				// SMTP-only configuration endpoints
				registry.GET("/smtp-config", registryHandler.GetSMTPConfig)
				registry.PUT("/smtp-config", registryHandler.UpdateSMTPConfig)
			}

			// Referral configuration (super admin only)
			referrals := admin.Group("/referrals")
			{
				referrals.GET("/config", referralHandler.GetAdminConfig)
				referrals.PUT("/config", referralHandler.UpdateConfig)
			}
		}

		// Admin routes for registry form management (admins and user managers)
		registryAdmin := api.Group("/admin/registry")
		registryAdmin.Use(middleware.AuthMiddleware(authService))
		registryAdmin.Use(middleware.RequirePermission(models.PermManageUsers))
		{
			registryAdmin.POST("/form-schema", registryHandler.CreateFormSchema)
			registryAdmin.GET("/form-schemas", registryHandler.ListFormSchemas)
			registryAdmin.GET("/form-schema/:id", registryHandler.GetFormSchema)
			registryAdmin.PUT("/form-schema/:id", registryHandler.UpdateFormSchema)
			registryAdmin.DELETE("/form-schema/:id", registryHandler.DeleteFormSchema)
		}

		// Registry routes (authenticated users)
		registry := api.Group("/registry")
		registry.Use(middleware.AuthMiddleware(authService))
		{
			registry.GET("/config", registryHandler.GetPublicConfiguration)
			registry.GET("/form-schema", registryHandler.GetActiveFormSchema)
			registry.POST("/submit", registryHandler.SubmitForm)
			registry.GET("/submissions", registryHandler.GetUserSubmissions)
			registry.GET("/submissions/:id", registryHandler.GetSubmission)
			registry.GET("/example-documents", registryHandler.GetExampleDocuments)
			registry.GET("/example-documents/download", registryHandler.GetExampleDocumentDownloadLink)
			registry.GET("/document-download", registryHandler.GetSubmissionDocumentDownloadLink)
		}

		// Referral routes (authenticated users)
		referrals := api.Group("/referrals")
		referrals.Use(middleware.AuthMiddleware(authService))
		{
			referrals.GET("/config", referralHandler.GetConfig)
			referrals.POST("/access", referralHandler.LogAccess)
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
