package routes

import (
	"fiscaflow/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes registers all API routes
func RegisterAPIRoutes(r *gin.Engine, userHandler *handlers.UserHandler, transactionHandler *handlers.TransactionHandler) {
	api := r.Group("/api/v1")

	// User routes (assume already registered)
	userHandler.RegisterRoutes(api)

	// Transaction routes (protected)
	api.Use(AuthMiddleware()) // Assume AuthMiddleware sets user_id in context
	transactionHandler.RegisterRoutes(api)
}

// AuthMiddleware is a placeholder for the actual authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This should set user_id in context after validating JWT
		// For now, assume user_id is set for testing
		// c.Set("user_id", uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		c.Next()
	}
}
