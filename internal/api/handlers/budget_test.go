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

	"fiscaflow/internal/domain/budget"
)

// MockBudgetService is a mock implementation of budget.Service
type MockBudgetService struct {
	mock.Mock
}

func (m *MockBudgetService) CreateBudget(ctx context.Context, userID uuid.UUID, req *budget.CreateBudgetRequest) (*budget.BudgetResponse, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*budget.BudgetResponse), args.Error(1)
}

func (m *MockBudgetService) GetBudget(ctx context.Context, userID, budgetID uuid.UUID) (*budget.BudgetResponse, error) {
	args := m.Called(ctx, userID, budgetID)
	return args.Get(0).(*budget.BudgetResponse), args.Error(1)
}

func (m *MockBudgetService) ListBudgets(ctx context.Context, userID uuid.UUID, offset, limit int) ([]budget.BudgetResponse, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).([]budget.BudgetResponse), args.Error(1)
}

func (m *MockBudgetService) UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, req *budget.UpdateBudgetRequest) (*budget.BudgetResponse, error) {
	args := m.Called(ctx, userID, budgetID, req)
	return args.Get(0).(*budget.BudgetResponse), args.Error(1)
}

func (m *MockBudgetService) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	args := m.Called(ctx, userID, budgetID)
	return args.Error(0)
}

func (m *MockBudgetService) AddBudgetCategory(ctx context.Context, userID, budgetID uuid.UUID, req *budget.CreateBudgetCategoryRequest) (*budget.BudgetCategoryResponse, error) {
	args := m.Called(ctx, userID, budgetID, req)
	return args.Get(0).(*budget.BudgetCategoryResponse), args.Error(1)
}

func (m *MockBudgetService) GetBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID) (*budget.BudgetCategoryResponse, error) {
	args := m.Called(ctx, userID, budgetID, categoryID)
	return args.Get(0).(*budget.BudgetCategoryResponse), args.Error(1)
}

func (m *MockBudgetService) ListBudgetCategories(ctx context.Context, userID, budgetID uuid.UUID) ([]budget.BudgetCategoryResponse, error) {
	args := m.Called(ctx, userID, budgetID)
	return args.Get(0).([]budget.BudgetCategoryResponse), args.Error(1)
}

func (m *MockBudgetService) UpdateBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID, req *budget.UpdateBudgetCategoryRequest) (*budget.BudgetCategoryResponse, error) {
	args := m.Called(ctx, userID, budgetID, categoryID, req)
	return args.Get(0).(*budget.BudgetCategoryResponse), args.Error(1)
}

func (m *MockBudgetService) DeleteBudgetCategory(ctx context.Context, userID, budgetID, categoryID uuid.UUID) error {
	args := m.Called(ctx, userID, budgetID, categoryID)
	return args.Error(0)
}

func (m *MockBudgetService) GetBudgetSummary(ctx context.Context, userID, budgetID uuid.UUID) (*budget.BudgetSummary, error) {
	args := m.Called(ctx, userID, budgetID)
	return args.Get(0).(*budget.BudgetSummary), args.Error(1)
}

func (m *MockBudgetService) UpdateBudgetFromTransaction(ctx context.Context, userID, budgetID, categoryID uuid.UUID, amount float64) error {
	args := m.Called(ctx, userID, budgetID, categoryID, amount)
	return args.Error(0)
}

