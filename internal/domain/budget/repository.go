package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for budget data access
type Repository interface {
	// Budget operations
	Create(ctx context.Context, budget *Budget) error
	GetByID(ctx context.Context, id uuid.UUID) (*Budget, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]Budget, error)
	Update(ctx context.Context, budget *Budget) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Budget category operations
	CreateCategory(ctx context.Context, budgetCategory *BudgetCategory) error
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*BudgetCategory, error)
	GetCategoriesByBudgetID(ctx context.Context, budgetID uuid.UUID) ([]BudgetCategory, error)
	UpdateCategory(ctx context.Context, budgetCategory *BudgetCategory) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Budget analysis operations
	GetBudgetSummary(ctx context.Context, budgetID uuid.UUID) (*BudgetSummary, error)
	UpdateSpentAmount(ctx context.Context, budgetID, categoryID uuid.UUID, amount float64) error
	GetActiveBudgetsByUser(ctx context.Context, userID uuid.UUID) ([]Budget, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new budget repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create creates a new budget
func (r *repository) Create(ctx context.Context, budget *Budget) error {
	budget.CreatedAt = time.Now()
	budget.UpdatedAt = time.Now()

	if budget.Currency == "" {
		budget.Currency = "USD"
	}

	return r.db.WithContext(ctx).Create(budget).Error
}

// GetByID retrieves a budget by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Budget, error) {
	var budget Budget
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&budget).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("budget not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}
	return &budget, nil
}

// GetByUserID retrieves budgets for a user with pagination
func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]Budget, error) {
	var budgets []Budget
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&budgets).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}

	return budgets, nil
}

// Update updates a budget
func (r *repository) Update(ctx context.Context, budget *Budget) error {
	budget.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Save(budget)
	if result.Error != nil {
		return fmt.Errorf("failed to update budget: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("budget not found")
	}

	return nil
}

// Delete deletes a budget
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&Budget{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete budget: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("budget not found")
	}

	return nil
}

// CreateCategory creates a new budget category
func (r *repository) CreateCategory(ctx context.Context, budgetCategory *BudgetCategory) error {
	budgetCategory.CreatedAt = time.Now()
	budgetCategory.UpdatedAt = time.Now()

	if budgetCategory.AlertThreshold == 0 {
		budgetCategory.AlertThreshold = 0.80 // Default 80% threshold
	}

	return r.db.WithContext(ctx).Create(budgetCategory).Error
}

// GetCategoryByID retrieves a budget category by ID
func (r *repository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*BudgetCategory, error) {
	var budgetCategory BudgetCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&budgetCategory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("budget category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get budget category: %w", err)
	}
	return &budgetCategory, nil
}

// GetCategoriesByBudgetID retrieves all categories for a budget
func (r *repository) GetCategoriesByBudgetID(ctx context.Context, budgetID uuid.UUID) ([]BudgetCategory, error) {
	var budgetCategories []BudgetCategory
	err := r.db.WithContext(ctx).
		Where("budget_id = ?", budgetID).
		Order("created_at ASC").
		Find(&budgetCategories).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get budget categories: %w", err)
	}

	return budgetCategories, nil
}

