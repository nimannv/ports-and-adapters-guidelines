package gorm

import (
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenDB connects to either an SQLite or PostgreSQL database via GORM and runs migrations.
func OpenDB(driver, dsn string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "postgres", "postgresql":
		dialector = postgres.Open(dsn)
	case "sqlite", "sqlite3", "":
		if dsn == "" {
			dsn = "tasks.db"
		}
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %q (must be 'sqlite' or 'postgres')", driver)
	}

	cfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(dialector, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate tables
	if err := db.AutoMigrate(&BoardModel{}, &ColumnModel{}, &TaskModel{}); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}

	return db, nil
}

