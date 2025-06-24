package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"fiscaflow/internal/domain/budget"
)

// MockService is a mock of Service interface.
type MockService struct {
	mock.Mock
}

// CreateBudget mocks base method.
func (m *MockService) CreateBudget(ctx context.Context, userID string, req *budget.CreateBudgetRequest) (*budget.Budget, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*budget.Budget), args.Error(1)
}

// GetBudget mocks base method.
func (m *MockService) GetBudget(ctx context.Context, budgetID, userID string) (*budget.Budget, error) {
	args := m.Called(ctx, budgetID, userID)
	return args.Get(0).(*budget.Budget), args.Error(1)
}

// ListBudgets mocks base method.
func (m *MockService) ListBudgets(ctx context.Context, userID string) ([]*budget.Budget, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*budget.Budget), args.Error(1)
}

// UpdateBudget mocks base method.
func (m *MockService) UpdateBudget(ctx context.Context, budgetID, userID string, req *budget.UpdateBudgetRequest) (*budget.Budget, error) {
	args := m.Called(ctx, budgetID, userID, req)
	return args.Get(0).(*budget.Budget), args.Error(1)
}

// DeleteBudget mocks base method.
func (m *MockService) DeleteBudget(ctx context.Context, budgetID, userID string) error {
	args := m.Called(ctx, budgetID, userID)
	return args.Error(0)
}

// GetBudgetSummary mocks base method.
func (m *MockService) GetBudgetSummary(ctx context.Context, budgetID, userID string) (*budget.BudgetSummary, error) {
	args := m.Called(ctx, budgetID, userID)
	return args.Get(0).(*budget.BudgetSummary), args.Error(1)
}
