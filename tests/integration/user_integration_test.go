package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"fiscaflow/internal/domain/user"
)

// TestUser is a SQLite-compatible version of the User model for integration tests
// Remove type TestUser struct { ... }

// TestUserSession is a SQLite-compatible version of the UserSession model
// DeviceInfo is stored as a JSON string
// Remove type TestUserSession struct { ... }

// TableName specifies the table name for TestUser
// Remove func (TestUser) TableName() string { ... }

// TableName specifies the table name for TestUserSession
// Remove func (TestUserSession) TableName() string { ... }

// NewTestDatabase creates a new test database using SQLite in memory
// Remove func NewTestDatabase(t *testing.T) *TestDatabase { ... }

// Cleanup cleans up the test database
// Remove func (td *TestDatabase) Cleanup() { ... }

// TestRepository is a test implementation of the user.Repository interface
type TestRepository struct {
	db *gorm.DB
}

func NewTestRepository(db *gorm.DB) user.Repository {
	return &TestRepository{db: db}
}

func (r *TestRepository) Create(ctx context.Context, u *user.User) error {
	testUser := &TestUser{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Phone:        u.Phone,
		Timezone:     u.Timezone,
		Status:       string(u.Status),
		Role:         string(u.Role),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(testUser).Error
}

func (r *TestRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var testUser TestUser
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&testUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	userID, err := uuid.Parse(testUser.ID)
	if err != nil {
		return nil, err
	}
	return &user.User{
		ID:           userID,
		Email:        testUser.Email,
		PasswordHash: testUser.PasswordHash,
		FirstName:    testUser.FirstName,
		LastName:     testUser.LastName,
		Phone:        testUser.Phone,
		Timezone:     testUser.Timezone,
		Status:       user.UserStatus(testUser.Status),
		Role:         user.UserRole(testUser.Role),
		CreatedAt:    testUser.CreatedAt,
		UpdatedAt:    testUser.UpdatedAt,
	}, nil
}

func (r *TestRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var testUser TestUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&testUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	userID, err := uuid.Parse(testUser.ID)
	if err != nil {
		return nil, err
	}
	return &user.User{
		ID:           userID,
		Email:        testUser.Email,
		PasswordHash: testUser.PasswordHash,
		FirstName:    testUser.FirstName,
		LastName:     testUser.LastName,
		Phone:        testUser.Phone,
		Timezone:     testUser.Timezone,
		Status:       user.UserStatus(testUser.Status),
		Role:         user.UserRole(testUser.Role),
		CreatedAt:    testUser.CreatedAt,
		UpdatedAt:    testUser.UpdatedAt,
	}, nil
}

func (r *TestRepository) Update(ctx context.Context, u *user.User) error {
	testUser := &TestUser{
		ID:           u.ID.String(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Phone:        u.Phone,
		Timezone:     u.Timezone,
		Status:       string(u.Status),
		Role:         string(u.Role),
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
	return r.db.WithContext(ctx).Save(testUser).Error
}

func (r *TestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id.String()).Delete(&TestUser{}).Error
}

func (r *TestRepository) List(ctx context.Context, offset, limit int) ([]user.User, error) {
	var testUsers []TestUser
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&testUsers).Error
	if err != nil {
		return nil, err
	}

	users := make([]user.User, len(testUsers))
	for i, testUser := range testUsers {
		userID, err := uuid.Parse(testUser.ID)
		if err != nil {
			return nil, err
		}

		users[i] = user.User{
			ID:        userID,
			Email:     testUser.Email,
			FirstName: testUser.FirstName,
			LastName:  testUser.LastName,
			Phone:     testUser.Phone,
			Timezone:  testUser.Timezone,
			Status:    user.UserStatus(testUser.Status),
			Role:      user.UserRole(testUser.Role),
			CreatedAt: testUser.CreatedAt,
			UpdatedAt: testUser.UpdatedAt,
		}
	}
	return users, nil
}

func (r *TestRepository) CreateSession(ctx context.Context, session *user.UserSession) error {
	testSession := &TestUserSession{
		ID:           session.ID.String(),
		UserID:       session.UserID.String(),
		RefreshToken: session.RefreshToken,
		DeviceInfo:   "", // Convert DeviceInfo to string if needed
		IPAddress:    session.IPAddress,
		UserAgent:    session.UserAgent,
		ExpiresAt:    session.ExpiresAt,
		CreatedAt:    session.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(testSession).Error
}

func (r *TestRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*user.UserSession, error) {
	var testSession TestUserSession
	err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&testSession).Error
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(testSession.UserID)
	if err != nil {
		return nil, err
	}

	sessionID, err := uuid.Parse(testSession.ID)
	if err != nil {
		return nil, err
	}

	return &user.UserSession{
		ID:           sessionID,
		UserID:       userID,
		RefreshToken: testSession.RefreshToken,
		DeviceInfo:   nil, // Convert string back to DeviceInfo if needed
		IPAddress:    testSession.IPAddress,
		UserAgent:    testSession.UserAgent,
		ExpiresAt:    testSession.ExpiresAt,
		CreatedAt:    testSession.CreatedAt,
	}, nil
}

func (r *TestRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", sessionID.String()).Delete(&TestUserSession{}).Error
}

func (r *TestRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID.String()).Delete(&TestUserSession{}).Error
}

