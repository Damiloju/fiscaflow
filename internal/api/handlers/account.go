package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"fiscaflow/internal/domain/transaction"
)

// AccountHandler handles account-related HTTP requests
type AccountHandler struct {
	Service transaction.Service
}

// NewAccountHandler creates a new AccountHandler
func NewAccountHandler(service transaction.Service) *AccountHandler {
	return &AccountHandler{Service: service}
}

// RegisterRoutes registers account routes
func (h *AccountHandler) RegisterRoutes(rg *gin.RouterGroup) {
	acc := rg.Group("/accounts")
	acc.POST("", h.CreateAccount)
	acc.GET("", h.ListAccounts)
	acc.GET(":id", h.GetAccount)
	acc.PUT(":id", h.UpdateAccount)
	acc.DELETE(":id", h.DeleteAccount)
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "CreateAccount")
	defer span.End()

	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	var req transaction.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.Service.CreateAccount(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, account)
}

// GetAccount handles GET /accounts/:id
func (h *AccountHandler) GetAccount(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "GetAccount")
	defer span.End()

	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	account, err := h.Service.GetAccount(ctx, userID, id)
	if err != nil {
		if errors.Is(err, transaction.ErrAccountNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

// ListAccounts handles GET /accounts
func (h *AccountHandler) ListAccounts(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "ListAccounts")
	defer span.End()

	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	accounts, err := h.Service.GetAccounts(ctx, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

// UpdateAccount handles PUT /accounts/:id
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "UpdateAccount")
	defer span.End()

	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	var req transaction.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.Service.UpdateAccount(ctx, userID, id, &req)
	if err != nil {
		if errors.Is(err, transaction.ErrAccountNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, account)
}

// DeleteAccount handles DELETE /accounts/:id
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "DeleteAccount")
	defer span.End()

	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account id"})
		return
	}

	if err := h.Service.DeleteAccount(ctx, userID, id); err != nil {
		if errors.Is(err, transaction.ErrAccountNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
