package budget

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Service defines the interface for budget business logic
type Service interface {
	// Budget operations
	CreateBudget(ctx context.Context, userID uuid.UUID, req *CreateBudgetRequest) (*BudgetResponse, error)
	GetBudget(ctx context.Context, userID, budgetID uuid.UUID) (*BudgetResponse, error)
	ListBudgets(ctx context.Context, userID uuid.UUID, offset, limit int) ([]BudgetResponse, error)
	UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, req *UpdateBudgetRequest) (*BudgetResponse, error)
	DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error

	// Budget category operations
	AddBudgetCategory(ctx context.Context, userID, budgetID uuid.UUID, req *CreateBudgetCategoryRequest) (*BudgetCategoryResponse, error)
	GetBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID) (*BudgetCategoryResponse, error)
	ListBudgetCategories(ctx context.Context, userID, budgetID uuid.UUID) ([]BudgetCategoryResponse, error)
	UpdateBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID, req *UpdateBudgetCategoryRequest) (*BudgetCategoryResponse, error)
	DeleteBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID) error

	// Budget analysis
	GetBudgetSummary(ctx context.Context, userID, budgetID uuid.UUID) (*BudgetSummary, error)
	UpdateBudgetFromTransaction(ctx context.Context, userID, budgetID, categoryID uuid.UUID, amount float64) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new budget service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateBudget creates a new budget for a user
func (s *service) CreateBudget(ctx context.Context, userID uuid.UUID, req *CreateBudgetRequest) (*BudgetResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.CreateBudget",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_name", req.Name),
			attribute.String("period_type", string(req.PeriodType)),
		),
	)
	defer span.End()

	// Validate request
	if err := s.validateCreateBudgetRequest(req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Create budget
	budget := &Budget{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		PeriodType:  req.PeriodType,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		TotalAmount: req.TotalAmount,
		Currency:    req.Currency,
		Settings:    req.Settings,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, budget); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	span.SetAttributes(attribute.String("budget_id", budget.ID.String()))
	return s.toBudgetResponse(budget), nil
}

// GetBudget retrieves a budget by ID
func (s *service) GetBudget(ctx context.Context, userID, budgetID uuid.UUID) (*BudgetResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.GetBudget",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
		),
	)
	defer span.End()

	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Check ownership
	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return nil, fmt.Errorf("unauthorized access to budget")
	}

	return s.toBudgetResponse(budget), nil
}

// ListBudgets retrieves budgets for a user
func (s *service) ListBudgets(ctx context.Context, userID uuid.UUID, offset, limit int) ([]BudgetResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.ListBudgets",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.Int("offset", offset),
			attribute.Int("limit", limit),
		),
	)
	defer span.End()

	budgets, err := s.repo.GetByUserID(ctx, userID, offset, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	responses := make([]BudgetResponse, len(budgets))
	for i, budget := range budgets {
		responses[i] = *s.toBudgetResponse(&budget)
	}

	span.SetAttributes(attribute.Int("budgets_count", len(responses)))
	return responses, nil
}

// UpdateBudget updates a budget
func (s *service) UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, req *UpdateBudgetRequest) (*BudgetResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.UpdateBudget",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
		),
	)
	defer span.End()

	// Get existing budget
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Check ownership
	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return nil, fmt.Errorf("unauthorized access to budget")
	}

	// Update fields
	if req.Name != nil {
		budget.Name = *req.Name
	}
	if req.Description != nil {
		budget.Description = *req.Description
	}
	if req.PeriodType != nil {
		budget.PeriodType = *req.PeriodType
	}
	if req.StartDate != nil {
		budget.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		budget.EndDate = req.EndDate
	}
	if req.TotalAmount != nil {
		budget.TotalAmount = *req.TotalAmount
	}
	if req.Currency != nil {
		budget.Currency = *req.Currency
	}
	if req.IsActive != nil {
		budget.IsActive = *req.IsActive
	}
	if req.Settings != nil {
		budget.Settings = *req.Settings
	}

	// Validate updated budget
	if err := s.validateBudget(budget); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if err := s.repo.Update(ctx, budget); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return s.toBudgetResponse(budget), nil
}

