package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
    "time" // Tambahkan import time untuk konfigurasi pool

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {
    // Gunakan nama variabel lingkungan yang sebenarnya
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

    // Tambahkan parameter parseTime dan loc yang penting untuk MySQL driver
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

    // Tambahkan konfigurasi pool koneksi di sini
    DB.SetMaxOpenConns(100)
    DB.SetMaxIdleConns(10)
    DB.SetConnMaxLifetime(time.Hour)

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func CloseDB() error {
	if DB == nil {
		return nil
	}
	log.Println("Database connection closed.")
	return DB.Close()
}