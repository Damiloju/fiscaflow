package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"fiscaflow/internal/domain/transaction"
)

type mockTransactionService struct {
	mock.Mock
}

func (m *mockTransactionService) CreateTransaction(ctx context.Context, userID uuid.UUID, req *transaction.CreateTransactionRequest) (*transaction.TransactionResponse, error) {
	args := m.Called(ctx, userID, req)
	if resp, ok := args.Get(0).(*transaction.TransactionResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTransactionService) GetTransaction(ctx context.Context, userID, transactionID uuid.UUID) (*transaction.TransactionResponse, error) {
	args := m.Called(ctx, userID, transactionID)
	if resp, ok := args.Get(0).(*transaction.TransactionResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTransactionService) GetTransactions(ctx context.Context, userID uuid.UUID, offset, limit int) ([]transaction.TransactionResponse, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).([]transaction.TransactionResponse), args.Error(1)
}
func (m *mockTransactionService) UpdateTransaction(ctx context.Context, userID, transactionID uuid.UUID, req *transaction.UpdateTransactionRequest) (*transaction.TransactionResponse, error) {
	args := m.Called(ctx, userID, transactionID, req)
	if resp, ok := args.Get(0).(*transaction.TransactionResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTransactionService) DeleteTransaction(ctx context.Context, userID, transactionID uuid.UUID) error {
	args := m.Called(ctx, userID, transactionID)
	return args.Error(0)
}

// Other methods omitted for brevity
func (m *mockTransactionService) CreateCategory(ctx context.Context, req *transaction.CreateCategoryRequest) (*transaction.Category, error) {
	return nil, nil
}
func (m *mockTransactionService) GetCategory(ctx context.Context, categoryID uuid.UUID) (*transaction.Category, error) {
	return nil, nil
}
func (m *mockTransactionService) GetCategories(ctx context.Context, offset, limit int) ([]transaction.Category, error) {
	return nil, nil
}
func (m *mockTransactionService) GetDefaultCategories(ctx context.Context) ([]transaction.Category, error) {
	return nil, nil
}
func (m *mockTransactionService) UpdateCategory(ctx context.Context, categoryID uuid.UUID, req *transaction.CreateCategoryRequest) (*transaction.Category, error) {
	return nil, nil
}
func (m *mockTransactionService) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	return nil
}
func (m *mockTransactionService) CreateAccount(ctx context.Context, userID uuid.UUID, req *transaction.CreateAccountRequest) (*transaction.Account, error) {
	return nil, nil
}
func (m *mockTransactionService) GetAccount(ctx context.Context, userID, accountID uuid.UUID) (*transaction.Account, error) {
	return nil, nil
}
func (m *mockTransactionService) GetAccounts(ctx context.Context, userID uuid.UUID) ([]transaction.Account, error) {
	return nil, nil
}
func (m *mockTransactionService) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, req *transaction.CreateAccountRequest) (*transaction.Account, error) {
	return nil, nil
}
func (m *mockTransactionService) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	return nil
}

func setupRouterWithTransactionHandler(svc transaction.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewTransactionHandler(svc)
	// Simulate auth middleware
	r.Use(func(c *gin.Context) {
		c.Set("user_id", uuid.MustParse("11111111-1111-1111-1111-111111111111"))
		c.Next()
	})
	h.RegisterRoutes(r.Group("/api/v1"))
	return r
}

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	svc := new(mockTransactionService)
	r := setupRouterWithTransactionHandler(svc)
	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	createReq := transaction.CreateTransactionRequest{
		AccountID:       uuid.New(),
		Amount:          50.0,
		Currency:        "USD",
		Description:     "Lunch",
		TransactionDate: time.Now(),
	}
	resp := &transaction.TransactionResponse{ID: uuid.New(), UserID: userID, Amount: 50.0, Description: "Lunch"}
	svc.On("CreateTransaction", mock.Anything, userID, mock.Anything).Return(resp, nil)
	body, _ := json.Marshal(createReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
