package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fiscaflow/internal/domain/user"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := userService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// GetUserIDFromContext extracts user ID from gin context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}

	return userID, true
}

// GetUserEmailFromContext extracts user email from gin context
func GetUserEmailFromContext(c *gin.Context) (string, bool) {
	emailInterface, exists := c.Get("user_email")
	if !exists {
		return "", false
	}

	email, ok := emailInterface.(string)
	if !ok {
		return "", false
	}

	return email, true
}

// GetUserRoleFromContext extracts user role from gin context
func GetUserRoleFromContext(c *gin.Context) (user.UserRole, bool) {
	roleInterface, exists := c.Get("user_role")
	if !exists {
		return "", false
	}

	role, ok := roleInterface.(user.UserRole)
	if !ok {
		return "", false
	}

	return role, true
}
