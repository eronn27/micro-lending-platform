package migrations

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "path/filepath"
    "sort"
    "strings"
)

// Migration represents a database migration file
type Migration struct {
    Version string  // e.g., "001", "002"
    Name    string  // e.g., "initial_schema", "add_users_table"
    Up      string  // SQL content to apply the migration
}

// MigrationHistory tracks applied migrations in database
type MigrationHistory struct {
    ID        int
    Version   string
    Name      string
    AppliedAt string
}

// RunMigrations executes all pending migrations in order
func RunMigrations(db *sql.DB, migrationsPath string) error {
    log.Println("Starting database migrations...")

    // Create migration history table if it doesn't exist
    if err := createMigrationTable(db); err != nil {
        return fmt.Errorf("failed to create migration table: %w", err)
    }

    // Get already applied migrations from database
    appliedMigrations, err := getAppliedMigrations(db)
    if err != nil {
        return fmt.Errorf("failed to get applied migrations: %w", err)
    }

    // Load migration files from directory
    migrations, err := loadMigrations(migrationsPath)
    if err != nil {
        return fmt.Errorf("failed to load migrations: %w", err)
    }

    if len(migrations) == 0 {
        log.Println("No migration files found")
        return nil
    }

    // Apply pending migrations in version order
    appliedCount := 0
    for _, migration := range migrations {
        if _, alreadyApplied := appliedMigrations[migration.Version]; alreadyApplied {
            log.Printf("Migration %s (%s) already applied, skipping", migration.Version, migration.Name)
            continue
        }

        log.Printf("Applying migration: %s - %s", migration.Version, migration.Name)
        
        // Apply migration SQL statements
        if err := applyMigration(db, migration); err != nil {
            return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
        }

        // Record successful migration in history table
        if err := recordMigration(db, migration); err != nil {
            return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
        }

        appliedCount++
        log.Printf("Successfully applied migration: %s - %s", migration.Version, migration.Name)
    }

    if appliedCount == 0 {
        log.Println("No new migrations to apply")
    } else {
        log.Printf("Applied %d migration(s) successfully", appliedCount)
    }

    return nil
}

// loadMigrations reads and parses all SQL files in migrations directory
func loadMigrations(migrationsPath string) ([]Migration, error) {
    files, err := ioutil.ReadDir(migrationsPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read migrations directory: %w", err)
    }

    var migrations []Migration
    for _, file := range files {
        if file.IsDir() {
            continue
        }

        filename := file.Name()
        // Look for files like "001_initial_schema.sql"
        if strings.HasSuffix(filename, ".sql") && strings.HasPrefix(filename, "00") {
            parts := strings.Split(filename, "_")
            if len(parts) < 2 {
                continue // Skip invalid filenames
            }

            version := parts[0]  // "001"
            name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".sql")  // "initial_schema"
            
            // Read SQL file content
            content, err := ioutil.ReadFile(filepath.Join(migrationsPath, filename))
            if err != nil {
                return nil, fmt.Errorf("failed to read migration file %s: %w", filename, err)
            }

            migrations = append(migrations, Migration{
                Version: version,
                Name:    name,
                Up:      string(content),
            })
        }
    }

    // Sort migrations by version number (001, 002, 003...)
    sort.Slice(migrations, func(i, j int) bool {
        return migrations[i].Version < migrations[j].Version
    })

    return migrations, nil
}

// createMigrationTable creates the migration tracking table
func createMigrationTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS migration_history (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        version VARCHAR(10) NOT NULL UNIQUE,  -- Migration version
        name VARCHAR(255) NOT NULL,           -- Migration name
        applied_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- When applied
    )
    `
    _, err := db.Exec(query)
    return err
}

// getAppliedMigrations retrieves already applied migrations from database
func getAppliedMigrations(db *sql.DB) (map[string]bool, error) {
    query := "SELECT version, name FROM migration_history ORDER BY version"
    rows, err := db.Query(query)
    if err != nil {
        // If table doesn't exist yet, return empty map (no migrations applied)
        return make(map[string]bool), nil
    }
    defer rows.Close()

    applied := make(map[string]bool)
    for rows.Next() {
        var version, name string
        if err := rows.Scan(&version, &name); err != nil {
            return nil, err
        }
        applied[version] = true
        log.Printf("Found applied migration: %s - %s", version, name)
    }

    return applied, nil
}

// applyMigration executes SQL statements from migration file
func applyMigration(db *sql.DB, migration Migration) error {
    // Split SQL by semicolons to handle multiple statements
    statements := strings.Split(migration.Up, ";")
    
    for i, stmt := range statements {
        stmt = strings.TrimSpace(stmt)
        if stmt == "" {
            continue
        }
        
        // Skip transaction commands since SQLite doesn't support nested transactions
        if strings.HasPrefix(strings.ToUpper(stmt), "BEGIN") || 
           strings.HasPrefix(strings.ToUpper(stmt), "COMMIT") {
            continue
        }
        
        log.Printf("Executing statement %d for migration %s", i+1, migration.Version)
        _, err := db.Exec(stmt)
        if err != nil {
            return fmt.Errorf("failed to execute statement %d: %w. SQL: %s", i+1, err, stmt)
        }
    }
    
    return nil
}

// recordMigration saves migration record to prevent re-application
func recordMigration(db *sql.DB, migration Migration) error {
    query := "INSERT INTO migration_history (version, name) VALUES (?, ?)"
    _, err := db.Exec(query, migration.Version, migration.Name)
    return err
}

// GetMigrationStatus shows current migration state (applied vs pending)
func GetMigrationStatus(db *sql.DB, migrationsPath string) ([]MigrationHistory, error) {
    applied, err := getAppliedMigrations(db)
    if err != nil {
        return nil, err
    }

    migrations, err := loadMigrations(migrationsPath)
    if err != nil {
        return nil, err
    }

    var status []MigrationHistory
    for _, migration := range migrations {
        appliedAt := "PENDING"
        if applied[migration.Version] {
            appliedAt = "APPLIED"
        }

        status = append(status, MigrationHistory{
            Version: migration.Version,
            Name:    migration.Name,
            AppliedAt: appliedAt,
        })
    }

    return status, nil
}
