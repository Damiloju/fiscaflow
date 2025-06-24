package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"fiscaflow/internal/api/handlers"
	"fiscaflow/internal/api/middleware"
	"fiscaflow/internal/config"
	"fiscaflow/internal/domain/transaction"
	"fiscaflow/internal/domain/user"
	"fiscaflow/internal/infrastructure/database"
)

// Server represents the API server
type Server struct {
	config             *config.Config
	logger             *zap.Logger
	userService        user.Service
	userHandler        *handlers.UserHandler
	transactionService transaction.Service
	transactionHandler *handlers.TransactionHandler
	categoryHandler    *handlers.CategoryHandler
	accountHandler     *handlers.AccountHandler
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
	transactionRepo := transaction.NewRepository(db.GetDB())

	// Initialize services
	userService := user.NewService(userRepo, cfg.JWT.Secret)
	transactionService := transaction.NewService(transactionRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, logger)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	categoryHandler := handlers.NewCategoryHandler(transactionService)
	accountHandler := handlers.NewAccountHandler(transactionService)

	return &Server{
		config:             cfg,
		logger:             logger,
		userService:        userService,
		userHandler:        userHandler,
		transactionService: transactionService,
		transactionHandler: transactionHandler,
		categoryHandler:    categoryHandler,
		accountHandler:     accountHandler,
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

	// Transaction routes (protected)
	transactions := v1.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware(s.userService))
	{
		transactions.POST("", s.transactionHandler.CreateTransaction)
		transactions.GET("", s.transactionHandler.ListTransactions)
		transactions.GET(":id", s.transactionHandler.GetTransaction)
		transactions.PUT(":id", s.transactionHandler.UpdateTransaction)
		transactions.DELETE(":id", s.transactionHandler.DeleteTransaction)
	}

	// Category routes
	categories := v1.Group("/categories")
	{
		// Public routes
		categories.GET("/default", s.categoryHandler.GetDefaultCategories)

		// Protected routes
		protectedCategories := categories.Group("")
		protectedCategories.Use(middleware.AuthMiddleware(s.userService))
		{
			protectedCategories.POST("", s.categoryHandler.CreateCategory)
			protectedCategories.GET("", s.categoryHandler.ListCategories)
			protectedCategories.GET(":id", s.categoryHandler.GetCategory)
			protectedCategories.PUT(":id", s.categoryHandler.UpdateCategory)
			protectedCategories.DELETE(":id", s.categoryHandler.DeleteCategory)
		}
	}

	// Account routes (protected)
	accounts := v1.Group("/accounts")
	accounts.Use(middleware.AuthMiddleware(s.userService))
	{
		accounts.POST("", s.accountHandler.CreateAccount)
		accounts.GET("", s.accountHandler.ListAccounts)
		accounts.GET(":id", s.accountHandler.GetAccount)
		accounts.PUT(":id", s.accountHandler.UpdateAccount)
		accounts.DELETE(":id", s.accountHandler.DeleteAccount)
	}

	s.logger.Info("API routes configured")
}
