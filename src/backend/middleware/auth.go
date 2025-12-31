package middleware

import (
	"strings"
	"teamitmivhs/work-order-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store claims in context untuk digunakan di controller
		c.Set("user_id", claims.ID)
		c.Set("user_name", claims.Name)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware checks if user is admin
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

// OperatorMiddleware checks if user is operator
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

// GetUserIDFromContext gets user ID from context
func GetUserIDFromContext(c *gin.Context) int {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	id, ok := userID.(int)
	if !ok {
		return 0
	}
	return id
}

// GetUserRoleFromContext gets user role from context
func GetUserRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("user_role")
	if !exists {
		return ""
	}
	roleStr, ok := role.(string)
	if !ok {
		return ""
	}
	return roleStr
}
