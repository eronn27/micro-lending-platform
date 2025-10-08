package handlers

import (
	"micro-lending-platform/backend/internal/auth"
	"micro-lending-platform/backend/internal/services"
	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all API routes and groups them by version and resource
func SetupRoutes(
	router *gin.Engine,
	authService *services.AuthService,
	clientService *services.ClientService,
	loanService *services.LoanService,
	paymentService *services.PaymentService,
	reportService *services.ReportService,
) {
	// Initialize handlers
	authHandler := NewAuthHandler(authService)
	clientHandler := NewClientHandler(clientService)
	loanHandler := NewLoanHandler(loanService)
	paymentHandler := NewPaymentHandler(paymentService)
	reportHandler := NewReportHandler(reportService)


	// API v1 group
	v1 := router.Group("/api/v1")
	{
		setupAuthRoutes(v1, authHandler)
		setupClientRoutes(v1, clientHandler)
		setupLoanRoutes(v1, loanHandler)
		setupPaymentRoutes(v1, paymentHandler)
		setupReportRoutes(v1, reportHandler)
	}

	// System routes
	setupSystemRoutes(router)
}

// setupAuthRoutes configures all authentication endpoints
func setupAuthRoutes(rg *gin.RouterGroup, h *AuthHandler) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
		authGroup.POST("/register", h.Register)
		authGroup.POST("/refresh", h.RefreshToken)

		protected := authGroup.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			protected.GET("/me", h.GetCurrentUser)
			protected.PUT("/profile", h.UpdateProfile)
			protected.PATCH("/password", h.ChangePassword)
			protected.POST("/logout", h.Logout)
		}
	}
}

// setupClientRoutes configures all client management endpoints
func setupClientRoutes(rg *gin.RouterGroup, h *ClientHandler) {
	clients := rg.Group("/clients")
	clients.Use(auth.AuthMiddleware())

	{
		clients.GET("", h.GetAllClients)
		clients.POST("", h.CreateClient)
		clients.POST("/simple", h.CreateSimpleClient)
		clients.GET("/:id", h.GetClientByID)
		clients.GET("/:id/details", h.GetClientWithDetails)
		clients.PUT("/:id", h.UpdateClient)
		clients.DELETE("/:id", h.DeleteClient)
		clients.PATCH("/:id/restore", h.RestoreClient)

		clients.GET("/search", h.SearchClients)
		clients.GET("/check-duplicate", h.CheckDuplicate)
		clients.GET("/control-number/:controlNumber", h.GetClientByControlNumber)

		clients.GET("/stats", h.GetClientStats)
		clients.GET("/export", h.ExportClients)

		clients.POST("/bulk", h.BulkCreateClients)
	}
}

// setupLoanRoutes configures loan management endpoints
func setupLoanRoutes(rg *gin.RouterGroup, h *LoanHandler) {
	loans := rg.Group("/loans")
	loans.Use(auth.AuthMiddleware())

	{
		loans.GET("", h.GetAllLoans)
		loans.POST("", h.CreateLoan)
		loans.GET("/:id", h.GetLoanByID)
		loans.PUT("/:id", h.UpdateLoan)
		loans.DELETE("/:id", h.DeleteLoan)

		loans.GET("/client/:clientId", h.GetLoansByClientID)
		loans.GET("/stats", h.GetLoanStats)
	}
}

// setupPaymentRoutes configures payment management endpoints
func setupPaymentRoutes(rg *gin.RouterGroup, h *PaymentHandler) {
	payments := rg.Group("/payments")
	payments.Use(auth.AuthMiddleware())

	{
		payments.GET("", h.GetAllPayments)
		payments.POST("", h.CreatePayment)
		payments.GET("/:id", h.GetPaymentByID)
		payments.PUT("/:id", h.UpdatePayment)
		payments.DELETE("/:id", h.DeletePayment)

		payments.GET("/loan/:loanId", h.GetPaymentsByLoanID)
	}
}

// setupReportRoutes configures reporting endpoints
func setupReportRoutes(rg *gin.RouterGroup, h *ReportHandler) {
	reports := rg.Group("/reports")
	reports.Use(auth.AuthMiddleware())

	{
		reports.GET("/weekly", h.GetWeeklyReport)
		reports.GET("/monthly", h.GetMonthlyReport)
	}
}

// setupSystemRoutes configures system-level endpoints
func setupSystemRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"service":   "micro-lending-platform",
			"version":   "1.0.0",
			"timestamp": gin.H{"iso": "", "unix": 0},
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "Micro Lending Platform API",
			"version":     "1.0.0",
			"description": "Backend API for micro-lending operations and client management",
			"endpoints": gin.H{
				"api_v1":        "/api/v1",
				"health":        "/health",
				"documentation": "TBD",
			},
			"contact": gin.H{
				"support": "TBD",
			},
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":   "Endpoint not found",
			"message": "The requested API endpoint does not exist",
			"path":    c.Request.URL.Path,
			"suggestions": gin.H{
				"api_base":    "/api/v1",
				"health":      "/health",
				"api_info":    "/",
			},
		})
	})
}

// AdminMiddleware for admin-only route protection
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
