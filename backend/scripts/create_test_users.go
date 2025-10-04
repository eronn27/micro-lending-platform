package main

import (
    "fmt"
    "log"
    "micro-lending-platform/backend/internal/config"
    "micro-lending-platform/backend/internal/database"
    "micro-lending-platform/backend/internal/models"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

func main() {
    log.Println("Creating test users...")

    // Load configuration
    cfg := config.Load()

    // Initialize database connection
    db, err := database.NewDB(cfg)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Create test users
    if err := createTestUsers(db.DB); err != nil {
        log.Fatal("Failed to create test users:", err)
    }

    log.Println("Test users created successfully!")
    log.Println("You can now login with any of these credentials:")
    log.Println("Admin: username='admin', password='admin123'")
    log.Println("Staff: username='staff1', password='staff123'")
    log.Println("Manager: username='manager1', password='manager123'")
}

func createTestUsers(db *gorm.DB) error {
    testUsers := []struct {
        username string
        password string
        isAdmin  bool
    }{
        {"admin", "admin123", true},
        {"staff1", "staff123", false},
        {"manager1", "manager123", true},
        {"john", "john123", false},
        {"sarah", "sarah123", false},
    }

    for _, userData := range testUsers {
        // Check if user already exists
        var existingUser models.User
        result := db.Where("username = ?", userData.username).First(&existingUser)
        
        if result.Error == nil {
            log.Printf("User %s already exists, skipping...", userData.username)
            continue
        }

        // Hash password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.password), bcrypt.DefaultCost)
        if err != nil {
            return fmt.Errorf("failed to hash password for %s: %w", userData.username, err)
        }

        // Create user
        user := models.User{
            Username:     userData.username,
            PasswordHash: string(hashedPassword),
            IsAdmin:      userData.isAdmin,
        }

        if err := db.Create(&user).Error; err != nil {
            return fmt.Errorf("failed to create user %s: %w", userData.username, err)
        }

        log.Printf("Created user: %s (Admin: %v)", userData.username, userData.isAdmin)
    }

    return nil
}
