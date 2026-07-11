package middleware

import (
	"strings"
	"teamitmivhs/work-order-backend/repository"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token dari header Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// Format yang diharapkan: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		memberRepo := repository.NewMemberRepository()
		member, err := memberRepo.GetMemberByID(claims.ID)
		if err != nil || member == nil {
			utils.Unauthorized(c, "User is no longer active")
			c.Abort()
			return
		}
		if member.AccountStatus != "" && member.AccountStatus != "active" {
			utils.Forbidden(c, "User account is not active")
			c.Abort()
			return
		}
		if member.MembershipStatus == "alumni" || member.MembershipStatus == "inactive" {
			utils.Forbidden(c, "User account is no longer active")
			c.Abort()
			return
		}

		// Ambil role terbaru dari DB agar token lama tidak tetap membawa privilege.
		c.Set("user_id", claims.ID)
		c.Set("user_name", member.Name)
		c.Set("user_role", member.Role)

		c.Next()
	}
}

// AdminMiddleware memastikan user yang login adalah Admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.Unauthorized(c, "Missing user role information")
			c.Abort()
			return
		}

		if role != "Admin" {
			utils.Forbidden(c, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OperatorMiddleware memastikan user adalah Operator atau Admin
func OperatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.Unauthorized(c, "Missing user role information")
			c.Abort()
			return
		}

		if role != "Operator" && role != "Admin" {
			utils.Forbidden(c, "Operator access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserIDFromContext mengambil user ID dari context.
// FIX: return (int, bool) — pola Go idiomatis yang aman.
// Sebelumnya return 0 sebagai sentinel value yang berbahaya
// karena 0 bisa lolos pengecekan if userID > 0 yang ceroboh.
func GetUserIDFromContext(c *gin.Context) (int, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int)
	if !ok {
		return 0, false
	}
	return id, true
}

// GetUserRoleFromContext mengambil role user dari context.
// FIX: return (string, bool) — konsisten dengan GetUserIDFromContext.
func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	if !ok {
		return "", false
	}
	return roleStr, true
}

// GetUserNameFromContext mengambil nama user dari context.
func GetUserNameFromContext(c *gin.Context) (string, bool) {
	name, exists := c.Get("user_name")
	if !exists {
		return "", false
	}
	nameStr, ok := name.(string)
	if !ok {
		return "", false
	}
	return nameStr, true
}
