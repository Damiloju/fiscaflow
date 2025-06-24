package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"fiscaflow/internal/domain/user"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService user.Service
	logger      *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService user.Service, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *UserHandler) Register(c *gin.Context) {
	ctx, span := otel.Tracer("user").Start(c.Request.Context(), "user.Register")
	defer span.End()

	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	span.SetAttributes(
		attribute.String("user.email", req.Email),
		attribute.String("user.first_name", req.FirstName),
		attribute.String("user.last_name", req.LastName),
	)

	userResponse, err := h.userService.Register(ctx, &req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		switch err {
		case user.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		default:
			h.logger.Error("Failed to register user", zap.Error(err), zap.String("email", req.Email))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		}
		return
	}

	span.SetAttributes(attribute.String("user.id", userResponse.ID.String()))
	span.SetStatus(codes.Ok, "user registered successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    userResponse,
	})
}

// Login handles user login
func (h *UserHandler) Login(c *gin.Context) {
	ctx, span := otel.Tracer("user").Start(c.Request.Context(), "user.Login")
	defer span.End()

	var req user.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	span.SetAttributes(attribute.String("user.email", req.Email))

	loginResponse, err := h.userService.Login(ctx, &req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		switch err {
		case user.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		case user.ErrUserInactive:
			c.JSON(http.StatusForbidden, gin.H{"error": "User account is inactive"})
		default:
			h.logger.Error("Failed to login user", zap.Error(err), zap.String("email", req.Email))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		}
		return
	}

	span.SetAttributes(
		attribute.String("user.id", loginResponse.User.ID.String()),
		attribute.String("user.role", string(loginResponse.User.Role)),
	)
	span.SetStatus(codes.Ok, "user logged in successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    loginResponse,
	})
}

// GetProfile retrieves user profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx, span := otel.Tracer("user").Start(c.Request.Context(), "user.GetProfile")
	defer span.End()

	userID, exists := c.Get("user_id")
	if !exists {
		span.RecordError(user.ErrInvalidToken)
		span.SetStatus(codes.Error, "user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		span.RecordError(user.ErrInvalidToken)
		span.SetStatus(codes.Error, "invalid user_id type")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	span.SetAttributes(attribute.String("user.id", userUUID.String()))

	userResponse, err := h.userService.GetProfile(ctx, userUUID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		switch err {
		case user.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			h.logger.Error("Failed to get user profile", zap.Error(err), zap.String("user_id", userUUID.String()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		}
		return
	}

	span.SetStatus(codes.Ok, "profile retrieved successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"user":    userResponse,
	})
}

// UpdateProfile updates user profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	ctx, span := otel.Tracer("user").Start(c.Request.Context(), "user.UpdateProfile")
	defer span.End()

	userID, exists := c.Get("user_id")
	if !exists {
		span.RecordError(user.ErrInvalidToken)
		span.SetStatus(codes.Error, "user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		span.RecordError(user.ErrInvalidToken)
		span.SetStatus(codes.Error, "invalid user_id type")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	span.SetAttributes(attribute.String("user.id", userUUID.String()))

	userResponse, err := h.userService.UpdateProfile(ctx, userUUID, &req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		switch err {
		case user.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			h.logger.Error("Failed to update user profile", zap.Error(err), zap.String("user_id", userUUID.String()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		}
		return
	}

	span.SetStatus(codes.Ok, "profile updated successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    userResponse,
	})
}

// RefreshToken refreshes access token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	ctx, span := otel.Tracer("user").Start(c.Request.Context(), "user.RefreshToken")
	defer span.End()

	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	span.SetAttributes(attribute.String("refresh_token", req.RefreshToken))

	loginResponse, err := h.userService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		switch err {
		case user.ErrInvalidRefreshToken:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		case user.ErrRefreshTokenExpired:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		default:
			h.logger.Error("Failed to refresh token", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		}
		return
	}

	span.SetAttributes(attribute.String("user.id", loginResponse.User.ID.String()))
	span.SetStatus(codes.Ok, "token refreshed successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data":    loginResponse,
	})
}

// Logout handles user logout
func (h *UserHandler) Logout(c *gin.Context) {
	ctx, span := otel.Tracer("user").Start(c.Request.Context(), "user.Logout")
	defer span.End()

	userID, exists := c.Get("user_id")
	if !exists {
		span.RecordError(user.ErrInvalidToken)
		span.SetStatus(codes.Error, "user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		span.RecordError(user.ErrInvalidToken)
		span.SetStatus(codes.Error, "invalid user_id type")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get session ID from header or query param
	sessionIDStr := c.GetHeader("X-Session-ID")
	if sessionIDStr == "" {
		sessionIDStr = c.Query("session_id")
	}

	if sessionIDStr == "" {
		span.RecordError(user.ErrInvalidToken)
		span.SetStatus(codes.Error, "session_id not provided")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid session_id format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID format"})
		return
	}

	span.SetAttributes(
		attribute.String("user.id", userUUID.String()),
		attribute.String("session.id", sessionID.String()),
	)

	if err := h.userService.Logout(ctx, userUUID, sessionID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.logger.Error("Failed to logout user", zap.Error(err), zap.String("user_id", userUUID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	span.SetStatus(codes.Ok, "user logged out successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// ListUsers handles listing users (admin only)
func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx, span := otel.Tracer("user").Start(c.Request.Context(), "user.ListUsers")
	defer span.End()

	// Get pagination parameters
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if limit > 100 {
		limit = 100
	}

	span.SetAttributes(
		attribute.Int("offset", offset),
		attribute.Int("limit", limit),
	)

	users, err := h.userService.(interface {
		List(ctx context.Context, offset, limit int) ([]user.User, error)
	}).List(ctx, offset, limit)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		h.logger.Error("Failed to list users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list users"})
		return
	}

	// Convert to responses
	userResponses := make([]*user.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = &user.UserResponse{
			ID:            u.ID,
			Email:         u.Email,
			FirstName:     u.FirstName,
			LastName:      u.LastName,
			Phone:         u.Phone,
			DateOfBirth:   u.DateOfBirth,
			Timezone:      u.Timezone,
			Locale:        u.Locale,
			Role:          u.Role,
			Status:        u.Status,
			EmailVerified: u.EmailVerified,
			PhoneVerified: u.PhoneVerified,
			LastLoginAt:   u.LastLoginAt,
			CreatedAt:     u.CreatedAt,
			UpdatedAt:     u.UpdatedAt,
		}
	}

	span.SetStatus(codes.Ok, "users listed successfully")

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"users":   userResponses,
		"pagination": gin.H{
			"offset": offset,
			"limit":  limit,
			"count":  len(userResponses),
		},
	})
}

// RegisterRoutes registers user-related routes
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/users/register", h.Register)
	rg.POST("/users/login", h.Login)
	rg.GET("/users/profile", h.GetProfile)
	rg.PUT("/users/profile", h.UpdateProfile)
	rg.POST("/users/refresh-token", h.RefreshToken)
	rg.POST("/users/logout", h.Logout)
}
