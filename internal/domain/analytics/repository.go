package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for analytics data access
type Repository interface {
	// Categorization rule operations
	CreateCategorizationRule(ctx context.Context, rule *CategorizationRule) error
	GetCategorizationRuleByID(ctx context.Context, id uuid.UUID) (*CategorizationRule, error)
	GetCategorizationRules(ctx context.Context, offset, limit int) ([]CategorizationRule, error)
	GetActiveCategorizationRules(ctx context.Context) ([]CategorizationRule, error)
	UpdateCategorizationRule(ctx context.Context, rule *CategorizationRule) error
	DeleteCategorizationRule(ctx context.Context, id uuid.UUID) error

	// Category operations
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)

	// Transaction operations for ML
	GetSimilarTransactions(ctx context.Context, description string, limit int) ([]Transaction, error)
	GetTransactionsByPeriod(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]Transaction, error)

	// Spending analysis operations
	CreateSpendingAnalysis(ctx context.Context, analysis *SpendingAnalysis) error
	GetSpendingAnalysisByID(ctx context.Context, id uuid.UUID) (*SpendingAnalysis, error)
	GetSpendingAnalysisByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*SpendingAnalysis, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new analytics repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateCategorizationRule creates a new categorization rule
func (r *repository) CreateCategorizationRule(ctx context.Context, rule *CategorizationRule) error {
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(rule).Error
}

// GetCategorizationRuleByID retrieves a categorization rule by ID
func (r *repository) GetCategorizationRuleByID(ctx context.Context, id uuid.UUID) (*CategorizationRule, error) {
	var rule CategorizationRule
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("categorization rule not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get categorization rule: %w", err)
	}
	return &rule, nil
}

// GetCategorizationRules retrieves categorization rules with pagination
func (r *repository) GetCategorizationRules(ctx context.Context, offset, limit int) ([]CategorizationRule, error) {
	var rules []CategorizationRule
	err := r.db.WithContext(ctx).
		Order("priority DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&rules).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get categorization rules: %w", err)
	}

	return rules, nil
}

// GetActiveCategorizationRules retrieves all active categorization rules
func (r *repository) GetActiveCategorizationRules(ctx context.Context) ([]CategorizationRule, error) {
	var rules []CategorizationRule
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("priority DESC, created_at DESC").
		Find(&rules).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active categorization rules: %w", err)
	}

	return rules, nil
}

// UpdateCategorizationRule updates a categorization rule
func (r *repository) UpdateCategorizationRule(ctx context.Context, rule *CategorizationRule) error {
	rule.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Save(rule)
	if result.Error != nil {
		return fmt.Errorf("failed to update categorization rule: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("categorization rule not found")
	}

	return nil
}

// DeleteCategorizationRule deletes a categorization rule
func (r *repository) DeleteCategorizationRule(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&CategorizationRule{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete categorization rule: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("categorization rule not found")
	}

	return nil
}

// GetCategoryByID retrieves a category by ID
func (r *repository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

// GetSimilarTransactions retrieves transactions similar to the given description
func (r *repository) GetSimilarTransactions(ctx context.Context, description string, limit int) ([]Transaction, error) {
	var transactions []Transaction

	// Simple similarity search using LIKE
	searchTerm := "%" + description + "%"

	err := r.db.WithContext(ctx).
		Where("description ILIKE ? OR merchant ILIKE ?", searchTerm, searchTerm).
		Order("created_at DESC").
		Limit(limit).
		Find(&transactions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get similar transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionsByPeriod retrieves transactions for a user within a date range
func (r *repository) GetTransactionsByPeriod(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND transaction_date >= ? AND transaction_date <= ?", userID, startDate, endDate).
		Order("transaction_date DESC").
		Find(&transactions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by period: %w", err)
	}

	return transactions, nil
}

// CreateSpendingAnalysis creates a new spending analysis
func (r *repository) CreateSpendingAnalysis(ctx context.Context, analysis *SpendingAnalysis) error {
	analysis.CreatedAt = time.Now()
	analysis.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(analysis).Error
}

// GetSpendingAnalysisByID retrieves a spending analysis by ID
func (r *repository) GetSpendingAnalysisByID(ctx context.Context, id uuid.UUID) (*SpendingAnalysis, error) {
	var analysis SpendingAnalysis
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&analysis).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("spending analysis not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get spending analysis: %w", err)
	}
	return &analysis, nil
}

// GetSpendingAnalysisByUser retrieves spending analysis for a user within a date range
func (r *repository) GetSpendingAnalysisByUser(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*SpendingAnalysis, error) {
	var analysis SpendingAnalysis
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND period_start >= ? AND period_end <= ?", userID, startDate, endDate).
		Order("created_at DESC").
		First(&analysis).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("spending analysis not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get spending analysis: %w", err)
	}
	return &analysis, nil
}

// Category represents a transaction category (imported from transaction domain)
type Category struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Color       string     `json:"color"`
	ParentID    *uuid.UUID `json:"parent_id" gorm:"type:uuid"`
	IsDefault   bool       `json:"is_default" gorm:"default:false"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Transaction represents a transaction (imported from transaction domain)
type Transaction struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	FamilyID   *uuid.UUID `json:"family_id" gorm:"type:uuid"`
	AccountID  uuid.UUID  `json:"account_id" gorm:"type:uuid;not null"`
	CategoryID *uuid.UUID `json:"category_id" gorm:"type:uuid"`

	Amount      float64 `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency    string  `json:"currency" gorm:"default:'USD'"`
	Description string  `json:"description" gorm:"not null"`
	Merchant    string  `json:"merchant"`
	Location    string  `json:"location" gorm:"type:jsonb"`

	TransactionDate time.Time  `json:"transaction_date" gorm:"not null"`
	PostedDate      *time.Time `json:"posted_date"`
	Status          string     `json:"status" gorm:"default:'pending'"`

	CategorizationSource     string   `json:"categorization_source" gorm:"default:'manual'"`
	CategorizationConfidence *float64 `json:"categorization_confidence"`

	Tags       []string `json:"tags" gorm:"type:text[]"`
	Notes      string   `json:"notes"`
	ReceiptURL string   `json:"receipt_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}

// TableName specifies the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}
