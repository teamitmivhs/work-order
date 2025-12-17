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

func InitDB() error {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "adminit2025"
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		dbPass = "databaseit2045"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "dbwoit"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	var err error
	maxRetries := 30
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Attempt %d: Failed to open database: %v", i+1, err)
			time.Sleep(retryDelay)
			continue
		}

		DB.SetMaxOpenConns(100)
		DB.SetMaxIdleConns(10)
		DB.SetConnMaxLifetime(time.Hour)

		err = DB.Ping()
		if err != nil {
			log.Printf("Attempt %d: Failed to connect to database: %v", i+1, err)
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