func TestBudgetHandler_CreateBudget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	budgetID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*MockBudgetService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful creation",
			requestBody: `{"name":"Monthly Budget","description":"My monthly budget","period_type":"monthly","start_date":"2024-06-01T00:00:00Z","total_amount":5000.00,"currency":"USD"}`,
			setupMock: func(mockService *MockBudgetService) {
				response := &budget.BudgetResponse{
					ID:          budgetID,
					UserID:      userID,
					Name:        "Monthly Budget",
					Description: "My monthly budget",
					PeriodType:  budget.PeriodTypeMonthly,
					StartDate:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					TotalAmount: 5000.00,
					Currency:    "USD",
					IsActive:    true,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				}
				mockService.On("CreateBudget", mock.Anything, userID, mock.AnythingOfType("*budget.CreateBudgetRequest")).
					Return(response, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"budget":{"id":"` + budgetID.String() + `","user_id":"` + userID.String() + `","name":"Monthly Budget","description":"My monthly budget","period_type":"monthly","start_date":"2024-06-01T00:00:00Z","total_amount":5000,"currency":"USD","is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","end_date":null,"family_id":null,"settings":""}}`,
		},
		{
			name:        "invalid request body",
			requestBody: `{"invalid":"json"`,
			setupMock: func(mockService *MockBudgetService) {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`,
		},
		{
			name:        "internal server error",
			requestBody: `{"name":"Monthly Budget","description":"My monthly budget","period_type":"monthly","start_date":"2024-06-01T00:00:00Z","total_amount":5000.00,"currency":"USD"}`,
			setupMock: func(mockService *MockBudgetService) {
				mockService.On("CreateBudget", mock.Anything, userID, mock.AnythingOfType("*budget.CreateBudgetRequest")).
					Return((*budget.BudgetResponse)(nil), fmt.Errorf("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBudgetService{}
			tt.setupMock(mockService)

			handler := NewBudgetHandler(mockService)

			// Use proper Gin test setup
			router := gin.New()
			router.POST("/budgets", func(c *gin.Context) {
				c.Set("user_id", userID)
				handler.CreateBudget(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/budgets", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestBudgetHandler_GetBudget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	budgetID := uuid.New()

	tests := []struct {
		name           string
		budgetID       string
		setupMock      func(*MockBudgetService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "successful retrieval",
			budgetID: budgetID.String(),
			setupMock: func(mockService *MockBudgetService) {
				response := &budget.BudgetResponse{
					ID:          budgetID,
					UserID:      userID,
					Name:        "Monthly Budget",
					Description: "My monthly budget",
					PeriodType:  budget.PeriodTypeMonthly,
					StartDate:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					TotalAmount: 5000.00,
					Currency:    "USD",
					IsActive:    true,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				}
				mockService.On("GetBudget", mock.Anything, userID, budgetID).
					Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"budget":{"id":"` + budgetID.String() + `","user_id":"` + userID.String() + `","name":"Monthly Budget","description":"My monthly budget","period_type":"monthly","start_date":"2024-06-01T00:00:00Z","total_amount":5000,"currency":"USD","is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","end_date":null,"family_id":null,"settings":""}}`,
		},
		{
			name:     "invalid budget ID",
			budgetID: "invalid-uuid",
			setupMock: func(mockService *MockBudgetService) {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid budget ID"}`,
		},
		{
			name:     "budget not found",
			budgetID: budgetID.String(),
			setupMock: func(mockService *MockBudgetService) {
				mockService.On("GetBudget", mock.Anything, userID, budgetID).
					Return((*budget.BudgetResponse)(nil), fmt.Errorf("budget not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"budget not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBudgetService{}
			tt.setupMock(mockService)

			handler := NewBudgetHandler(mockService)

			// Use proper Gin test setup
			router := gin.New()
			router.GET("/budgets/:id", func(c *gin.Context) {
				c.Set("user_id", userID)
				handler.GetBudget(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/budgets/"+tt.budgetID, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestBudgetHandler_ListBudgets(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	budgetID1 := uuid.New()
	budgetID2 := uuid.New()

	tests := []struct {
		name           string
		setupMock      func(*MockBudgetService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful retrieval",
			setupMock: func(mockService *MockBudgetService) {
				budgets := []budget.BudgetResponse{
					{
						ID:          budgetID1,
						UserID:      userID,
						Name:        "Monthly Budget",
						Description: "My monthly budget",
						PeriodType:  budget.PeriodTypeMonthly,
						StartDate:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
						TotalAmount: 5000.00,
						Currency:    "USD",
						IsActive:    true,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
					{
						ID:          budgetID2,
						UserID:      userID,
						Name:        "Yearly Budget",
						Description: "My yearly budget",
						PeriodType:  budget.PeriodTypeYearly,
						StartDate:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						TotalAmount: 60000.00,
						Currency:    "USD",
						IsActive:    true,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
				}
				mockService.On("ListBudgets", mock.Anything, userID, 0, 20).
					Return(budgets, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"budgets":[{"id":"` + budgetID1.String() + `","user_id":"` + userID.String() + `","name":"Monthly Budget","description":"My monthly budget","period_type":"monthly","start_date":"2024-06-01T00:00:00Z","total_amount":5000,"currency":"USD","is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","end_date":null,"family_id":null,"settings":""},{"id":"` + budgetID2.String() + `","user_id":"` + userID.String() + `","name":"Yearly Budget","description":"My yearly budget","period_type":"yearly","start_date":"2024-01-01T00:00:00Z","total_amount":60000,"currency":"USD","is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","end_date":null,"family_id":null,"settings":""}]}`,
		},
		{
			name: "internal server error",
			setupMock: func(mockService *MockBudgetService) {
				mockService.On("ListBudgets", mock.Anything, userID, 0, 20).
					Return([]budget.BudgetResponse{}, fmt.Errorf("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBudgetService{}
			tt.setupMock(mockService)

			handler := NewBudgetHandler(mockService)

			// Use proper Gin test setup
			router := gin.New()
			router.GET("/budgets", func(c *gin.Context) {
				c.Set("user_id", userID)
				handler.ListBudgets(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/budgets", nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestBudgetHandler_UpdateBudget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	budgetID := uuid.New()

	tests := []struct {
		name           string
		budgetID       string
		requestBody    string
		setupMock      func(*MockBudgetService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "successful update",
			budgetID:    budgetID.String(),
			requestBody: `{"name":"Updated Budget","description":"Updated description","total_amount":6000.00}`,
			setupMock: func(mockService *MockBudgetService) {
				response := &budget.BudgetResponse{
					ID:          budgetID,
					UserID:      userID,
					Name:        "Updated Budget",
					Description: "Updated description",
					PeriodType:  budget.PeriodTypeMonthly,
					StartDate:   time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					TotalAmount: 6000.00,
					Currency:    "USD",
					IsActive:    true,
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
				}
				mockService.On("UpdateBudget", mock.Anything, userID, budgetID, mock.AnythingOfType("*budget.UpdateBudgetRequest")).
					Return(response, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"budget":{"id":"` + budgetID.String() + `","user_id":"` + userID.String() + `","name":"Updated Budget","description":"Updated description","period_type":"monthly","start_date":"2024-06-01T00:00:00Z","total_amount":6000,"currency":"USD","is_active":true,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","end_date":null,"family_id":null,"settings":""}}`,
		},
		{
			name:        "invalid budget ID",
			budgetID:    "invalid-uuid",
			requestBody: `{"name":"Updated Budget"}`,
			setupMock: func(mockService *MockBudgetService) {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid budget ID"}`,
		},
		{
			name:        "invalid request body",
			budgetID:    budgetID.String(),
			requestBody: `{"name":123}`,
			setupMock: func(mockService *MockBudgetService) {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request body"}`,
		},
		{
			name:        "budget not found",
			budgetID:    budgetID.String(),
			requestBody: `{"name":"Updated Budget"}`,
			setupMock: func(mockService *MockBudgetService) {
				mockService.On("UpdateBudget", mock.Anything, userID, budgetID, mock.AnythingOfType("*budget.UpdateBudgetRequest")).
					Return((*budget.BudgetResponse)(nil), fmt.Errorf("budget not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"budget not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBudgetService{}
			tt.setupMock(mockService)

			handler := NewBudgetHandler(mockService)

			// Use proper Gin test setup
			router := gin.New()
			router.PUT("/budgets/:id", func(c *gin.Context) {
				c.Set("user_id", userID)
				handler.UpdateBudget(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/budgets/"+tt.budgetID, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}

func TestBudgetHandler_DeleteBudget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New()
	budgetID := uuid.New()

	tests := []struct {
		name           string
		budgetID       string
		setupMock      func(*MockBudgetService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "successful deletion",
			budgetID: budgetID.String(),
			setupMock: func(mockService *MockBudgetService) {
				mockService.On("DeleteBudget", mock.Anything, userID, budgetID).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:     "invalid budget ID",
			budgetID: "invalid-uuid",
			setupMock: func(mockService *MockBudgetService) {
				// No mock setup needed for this case
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid budget ID"}`,
		},
		{
			name:     "budget not found",
			budgetID: budgetID.String(),
			setupMock: func(mockService *MockBudgetService) {
				mockService.On("DeleteBudget", mock.Anything, userID, budgetID).
					Return(fmt.Errorf("budget not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"budget not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockBudgetService{}
			tt.setupMock(mockService)

			handler := NewBudgetHandler(mockService)

			// Use proper Gin test setup
			router := gin.New()
			router.DELETE("/budgets/:id", func(c *gin.Context) {
				c.Set("user_id", userID)
				handler.DeleteBudget(c)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/budgets/"+tt.budgetID, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
