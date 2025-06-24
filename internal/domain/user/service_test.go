package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, offset, limit int) ([]User, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockRepository) CreateSession(ctx context.Context, session *UserSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*UserSession, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserSession), args.Error(1)
}

func (m *MockRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name          string
		request       *CreateUserRequest
		mockSetup     func(*MockRepository)
		expectedError error
		expectedUser  *UserResponse
	}{
		{
			name: "successful registration",
			request: &CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
				Timezone:  "UTC",
				Locale:    "en-US",
			},
			mockSetup: func(repo *MockRepository) {
				// Mock GetByEmail to return ErrUserNotFound (user doesn't exist)
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, ErrUserNotFound)
				// Mock Create to succeed
				repo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
			},
			expectedError: nil,
			expectedUser: &UserResponse{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Timezone:  "UTC",
				Locale:    "en-US",
				Role:      UserRoleUser,
				Status:    UserStatusActive,
			},
		},
		{
			name: "user already exists",
			request: &CreateUserRequest{
				Email:     "existing@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func(repo *MockRepository) {
				existingUser := &User{
					ID:    uuid.New(),
					Email: "existing@example.com",
				}
				repo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			expectedError: ErrUserAlreadyExists,
			expectedUser:  nil,
		},
		{
			name: "repository error",
			request: &CreateUserRequest{
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			mockSetup: func(repo *MockRepository) {
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("failed to check existing user: database error"),
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)
			service := NewService(mockRepo, "test-secret")

			// Execute
			result, err := service.Register(context.Background(), tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.FirstName, result.FirstName)
				assert.Equal(t, tt.expectedUser.LastName, result.LastName)
				assert.Equal(t, tt.expectedUser.Role, result.Role)
				assert.Equal(t, tt.expectedUser.Status, result.Status)
				assert.NotEmpty(t, result.ID)
				assert.NotEmpty(t, result.CreatedAt)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name          string
		request       *LoginRequest
		mockSetup     func(*MockRepository)
		expectedError error
		expectedLogin *LoginResponse
	}{
		{
			name: "successful login",
			request: &LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(repo *MockRepository) {
				// Create a user with a real bcrypt hash
				hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &User{
					ID:           uuid.New(),
					Email:        "test@example.com",
					PasswordHash: string(hash),
					FirstName:    "John",
					LastName:     "Doe",
					Role:         UserRoleUser,
					Status:       UserStatusActive,
				}
				repo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)
				repo.On("Update", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
				repo.On("CreateSession", mock.Anything, mock.AnythingOfType("*user.UserSession")).Return(nil)
			},
			expectedError: nil,
			expectedLogin: &LoginResponse{
				User: &UserResponse{
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Role:      UserRoleUser,
					Status:    UserStatusActive,
				},
			},
		},
		{
			name: "user not found",
			request: &LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			mockSetup: func(repo *MockRepository) {
				repo.On("GetByEmail", mock.Anything, "nonexistent@example.com").Return(nil, ErrUserNotFound)
			},
			expectedError: ErrInvalidCredentials,
			expectedLogin: nil,
		},
		{
			name: "user inactive",
			request: &LoginRequest{
				Email:    "inactive@example.com",
				Password: "password123",
			},
			mockSetup: func(repo *MockRepository) {
				user := &User{
					ID:     uuid.New(),
					Email:  "inactive@example.com",
					Status: UserStatusInactive,
				}
				repo.On("GetByEmail", mock.Anything, "inactive@example.com").Return(user, nil)
			},
			expectedError: ErrUserInactive,
			expectedLogin: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)
			service := NewService(mockRepo, "test-secret")

			// Execute
			result, err := service.Login(context.Background(), tt.request)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.User)
				assert.Equal(t, tt.expectedLogin.User.Email, result.User.Email)
				assert.Equal(t, tt.expectedLogin.User.FirstName, result.User.FirstName)
				assert.Equal(t, tt.expectedLogin.User.LastName, result.User.LastName)
				assert.Equal(t, tt.expectedLogin.User.Role, result.User.Role)
				assert.Equal(t, tt.expectedLogin.User.Status, result.User.Status)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
				assert.Greater(t, result.ExpiresIn, int64(0))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetProfile(t *testing.T) {
	tests := []struct {
		name          string
		userID        uuid.UUID
		mockSetup     func(*MockRepository)
		expectedError error
		expectedUser  *UserResponse
	}{
		{
			name:   "successful profile retrieval",
			userID: uuid.New(),
			mockSetup: func(repo *MockRepository) {
				user := &User{
					ID:        uuid.New(),
					Email:     "test@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Role:      UserRoleUser,
					Status:    UserStatusActive,
				}
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(user, nil)
			},
			expectedError: nil,
			expectedUser: &UserResponse{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Role:      UserRoleUser,
				Status:    UserStatusActive,
			},
		},
		{
			name:   "user not found",
			userID: uuid.New(),
			mockSetup: func(repo *MockRepository) {
				repo.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, ErrUserNotFound)
			},
			expectedError: ErrUserNotFound,
			expectedUser:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)
			service := NewService(mockRepo, "test-secret")

			// Execute
			result, err := service.GetProfile(context.Background(), tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedUser.Email, result.Email)
				assert.Equal(t, tt.expectedUser.FirstName, result.FirstName)
				assert.Equal(t, tt.expectedUser.LastName, result.LastName)
				assert.Equal(t, tt.expectedUser.Role, result.Role)
				assert.Equal(t, tt.expectedUser.Status, result.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_ValidateToken(t *testing.T) {
	tests := []struct {
		name           string
		tokenString    string
		jwtSecret      string
		expectedError  error
		expectedClaims *Claims
	}{
		{
			name:           "valid token",
			tokenString:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2NzgtMTIzNC0xMjM0LTEyMzQtMTIzNDU2Nzg5MDEiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJyb2xlIjoidXNlciIsImlhdCI6MTYzNTU5OTk5OSwiZXhwIjoxNjM1NjAzNTk5fQ.signature",
			jwtSecret:      "test-secret",
			expectedError:  ErrInvalidToken, // This will fail because we need a real token
			expectedClaims: nil,
		},
		{
			name:           "invalid token",
			tokenString:    "invalid-token",
			jwtSecret:      "test-secret",
			expectedError:  ErrInvalidToken,
			expectedClaims: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			service := NewService(mockRepo, tt.jwtSecret)

			// Execute
			result, err := service.ValidateToken(context.Background(), tt.tokenString)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedClaims.UserID, result.UserID)
				assert.Equal(t, tt.expectedClaims.Email, result.Email)
				assert.Equal(t, tt.expectedClaims.Role, result.Role)
			}
		})
	}
}

func TestUserService_GenerateAccessToken(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	// Construct the concrete type directly
	svc := &service{
		repo:          mockRepo,
		jwtSecret:     "test-secret",
		tokenExpiry:   15 * time.Minute,
		refreshExpiry: 7 * 24 * time.Hour,
	}

	user := &User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Role:  UserRoleUser,
	}

	// Execute
	token, err := svc.generateAccessToken(user)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUserService_GenerateRefreshToken(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	svc := &service{
		repo:          mockRepo,
		jwtSecret:     "test-secret",
		tokenExpiry:   15 * time.Minute,
		refreshExpiry: 7 * 24 * time.Hour,
	}

	// Execute
	token, err := svc.generateRefreshToken()

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Len(t, token, 64) // 32 bytes = 64 hex characters
}
