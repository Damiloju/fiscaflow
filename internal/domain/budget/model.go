package budget

import (
	"time"

	"github.com/google/uuid"
)

// Budget represents a user's budget
type Budget struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	FamilyID    *uuid.UUID `json:"family_id" gorm:"type:uuid"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`

	PeriodType PeriodType `json:"period_type" gorm:"not null"`
	StartDate  time.Time  `json:"start_date" gorm:"not null"`
	EndDate    *time.Time `json:"end_date"`

	TotalAmount float64 `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	Currency    string  `json:"currency" gorm:"default:'USD'"`

	IsActive bool   `json:"is_active" gorm:"default:true"`
	Settings string `json:"settings" gorm:"type:jsonb;default:'{}'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BudgetCategory represents a category allocation within a budget
type BudgetCategory struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BudgetID   uuid.UUID `json:"budget_id" gorm:"type:uuid;not null"`
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null"`

	AllocatedAmount float64 `json:"allocated_amount" gorm:"type:decimal(15,2);not null"`
	SpentAmount     float64 `json:"spent_amount" gorm:"type:decimal(15,2);default:0.00"`

	AlertThreshold float64 `json:"alert_threshold" gorm:"type:decimal(3,2);default:0.80"`
	IsActive       bool    `json:"is_active" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PeriodType represents the budget period type
type PeriodType string

const (
	PeriodTypeMonthly   PeriodType = "monthly"
	PeriodTypeQuarterly PeriodType = "quarterly"
	PeriodTypeYearly    PeriodType = "yearly"
	PeriodTypeCustom    PeriodType = "custom"
)

// CreateBudgetRequest represents a request to create a new budget
type CreateBudgetRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	PeriodType  PeriodType `json:"period_type" binding:"required"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	TotalAmount float64    `json:"total_amount" binding:"required"`
	Currency    string     `json:"currency"`
	Settings    string     `json:"settings"`
}

// UpdateBudgetRequest represents a request to update a budget
type UpdateBudgetRequest struct {
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	PeriodType  *PeriodType `json:"period_type"`
	StartDate   *time.Time  `json:"start_date"`
	EndDate     *time.Time  `json:"end_date"`
	TotalAmount *float64    `json:"total_amount"`
	Currency    *string     `json:"currency"`
	IsActive    *bool       `json:"is_active"`
	Settings    *string     `json:"settings"`
}

// CreateBudgetCategoryRequest represents a request to create a budget category
type CreateBudgetCategoryRequest struct {
	CategoryID      uuid.UUID `json:"category_id" binding:"required"`
	AllocatedAmount float64   `json:"allocated_amount" binding:"required"`
	AlertThreshold  float64   `json:"alert_threshold"`
}

// UpdateBudgetCategoryRequest represents a request to update a budget category
type UpdateBudgetCategoryRequest struct {
	AllocatedAmount *float64 `json:"allocated_amount"`
	AlertThreshold  *float64 `json:"alert_threshold"`
	IsActive        *bool    `json:"is_active"`
}

// BudgetResponse represents a budget response
type BudgetResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	FamilyID    *uuid.UUID `json:"family_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	PeriodType  PeriodType `json:"period_type"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	TotalAmount float64    `json:"total_amount"`
	Currency    string     `json:"currency"`
	IsActive    bool       `json:"is_active"`
	Settings    string     `json:"settings"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// BudgetCategoryResponse represents a budget category response
type BudgetCategoryResponse struct {
	ID              uuid.UUID `json:"id"`
	BudgetID        uuid.UUID `json:"budget_id"`
	CategoryID      uuid.UUID `json:"category_id"`
	AllocatedAmount float64   `json:"allocated_amount"`
	SpentAmount     float64   `json:"spent_amount"`
	AlertThreshold  float64   `json:"alert_threshold"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BudgetSummary represents a budget summary with spending analysis
type BudgetSummary struct {
	Budget           *BudgetResponse          `json:"budget"`
	Categories       []BudgetCategoryResponse `json:"categories"`
	TotalAllocated   float64                  `json:"total_allocated"`
	TotalSpent       float64                  `json:"total_spent"`
	RemainingAmount  float64                  `json:"remaining_amount"`
	SpendingProgress float64                  `json:"spending_progress"` // Percentage spent
	Alerts           []BudgetAlert            `json:"alerts"`
}

// BudgetAlert represents a budget alert
type BudgetAlert struct {
	CategoryID      uuid.UUID `json:"category_id"`
	CategoryName    string    `json:"category_name"`
	AllocatedAmount float64   `json:"allocated_amount"`
	SpentAmount     float64   `json:"spent_amount"`
	Threshold       float64   `json:"threshold"`
	AlertType       string    `json:"alert_type"` // "warning", "critical", "over_budget"
	Message         string    `json:"message"`
}

// TableName specifies the table name for Budget
func (Budget) TableName() string {
	return "budgets"
}

// TableName specifies the table name for BudgetCategory
func (BudgetCategory) TableName() string {
	return "budget_categories"
}
