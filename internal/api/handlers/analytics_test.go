package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"fiscaflow/internal/domain/analytics"
)

// MockAnalyticsService is a mock implementation of analytics.Service
type MockAnalyticsService struct {
	mock.Mock
}

func (m *MockAnalyticsService) CategorizeTransaction(ctx context.Context, req *analytics.CategorizationRequest) (*analytics.CategorizationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*analytics.CategorizationResponse), args.Error(1)
}

func (m *MockAnalyticsService) AnalyzeSpending(ctx context.Context, userID uuid.UUID, req *analytics.SpendingAnalysisRequest) (*analytics.SpendingAnalysisResponse, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*analytics.SpendingAnalysisResponse), args.Error(1)
}

func (m *MockAnalyticsService) GetSpendingInsights(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]analytics.SpendingInsight, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	return args.Get(0).([]analytics.SpendingInsight), args.Error(1)
}

func (m *MockAnalyticsService) CreateCategorizationRule(ctx context.Context, req *analytics.CreateCategorizationRuleRequest) (*analytics.CategorizationRuleResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*analytics.CategorizationRuleResponse), args.Error(1)
}

func (m *MockAnalyticsService) GetCategorizationRule(ctx context.Context, id uuid.UUID) (*analytics.CategorizationRuleResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*analytics.CategorizationRuleResponse), args.Error(1)
}

func (m *MockAnalyticsService) ListCategorizationRules(ctx context.Context, offset, limit int) ([]analytics.CategorizationRuleResponse, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]analytics.CategorizationRuleResponse), args.Error(1)
}

func (m *MockAnalyticsService) UpdateCategorizationRule(ctx context.Context, id uuid.UUID, req *analytics.UpdateCategorizationRuleRequest) (*analytics.CategorizationRuleResponse, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*analytics.CategorizationRuleResponse), args.Error(1)
}

