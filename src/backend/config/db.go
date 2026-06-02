package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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
