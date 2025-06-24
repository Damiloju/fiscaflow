package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"

	"fiscaflow/internal/domain/transaction"
)

// TransactionHandler handles transaction-related HTTP requests
type TransactionHandler struct {
	Service transaction.Service
}

// NewTransactionHandler creates a new TransactionHandler
func NewTransactionHandler(service transaction.Service) *TransactionHandler {
	return &TransactionHandler{Service: service}
}

// RegisterRoutes registers transaction routes
func (h *TransactionHandler) RegisterRoutes(rg *gin.RouterGroup) {
	tr := rg.Group("/transactions")
	tr.POST("", h.CreateTransaction)
	tr.GET("", h.ListTransactions)
	tr.GET(":id", h.GetTransaction)
	tr.PUT(":id", h.UpdateTransaction)
	tr.DELETE(":id", h.DeleteTransaction)
}

// CreateTransaction handles POST /transactions
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "CreateTransaction")
	defer span.End()

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	var req transaction.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Service.CreateTransaction(ctx, uid, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

// GetTransaction handles GET /transactions/:id
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "GetTransaction")
	defer span.End()

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	resp, err := h.Service.GetTransaction(ctx, uid, id)
	if err != nil {
		if errors.Is(err, transaction.ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ListTransactions handles GET /transactions
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "ListTransactions")
	defer span.End()

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	resp, err := h.Service.GetTransactions(ctx, uid, offset, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateTransaction handles PUT /transactions/:id
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "UpdateTransaction")
	defer span.End()

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	var req transaction.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.Service.UpdateTransaction(ctx, uid, id, &req)
	if err != nil {
		if errors.Is(err, transaction.ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// DeleteTransaction handles DELETE /transactions/:id
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "DeleteTransaction")
	defer span.End()

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction id"})
		return
	}

	if err := h.Service.DeleteTransaction(ctx, uid, id); err != nil {
		if errors.Is(err, transaction.ErrTransactionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
