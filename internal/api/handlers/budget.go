package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fiscaflow/internal/domain/budget"
)

// BudgetHandler handles budget-related HTTP requests
type BudgetHandler struct {
	budgetService budget.Service
}

// NewBudgetHandler creates a new budget handler
func NewBudgetHandler(budgetService budget.Service) *BudgetHandler {
	return &BudgetHandler{
		budgetService: budgetService,
	}
}

// CreateBudget handles POST /api/v1/budgets
func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req budget.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	budgetResponse, err := h.budgetService.CreateBudget(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"budget": budgetResponse})
}

// GetBudget handles GET /api/v1/budgets/:id
func (h *BudgetHandler) GetBudget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	budgetResponse, err := h.budgetService.GetBudget(c.Request.Context(), userUUID, budgetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budget": budgetResponse})
}

// ListBudgets handles GET /api/v1/budgets
func (h *BudgetHandler) ListBudgets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if limit > 100 {
		limit = 100
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	budgets, err := h.budgetService.ListBudgets(c.Request.Context(), userUUID, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budgets": budgets})
}

// UpdateBudget handles PUT /api/v1/budgets/:id
func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	var req budget.UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	budgetResponse, err := h.budgetService.UpdateBudget(c.Request.Context(), userUUID, budgetID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budget": budgetResponse})
}

// DeleteBudget handles DELETE /api/v1/budgets/:id
func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.budgetService.DeleteBudget(c.Request.Context(), userUUID, budgetID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetBudgetSummary handles GET /api/v1/budgets/:id/summary
func (h *BudgetHandler) GetBudgetSummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	summary, err := h.budgetService.GetBudgetSummary(c.Request.Context(), userUUID, budgetID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

// AddBudgetCategory handles POST /api/v1/budgets/:id/categories
func (h *BudgetHandler) AddBudgetCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	var req budget.CreateBudgetCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	categoryResponse, err := h.budgetService.AddBudgetCategory(c.Request.Context(), userUUID, budgetID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"budget_category": categoryResponse})
}

// GetBudgetCategory handles GET /api/v1/budgets/:id/categories/:categoryId
func (h *BudgetHandler) GetBudgetCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	categoryResponse, err := h.budgetService.GetBudgetCategory(c.Request.Context(), userUUID, budgetID, categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budget_category": categoryResponse})
}

// ListBudgetCategories handles GET /api/v1/budgets/:id/categories
func (h *BudgetHandler) ListBudgetCategories(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	categories, err := h.budgetService.ListBudgetCategories(c.Request.Context(), userUUID, budgetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budget_categories": categories})
}

// UpdateBudgetCategory handles PUT /api/v1/budgets/:id/categories/:categoryId
func (h *BudgetHandler) UpdateBudgetCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req budget.UpdateBudgetCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	categoryResponse, err := h.budgetService.UpdateBudgetCategory(c.Request.Context(), userUUID, budgetID, categoryID, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"budget_category": categoryResponse})
}

// DeleteBudgetCategory handles DELETE /api/v1/budgets/:id/categories/:categoryId
func (h *BudgetHandler) DeleteBudgetCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid budget ID"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("categoryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.budgetService.DeleteBudgetCategory(c.Request.Context(), userUUID, budgetID, categoryID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RegisterRoutes registers all budget routes
func (h *BudgetHandler) RegisterRoutes(api *gin.RouterGroup) {
	budgets := api.Group("/budgets")
	{
		budgets.POST("", h.CreateBudget)
		budgets.GET("", h.ListBudgets)
		budgets.GET("/:id", h.GetBudget)
		budgets.PUT("/:id", h.UpdateBudget)
		budgets.DELETE("/:id", h.DeleteBudget)
		budgets.GET("/:id/summary", h.GetBudgetSummary)

		// Budget categories
		budgets.POST("/:id/categories", h.AddBudgetCategory)
		budgets.GET("/:id/categories", h.ListBudgetCategories)
		budgets.GET("/:id/categories/:categoryId", h.GetBudgetCategory)
		budgets.PUT("/:id/categories/:categoryId", h.UpdateBudgetCategory)
		budgets.DELETE("/:id/categories/:categoryId", h.DeleteBudgetCategory)
	}
}