// DeleteBudget deletes a budget
func (s *service) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	ctx, span := otel.Tracer("").Start(ctx, "budget.DeleteBudget",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
		),
	)
	defer span.End()

	// Get existing budget to check ownership
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Check ownership
	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return fmt.Errorf("unauthorized access to budget")
	}

	if err := s.repo.Delete(ctx, budgetID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// AddBudgetCategory adds a category to a budget
func (s *service) AddBudgetCategory(ctx context.Context, userID, budgetID uuid.UUID, req *CreateBudgetCategoryRequest) (*BudgetCategoryResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.AddBudgetCategory",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
			attribute.String("category_id", req.CategoryID.String()),
		),
	)
	defer span.End()

	// Verify budget ownership
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return nil, fmt.Errorf("unauthorized access to budget")
	}

	// Create budget category
	budgetCategory := &BudgetCategory{
		BudgetID:        budgetID,
		CategoryID:      req.CategoryID,
		AllocatedAmount: req.AllocatedAmount,
		AlertThreshold:  req.AlertThreshold,
		IsActive:        true,
	}

	if err := s.repo.CreateCategory(ctx, budgetCategory); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("budget_category_id", budgetCategory.ID.String()))
	return s.toBudgetCategoryResponse(budgetCategory), nil
}

// GetBudgetCategory retrieves a budget category
func (s *service) GetBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID) (*BudgetCategoryResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.GetBudgetCategory",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
			attribute.String("category_id", categoryID.String()),
		),
	)
	defer span.End()

	// Verify budget ownership
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return nil, fmt.Errorf("unauthorized access to budget")
	}

	budgetCategory, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Verify the category belongs to the budget
	if budgetCategory.BudgetID != budgetID {
		span.SetStatus(codes.Error, "category does not belong to budget")
		return nil, fmt.Errorf("category does not belong to budget")
	}

	return s.toBudgetCategoryResponse(budgetCategory), nil
}

// ListBudgetCategories retrieves all categories for a budget
func (s *service) ListBudgetCategories(ctx context.Context, userID, budgetID uuid.UUID) ([]BudgetCategoryResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.ListBudgetCategories",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
		),
	)
	defer span.End()

	// Verify budget ownership
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return nil, fmt.Errorf("unauthorized access to budget")
	}

	categories, err := s.repo.GetCategoriesByBudgetID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	responses := make([]BudgetCategoryResponse, len(categories))
	for i, category := range categories {
		responses[i] = *s.toBudgetCategoryResponse(&category)
	}

	span.SetAttributes(attribute.Int("categories_count", len(responses)))
	return responses, nil
}

// UpdateBudgetCategory updates a budget category
func (s *service) UpdateBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID, req *UpdateBudgetCategoryRequest) (*BudgetCategoryResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.UpdateBudgetCategory",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
			attribute.String("category_id", categoryID.String()),
		),
	)
	defer span.End()

	// Verify budget ownership
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return nil, fmt.Errorf("unauthorized access to budget")
	}

	budgetCategory, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Verify the category belongs to the budget
	if budgetCategory.BudgetID != budgetID {
		span.SetStatus(codes.Error, "category does not belong to budget")
		return nil, fmt.Errorf("category does not belong to budget")
	}

	// Update fields
	if req.AllocatedAmount != nil {
		budgetCategory.AllocatedAmount = *req.AllocatedAmount
	}
	if req.AlertThreshold != nil {
		budgetCategory.AlertThreshold = *req.AlertThreshold
	}
	if req.IsActive != nil {
		budgetCategory.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateCategory(ctx, budgetCategory); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return s.toBudgetCategoryResponse(budgetCategory), nil
}

