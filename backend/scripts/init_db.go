package main

import (
    "database/sql"
    "fmt"
    "log"
    "micro-lending-platform/backend/internal/config"
    "micro-lending-platform/backend/internal/database"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/internal/auth"
    "micro-lending-platform/backend/migrations"
    "os"
    "path/filepath"
    "gorm.io/gorm"
)

func main() {
    log.Println("Starting database initialization...")

    // Load configuration
    cfg := config.Load()
    log.Printf("Configuration loaded: DB=%s, Port=%s", cfg.DBPath, cfg.ServerPort)

    // Create data directory if it doesn't exist
    if err := createDataDirectory(cfg.DBPath); err != nil {
        log.Fatal("Failed to create data directory:", err)
    }

    // Initialize database connection
    db, err := database.NewDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    log.Println("Database connection established successfully")

    // Get SQL database instance for migrations
    sqlDB, err := db.DB.DB()
    if err != nil {
        log.Fatal("Failed to get SQL database instance:", err)
    }

    // Check if migrations table exists to determine if this is a fresh install
    isFreshInstall, err := isFreshInstall(sqlDB)
    if err != nil {
        log.Fatal("Failed to check database state:", err)
    }

    if isFreshInstall {
        log.Println("Fresh installation detected. Running initial setup...")
    } else {
        log.Println("Existing database detected. Running migrations...")
    }

    // Run SQL migrations
    migrationsPath := "./migrations" // Fixed path
    if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
        log.Fatalf("Migrations directory not found: %s", migrationsPath)
    }

    log.Printf("Running migrations from: %s", migrationsPath)
    if err := migrations.RunMigrations(sqlDB, migrationsPath); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    // Run GORM auto-migrations for any model changes
    log.Println("Running GORM auto-migrations...")
    if err := db.AutoMigrate(); err != nil {
        log.Fatal("Failed to run auto migrations:", err)
    }

    // Verify migrations were applied successfully
    if err := verifyMigrations(sqlDB, migrationsPath); err != nil {
        log.Fatal("Migration verification failed:", err)
    }

    // Create default admin user if it doesn't exist
    if err := createDefaultAdmin(db.DB); err != nil {
        log.Fatal("Failed to create default admin user:", err)
    }

    // Display migration status
    if err := displayMigrationStatus(sqlDB, migrationsPath); err != nil {
        log.Printf("Warning: Failed to display migration status: %v", err)
    }

    log.Println("Database initialization completed successfully!")
    log.Printf("Database file: %s", cfg.DBPath)
    log.Println("Default admin credentials: username='admin', password='admin123'")
}

// createDataDirectory creates the directory for the database file
func createDataDirectory(dbPath string) error {
    dir := filepath.Dir(dbPath)
    if dir != "." && dir != "/" {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return fmt.Errorf("failed to create directory %s: %w", dir, err)
        }
        log.Printf("Created data directory: %s", dir)
    }
    return nil
}

// isFreshInstall checks if this is a fresh database installation
func isFreshInstall(db *sql.DB) (bool, error) {
    // Check if any tables exist
    query := `
    SELECT COUNT(*) FROM sqlite_master 
    WHERE type='table' AND name NOT LIKE 'sqlite_%'
    `
    var tableCount int
    err := db.QueryRow(query).Scan(&tableCount)
    if err != nil {
        return false, err
    }

    return tableCount == 0, nil
}

// verifyMigrations verifies that migrations were applied correctly
func verifyMigrations(db *sql.DB, migrationsPath string) error {
    // Check if essential tables were created
    essentialTables := []string{"clients", "loans", "payments", "users", "migration_history"}
    
    for _, table := range essentialTables {
        query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
        var count int
        err := db.QueryRow(query, table).Scan(&count)
        if err != nil {
            return fmt.Errorf("failed to check table %s: %w", table, err)
        }
        if count == 0 {
            return fmt.Errorf("essential table %s was not created", table)
        }
    }

    log.Printf("Verified %d essential tables exist", len(essentialTables))
    return nil
}

// createDefaultAdmin ensures the default admin user exists
func createDefaultAdmin(db *gorm.DB) error {
    // Check if admin user already exists
    var count int64
    result := db.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
    if result.Error != nil {
        return fmt.Errorf("failed to check admin user: %w", result.Error)
    }

    if count > 0 {
        log.Println("Default admin user already exists")
        return nil
    }

    // Create default admin user
    adminUser := &models.User{
        Username: "admin",
        IsAdmin:  true,
    }

    // Hash password
    hashedPassword, err := auth.HashPassword("admin123")
    if err != nil {
        return fmt.Errorf("failed to hash admin password: %w", err)
    }
    adminUser.PasswordHash = hashedPassword

    // Create user
    result = db.Create(adminUser)
    if result.Error != nil {
        return fmt.Errorf("failed to create admin user: %w", result.Error)
    }

    log.Println("Default admin user created successfully")
    return nil
}

// displayMigrationStatus shows the current migration status
func displayMigrationStatus(db *sql.DB, migrationsPath string) error {
    status, err := migrations.GetMigrationStatus(db, migrationsPath)
    if err != nil {
        return err
    }

    log.Println("Migration Status:")
    log.Println("=================")
    for _, migration := range status {
        log.Printf("%s - %s: %s", migration.Version, migration.Name, migration.AppliedAt)
    }
    log.Println("=================")

    return nil
}
