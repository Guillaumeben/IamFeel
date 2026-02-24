package db

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"

    _ "modernc.org/sqlite"
)

// DB wraps the database connection
type DB struct {
    conn *sql.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
    // Ensure the directory exists
    dir := filepath.Dir(dbPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create database directory: %w", err)
    }

    // Open database connection
    conn, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Set connection pool settings
    conn.SetMaxOpenConns(1) // SQLite works best with single connection
    conn.SetMaxIdleConns(1)

    // Enable foreign keys
    if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
    }

    db := &DB{conn: conn}

    // Run migrations
    if err := db.migrate(); err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to run migrations: %w", err)
    }

    return db, nil
}

// migrate runs the database migrations
func (db *DB) migrate() error {
    // Execute schema
    if _, err := db.conn.Exec(Schema); err != nil {
        return fmt.Errorf("failed to execute schema: %w", err)
    }

    // Run column migrations (add new columns if they don't exist)
    if err := db.migrateColumns(); err != nil {
        return fmt.Errorf("failed to migrate columns: %w", err)
    }

    return nil
}

// migrateColumns adds new columns to existing tables if they don't exist
func (db *DB) migrateColumns() error {
    // Check and add performance_notes column
    if !db.columnExists("training_sessions", "performance_notes") {
        if _, err := db.conn.Exec("ALTER TABLE training_sessions ADD COLUMN performance_notes TEXT"); err != nil {
            return fmt.Errorf("failed to add performance_notes column: %w", err)
        }
    }

    // Check and add skipped column
    if !db.columnExists("training_sessions", "skipped") {
        if _, err := db.conn.Exec("ALTER TABLE training_sessions ADD COLUMN skipped BOOLEAN DEFAULT 0"); err != nil {
            return fmt.Errorf("failed to add skipped column: %w", err)
        }
    }

    // Check and add skip_reason column
    if !db.columnExists("training_sessions", "skip_reason") {
        if _, err := db.conn.Exec("ALTER TABLE training_sessions ADD COLUMN skip_reason TEXT"); err != nil {
            return fmt.Errorf("failed to add skip_reason column: %w", err)
        }
    }

    return nil
}

// columnExists checks if a column exists in a table
func (db *DB) columnExists(tableName, columnName string) bool {
    query := fmt.Sprintf("SELECT COUNT(*) FROM pragma_table_info('%s') WHERE name='%s'", tableName, columnName)
    var count int
    err := db.conn.QueryRow(query).Scan(&count)
    return err == nil && count > 0
}

// Close closes the database connection
func (db *DB) Close() error {
    if db.conn != nil {
        return db.conn.Close()
    }
    return nil
}

// Conn returns the underlying database connection
func (db *DB) Conn() *sql.DB {
    return db.conn
}

// Health checks if the database is healthy
func (db *DB) Health() error {
    return db.conn.Ping()
}
