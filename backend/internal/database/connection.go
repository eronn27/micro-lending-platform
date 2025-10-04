package database

import (
    "fmt"
    "log"
    "micro-lending-platform/backend/internal/config"
    "os"
    "path/filepath"
    "time"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB  // Global database instance

type Database struct {
    *gorm.DB
}

// NewDB creates and configures a new database connection
func NewDB(cfg *config.Config) (*Database, error) {
    // Create directory for database file if it doesn't exist
    if err := createDataDir(cfg.DBPath); err != nil {
        return nil, fmt.Errorf("failed to create data directory: %w", err)
    }

    // Configure GORM based on environment
    gormConfig := &gorm.Config{}
    if cfg.Environment == "development" {
        gormConfig.Logger = logger.Default.LogMode(logger.Info)  // Verbose logging
    } else {
        gormConfig.Logger = logger.Default.LogMode(logger.Error)  // Only errors
    }

    // Open SQLite database connection
    db, err := gorm.Open(sqlite.Open(cfg.DBPath), gormConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Configure connection pool for better performance
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get SQL DB: %w", err)
    }

    // Connection pool settings
    sqlDB.SetMaxOpenConns(25)        // Maximum open connections
    sqlDB.SetMaxIdleConns(5)         // Maximum idle connections  
    sqlDB.SetConnMaxLifetime(30 * time.Minute)  // Maximum connection lifetime

    // Enable foreign key constraints (SQLite has them off by default)
    if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
        log.Printf("Warning: Failed to enable foreign keys: %v", err)
    }

    log.Printf("Database connected successfully: %s", cfg.DBPath)
    
    DB = db  // Set global instance
    return &Database{db}, nil
}

// NewDBWithMigrations creates connection and runs migrations automatically
func NewDBWithMigrations(cfg *config.Config, migrationsPath string) (*Database, error) {
    db, err := NewDB(cfg)
    if err != nil {
        return nil, err
    }

    // Run migrations immediately after connection
    if err := db.RunAllMigrations(migrationsPath); err != nil {
        db.Close()  // Close connection if migrations fail
        return nil, fmt.Errorf("migrations failed: %w", err)
    }

    return db, nil
}

// createDataDir ensures the database directory exists
func createDataDir(dbPath string) error {
    dir := filepath.Dir(dbPath)
    if dir != "." && dir != "/" {
        if err := os.MkdirAll(dir, 0755); err != nil {
            return err
        }
        log.Printf("Created data directory: %s", dir)
    }
    return nil
}

// GetDB returns the global database instance (singleton pattern)
func GetDB() *gorm.DB {
    return DB
}

// Close safely closes the database connection
func (d *Database) Close() error {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

// HealthCheck verifies database is responsive
func (d *Database) HealthCheck() error {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return fmt.Errorf("failed to get SQL DB: %w", err)
    }

    // Basic ping test
    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("database ping failed: %w", err)
    }

    // Execute a simple query to verify functionality
    var result int
    if err := sqlDB.QueryRow("SELECT 1").Scan(&result); err != nil {
        return fmt.Errorf("database query failed: %w", err)
    }

    if result != 1 {
        return fmt.Errorf("unexpected health check result: %d", result)
    }

    return nil
}
