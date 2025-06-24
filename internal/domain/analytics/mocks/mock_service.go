package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"fiscaflow/internal/domain/analytics"
)

// MockService is a mock of Service interface.
type MockService struct {
	mock.Mock
}

// CategorizeTransaction mocks base method.
func (m *MockService) CategorizeTransaction(ctx context.Context, req *analytics.CategorizationRequest) (*analytics.CategorizationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*analytics.CategorizationResponse), args.Error(1)
}

// CreateCategorizationRule mocks base method.
func (m *MockService) CreateCategorizationRule(ctx context.Context, req *analytics.CreateCategorizationRuleRequest) (*analytics.CategorizationRuleResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*analytics.CategorizationRuleResponse), args.Error(1)
}

// GetCategorizationRule mocks base method.
func (m *MockService) GetCategorizationRule(ctx context.Context, id uuid.UUID) (*analytics.CategorizationRuleResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*analytics.CategorizationRuleResponse), args.Error(1)
}

// ListCategorizationRules mocks base method.
func (m *MockService) ListCategorizationRules(ctx context.Context, offset, limit int) ([]analytics.CategorizationRuleResponse, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]analytics.CategorizationRuleResponse), args.Error(1)
}

// UpdateCategorizationRule mocks base method.
func (m *MockService) UpdateCategorizationRule(ctx context.Context, id uuid.UUID, req *analytics.UpdateCategorizationRuleRequest) (*analytics.CategorizationRuleResponse, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*analytics.CategorizationRuleResponse), args.Error(1)
}

// DeleteCategorizationRule mocks base method.
func (m *MockService) DeleteCategorizationRule(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// AnalyzeSpending mocks base method.
func (m *MockService) AnalyzeSpending(ctx context.Context, userID uuid.UUID, req *analytics.SpendingAnalysisRequest) (*analytics.SpendingAnalysisResponse, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*analytics.SpendingAnalysisResponse), args.Error(1)
}

// GetSpendingInsights mocks base method.
func (m *MockService) GetSpendingInsights(ctx context.Context, userID uuid.UUID, periodStart, periodEnd time.Time) ([]analytics.SpendingInsight, error) {
	args := m.Called(ctx, userID, periodStart, periodEnd)
	return args.Get(0).([]analytics.SpendingInsight), args.Error(1)
}
