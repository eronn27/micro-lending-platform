package database

import (
    "fmt"
    "log"
    "micro-lending-platform/backend/internal/models"
    "micro-lending-platform/backend/migrations"
    "os"
)

// AutoMigrate runs GORM auto-migration only if no tables exist
// This provides fallback if SQL migrations haven't been run
func (d *Database) AutoMigrate() error {
    log.Println("Starting GORM auto-migration...")
    
    // Check if tables already exist from SQL migrations
    var tableCount int64
    if err := d.DB.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='clients'").Scan(&tableCount).Error; err != nil {
        return fmt.Errorf("failed to check existing tables: %w", err)
    }

    // Skip GORM auto-migration if SQL migrations have already created tables
    if tableCount > 0 {
        log.Println("Tables already exist from SQL migrations, skipping GORM auto-migration")
        return nil
    }

    // Define all models that need database tables
    modelsToMigrate := []interface{}{
        &models.Client{}, &models.IncomeInfo{}, &models.Loan{}, &models.Payment{},
        &models.CoMaker{}, &models.FamilyMember{}, &models.Document{}, &models.User{},
    }

    // Create tables for each model
    for i, model := range modelsToMigrate {
        log.Printf("Migrating model %d/%d...", i+1, len(modelsToMigrate))
        
        if err := d.DB.AutoMigrate(model); err != nil {
            return fmt.Errorf("failed to auto-migrate model %T: %w", model, err)
        }
    }

    log.Println("GORM auto-migration completed successfully")
    return nil
}

// RunSQLMigrations runs the SQL-based migration system
func (d *Database) RunSQLMigrations(migrationsPath string) error {
    log.Println("Starting SQL migrations...")
    
    // Verify migrations directory exists
    if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
        return fmt.Errorf("migrations directory does not exist: %s", migrationsPath)
    }

    // Get raw SQL database connection (bypassing GORM)
    sqlDB, err := d.DB.DB()
    if err != nil {
        return fmt.Errorf("failed to get SQL database instance: %w", err)
    }

    // Run the migration engine
    if err := migrations.RunMigrations(sqlDB, migrationsPath); err != nil {
        return fmt.Errorf("failed to run SQL migrations: %w", err)
    }

    log.Println("SQL migrations completed successfully")
    return nil
}

// RunAllMigrations orchestrates the complete migration process
func (d *Database) RunAllMigrations(migrationsPath string) error {
    log.Println("Running all database migrations...")

    // Priority 1: Run SQL migrations (preferred method)
    if err := d.RunSQLMigrations(migrationsPath); err != nil {
        return fmt.Errorf("SQL migrations failed: %w", err)
    }

    // Priority 2: Run GORM auto-migrations as fallback
    if err := d.AutoMigrate(); err != nil {
        log.Printf("GORM auto-migration had issues (normal if tables already exist): %v", err)
        // Don't fail here - it's normal if SQL migrations already created tables
    }

    log.Println("All database migrations completed successfully")
    return nil
}


// VerifyDatabaseConnection checks database connectivity
func (d *Database) VerifyDatabaseConnection() error {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return fmt.Errorf("failed to get SQL database: %w", err)
    }

    // Test database connection
    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("database ping failed: %w", err)
    }

    log.Printf("Database connection verified successfully")
    return nil
}
