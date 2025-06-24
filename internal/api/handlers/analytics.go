package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"fiscaflow/internal/domain/analytics"
)

// AnalyticsHandler handles analytics-related HTTP requests
type AnalyticsHandler struct {
	analyticsService analytics.Service
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(analyticsService analytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// CategorizeTransaction handles POST /api/v1/analytics/categorize
func (h *AnalyticsHandler) CategorizeTransaction(c *gin.Context) {
	var req analytics.CategorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.analyticsService.CategorizeTransaction(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categorization": response})
}

// CreateCategorizationRule handles POST /api/v1/analytics/categorization-rules
func (h *AnalyticsHandler) CreateCategorizationRule(c *gin.Context) {
	var req analytics.CreateCategorizationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.analyticsService.CreateCategorizationRule(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"categorization_rule": response})
}

// GetCategorizationRule handles GET /api/v1/analytics/categorization-rules/:id
func (h *AnalyticsHandler) GetCategorizationRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	response, err := h.analyticsService.GetCategorizationRule(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categorization_rule": response})
}

// ListCategorizationRules handles GET /api/v1/analytics/categorization-rules
func (h *AnalyticsHandler) ListCategorizationRules(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if limit > 100 {
		limit = 100
	}

	rules, err := h.analyticsService.ListCategorizationRules(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categorization_rules": rules})
}

// UpdateCategorizationRule handles PUT /api/v1/analytics/categorization-rules/:id
func (h *AnalyticsHandler) UpdateCategorizationRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	var req analytics.UpdateCategorizationRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := h.analyticsService.UpdateCategorizationRule(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categorization_rule": response})
}

// DeleteCategorizationRule handles DELETE /api/v1/analytics/categorization-rules/:id
func (h *AnalyticsHandler) DeleteCategorizationRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule ID"})
		return
	}

	if err := h.analyticsService.DeleteCategorizationRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AnalyzeSpending handles POST /api/v1/analytics/spending
func (h *AnalyticsHandler) AnalyzeSpending(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req analytics.SpendingAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	response, err := h.analyticsService.AnalyzeSpending(c.Request.Context(), userUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analysis": response})
}

// GetSpendingInsights handles GET /api/v1/analytics/spending/insights
func (h *AnalyticsHandler) GetSpendingInsights(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format (YYYY-MM-DD)"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format (YYYY-MM-DD)"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID"})
		return
	}

	insights, err := h.analyticsService.GetSpendingInsights(c.Request.Context(), userUUID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"insights": insights})
}

// RegisterRoutes registers all analytics routes
func (h *AnalyticsHandler) RegisterRoutes(api *gin.RouterGroup) {
	analytics := api.Group("/analytics")
	{
		// Categorization
		analytics.POST("/categorize", h.CategorizeTransaction)

		// Categorization rules
		analytics.POST("/categorization-rules", h.CreateCategorizationRule)
		analytics.GET("/categorization-rules", h.ListCategorizationRules)
		analytics.GET("/categorization-rules/:id", h.GetCategorizationRule)
		analytics.PUT("/categorization-rules/:id", h.UpdateCategorizationRule)
		analytics.DELETE("/categorization-rules/:id", h.DeleteCategorizationRule)

		// Spending analysis
		analytics.POST("/spending", h.AnalyzeSpending)
		analytics.GET("/spending/insights", h.GetSpendingInsights)
	}
}
