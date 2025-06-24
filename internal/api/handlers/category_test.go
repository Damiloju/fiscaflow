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

type mockCategoryService struct {
	mock.Mock
}

func (m *mockCategoryService) CreateCategory(ctx context.Context, req *transaction.CreateCategoryRequest) (*transaction.Category, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*transaction.Category), args.Error(1)
}
func (m *mockCategoryService) GetCategory(ctx context.Context, categoryID uuid.UUID) (*transaction.Category, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).(*transaction.Category), args.Error(1)
}
func (m *mockCategoryService) GetCategories(ctx context.Context, offset, limit int) ([]transaction.Category, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]transaction.Category), args.Error(1)
}
func (m *mockCategoryService) GetDefaultCategories(ctx context.Context) ([]transaction.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]transaction.Category), args.Error(1)
}
func (m *mockCategoryService) UpdateCategory(ctx context.Context, categoryID uuid.UUID, req *transaction.CreateCategoryRequest) (*transaction.Category, error) {
	args := m.Called(ctx, categoryID, req)
	return args.Get(0).(*transaction.Category), args.Error(1)
}
func (m *mockCategoryService) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	args := m.Called(ctx, categoryID)
	return args.Error(0)
}

// Unused service methods for interface compliance
func (m *mockCategoryService) CreateTransaction(context.Context, uuid.UUID, *transaction.CreateTransactionRequest) (*transaction.TransactionResponse, error) {
	return nil, nil
}
func (m *mockCategoryService) GetTransaction(context.Context, uuid.UUID, uuid.UUID) (*transaction.TransactionResponse, error) {
	return nil, nil
}
func (m *mockCategoryService) GetTransactions(context.Context, uuid.UUID, int, int) ([]transaction.TransactionResponse, error) {
	return nil, nil
}
func (m *mockCategoryService) UpdateTransaction(context.Context, uuid.UUID, uuid.UUID, *transaction.UpdateTransactionRequest) (*transaction.TransactionResponse, error) {
	return nil, nil
}
func (m *mockCategoryService) DeleteTransaction(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}
func (m *mockCategoryService) CreateAccount(context.Context, uuid.UUID, *transaction.CreateAccountRequest) (*transaction.Account, error) {
	return nil, nil
}
func (m *mockCategoryService) GetAccount(context.Context, uuid.UUID, uuid.UUID) (*transaction.Account, error) {
	return nil, nil
}
func (m *mockCategoryService) GetAccounts(context.Context, uuid.UUID) ([]transaction.Account, error) {
	return nil, nil
}
func (m *mockCategoryService) UpdateAccount(context.Context, uuid.UUID, uuid.UUID, *transaction.CreateAccountRequest) (*transaction.Account, error) {
	return nil, nil
}
func (m *mockCategoryService) DeleteAccount(context.Context, uuid.UUID, uuid.UUID) error {
	return nil
}

func TestCategoryHandler_CreateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCategoryService)
	h := NewCategoryHandler(mockSvc)
	r := gin.Default()
	r.POST("/categories", h.CreateCategory)

	cat := &transaction.Category{ID: uuid.New(), Name: "Test"}
	mockSvc.On("CreateCategory", mock.Anything, mock.Anything).Return(cat, nil)

	body, _ := json.Marshal(transaction.CreateCategoryRequest{Name: "Test"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var resp transaction.Category
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, cat.ID, resp.ID)
}

func TestCategoryHandler_CreateCategory_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCategoryService)
	h := NewCategoryHandler(mockSvc)
	r := gin.Default()
	r.POST("/categories", h.CreateCategory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewReader([]byte("bad json")))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryHandler_GetCategory_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCategoryService)
	h := NewCategoryHandler(mockSvc)
	r := gin.Default()
	r.GET("/categories/:id", h.GetCategory)

	mockSvc.On("GetCategory", mock.Anything, mock.Anything).Return(&transaction.Category{}, transaction.ErrCategoryNotFound)
	w := httptest.NewRecorder()
	id := uuid.New().String()
	req, _ := http.NewRequest("GET", "/categories/"+id, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCategoryHandler_ListCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCategoryService)
	h := NewCategoryHandler(mockSvc)
	r := gin.Default()
	r.GET("/categories", h.ListCategories)

	mockSvc.On("GetCategories", mock.Anything, 0, 50).Return([]transaction.Category{{ID: uuid.New(), Name: "Test"}}, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_GetDefaultCategories(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCategoryService)
	h := NewCategoryHandler(mockSvc)
	r := gin.Default()
	r.GET("/categories/default", h.GetDefaultCategories)

	mockSvc.On("GetDefaultCategories", mock.Anything).Return([]transaction.Category{{ID: uuid.New(), Name: "Default"}}, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/default", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_UpdateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCategoryService)
	h := NewCategoryHandler(mockSvc)
	r := gin.Default()
	r.PUT("/categories/:id", h.UpdateCategory)

	cat := &transaction.Category{ID: uuid.New(), Name: "Updated"}
	mockSvc.On("UpdateCategory", mock.Anything, mock.Anything, mock.Anything).Return(cat, nil)
	id := cat.ID.String()
	body, _ := json.Marshal(transaction.CreateCategoryRequest{Name: "Updated"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/categories/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCategoryHandler_DeleteCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mockCategoryService)
	h := NewCategoryHandler(mockSvc)
	r := gin.Default()
	r.DELETE("/categories/:id", h.DeleteCategory)

	mockSvc.On("DeleteCategory", mock.Anything, mock.Anything).Return(nil)
	id := uuid.New().String()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/"+id, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
