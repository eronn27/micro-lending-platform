package main

import (
    "log"
    "micro-lending-platform/backend/internal/config"
    "micro-lending-platform/backend/internal/database"
    "micro-lending-platform/backend/internal/handlers"
    "micro-lending-platform/backend/internal/repositories"
    "micro-lending-platform/backend/internal/services"
    "path/filepath"
    "github.com/gin-gonic/gin"
)

// main is the entry point of the micro-lending platform backend API
func main() {
    // Load configuration from environment variables or default values
    cfg := config.Load()
    
    // Initialize database with migrations
    migrationsPath := filepath.Join(".", "migrations")
    db, err := database.NewDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close() // Ensure database connection is closed when program exits

    // Run all database migrations to ensure schema is up-to-date
    if err := db.RunAllMigrations(migrationsPath); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }


    // Configure Gin router based on environment
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode) // Disable debug logging in production
    }
    
    // Create a new Gin router with default middleware
    router := gin.Default()
    
    // Add essential middleware to all routes
    router.Use(gin.Recovery()) // Recovery middleware recovers from any panics
    
    // Logger middleware prints HTTP request logs (only in development)
    if cfg.Environment == "development" {
        router.Use(gin.Logger())
    }

    // CORS (Cross-Origin Resource Sharing) middleware
    // This allows the frontend (running on localhost:3000) to communicate with the backend
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS") // Added PATCH
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Accept")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        
        // Handle OPTIONS method explicitly (preflight requests)
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204) // No Content - successful preflight
            return
        }
        
        c.Next() // Continue to the next middleware/handler
    })

    // Initialize repositories
    userRepo := repositories.NewUserRepository(db.DB)
    clientRepo := repositories.NewClientRepository(db.DB)
    loanRepo := repositories.NewLoanRepository(db.DB)
    paymentRepo := repositories.NewPaymentRepository(db.DB)
    reportRepo := repositories.NewReportRepository(db.DB)

    // Initialize services
    authService := services.NewAuthService(userRepo)
    clientService := services.NewClientService(clientRepo)
    loanService := services.NewLoanService(loanRepo, clientRepo)
    paymentService := services.NewPaymentService(paymentRepo, loanRepo)
    reportService := services.NewReportService(reportRepo) 

    // Setup routes with all services
    handlers.SetupRoutes(router, authService, clientService, loanService, paymentService, reportService)

    // Start the HTTP server on the configured port
    log.Printf("Server starting on port %s", cfg.ServerPort)
    if err := router.Run(":" + cfg.ServerPort); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