func TestUserIntegration_RegisterAndLogin(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repository
	userRepo := NewTestRepository(db.DB)

	// Test user registration
	testUser := &user.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "+1234567890",
		Timezone:  "America/New_York",
		Status:    user.UserStatusActive,
		Role:      user.UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test creating user
	err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Test retrieving user by ID
	retrievedUser, err := userRepo.GetByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, retrievedUser.ID)
	assert.Equal(t, testUser.Email, retrievedUser.Email)
	assert.Equal(t, testUser.FirstName, retrievedUser.FirstName)
	assert.Equal(t, testUser.LastName, retrievedUser.LastName)

	// Test retrieving user by email
	retrievedUserByEmail, err := userRepo.GetByEmail(context.Background(), testUser.Email)
	require.NoError(t, err)
	assert.Equal(t, testUser.ID, retrievedUserByEmail.ID)
	assert.Equal(t, testUser.Email, retrievedUserByEmail.Email)

	// Test updating user
	testUser.FirstName = "Updated"
	testUser.UpdatedAt = time.Now()
	err = userRepo.Update(context.Background(), testUser)
	require.NoError(t, err)

	// Verify update
	updatedUser, err := userRepo.GetByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", updatedUser.FirstName)

	// Test listing users
	users, err := userRepo.List(context.Background(), 0, 10)
	require.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, testUser.ID, users[0].ID)

	// Test session management
	testSession := &user.UserSession{
		ID:           uuid.New(),
		UserID:       testUser.ID,
		RefreshToken: "refresh_token_123",
		DeviceInfo:   nil,
		IPAddress:    "192.168.1.1",
		UserAgent:    "Mozilla/5.0",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	// Test creating session
	err = userRepo.CreateSession(context.Background(), testSession)
	require.NoError(t, err)

	// Test retrieving session by refresh token
	retrievedSession, err := userRepo.GetSessionByRefreshToken(context.Background(), testSession.RefreshToken)
	require.NoError(t, err)
	assert.Equal(t, testSession.ID, retrievedSession.ID)
	assert.Equal(t, testSession.UserID, retrievedSession.UserID)
	assert.Equal(t, testSession.RefreshToken, retrievedSession.RefreshToken)

	// Test revoking session
	err = userRepo.RevokeSession(context.Background(), testSession.ID)
	require.NoError(t, err)

	// Verify session is revoked
	_, err = userRepo.GetSessionByRefreshToken(context.Background(), testSession.RefreshToken)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserIntegration_ProfileManagement(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repository
	userRepo := NewTestRepository(db.DB)

	// Create a test user
	testUser := &user.User{
		ID:        uuid.New(),
		Email:     "profile@example.com",
		FirstName: "Jane",
		LastName:  "Smith",
		Phone:     "+1234567890",
		Timezone:  "UTC",
		Status:    user.UserStatusActive,
		Role:      user.UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Test profile updates
	testUser.FirstName = "Updated Profile"
	testUser.Phone = "+0987654321"
	testUser.Timezone = "America/New_York"
	testUser.UpdatedAt = time.Now()

	err = userRepo.Update(context.Background(), testUser)
	require.NoError(t, err)

	// Verify profile updates
	updatedUser, err := userRepo.GetByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Profile", updatedUser.FirstName)
	assert.Equal(t, "+0987654321", updatedUser.Phone)
	assert.Equal(t, "America/New_York", updatedUser.Timezone)
}

func TestUserIntegration_TokenManagement(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repository
	userRepo := NewTestRepository(db.DB)

	// Create a test user
	testUser := &user.User{
		ID:        uuid.New(),
		Email:     "token@example.com",
		FirstName: "Token",
		LastName:  "User",
		Status:    user.UserStatusActive,
		Role:      user.UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create multiple sessions
	session1 := &user.UserSession{
		ID:           uuid.New(),
		UserID:       testUser.ID,
		RefreshToken: "refresh_token_1",
		DeviceInfo:   nil,
		IPAddress:    "192.168.1.1",
		UserAgent:    "Mozilla/5.0",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	session2 := &user.UserSession{
		ID:           uuid.New(),
		UserID:       testUser.ID,
		RefreshToken: "refresh_token_2",
		DeviceInfo:   nil,
		IPAddress:    "192.168.1.2",
		UserAgent:    "Chrome/90.0",
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
	}

	err = userRepo.CreateSession(context.Background(), session1)
	require.NoError(t, err)

	err = userRepo.CreateSession(context.Background(), session2)
	require.NoError(t, err)

	// Test revoking all user sessions
	err = userRepo.RevokeAllUserSessions(context.Background(), testUser.ID)
	require.NoError(t, err)

	// Verify all sessions are revoked
	_, err = userRepo.GetSessionByRefreshToken(context.Background(), session1.RefreshToken)
	assert.Error(t, err)

	_, err = userRepo.GetSessionByRefreshToken(context.Background(), session2.RefreshToken)
	assert.Error(t, err)
}

func TestUserIntegration_DuplicateRegistration(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repository
	userRepo := NewTestRepository(db.DB)

	// Create first user
	user1 := &user.User{
		ID:        uuid.New(),
		Email:     "duplicate@example.com",
		FirstName: "First",
		LastName:  "User",
		Status:    user.UserStatusActive,
		Role:      user.UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := userRepo.Create(context.Background(), user1)
	require.NoError(t, err)

	// Try to create second user with same email
	user2 := &user.User{
		ID:        uuid.New(),
		Email:     "duplicate@example.com", // Same email
		FirstName: "Second",
		LastName:  "User",
		Status:    user.UserStatusActive,
		Role:      user.UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = userRepo.Create(context.Background(), user2)
	assert.Error(t, err) // Should fail due to unique constraint
}

func TestUserIntegration_RepositoryOperations(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repository
	userRepo := NewTestRepository(db.DB)

	// Test empty list
	users, err := userRepo.List(context.Background(), 0, 10)
	require.NoError(t, err)
	assert.Len(t, users, 0)

	// Create multiple users
	user1 := &user.User{
		ID:        uuid.New(),
		Email:     "user1@example.com",
		FirstName: "User",
		LastName:  "One",
		Status:    user.UserStatusActive,
		Role:      user.UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user2 := &user.User{
		ID:        uuid.New(),
		Email:     "user2@example.com",
		FirstName: "User",
		LastName:  "Two",
		Status:    user.UserStatusActive,
		Role:      user.UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = userRepo.Create(context.Background(), user1)
	require.NoError(t, err)

	err = userRepo.Create(context.Background(), user2)
	require.NoError(t, err)

	// Test listing with pagination
	users, err = userRepo.List(context.Background(), 0, 1)
	require.NoError(t, err)
	assert.Len(t, users, 1)

	users, err = userRepo.List(context.Background(), 0, 10)
	require.NoError(t, err)
	assert.Len(t, users, 2)

	// Test deleting user
	err = userRepo.Delete(context.Background(), user1.ID)
	require.NoError(t, err)

	// Verify user is deleted
	_, err = userRepo.GetByID(context.Background(), user1.ID)
	assert.Error(t, err)

	// Verify other user still exists
	remainingUser, err := userRepo.GetByID(context.Background(), user2.ID)
	require.NoError(t, err)
	assert.Equal(t, user2.ID, remainingUser.ID)
}