// DeleteBudgetCategory deletes a budget category
func (s *service) DeleteBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID) error {
	ctx, span := otel.Tracer("").Start(ctx, "budget.DeleteBudgetCategory",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
			attribute.String("category_id", categoryID.String()),
		),
	)
	defer span.End()

	// Verify budget ownership
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return fmt.Errorf("unauthorized access to budget")
	}

	budgetCategory, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	// Verify the category belongs to the budget
	if budgetCategory.BudgetID != budgetID {
		span.SetStatus(codes.Error, "category does not belong to budget")
		return fmt.Errorf("category does not belong to budget")
	}

	if err := s.repo.DeleteCategory(ctx, categoryID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// GetBudgetSummary retrieves a comprehensive budget summary
func (s *service) GetBudgetSummary(ctx context.Context, userID, budgetID uuid.UUID) (*BudgetSummary, error) {
	ctx, span := otel.Tracer("").Start(ctx, "budget.GetBudgetSummary",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
		),
	)
	defer span.End()

	// Verify budget ownership
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return nil, fmt.Errorf("unauthorized access to budget")
	}

	summary, err := s.repo.GetBudgetSummary(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(
		attribute.Float64("total_allocated", summary.TotalAllocated),
		attribute.Float64("total_spent", summary.TotalSpent),
		attribute.Float64("spending_progress", summary.SpendingProgress),
		attribute.Int("alerts_count", len(summary.Alerts)),
	)

	return summary, nil
}

// UpdateBudgetFromTransaction updates budget spending when a transaction is created/updated
func (s *service) UpdateBudgetFromTransaction(ctx context.Context, userID, budgetID, categoryID uuid.UUID, amount float64) error {
	ctx, span := otel.Tracer("").Start(ctx, "budget.UpdateBudgetFromTransaction",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("budget_id", budgetID.String()),
			attribute.String("category_id", categoryID.String()),
			attribute.Float64("amount", amount),
		),
	)
	defer span.End()

	// Verify budget ownership
	budget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if budget.UserID != userID {
		span.SetStatus(codes.Error, "unauthorized access to budget")
		return fmt.Errorf("unauthorized access to budget")
	}

	// Update the spent amount for the category
	if err := s.repo.UpdateSpentAmount(ctx, budgetID, categoryID, amount); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// Helper methods

func (s *service) validateCreateBudgetRequest(req *CreateBudgetRequest) error {
	if req.Name == "" {
		return fmt.Errorf("budget name is required")
	}
	if req.TotalAmount <= 0 {
		return fmt.Errorf("total amount must be positive")
	}
	if req.StartDate.IsZero() {
		return fmt.Errorf("start date is required")
	}
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return fmt.Errorf("end date must be after start date")
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}
	return nil
}

func (s *service) validateBudget(budget *Budget) error {
	if budget.Name == "" {
		return fmt.Errorf("budget name is required")
	}
	if budget.TotalAmount <= 0 {
		return fmt.Errorf("total amount must be positive")
	}
	if budget.StartDate.IsZero() {
		return fmt.Errorf("start date is required")
	}
	if budget.EndDate != nil && budget.EndDate.Before(budget.StartDate) {
		return fmt.Errorf("end date must be after start date")
	}
	return nil
}

func (s *service) toBudgetResponse(budget *Budget) *BudgetResponse {
	return &BudgetResponse{
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
}

func (s *service) toBudgetCategoryResponse(budgetCategory *BudgetCategory) *BudgetCategoryResponse {
	return &BudgetCategoryResponse{
		ID:              budgetCategory.ID,
		BudgetID:        budgetCategory.BudgetID,
		CategoryID:      budgetCategory.CategoryID,
		AllocatedAmount: budgetCategory.AllocatedAmount,
		SpentAmount:     budgetCategory.SpentAmount,
		AlertThreshold:  budgetCategory.AlertThreshold,
		IsActive:        budgetCategory.IsActive,
		CreatedAt:       budgetCategory.CreatedAt,
		UpdatedAt:       budgetCategory.UpdatedAt,
	}
}
