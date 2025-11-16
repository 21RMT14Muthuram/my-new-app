package database

import (
    "database/sql"
    "fmt"
    "log"
    _ "modernc.org/sqlite" // Pure-Go SQLite driver
)

var DB *sql.DB

func Connect() error {
    var err error
    // modernc.org/sqlite connection string
    DB, err = sql.Open("sqlite", "../app.db")
    if err != nil {
        return err
    }
    
    // Test the connection
    err = DB.Ping()
    if err != nil {
        return err
    }
    
    // Enable WAL mode for better performance
    _, err = DB.Exec("PRAGMA journal_mode=WAL;")
    if err != nil {
        log.Printf("Warning: Could not enable WAL mode: %v", err)
    }
    
    // Enable foreign keys
    _, err = DB.Exec("PRAGMA foreign_keys=ON;")
    if err != nil {
        log.Printf("Warning: Could not enable foreign keys: %v", err)
    }
    
    // Enable busy timeout
    _, err = DB.Exec("PRAGMA busy_timeout=5000;")
    if err != nil {
        log.Printf("Warning: Could not set busy timeout: %v", err)
    }
    
    // Create tables if they don't exist
    err = createTables()
    if err != nil {
        return err
    }
    
    log.Println("Connected to SQLite database successfully using modernc.org/sqlite")
    return nil
}

func createTables() error {
    createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        is_verified BOOLEAN DEFAULT 0,
        otp_code TEXT,
        otp_expires_at DATETIME,
        verified_at DATETIME,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    
    _, err := DB.Exec(createUsersTable)
    if err != nil {
        return fmt.Errorf("failed to create users table: %v", err)
    }
    
    return nil
}

func Migrate() error {
    // Check if new columns exist, if not add them
    columns := []string{"is_verified", "otp_code", "otp_expires_at", "verified_at", "updated_at"}
    
    for _, column := range columns {
        var exists bool
        err := DB.QueryRow(`
            SELECT COUNT(*) > 0 
            FROM pragma_table_info('users') 
            WHERE name = ?
        `, column).Scan(&exists)
        
        if err != nil {
            return fmt.Errorf("failed to check column %s: %v", column, err)
        }
        
        if !exists {
            var alterSQL string
            switch column {
            case "is_verified":
                alterSQL = "ALTER TABLE users ADD COLUMN is_verified BOOLEAN DEFAULT 0"
            case "otp_code":
                alterSQL = "ALTER TABLE users ADD COLUMN otp_code TEXT"
            case "otp_expires_at":
                alterSQL = "ALTER TABLE users ADD COLUMN otp_expires_at DATETIME"
            case "verified_at":
                alterSQL = "ALTER TABLE users ADD COLUMN verified_at DATETIME"
            case "updated_at":
                alterSQL = "ALTER TABLE users ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP"
            }
            
            _, err := DB.Exec(alterSQL)
            if err != nil {
                return fmt.Errorf("failed to add column %s: %v", column, err)
            }
            log.Printf("Added column %s to users table", column)
        }
    }
    
    log.Println("Database migration completed successfully")
    return nil
}