// UpdateCategory updates a budget category
func (r *repository) UpdateCategory(ctx context.Context, budgetCategory *BudgetCategory) error {
	budgetCategory.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Save(budgetCategory)
	if result.Error != nil {
		return fmt.Errorf("failed to update budget category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("budget category not found")
	}

	return nil
}

// DeleteCategory deletes a budget category
func (r *repository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&BudgetCategory{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete budget category: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("budget category not found")
	}

	return nil
}

// GetBudgetSummary retrieves a comprehensive budget summary
func (r *repository) GetBudgetSummary(ctx context.Context, budgetID uuid.UUID) (*BudgetSummary, error) {
	// Get the budget
	budget, err := r.GetByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}

	// Get budget categories
	categories, err := r.GetCategoriesByBudgetID(ctx, budgetID)
	if err != nil {
		return nil, err
	}

	// Calculate totals
	var totalAllocated, totalSpent float64
	for _, category := range categories {
		totalAllocated += category.AllocatedAmount
		totalSpent += category.SpentAmount
	}

	// Calculate spending progress
	var spendingProgress float64
	if totalAllocated > 0 {
		spendingProgress = (totalSpent / totalAllocated) * 100
	}

	// Generate alerts
	alerts := r.generateAlerts(categories)

	// Convert to response types
	budgetResponse := &BudgetResponse{
		ID:          budget.ID,
		UserID:      budget.UserID,
		FamilyID:    budget.FamilyID,
		Name:        budget.Name,
		Description: budget.Description,
		PeriodType:  budget.PeriodType,
		StartDate:   budget.StartDate,
		EndDate:     budget.EndDate,
		TotalAmount: budget.TotalAmount,
		Currency:    budget.Currency,
		IsActive:    budget.IsActive,
		Settings:    budget.Settings,
		CreatedAt:   budget.CreatedAt,
		UpdatedAt:   budget.UpdatedAt,
	}

	categoryResponses := make([]BudgetCategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = BudgetCategoryResponse(category)
	}

	return &BudgetSummary{
		Budget:           budgetResponse,
		Categories:       categoryResponses,
		TotalAllocated:   totalAllocated,
		TotalSpent:       totalSpent,
		RemainingAmount:  totalAllocated - totalSpent,
		SpendingProgress: spendingProgress,
		Alerts:           alerts,
	}, nil
}

// UpdateSpentAmount updates the spent amount for a budget category
func (r *repository) UpdateSpentAmount(ctx context.Context, budgetID, categoryID uuid.UUID, amount float64) error {
	// This would typically be called when a transaction is created/updated
	// For now, we'll just update the spent amount directly
	result := r.db.WithContext(ctx).
		Model(&BudgetCategory{}).
		Where("budget_id = ? AND category_id = ?", budgetID, categoryID).
		Update("spent_amount", amount)

	if result.Error != nil {
		return fmt.Errorf("failed to update spent amount: %w", result.Error)
	}

	return nil
}

// GetActiveBudgetsByUser retrieves active budgets for a user
func (r *repository) GetActiveBudgetsByUser(ctx context.Context, userID uuid.UUID) ([]Budget, error) {
	var budgets []Budget
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&budgets).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active budgets: %w", err)
	}

	return budgets, nil
}

// generateAlerts generates budget alerts based on spending thresholds
func (r *repository) generateAlerts(categories []BudgetCategory) []BudgetAlert {
	var alerts []BudgetAlert

	for _, category := range categories {
		if category.AllocatedAmount <= 0 {
			continue
		}

		spendingRatio := category.SpentAmount / category.AllocatedAmount

		var alertType string
		var message string

		if spendingRatio >= 1.0 {
			// Over budget
			alertType = "over_budget"
			message = fmt.Sprintf("You've exceeded your budget for this category by $%.2f",
				category.SpentAmount-category.AllocatedAmount)
		} else if spendingRatio >= category.AlertThreshold {
			// Warning threshold reached
			alertType = "warning"
			message = fmt.Sprintf("You've used %.1f%% of your budget for this category",
				spendingRatio*100)
		} else if spendingRatio >= 0.9 {
			// Critical threshold (90%)
			alertType = "critical"
			message = fmt.Sprintf("You're approaching your budget limit (%.1f%% used)",
				spendingRatio*100)
		}

		if alertType != "" {
			alerts = append(alerts, BudgetAlert{
				CategoryID:      category.CategoryID,
				AllocatedAmount: category.AllocatedAmount,
				SpentAmount:     category.SpentAmount,
				Threshold:       category.AlertThreshold,
				AlertType:       alertType,
				Message:         message,
			})
		}
	}

	return alerts
}
