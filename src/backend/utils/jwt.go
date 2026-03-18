package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// FIX: jwtSecret tidak lagi di-evaluate saat package init.
// Sebelumnya var jwtSecret = []byte(getJWTSecret()) dipanggil sekali saat startup,
// sehingga jika JWT_SECRET belum di-set saat init (misalnya lewat secret manager),
// akan pakai default hardcoded tanpa peringatan apapun.
//
// Sekarang menggunakan sync.Once agar secret dibaca saat pertama kali dibutuhkan,
// sekaligus hanya dibaca sekali (thread-safe, efisien).
var (
	jwtSecret     []byte
	jwtSecretOnce sync.Once
)

func getJWTSecret() []byte {
	jwtSecretOnce.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			// FIX: log WARNING yang jelas agar tidak luput di production
			log.Println("[WARNING] JWT_SECRET environment variable is not set. " +
				"Using default secret — this is UNSAFE in production!")
			secret = "your-secret-key-change-in-production-12345"
		}
		jwtSecret = []byte(secret)
	})
	return jwtSecret
}

// GenerateToken menghasilkan JWT token untuk user dengan expiry 24 jam
func GenerateToken(id int, name, role string) (string, error) {
	if id <= 0 {
		return "", fmt.Errorf("invalid user ID: %d", id)
	}
	if name == "" {
		return "", fmt.Errorf("user name cannot be empty")
	}
	if role == "" {
		return "", fmt.Errorf("user role cannot be empty")
	}

	now := time.Now()
	claims := &Claims{
		ID:   id,
		Name: name,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ValidateToken memvalidasi dan mem-parse JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token string cannot be empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma signing sesuai ekspektasi
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return getJWTSecret(), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
