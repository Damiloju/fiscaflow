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

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	Service transaction.Service
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(service transaction.Service) *CategoryHandler {
	return &CategoryHandler{Service: service}
}

// RegisterRoutes registers category routes
func (h *CategoryHandler) RegisterRoutes(rg *gin.RouterGroup) {
	cat := rg.Group("/categories")
	cat.POST("", h.CreateCategory)
	cat.GET("", h.ListCategories)
	cat.GET(":id", h.GetCategory)
	cat.PUT(":id", h.UpdateCategory)
	cat.DELETE(":id", h.DeleteCategory)
	cat.GET("/default", h.GetDefaultCategories)
}

// CreateCategory handles POST /categories
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "CreateCategory")
	defer span.End()

	var req transaction.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.Service.CreateCategory(ctx, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}

// GetCategory handles GET /categories/:id
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "GetCategory")
	defer span.End()

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	category, err := h.Service.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, transaction.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// ListCategories handles GET /categories
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "ListCategories")
	defer span.End()

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	categories, err := h.Service.GetCategories(ctx, offset, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetDefaultCategories handles GET /categories/default
func (h *CategoryHandler) GetDefaultCategories(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "GetDefaultCategories")
	defer span.End()

	categories, err := h.Service.GetDefaultCategories(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// UpdateCategory handles PUT /categories/:id
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "UpdateCategory")
	defer span.End()

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var req transaction.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := h.Service.UpdateCategory(ctx, id, &req)
	if err != nil {
		if errors.Is(err, transaction.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// DeleteCategory handles DELETE /categories/:id
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	ctx, span := otel.Tracer("api").Start(c.Request.Context(), "DeleteCategory")
	defer span.End()

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	if err := h.Service.DeleteCategory(ctx, id); err != nil {
		if errors.Is(err, transaction.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
