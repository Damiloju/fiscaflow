package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"fiscaflow/internal/domain/transaction"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountService struct {
	mock.Mock
}

func (m *mockAccountService) CreateAccount(ctx context.Context, userID uuid.UUID, req *transaction.CreateAccountRequest) (*transaction.Account, error) {
	args := m.Called(ctx, userID, req)
	return args.Get(0).(*transaction.Account), args.Error(1)
}
func (m *mockAccountService) GetAccount(ctx context.Context, userID, accountID uuid.UUID) (*transaction.Account, error) {
	args := m.Called(ctx, userID, accountID)
	return args.Get(0).(*transaction.Account), args.Error(1)
}
func (m *mockAccountService) GetAccounts(ctx context.Context, userID uuid.UUID) ([]transaction.Account, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]transaction.Account), args.Error(1)
}
func (m *mockAccountService) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, req *transaction.CreateAccountRequest) (*transaction.Account, error) {
	args := m.Called(ctx, userID, accountID, req)
	return args.Get(0).(*transaction.Account), args.Error(1)
}
func (m *mockAccountService) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	args := m.Called(ctx, userID, accountID)
	return args.Error(0)
}

// Unused service methods for interface compliance
func (m *mockAccountService) CreateTransaction(context.Context, uuid.UUID, *transaction.CreateTransactionRequest) (*transaction.TransactionResponse, error) {
	return nil, nil
}
func (m *mockAccountService) GetTransaction(context.Context, uuid.UUID, uuid.UUID) (*transaction.TransactionResponse, error) {
	return nil, nil
}
func (m *mockAccountService) GetTransactions(context.Context, uuid.UUID, int, int) ([]transaction.TransactionResponse, error) {
	return nil, nil
}
func (m *mockAccountService) UpdateTransaction(context.Context, uuid.UUID, uuid.UUID, *transaction.UpdateTransactionRequest) (*transaction.TransactionResponse, error) {
	return nil, nil
}
func (m *mockAccountService) DeleteTransaction(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}
func (m *mockAccountService) CreateCategory(context.Context, *transaction.CreateCategoryRequest) (*transaction.Category, error) {
	return nil, nil
}
func (m *mockAccountService) GetCategory(context.Context, uuid.UUID) (*transaction.Category, error) {
	return nil, nil
}
func (m *mockAccountService) GetCategories(context.Context, int, int) ([]transaction.Category, error) {
	return nil, nil
}
func (m *mockAccountService) GetDefaultCategories(context.Context) ([]transaction.Category, error) {
	return nil, nil
}
func (m *mockAccountService) UpdateCategory(context.Context, uuid.UUID, *transaction.CreateCategoryRequest) (*transaction.Category, error) {
	return nil, nil
}
func (m *mockAccountService) DeleteCategory(context.Context, uuid.UUID) error {
	return nil
}

func TestAccountHandler_CreateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockAccountService)
	h := NewAccountHandler(mockSvc)
	r := gin.Default()
	userID := uuid.New()
	r.POST("/accounts", func(c *gin.Context) {
		c.Set("user_id", userID)
		h.CreateAccount(c)
	})

	accReq := &transaction.CreateAccountRequest{
		Name:        "Test Account",
		Type:        transaction.AccountTypeChecking,
		Institution: "Test Bank",
		Balance:     1000,
		Currency:    "USD",
	}
	acc := &transaction.Account{ID: uuid.New(), Name: "Test Account", Type: transaction.AccountTypeChecking, Institution: "Test Bank", Balance: 1000, Currency: "USD"}
	mockSvc.On("CreateAccount", mock.Anything, userID, accReq).Return(acc, nil)

	body, _ := json.Marshal(accReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp transaction.Account
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, acc.ID, resp.ID)
}

func TestAccountHandler_CreateAccount_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockAccountService)
	h := NewAccountHandler(mockSvc)
	r := gin.Default()
	r.POST("/accounts", h.CreateAccount)

	body, _ := json.Marshal(transaction.CreateAccountRequest{Name: "Test Account"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAccountHandler_GetAccount_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockAccountService)
	h := NewAccountHandler(mockSvc)
	r := gin.Default()
	r.GET("/accounts/:id", func(c *gin.Context) {
		c.Set("user_id", uuid.New())
		h.GetAccount(c)
	})

	mockSvc.On("GetAccount", mock.Anything, mock.Anything, mock.Anything).Return(&transaction.Account{}, transaction.ErrAccountNotFound)
	id := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/accounts/"+id, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAccountHandler_ListAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockAccountService)
	h := NewAccountHandler(mockSvc)
	r := gin.Default()
	r.GET("/accounts", func(c *gin.Context) {
		c.Set("user_id", uuid.New())
		h.ListAccounts(c)
	})

	mockSvc.On("GetAccounts", mock.Anything, mock.Anything).Return([]transaction.Account{{ID: uuid.New(), Name: "Test Account"}}, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/accounts", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAccountHandler_UpdateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockAccountService)
	h := NewAccountHandler(mockSvc)
	r := gin.Default()
	userID := uuid.New()
	acc := &transaction.Account{ID: uuid.New(), Name: "Updated Account", Type: transaction.AccountTypeChecking, Institution: "Test Bank", Balance: 2000, Currency: "USD"}
	accReq := &transaction.CreateAccountRequest{
		Name:        "Updated Account",
		Type:        transaction.AccountTypeChecking,
		Institution: "Test Bank",
		Balance:     2000,
		Currency:    "USD",
	}
	r.PUT("/accounts/:id", func(c *gin.Context) {
		c.Set("user_id", userID)
		h.UpdateAccount(c)
	})

	mockSvc.On("UpdateAccount", mock.Anything, userID, acc.ID, accReq).Return(acc, nil)
	id := acc.ID.String()
	body, _ := json.Marshal(accReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/accounts/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAccountHandler_DeleteAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockAccountService)
	h := NewAccountHandler(mockSvc)
	r := gin.Default()
	r.DELETE("/accounts/:id", func(c *gin.Context) {
		c.Set("user_id", uuid.New())
		h.DeleteAccount(c)
	})

	mockSvc.On("DeleteAccount", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	id := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/accounts/"+id, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