func (m *MockAnalyticsService) DeleteCategorizationRule(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestAnalyticsHandler_CategorizeTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*MockAnalyticsService) uuid.UUID
		expectedStatus int
		expectedBody   func(uuid.UUID) string
	}{
		{
			name:        "successful categorization",
			requestBody: `{"description":"grocery store","amount":150.50,"merchant":"walmart"}`,
			setupMock: func(mockService *MockAnalyticsService) uuid.UUID {
				categoryID := uuid.New()
				response := &analytics.CategorizationResponse{
					CategoryID:           categoryID,
					CategoryName:         "Food & Groceries",
					Confidence:           0.85,
					CategorizationSource: "rule",
					MatchedPattern:       "walmart",
				}
				mockService.On("CategorizeTransaction", mock.Anything, mock.AnythingOfType("*analytics.CategorizationRequest")).
					Return(response, nil)
				return categoryID
			},
			expectedStatus: http.StatusOK,
			expectedBody: func(categoryID uuid.UUID) string {
				return `{"categorization":{"category_id":"` + categoryID.String() + `","category_name":"Food & Groceries","confidence":0.85,"categorization_source":"rule","matched_pattern":"walmart"}}`
			},
		},
		{
			name:        "invalid request body",
			requestBody: `{"invalid":"json"`,
			setupMock: func(mockService *MockAnalyticsService) uuid.UUID {
				return uuid.Nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(_ uuid.UUID) string {
				return `{"error":"invalid request body"}`
			},
		},
		{
			name:        "internal server error",
			requestBody: `{"description":"grocery store","amount":150.50,"merchant":"walmart"}`,
			setupMock: func(mockService *MockAnalyticsService) uuid.UUID {
				mockService.On("CategorizeTransaction", mock.Anything, mock.AnythingOfType("*analytics.CategorizationRequest")).
					Return((*analytics.CategorizationResponse)(nil), fmt.Errorf("database error"))
				return uuid.Nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: func(_ uuid.UUID) string {
				return `{"error":"database error"}`
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAnalyticsService{}
			categoryID := tt.setupMock(mockService)

			handler := NewAnalyticsHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/analytics/categorize", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")
			// Set user_id for all except invalid request body (simulate auth)
			if tt.name != "invalid request body" {
				c.Set("user_id", uuid.New())
			}

			handler.CategorizeTransaction(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody(categoryID), w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandler_AnalyzeSpending(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*MockAnalyticsService) uuid.UUID
		expectedStatus int
		expectedBody   func(uuid.UUID) string
	}{
		{
			name:        "successful analysis",
			requestBody: `{"start_date":"2024-01-01T00:00:00Z","end_date":"2024-12-31T23:59:59Z","group_by":"category"}`,
			setupMock: func(mockService *MockAnalyticsService) uuid.UUID {
				categoryID := uuid.New()
				response := &analytics.SpendingAnalysisResponse{
					PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					PeriodEnd:   time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
					TotalSpent:  1000.0,
					TotalIncome: 1500.0,
					NetAmount:   500.0,
					CategoryBreakdown: []analytics.CategorySpending{
						{
							CategoryID:       categoryID,
							CategoryName:     "Food",
							Amount:           500.0,
							Percentage:       50.0,
							TransactionCount: 10,
						},
					},
					Insights:       nil,
					SpendingTrends: nil,
					TopCategories:  nil,
				}
				mockService.On("AnalyzeSpending", mock.Anything, userID, mock.AnythingOfType("*analytics.SpendingAnalysisRequest")).
					Return(response, nil)
				return categoryID
			},
			expectedStatus: http.StatusOK,
			expectedBody: func(categoryID uuid.UUID) string {
				return `{"analysis":{"period_start":"2024-01-01T00:00:00Z","period_end":"2024-12-31T23:59:59Z","total_spent":1000,"total_income":1500,"net_amount":500,"category_breakdown":[{"category_id":"` + categoryID.String() + `","category_name":"Food","amount":500,"percentage":50,"transaction_count":10}],"insights":null,"spending_trends":null,"top_categories":null}}`
			},
		},
		{
			name:        "invalid request body",
			requestBody: `{"invalid":"json"`,
			setupMock: func(mockService *MockAnalyticsService) uuid.UUID {
				return uuid.Nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(_ uuid.UUID) string {
				return `{"error":"invalid request body"}`
			},
		},
		{
			name:        "unauthorized",
			requestBody: `{"start_date":"2024-01-01T00:00:00Z","end_date":"2024-12-31T23:59:59Z"}`,
			setupMock: func(mockService *MockAnalyticsService) uuid.UUID {
				return uuid.Nil
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: func(_ uuid.UUID) string {
				return `{"error":"unauthorized"}`
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAnalyticsService{}
			categoryID := tt.setupMock(mockService)

			handler := NewAnalyticsHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/analytics/spending", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")
			if tt.name != "unauthorized" {
				c.Set("user_id", userID)
			}

			handler.AnalyzeSpending(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody(categoryID), w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandler_CreateCategorizationRule(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*MockAnalyticsService) (uuid.UUID, uuid.UUID)
		expectedStatus int
		expectedBody   func(uuid.UUID, uuid.UUID) string
	}{
		{
			name:        "successful creation",
			requestBody: `{"category_id":"123e4567-e89b-12d3-a456-426614174000","pattern":"grocery","pattern_type":"keyword","priority":1}`,
			setupMock: func(mockService *MockAnalyticsService) (uuid.UUID, uuid.UUID) {
				categoryID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				ruleID := uuid.New()
				response := &analytics.CategorizationRuleResponse{
					ID:           ruleID,
					CategoryID:   categoryID,
					CategoryName: "Food & Groceries",
					Pattern:      "grocery",
					PatternType:  "keyword",
					Priority:     1,
					IsActive:     true,
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				}
				mockService.On("CreateCategorizationRule", mock.Anything, mock.AnythingOfType("*analytics.CreateCategorizationRuleRequest")).
					Return(response, nil)
				return ruleID, categoryID
			},
			expectedStatus: http.StatusCreated,
			expectedBody: func(ruleID, categoryID uuid.UUID) string {
				return `{"categorization_rule":{"id":"` + ruleID.String() + `","category_id":"` + categoryID.String() + `","category_name":"Food & Groceries","pattern":"grocery","pattern_type":"keyword","priority":1,"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`
			},
		},
		{
			name:        "invalid request body",
			requestBody: `{"invalid":"json"`,
			setupMock: func(mockService *MockAnalyticsService) (uuid.UUID, uuid.UUID) {
				return uuid.Nil, uuid.Nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: func(_, _ uuid.UUID) string {
				return `{"error":"invalid request body"}`
			},
		},
		{
			name:        "internal server error",
			requestBody: `{"category_id":"123e4567-e89b-12d3-a456-426614174000","pattern":"grocery","pattern_type":"keyword"}`,
			setupMock: func(mockService *MockAnalyticsService) (uuid.UUID, uuid.UUID) {
				mockService.On("CreateCategorizationRule", mock.Anything, mock.AnythingOfType("*analytics.CreateCategorizationRuleRequest")).
					Return((*analytics.CategorizationRuleResponse)(nil), fmt.Errorf("database error"))
				return uuid.Nil, uuid.Nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: func(_, _ uuid.UUID) string {
				return `{"error":"database error"}`
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAnalyticsService{}
			ruleID, categoryID := tt.setupMock(mockService)

			handler := NewAnalyticsHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/analytics/categorization-rules", bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("user_id", uuid.New())

			handler.CreateCategorizationRule(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody(ruleID, categoryID), w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandler_GetCategorizationRule(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ruleID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name           string
		ruleID         string
		setupMock      func(*MockAnalyticsService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful retrieval",
			ruleID: ruleID.String(),
			setupMock: func(mockService *MockAnalyticsService) {
				response := &analytics.CategorizationRuleResponse{
					ID:           ruleID,
					CategoryID:   categoryID,
					CategoryName: "Food & Groceries",
					Pattern:      "grocery",
					PatternType:  "keyword",
					Priority:     1,
					IsActive:     true,
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				}
				mockService.On("GetCategorizationRule", mock.Anything, ruleID).
					Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"categorization_rule":{"id":"` + ruleID.String() + `","category_id":"` + categoryID.String() + `","category_name":"Food & Groceries","pattern":"grocery","pattern_type":"keyword","priority":1,"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name:   "invalid rule ID",
			ruleID: "invalid-uuid",
			setupMock: func(mockService *MockAnalyticsService) {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid rule ID"}`,
		},
		{
			name:   "rule not found",
			ruleID: ruleID.String(),
			setupMock: func(mockService *MockAnalyticsService) {
				mockService.On("GetCategorizationRule", mock.Anything, ruleID).
					Return((*analytics.CategorizationRuleResponse)(nil), fmt.Errorf("rule not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"rule not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAnalyticsService{}
			tt.setupMock(mockService)

			handler := NewAnalyticsHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/analytics/categorization-rules/"+tt.ruleID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.ruleID}}
			c.Set("user_id", uuid.New())

			handler.GetCategorizationRule(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandler_ListCategorizationRules(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	ruleID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name           string
		setupMock      func(*MockAnalyticsService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful retrieval",
			setupMock: func(mockService *MockAnalyticsService) {
				response := []analytics.CategorizationRuleResponse{
					{
						ID:           ruleID,
						CategoryID:   categoryID,
						CategoryName: "Food & Groceries",
						Pattern:      "grocery",
						PatternType:  "keyword",
						Priority:     1,
						IsActive:     true,
						CreatedAt:    time.Time{},
						UpdatedAt:    time.Time{},
					},
				}
				mockService.On("ListCategorizationRules", mock.Anything, 0, 20).
					Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"categorization_rules":[{"id":"` + ruleID.String() + `","category_id":"` + categoryID.String() + `","category_name":"Food & Groceries","pattern":"grocery","pattern_type":"keyword","priority":1,"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAnalyticsService{}
			tt.setupMock(mockService)

			handler := NewAnalyticsHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/analytics/categorization-rules", nil)
			c.Set("user_id", userID)

			handler.ListCategorizationRules(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandler_UpdateCategorizationRule(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ruleID := uuid.New()
	categoryID := uuid.New()

	tests := []struct {
		name           string
		ruleID         string
		requestBody    string
		setupMock      func(*MockAnalyticsService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful update",
			ruleID:      ruleID.String(),
			requestBody: `{"pattern":"updated-grocery","pattern_type":"keyword","priority":2,"is_active":true}`,
			setupMock: func(mockService *MockAnalyticsService) {
				response := &analytics.CategorizationRuleResponse{
					ID:           ruleID,
					CategoryID:   categoryID,
					CategoryName: "Updated Food & Groceries",
					Pattern:      "updated-grocery",
					PatternType:  "keyword",
					Priority:     2,
					IsActive:     true,
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				}
				mockService.On("UpdateCategorizationRule", mock.Anything, ruleID, mock.AnythingOfType("*analytics.UpdateCategorizationRuleRequest")).
					Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"categorization_rule":{"id":"` + ruleID.String() + `","category_id":"` + categoryID.String() + `","category_name":"Updated Food & Groceries","pattern":"updated-grocery","pattern_type":"keyword","priority":2,"is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name:           "invalid rule ID",
			ruleID:         "invalid-uuid",
			requestBody:    `{"pattern":"updated-grocery"}`,
			setupMock:      func(mockService *MockAnalyticsService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid rule ID"}`,
		},
		{
			name:           "invalid request body",
			ruleID:         ruleID.String(),
			requestBody:    `{"pattern":123}`,
			setupMock:      func(mockService *MockAnalyticsService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`,
		},
		{
			name:        "rule not found",
			ruleID:      ruleID.String(),
			requestBody: `{"pattern":"updated-grocery"}`,
			setupMock: func(mockService *MockAnalyticsService) {
				mockService.On("UpdateCategorizationRule", mock.Anything, ruleID, mock.AnythingOfType("*analytics.UpdateCategorizationRuleRequest")).
					Return((*analytics.CategorizationRuleResponse)(nil), fmt.Errorf("rule not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"rule not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAnalyticsService{}
			tt.setupMock(mockService)

			handler := NewAnalyticsHandler(mockService)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("PUT", "/analytics/categorization-rules/"+tt.ruleID, bytes.NewBufferString(tt.requestBody))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.ruleID}}
			c.Set("user_id", uuid.New())

			handler.UpdateCategorizationRule(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestAnalyticsHandler_DeleteCategorizationRule(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ruleID := uuid.New()

	tests := []struct {
		name           string
		ruleID         string
		setupMock      func(*MockAnalyticsService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successful deletion",
			ruleID: ruleID.String(),
			setupMock: func(mockService *MockAnalyticsService) {
				mockService.On("DeleteCategorizationRule", mock.Anything, ruleID).
					Return(nil).Once()
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:   "invalid rule ID",
			ruleID: "invalid-uuid",
			setupMock: func(mockService *MockAnalyticsService) {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid rule ID"}`,
		},
		{
			name:   "rule not found",
			ruleID: ruleID.String(),
			setupMock: func(mockService *MockAnalyticsService) {
				mockService.On("DeleteCategorizationRule", mock.Anything, ruleID).
					Return(fmt.Errorf("rule not found")).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"rule not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockAnalyticsService{}
			tt.setupMock(mockService)

			handler := NewAnalyticsHandler(mockService)

			// Use proper Gin test setup
			router := gin.New()
			router.DELETE("/analytics/categorization-rules/:id", handler.DeleteCategorizationRule)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/analytics/categorization-rules/"+tt.ruleID, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
