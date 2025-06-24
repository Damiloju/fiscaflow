package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"
	"fiscaflow/internal/api/handlers"
	"fiscaflow/internal/api/middleware"
	"fiscaflow/internal/domain/transaction"
	"fiscaflow/internal/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServerWithAuth(t *testing.T) (*gin.Engine, string) {
	db := NewTestDatabase(t)
	userRepo := NewTestRepository(db.DB)
	transactionRepo := NewTestTransactionRepository(db.DB)
	userService := user.NewService(userRepo, "test-secret")
	transactionService := transaction.NewService(transactionRepo)

	userHandler := handlers.NewUserHandler(userService, nil)
	categoryHandler := handlers.NewCategoryHandler(transactionService)
	accountHandler := handlers.NewAccountHandler(transactionService)

	r := gin.Default()
	api := r.Group("/api/v1")
	userHandler.RegisterRoutes(api)
	api.Use(middleware.AuthMiddleware(userService))
	categoryHandler.RegisterRoutes(api)
	accountHandler.RegisterRoutes(api)

	// Register and login test user
	registerReq := user.CreateUserRequest{
		Email:     "testcatacc@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	_, err := userService.Register(context.Background(), &registerReq)
	if err != nil && err != user.ErrUserAlreadyExists {
		t.Fatalf("failed to register test user: %v", err)
	}
	loginReq := user.LoginRequest{Email: registerReq.Email, Password: registerReq.Password}
	loginResp, err := userService.Login(context.Background(), &loginReq)
	require.NoError(t, err)
	return r, loginResp.AccessToken
}

func TestCategoryAPI_CRUD(t *testing.T) {
	r, token := setupTestServerWithAuth(t)

	// Create category
	catReq := transaction.CreateCategoryRequest{Name: "Groceries", Description: "Food", Icon: "cart", Color: "#fff"}
	body, _ := json.Marshal(catReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/categories", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createdCat transaction.Category
	_ = json.Unmarshal(w.Body.Bytes(), &createdCat)
	assert.Equal(t, catReq.Name, createdCat.Name)

	// Get category
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/categories/"+createdCat.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Update category
	updateReq := transaction.CreateCategoryRequest{Name: "Updated", Description: "Updated desc", Icon: "cart", Color: "#000"}
	body, _ = json.Marshal(updateReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/categories/"+createdCat.ID.String(), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// List categories
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/categories", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Delete category
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/categories/"+createdCat.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestAccountAPI_CRUD(t *testing.T) {
	r, token := setupTestServerWithAuth(t)

	// Create account
	accReq := transaction.CreateAccountRequest{Name: "Checking", Type: transaction.AccountTypeChecking, Institution: "Test Bank", Balance: 1000, Currency: "USD"}
	body, _ := json.Marshal(accReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/accounts", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createdAcc transaction.Account
	_ = json.Unmarshal(w.Body.Bytes(), &createdAcc)
	assert.Equal(t, accReq.Name, createdAcc.Name)

	// Get account
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/accounts/"+createdAcc.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Update account
	updateReq := transaction.CreateAccountRequest{Name: "Updated Account", Type: transaction.AccountTypeChecking, Institution: "Test Bank", Balance: 2000, Currency: "USD"}
	body, _ = json.Marshal(updateReq)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/accounts/"+createdAcc.ID.String(), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// List accounts
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/accounts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Delete account
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/accounts/"+createdAcc.ID.String(), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
