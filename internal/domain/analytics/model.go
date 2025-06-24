package analytics

import (
	"time"

	"github.com/google/uuid"
)

// CategorizationModel represents a machine learning model for transaction categorization
type CategorizationModel struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Version   string    `json:"version" gorm:"not null"`
	ModelType string    `json:"model_type" gorm:"not null"` // "keyword", "mlp", "transformer"
	Accuracy  float64   `json:"accuracy" gorm:"type:decimal(5,4)"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	ModelData string    `json:"model_data" gorm:"type:jsonb"` // Serialized model
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategorizationRule represents a rule-based categorization rule
type CategorizationRule struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CategoryID  uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`
	Pattern     string    `json:"pattern" gorm:"not null"`      // Regex pattern or keyword
	PatternType string    `json:"pattern_type" gorm:"not null"` // "regex", "keyword", "exact"
	Priority    int       `json:"priority" gorm:"default:0"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SpendingAnalysis represents spending analysis for a user
type SpendingAnalysis struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	PeriodStart       time.Time `json:"period_start" gorm:"not null"`
	PeriodEnd         time.Time `json:"period_end" gorm:"not null"`
	TotalSpent        float64   `json:"total_spent" gorm:"type:decimal(15,2)"`
	TotalIncome       float64   `json:"total_income" gorm:"type:decimal(15,2)"`
	NetAmount         float64   `json:"net_amount" gorm:"type:decimal(15,2)"`
	CategoryBreakdown string    `json:"category_breakdown" gorm:"type:jsonb"`
	TopCategories     string    `json:"top_categories" gorm:"type:jsonb"`
	SpendingTrends    string    `json:"spending_trends" gorm:"type:jsonb"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// CategorizationRequest represents a request to categorize a transaction
type CategorizationRequest struct {
	Description string  `json:"description" binding:"required"`
	Merchant    string  `json:"merchant"`
	Amount      float64 `json:"amount"`
	Location    string  `json:"location"`
}

// CategorizationResponse represents the categorization result
type CategorizationResponse struct {
	CategoryID            uuid.UUID            `json:"category_id"`
	CategoryName          string               `json:"category_name"`
	Confidence            float64              `json:"confidence"`
	CategorizationSource  string               `json:"categorization_source"` // "rule", "ml", "manual"
	MatchedPattern        string               `json:"matched_pattern,omitempty"`
	AlternativeCategories []CategorySuggestion `json:"alternative_categories,omitempty"`
}

// CategorySuggestion represents a suggested category with confidence
type CategorySuggestion struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Confidence   float64   `json:"confidence"`
	Reason       string    `json:"reason"`
}

// SpendingInsight represents a spending insight
type SpendingInsight struct {
	Type        string                 `json:"type"` // "trend", "anomaly", "pattern", "recommendation"
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"` // "low", "medium", "high"
	Data        map[string]interface{} `json:"data"`
	CreatedAt   time.Time              `json:"created_at"`
}

// CreateCategorizationRuleRequest represents a request to create a categorization rule
type CreateCategorizationRuleRequest struct {
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Pattern     string    `json:"pattern" binding:"required"`
	PatternType string    `json:"pattern_type" binding:"required"`
	Priority    int       `json:"priority"`
}

// UpdateCategorizationRuleRequest represents a request to update a categorization rule
type UpdateCategorizationRuleRequest struct {
	Pattern     *string `json:"pattern"`
	PatternType *string `json:"pattern_type"`
	Priority    *int    `json:"priority"`
	IsActive    *bool   `json:"is_active"`
}

// CategorizationRuleResponse represents a categorization rule response
type CategorizationRuleResponse struct {
	ID           uuid.UUID `json:"id"`
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Pattern      string    `json:"pattern"`
	PatternType  string    `json:"pattern_type"`
	Priority     int       `json:"priority"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SpendingAnalysisRequest represents a request for spending analysis
type SpendingAnalysisRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	GroupBy   string    `json:"group_by"` // "day", "week", "month", "category"
}

// SpendingAnalysisResponse represents a spending analysis response
type SpendingAnalysisResponse struct {
	PeriodStart       time.Time          `json:"period_start"`
	PeriodEnd         time.Time          `json:"period_end"`
	TotalSpent        float64            `json:"total_spent"`
	TotalIncome       float64            `json:"total_income"`
	NetAmount         float64            `json:"net_amount"`
	CategoryBreakdown []CategorySpending `json:"category_breakdown"`
	TopCategories     []CategorySpending `json:"top_categories"`
	SpendingTrends    []SpendingTrend    `json:"spending_trends"`
	Insights          []SpendingInsight  `json:"insights"`
}

// CategorySpending represents spending for a category
type CategorySpending struct {
	CategoryID       uuid.UUID `json:"category_id"`
	CategoryName     string    `json:"category_name"`
	Amount           float64   `json:"amount"`
	Percentage       float64   `json:"percentage"`
	TransactionCount int       `json:"transaction_count"`
}

// SpendingTrend represents a spending trend
type SpendingTrend struct {
	Period string  `json:"period"`
	Amount float64 `json:"amount"`
	Change float64 `json:"change"` // Percentage change from previous period
	Trend  string  `json:"trend"`  // "increasing", "decreasing", "stable"
}

// TableName specifies the table name for CategorizationModel
func (CategorizationModel) TableName() string {
	return "categorization_models"
}

// TableName specifies the table name for CategorizationRule
func (CategorizationRule) TableName() string {
	return "categorization_rules"
}

// TableName specifies the table name for SpendingAnalysis
func (SpendingAnalysis) TableName() string {
	return "spending_analyses"
}
