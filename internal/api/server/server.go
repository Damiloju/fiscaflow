package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"fiscaflow/internal/api/handlers"
	"fiscaflow/internal/api/middleware"
	"fiscaflow/internal/config"
	"fiscaflow/internal/domain/user"
	"fiscaflow/internal/infrastructure/database"
)

// Server represents the API server
type Server struct {
	config      *config.Config
	logger      *zap.Logger
	userService user.Service
	userHandler *handlers.UserHandler
}

// New creates a new API server instance
func New(cfg *config.Config, logger *zap.Logger) *Server {
	// Initialize database
	dbConfig := &database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}

	db, err := database.NewDatabase(dbConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto-migrate database
	if err := db.AutoMigrate(); err != nil {
		logger.Fatal("Failed to run database migrations", zap.Error(err))
	}

	// Initialize repositories
	userRepo := user.NewRepository(db.GetDB())

	// Initialize services
	userService := user.NewService(userRepo, cfg.JWT.Secret)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, logger)

	return &Server{
		config:      cfg,
		logger:      logger,
		userService: userService,
		userHandler: userHandler,
	}
}

// SetupRoutes configures all API routes
func (s *Server) SetupRoutes(router *gin.Engine) {
	// API v1 group
	v1 := router.Group("/api/v1")

	// User routes
	users := v1.Group("/users")
	{
		users.POST("/register", s.userHandler.Register)
		users.POST("/login", s.userHandler.Login)
		users.POST("/logout", s.userHandler.Logout)
		users.POST("/refresh", s.userHandler.RefreshToken)

		// Protected routes
		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware(s.userService))
		{
			protected.GET("/profile", s.userHandler.GetProfile)
			protected.PUT("/profile", s.userHandler.UpdateProfile)
		}
	}

	s.logger.Info("API routes configured")
}
