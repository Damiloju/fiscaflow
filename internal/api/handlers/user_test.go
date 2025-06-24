package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"fiscaflow/internal/domain/user"
)

// MockUserService is a mock implementation of the user.Service interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserResponse), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.LoginResponse), args.Error(1)
}

func (m *MockUserService) GetProfile(ctx context.Context, userID uuid.UUID) (*user.UserResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserResponse), args.Error(1)
}

func (m *MockUserService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.UserResponse), args.Error(1)
}

func (m *MockUserService) RefreshToken(ctx context.Context, refreshToken string) (*user.LoginResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.LoginResponse), args.Error(1)
}

func (m *MockUserService) Logout(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	args := m.Called(ctx, userID, sessionID)
	return args.Error(0)
}

func (m *MockUserService) ValidateToken(ctx context.Context, tokenString string) (*user.Claims, error) {
	args := m.Called(ctx, tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.Claims), args.Error(1)
}

func setupTestRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())

	// Setup routes
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/register", handler.Register)
			users.POST("/login", handler.Login)
			users.POST("/refresh-token", handler.RefreshToken)
			users.POST("/logout", handler.Logout)

			// Protected routes
			protected := users.Group("")
			protected.Use(func(c *gin.Context) {
				// Mock authentication middleware for testing
				userID := c.GetHeader("X-User-ID")
				if userID != "" {
					if id, err := uuid.Parse(userID); err == nil {
						c.Set("user_id", id)
					}
				}
				c.Next()
			})
			{
				protected.GET("/profile", handler.GetProfile)
				protected.PUT("/profile", handler.UpdateProfile)
			}
		}
	}

	return router
}

func TestUserHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful registration",
			requestBody: user.CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Timezone:  "UTC",
				Locale:    "en-US",
			},
			mockSetup: func(service *MockUserService) {
				expectedResponse := &user.UserResponse{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Timezone:  "UTC",
					Locale:    "en-US",
					Role:      user.UserRoleUser,
					Status:    user.UserStatusActive,
				}
				service.On("Register", mock.Anything, mock.AnythingOfType("*user.CreateUserRequest")).Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "User registered successfully",
			},
		},
		{
			name: "user already exists",
			requestBody: user.CreateUserRequest{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func(service *MockUserService) {
				service.On("Register", mock.Anything, mock.AnythingOfType("*user.CreateUserRequest")).Return(nil, user.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"error": "User already exists",
			},
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
				// Missing required fields
			},
			mockSetup: func(service *MockUserService) {
				// No mock setup needed for invalid request
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			logger := zap.NewNop()
			handler := NewUserHandler(mockService, logger)
			router := setupTestRouter(handler)

			// Create request
			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful login",
			requestBody: user.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(service *MockUserService) {
				expectedResponse := &user.LoginResponse{
					User: &user.UserResponse{
						ID:        uuid.New(),
						Email:     "test@example.com",
						FirstName: "John",
						LastName:  "Doe",
						Role:      user.UserRoleUser,
						Status:    user.UserStatusActive,
					},
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					ExpiresIn:    900,
				}
				service.On("Login", mock.Anything, mock.AnythingOfType("*user.LoginRequest")).Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Login successful",
			},
		},
		{
			name: "invalid credentials",
			requestBody: user.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(service *MockUserService) {
				service.On("Login", mock.Anything, mock.AnythingOfType("*user.LoginRequest")).Return(nil, user.ErrInvalidCredentials)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			logger := zap.NewNop()
			handler := NewUserHandler(mockService, logger)
			router := setupTestRouter(handler)

			// Create request
			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	tests := []struct {
		name           string
		userID         uuid.UUID
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful profile retrieval",
			userID: uuid.New(),
			mockSetup: func(service *MockUserService) {
				expectedResponse := &user.UserResponse{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Role:      user.UserRoleUser,
					Status:    user.UserStatusActive,
				}
				service.On("GetProfile", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Profile retrieved successfully",
			},
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(service *MockUserService) {
				service.On("GetProfile", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, user.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "User not found",
			},
		},
		{
			name:   "unauthorized - no user_id in context",
			userID: uuid.Nil,
			mockSetup: func(service *MockUserService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Unauthorized",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			logger := zap.NewNop()
			handler := NewUserHandler(mockService, logger)
			router := setupTestRouter(handler)

			// Create request
			req, _ := http.NewRequest("GET", "/api/v1/users/profile", nil)
			if tt.userID != uuid.Nil {
				req.Header.Set("X-User-ID", tt.userID.String())
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	tests := []struct {
		name           string
		userID         uuid.UUID
		requestBody    interface{}
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:   "successful profile update",
			userID: uuid.New(),
			requestBody: user.UpdateUserRequest{
				FirstName: "Jane",
				LastName:  "Smith",
				Phone:     "+1234567890",
			},
			mockSetup: func(service *MockUserService) {
				expectedResponse := &user.UserResponse{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "Jane",
					LastName:  "Smith",
					Phone:     "+1234567890",
					Role:      user.UserRoleUser,
					Status:    user.UserStatusActive,
				}
				service.On("UpdateProfile", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*user.UpdateUserRequest")).Return(expectedResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Profile updated successfully",
			},
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			requestBody: user.UpdateUserRequest{
				FirstName: "Jane",
			},
			mockSetup: func(service *MockUserService) {
				service.On("UpdateProfile", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*user.UpdateUserRequest")).Return(nil, user.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "User not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			tt.mockSetup(mockService)

			logger := zap.NewNop()
			handler := NewUserHandler(mockService, logger)
			router := setupTestRouter(handler)

			// Create request
			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/api/v1/users/profile", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", tt.userID.String())

			// Create response recorder
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, response[key])
			}

			mockService.AssertExpectations(t)
		})
	}
}
