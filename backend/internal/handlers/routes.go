package handlers

import (
    "micro-lending-platform/backend/internal/auth"
    "micro-lending-platform/backend/internal/services"
    "github.com/gin-gonic/gin"
)

// SetupRoutes initializes all API routes and groups them by version and resource
// This centralizes route configuration and makes the API structure clear and maintainable
func SetupRoutes(router *gin.Engine, authService *services.AuthService, clientService *services.ClientService) {
    // Initialize handlers with their respective services
    authHandler := NewAuthHandler(authService)
    clientHandler := NewClientHandler(clientService)

    // API v1 group - all routes under /api/v1 prefix
    // Versioning allows for future API changes without breaking existing clients
    v1 := router.Group("/api/v1")
    {
        setupAuthRoutes(v1, authHandler)      // Authentication and user management routes
        setupClientRoutes(v1, clientHandler)  // Client management routes
        setupLoanRoutes(v1)                   // Loan management routes (placeholder for future)
        setupReportRoutes(v1)                 // Reporting routes (placeholder for future)
    }

    // System routes (outside API versioning)
    // These are infrastructure endpoints that shouldn't change with API versions
    setupSystemRoutes(router)
}

// setupAuthRoutes configures all authentication and user management endpoints
func setupAuthRoutes(rg *gin.RouterGroup, h *AuthHandler) {
    // Auth group for all authentication-related endpoints
    authe := rg.Group("/auth")
    {
        // Public authentication endpoints (no authentication required)
        authe.POST("/login", h.Login)          // User login - validates credentials and returns JWT
        
        // The following endpoints are placeholders for future implementation
        authe.POST("/register", h.Register)    // User registration (future feature)
        authe.POST("/refresh", h.RefreshToken) // Token refresh (future feature)
        
        // Protected authentication endpoints (require valid JWT)
        // AuthMiddleware validates JWT tokens for all routes in this subgroup
        protected := authe.Group("")
        protected.Use(auth.AuthMiddleware())
        {
            protected.GET("/me", h.GetCurrentUser)         // Get current authenticated user info
            protected.PUT("/profile", h.UpdateProfile)     // Update user profile (future feature)
            protected.PATCH("/password", h.ChangePassword) // Change user password (future feature)
            protected.POST("/logout", h.Logout)            // Logout - token invalidation (future feature)
        }
    }
}

// setupClientRoutes configures all client management endpoints
// Clients are the core entity in the micro-lending platform
func setupClientRoutes(rg *gin.RouterGroup, h *ClientHandler) {
    // Clients group for all client-related operations
    clients := rg.Group("/clients")
    clients.Use(auth.AuthMiddleware()) // All client routes require authentication
    
    {
        // Client CRUD operations
        clients.GET("", h.GetAllClients)                           // Get paginated list of clients with search
        clients.POST("", h.CreateClient)                           // Create new client record with complete data
        clients.POST("/simple", h.CreateSimpleClient)              // Create simple client (backward compatibility)
        clients.GET("/:id", h.GetClientByID)                       // Get specific client by ID
        clients.GET("/:id/details", h.GetClientWithDetails)        // Get client with full related data
        clients.PUT("/:id", h.UpdateClient)                        // Update client information
        clients.DELETE("/:id", h.DeleteClient)                     // Soft delete client
        clients.PATCH("/:id/restore", h.RestoreClient)             // Restore soft-deleted client (future feature)
        
        // Client search and utility endpoints
        clients.GET("/search", h.SearchClients)                    // Advanced client search
        clients.GET("/check-duplicate", h.CheckDuplicate)          // Check for potential duplicate clients
        clients.GET("/control-number/:controlNumber", h.GetClientByControlNumber) // Get client by control number
        
        // Client statistics and reporting
        clients.GET("/stats", h.GetClientStats)                    // Get client statistics dashboard data
        clients.GET("/export", h.ExportClients)                    // Export clients data
        
        // Bulk operations
        clients.POST("/bulk", h.BulkCreateClients)                 // Bulk create multiple clients
    }
}

// setupLoanRoutes configures loan management endpoints (placeholder for future implementation)
// Loans represent the financial transactions between the platform and clients
func setupLoanRoutes(rg *gin.RouterGroup) {
    loans := rg.Group("/loans")
    loans.Use(auth.AuthMiddleware()) // All loan routes will require authentication
    
    {
        // Placeholder routes - to be implemented in future iterations
        // Using 501 status code (Not Implemented) to indicate planned functionality
        loans.GET("", func(c *gin.Context) {
            c.JSON(501, gin.H{
                "message": "Get loans endpoint - coming soon",
                "status":  "not_implemented",
            })
        })
        
        loans.POST("", func(c *gin.Context) {
            c.JSON(501, gin.H{
                "message": "Create loan endpoint - coming soon", 
                "status":  "not_implemented",
            })
        })
        
        loans.GET("/:id", func(c *gin.Context) {
            c.JSON(501, gin.H{
                "message": "Get loan by ID endpoint - coming soon",
                "status":  "not_implemented",
            })
        })
        
        loans.PUT("/:id", func(c *gin.Context) {
            c.JSON(501, gin.H{
                "message": "Update loan endpoint - coming soon",
                "status":  "not_implemented",
            })
        })
    }
}

// setupReportRoutes configures reporting endpoints (placeholder for future implementation)
// Reports provide analytics and business intelligence for platform management
func setupReportRoutes(rg *gin.RouterGroup) {
    reports := rg.Group("/reports")
    reports.Use(auth.AuthMiddleware()) // All report routes will require authentication
    reports.Use(AdminMiddleware()) // Typically reports require admin privileges
    
    {
        // Placeholder routes - to be implemented in future iterations
        reports.GET("/clients", func(c *gin.Context) {
            c.JSON(501, gin.H{
                "message": "Client reports endpoint - coming soon",
                "status":  "not_implemented",
            })
        })
        
        reports.GET("/loans", func(c *gin.Context) {
            c.JSON(501, gin.H{
                "message": "Loan reports endpoint - coming soon",
                "status":  "not_implemented", 
            })
        })
        
        reports.GET("/payments", func(c *gin.Context) {
            c.JSON(501, gin.H{
                "message": "Payment reports endpoint - coming soon",
                "status":  "not_implemented",
            })
        })
    }
}

// setupSystemRoutes configures system-level endpoints for infrastructure and monitoring
// These endpoints are essential for deployment, monitoring, and API discovery
func setupSystemRoutes(router *gin.Engine) {
    // Health check endpoint - used by load balancers, container orchestration, and monitoring systems
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status":    "healthy",
            "service":   "micro-lending-platform",
            "version":   "1.0.0",
            "timestamp": gin.H{"iso": "", "unix": 0}, // Would be populated in actual implementation
        })
    })
    
    // API info endpoint - provides basic API information and discovery
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "name":        "Micro Lending Platform API",
            "version":     "1.0.0", 
            "description": "Backend API for micro-lending operations and client management",
            "endpoints": gin.H{
                "api_v1":        "/api/v1",
                "health":        "/health",
                "documentation": "TBD", // Would point to Swagger/OpenAPI docs if available
            },
            "contact": gin.H{
                "support": "TBD", // Would contain support contact information
            },
        })
    })
    
    // 404 handler for undefined routes - provides helpful error messages
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

// AdminMiddleware is a placeholder for admin-only route protection
// This would be implemented to restrict certain endpoints to admin users only
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // For now, this is a pass-through middleware
        // Future implementation would check user roles from JWT claims
        c.Next()
    }
}
