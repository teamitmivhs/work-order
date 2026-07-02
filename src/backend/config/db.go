package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// requireEnv membaca env variable dan return error jika tidak di-set.
// Dipakai untuk konfigurasi kritis yang TIDAK boleh punya default hardcoded.
func requireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return val, nil
}

// getEnvOrDefault membaca env variable dengan nilai default.
// Hanya dipakai untuk konfigurasi non-sensitif (host, port, dbname).
func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func InitDB() error {
	// FIX: DB_USER dan DB_PASSWORD tidak lagi punya fallback hardcoded.
	// Credential sensitif harus selalu di-set via environment / secret manager.
	dbUser, err := requireEnv("DB_USER")
	if err != nil {
		return err
	}

	dbPass, err := requireEnv("DB_PASSWORD")
	if err != nil {
		return err
	}

	// Host, port, dan nama DB boleh punya default karena tidak sensitif
	dbHost := getEnvOrDefault("DB_HOST", "127.0.0.1")
	dbPort := getEnvOrDefault("DB_PORT", "3306")
	dbName := getEnvOrDefault("DB_NAME", "dbwoit")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	maxRetries := 30
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Attempt %d/%d: Failed to open database: %v", i+1, maxRetries, err)
			time.Sleep(retryDelay)
			continue
		}

		DB.SetMaxOpenConns(100)
		DB.SetMaxIdleConns(10)
		DB.SetConnMaxLifetime(time.Hour)

		if err = DB.Ping(); err != nil {
			log.Printf("Attempt %d/%d: Failed to ping database: %v", i+1, maxRetries, err)
			DB.Close()
			time.Sleep(retryDelay)
			continue
		}

		if err = runMigrations(dbName); err != nil {
			DB.Close()
			return fmt.Errorf("failed to run database migrations: %w", err)
		}

		log.Println("Database connected successfully")
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

func GetDB() *sql.DB {
	return DB
}

func CloseDB() error {
	if DB == nil {
		return nil
	}
	log.Println("Database connection closed.")
	return DB.Close()
}

func runMigrations(dbName string) error {
	if DB == nil {
		return nil
	}

	migrations := []struct {
		table  string
		column string
		sql    string
	}{
		{"members", "Division", "ALTER TABLE members ADD COLUMN Division varchar(50) DEFAULT NULL AFTER Role"},
		{"members", "AccountStatus", "ALTER TABLE members ADD COLUMN AccountStatus enum('pending','active','rejected','disabled') NOT NULL DEFAULT 'active' AFTER Avatar"},
		{"members", "MembershipStatus", "ALTER TABLE members ADD COLUMN MembershipStatus enum('active','alumni','inactive') NOT NULL DEFAULT 'active' AFTER AccountStatus"},
		{"members", "BatchYear", "ALTER TABLE members ADD COLUMN BatchYear varchar(20) DEFAULT NULL AFTER MembershipStatus"},
		{"members", "GraduationYear", "ALTER TABLE members ADD COLUMN GraduationYear int DEFAULT NULL AFTER BatchYear"},
		{"members", "CanHandleWorkOrder", "ALTER TABLE members ADD COLUMN CanHandleWorkOrder tinyint(1) NOT NULL DEFAULT 1 AFTER GraduationYear"},
		{"members", "RegisteredAt", "ALTER TABLE members ADD COLUMN RegisteredAt datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER CanHandleWorkOrder"},
		{"members", "ApprovedAt", "ALTER TABLE members ADD COLUMN ApprovedAt datetime DEFAULT NULL AFTER RegisteredAt"},
		{"members", "ApprovedBy", "ALTER TABLE members ADD COLUMN ApprovedBy int DEFAULT NULL AFTER ApprovedAt"},
		{"orders", "TrackingCode", "ALTER TABLE orders ADD COLUMN TrackingCode varchar(20) DEFAULT NULL AFTER OrderNumber"},
		{"orders", "StartedAt", "ALTER TABLE orders ADD COLUMN StartedAt datetime DEFAULT NULL AFTER TimeSort"},
		{"orders", "Rating", "ALTER TABLE orders ADD COLUMN Rating tinyint DEFAULT NULL AFTER Notes"},
		{"orders", "NotesQuality", "ALTER TABLE orders ADD COLUMN NotesQuality tinyint DEFAULT NULL AFTER Rating"},
		{"orders", "DocumentationPhoto", "ALTER TABLE orders ADD COLUMN DocumentationPhoto varchar(255) DEFAULT NULL AFTER NotesQuality"},
	}

	for _, migration := range migrations {
		exists, err := columnExists(dbName, migration.table, migration.column)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if _, err := DB.Exec(migration.sql); err != nil {
			return fmt.Errorf("add %s.%s: %w", migration.table, migration.column, err)
		}
	}

	if _, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS shift_day_counter (
			ID tinyint NOT NULL,
			LastDate date NOT NULL,
			LastDay int NOT NULL,
			LastMonth int NOT NULL,
			LastYear int NOT NULL,
			RolloverCount int NOT NULL DEFAULT 0,
			UpdatedAt datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (ID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
	`); err != nil {
		return fmt.Errorf("create shift_day_counter: %w", err)
	}

	if _, err := DB.Exec(`
		UPDATE members
		SET Division = CASE
		        WHEN LOWER(TRIM(Role)) = 'soundman' THEN 'Soundman'
		        WHEN LOWER(TRIM(Role)) = 'programmer' THEN 'Programmer'
		        WHEN LOWER(TRIM(Role)) = 'maintenance' THEN 'Maintenance'
		        WHEN LOWER(TRIM(Role)) IN ('data analyst', 'data-analyst', 'data_analyst') THEN 'Data Analyst'
		        ELSE Division
		    END,
		    Role = CASE
		        WHEN LOWER(TRIM(Role)) IN ('soundman', 'programmer', 'maintenance', 'data analyst', 'data-analyst', 'data_analyst') THEN 'Operator'
		        ELSE Role
		    END
		WHERE Division IS NULL OR Division = ''
	`); err != nil {
		return err
	}

	if _, err := DB.Exec(`
		UPDATE members
		SET Role = 'Operator'
		WHERE LOWER(TRIM(Role)) IN ('staff', 'technician')
	`); err != nil {
		return err
	}

	if _, err := DB.Exec(`
		UPDATE orders
		SET TrackingCode = CONCAT('WO-', LPAD(ID, 6, '0'))
		WHERE TrackingCode IS NULL OR TrackingCode = ''
	`); err != nil {
		return err
	}

	if _, err := DB.Exec(`
		UPDATE members
		SET AccountStatus = 'active',
		    MembershipStatus = 'active',
		    CanHandleWorkOrder = CASE WHEN Role = 'Guest' THEN 0 ELSE CanHandleWorkOrder END
		WHERE AccountStatus IS NULL OR AccountStatus = ''
	`); err != nil && !strings.Contains(err.Error(), "Data truncated") {
		return err
	}

	_, err := DB.Exec("UPDATE members SET CanHandleWorkOrder = 0 WHERE Role = 'Guest'")
	if err != nil {
		return err
	}

	return nil
}

func columnExists(dbName, tableName, columnName string) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = ?
	`, dbName, tableName, columnName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
