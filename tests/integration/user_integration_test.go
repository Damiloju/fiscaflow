package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"fiscaflow/internal/domain/user"
)

// TestDatabase represents a test database setup
type TestDatabase struct {
	DB *gorm.DB
}

// TestUser is a SQLite-compatible version of the User model for integration tests
type TestUser struct {
	ID               string     `json:"id" gorm:"type:text;primary_key"`
	Email            string     `json:"email" gorm:"unique;not null"`
	PasswordHash     string     `json:"-" gorm:"not null"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Phone            string     `json:"phone"`
	DateOfBirth      *time.Time `json:"date_of_birth"`
	Timezone         string     `json:"timezone"`
	Locale           string     `json:"locale"`
	Role             string     `json:"role"`
	Status           string     `json:"status"`
	EmailVerified    bool       `json:"email_verified"`
	PhoneVerified    bool       `json:"phone_verified"`
	TwoFactorEnabled bool       `json:"two_factor_enabled"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// TestUserSession is a SQLite-compatible version of the UserSession model
type TestUserSession struct {
	ID              string     `json:"id" gorm:"type:text;primary_key"`
	UserID          string     `json:"user_id" gorm:"type:text;not null"`
	RefreshToken    string     `json:"refresh_token" gorm:"unique;not null"`
	AccessTokenHash string     `json:"access_token_hash"`
	DeviceInfo      string     `json:"device_info"`
	IPAddress       string     `json:"ip_address"`
	UserAgent       string     `json:"user_agent"`
	ExpiresAt       time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt       *time.Time `json:"revoked_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// TableName specifies the table name for TestUser
func (TestUser) TableName() string {
	return "users"
}

// TableName specifies the table name for TestUserSession
func (TestUserSession) TableName() string {
	return "user_sessions"
}

// NewTestDatabase creates a new test database using SQLite in memory
func NewTestDatabase(t *testing.T) *TestDatabase {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema with SQLite-compatible models
	err = db.AutoMigrate(&TestUser{}, &TestUserSession{})
	require.NoError(t, err)

	return &TestDatabase{DB: db}
}

// Cleanup cleans up the test database
func (td *TestDatabase) Cleanup() {
	td.DB.Exec("DELETE FROM user_sessions")
	td.DB.Exec("DELETE FROM users")
}

// TestRepository is a test implementation of the user.Repository interface
type TestRepository struct {
	db *gorm.DB
}

func NewTestRepository(db *gorm.DB) user.Repository {
	return &TestRepository{db: db}
}

func (r *TestRepository) Create(ctx context.Context, u *user.User) error {
	testUser := &TestUser{
		ID:               u.ID.String(),
		Email:            u.Email,
		PasswordHash:     u.PasswordHash,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Phone:            u.Phone,
		DateOfBirth:      u.DateOfBirth,
		Timezone:         u.Timezone,
		Locale:           u.Locale,
		Role:             string(u.Role),
		Status:           string(u.Status),
		EmailVerified:    u.EmailVerified,
		PhoneVerified:    u.PhoneVerified,
		TwoFactorEnabled: u.TwoFactorEnabled,
		LastLoginAt:      u.LastLoginAt,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
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

	userID, _ := uuid.Parse(testUser.ID)
	return &user.User{
		ID:               userID,
		Email:            testUser.Email,
		PasswordHash:     testUser.PasswordHash,
		FirstName:        testUser.FirstName,
		LastName:         testUser.LastName,
		Phone:            testUser.Phone,
		DateOfBirth:      testUser.DateOfBirth,
		Timezone:         testUser.Timezone,
		Locale:           testUser.Locale,
		Role:             user.UserRole(testUser.Role),
		Status:           user.UserStatus(testUser.Status),
		EmailVerified:    testUser.EmailVerified,
		PhoneVerified:    testUser.PhoneVerified,
		TwoFactorEnabled: testUser.TwoFactorEnabled,
		LastLoginAt:      testUser.LastLoginAt,
		CreatedAt:        testUser.CreatedAt,
		UpdatedAt:        testUser.UpdatedAt,
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

	userID, _ := uuid.Parse(testUser.ID)
	return &user.User{
		ID:               userID,
		Email:            testUser.Email,
		PasswordHash:     testUser.PasswordHash,
		FirstName:        testUser.FirstName,
		LastName:         testUser.LastName,
		Phone:            testUser.Phone,
		DateOfBirth:      testUser.DateOfBirth,
		Timezone:         testUser.Timezone,
		Locale:           testUser.Locale,
		Role:             user.UserRole(testUser.Role),
		Status:           user.UserStatus(testUser.Status),
		EmailVerified:    testUser.EmailVerified,
		PhoneVerified:    testUser.PhoneVerified,
		TwoFactorEnabled: testUser.TwoFactorEnabled,
		LastLoginAt:      testUser.LastLoginAt,
		CreatedAt:        testUser.CreatedAt,
		UpdatedAt:        testUser.UpdatedAt,
	}, nil
}

func (r *TestRepository) Update(ctx context.Context, u *user.User) error {
	testUser := &TestUser{
		ID:               u.ID.String(),
		Email:            u.Email,
		PasswordHash:     u.PasswordHash,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Phone:            u.Phone,
		DateOfBirth:      u.DateOfBirth,
		Timezone:         u.Timezone,
		Locale:           u.Locale,
		Role:             string(u.Role),
		Status:           string(u.Status),
		EmailVerified:    u.EmailVerified,
		PhoneVerified:    u.PhoneVerified,
		TwoFactorEnabled: u.TwoFactorEnabled,
		LastLoginAt:      u.LastLoginAt,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
	return r.db.WithContext(ctx).Save(testUser).Error
}

func (r *TestRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&TestUser{}, "id = ?", id.String()).Error
}

func (r *TestRepository) List(ctx context.Context, offset, limit int) ([]user.User, error) {
	var testUsers []TestUser
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&testUsers).Error
	if err != nil {
		return nil, err
	}

	users := make([]user.User, len(testUsers))
	for i, testUser := range testUsers {
		userID, _ := uuid.Parse(testUser.ID)
		users[i] = user.User{
			ID:               userID,
			Email:            testUser.Email,
			PasswordHash:     testUser.PasswordHash,
			FirstName:        testUser.FirstName,
			LastName:         testUser.LastName,
			Phone:            testUser.Phone,
			DateOfBirth:      testUser.DateOfBirth,
			Timezone:         testUser.Timezone,
			Locale:           testUser.Locale,
			Role:             user.UserRole(testUser.Role),
			Status:           user.UserStatus(testUser.Status),
			EmailVerified:    testUser.EmailVerified,
			PhoneVerified:    testUser.PhoneVerified,
			TwoFactorEnabled: testUser.TwoFactorEnabled,
			LastLoginAt:      testUser.LastLoginAt,
			CreatedAt:        testUser.CreatedAt,
			UpdatedAt:        testUser.UpdatedAt,
		}
	}
	return users, nil
}

func (r *TestRepository) CreateSession(ctx context.Context, session *user.UserSession) error {
	testSession := &TestUserSession{
		ID:              session.ID.String(),
		UserID:          session.UserID.String(),
		RefreshToken:    session.RefreshToken,
		AccessTokenHash: session.AccessTokenHash,
		DeviceInfo:      session.DeviceInfo,
		IPAddress:       session.IPAddress,
		UserAgent:       session.UserAgent,
		ExpiresAt:       session.ExpiresAt,
		RevokedAt:       session.RevokedAt,
		CreatedAt:       session.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(testSession).Error
}

func (r *TestRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*user.UserSession, error) {
	var testSession TestUserSession
	err := r.db.WithContext(ctx).Where("refresh_token = ? AND revoked_at IS NULL", refreshToken).First(&testSession).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.ErrSessionNotFound
		}
		return nil, err
	}

	sessionID, _ := uuid.Parse(testSession.ID)
	userID, _ := uuid.Parse(testSession.UserID)

	return &user.UserSession{
		ID:              sessionID,
		UserID:          userID,
		RefreshToken:    testSession.RefreshToken,
		AccessTokenHash: testSession.AccessTokenHash,
		DeviceInfo:      testSession.DeviceInfo,
		IPAddress:       testSession.IPAddress,
		UserAgent:       testSession.UserAgent,
		ExpiresAt:       testSession.ExpiresAt,
		RevokedAt:       testSession.RevokedAt,
		CreatedAt:       testSession.CreatedAt,
	}, nil
}

func (r *TestRepository) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&TestUserSession{}).Where("id = ?", sessionID.String()).Update("revoked_at", time.Now()).Error
}

func (r *TestRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&TestUserSession{}).Where("user_id = ?", userID.String()).Update("revoked_at", time.Now()).Error
}

func TestUserIntegration_RegisterAndLogin(t *testing.T) {
	// Setup
	testDB := NewTestDatabase(t)
	defer testDB.Cleanup()

	userRepo := NewTestRepository(testDB.DB)
	userService := user.NewService(userRepo, "test-secret")

	ctx := context.Background()

	// Test registration
	t.Run("register new user", func(t *testing.T) {
		req := &user.CreateUserRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
			Timezone:  "UTC",
			Locale:    "en-US",
		}

		userResponse, err := userService.Register(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, userResponse)
		assert.Equal(t, req.Email, userResponse.Email)
		assert.Equal(t, req.FirstName, userResponse.FirstName)
		assert.Equal(t, req.LastName, userResponse.LastName)
		assert.Equal(t, user.UserRoleUser, userResponse.Role)
		assert.Equal(t, user.UserStatusActive, userResponse.Status)
		assert.NotEmpty(t, userResponse.ID)
		assert.NotEmpty(t, userResponse.CreatedAt)
	})

	// Test login
	t.Run("login with correct credentials", func(t *testing.T) {
		req := &user.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		loginResponse, err := userService.Login(ctx, req)
		require.NoError(t, err)
		assert.NotNil(t, loginResponse)
		assert.NotNil(t, loginResponse.User)
		assert.Equal(t, "test@example.com", loginResponse.User.Email)
		assert.NotEmpty(t, loginResponse.AccessToken)
		assert.NotEmpty(t, loginResponse.RefreshToken)
		assert.Greater(t, loginResponse.ExpiresIn, int64(0))
	})

	// Test login with wrong password
	t.Run("login with wrong password", func(t *testing.T) {
		req := &user.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		loginResponse, err := userService.Login(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, user.ErrInvalidCredentials, err)
		assert.Nil(t, loginResponse)
	})

	// Test login with non-existent user
	t.Run("login with non-existent user", func(t *testing.T) {
		req := &user.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		loginResponse, err := userService.Login(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, user.ErrInvalidCredentials, err)
		assert.Nil(t, loginResponse)
	})
}

func TestUserIntegration_ProfileManagement(t *testing.T) {
	// Setup
	testDB := NewTestDatabase(t)
	defer testDB.Cleanup()

	userRepo := NewTestRepository(testDB.DB)
	userService := user.NewService(userRepo, "test-secret")

	ctx := context.Background()

	// Create a user first
	createReq := &user.CreateUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	userResponse, err := userService.Register(ctx, createReq)
	require.NoError(t, err)
	userID := userResponse.ID

	// Test get profile
	t.Run("get user profile", func(t *testing.T) {
		profile, err := userService.GetProfile(ctx, userID)
		require.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, userID, profile.ID)
		assert.Equal(t, "test@example.com", profile.Email)
		assert.Equal(t, "John", profile.FirstName)
		assert.Equal(t, "Doe", profile.LastName)
	})

	// Test update profile
	t.Run("update user profile", func(t *testing.T) {
		updateReq := &user.UpdateUserRequest{
			FirstName: "Jane",
			LastName:  "Smith",
			Phone:     "+1234567890",
		}

		updatedProfile, err := userService.UpdateProfile(ctx, userID, updateReq)
		require.NoError(t, err)
		assert.NotNil(t, updatedProfile)
		assert.Equal(t, "Jane", updatedProfile.FirstName)
		assert.Equal(t, "Smith", updatedProfile.LastName)
		assert.Equal(t, "+1234567890", updatedProfile.Phone)

		// Verify the update persisted
		profile, err := userService.GetProfile(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "Jane", profile.FirstName)
		assert.Equal(t, "Smith", profile.LastName)
		assert.Equal(t, "+1234567890", profile.Phone)
	})

	// Test get non-existent user
	t.Run("get non-existent user profile", func(t *testing.T) {
		nonExistentID := uuid.New()
		profile, err := userService.GetProfile(ctx, nonExistentID)
		assert.Error(t, err)
		assert.Equal(t, user.ErrUserNotFound, err)
		assert.Nil(t, profile)
	})
}

func TestUserIntegration_TokenManagement(t *testing.T) {
	// Setup
	testDB := NewTestDatabase(t)
	defer testDB.Cleanup()

	userRepo := NewTestRepository(testDB.DB)
	userService := user.NewService(userRepo, "test-secret")

	ctx := context.Background()

	// Create a user and login to get tokens
	createReq := &user.CreateUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	_, err := userService.Register(ctx, createReq)
	require.NoError(t, err)

	loginReq := &user.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginResponse, err := userService.Login(ctx, loginReq)
	require.NoError(t, err)
	refreshToken := loginResponse.RefreshToken

	// Test refresh token
	t.Run("refresh access token", func(t *testing.T) {
		newLoginResponse, err := userService.RefreshToken(ctx, refreshToken)
		require.NoError(t, err)
		assert.NotNil(t, newLoginResponse)
		assert.NotEmpty(t, newLoginResponse.AccessToken)
		assert.Equal(t, refreshToken, newLoginResponse.RefreshToken) // Same refresh token
		assert.Greater(t, newLoginResponse.ExpiresIn, int64(0))
	})

	// Test invalid refresh token
	t.Run("refresh with invalid token", func(t *testing.T) {
		newLoginResponse, err := userService.RefreshToken(ctx, "invalid-token")
		assert.Error(t, err)
		assert.Equal(t, user.ErrInvalidRefreshToken, err)
		assert.Nil(t, newLoginResponse)
	})

	// Test token validation
	t.Run("validate access token", func(t *testing.T) {
		claims, err := userService.ValidateToken(ctx, loginResponse.AccessToken)
		// The token should be valid since it was just generated
		require.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "test@example.com", claims.Email)
		assert.Equal(t, user.UserRoleUser, claims.Role)
	})
}

func TestUserIntegration_DuplicateRegistration(t *testing.T) {
	// Setup
	testDB := NewTestDatabase(t)
	defer testDB.Cleanup()

	userRepo := NewTestRepository(testDB.DB)
	userService := user.NewService(userRepo, "test-secret")

	ctx := context.Background()

	// Register first user
	req := &user.CreateUserRequest{
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
	}

	userResponse, err := userService.Register(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, userResponse)

	// Try to register with same email
	t.Run("duplicate registration", func(t *testing.T) {
		duplicateReq := &user.CreateUserRequest{
			Email:     "test@example.com",
			Password:  "differentpassword",
			FirstName: "Jane",
			LastName:  "Smith",
		}

		duplicateResponse, err := userService.Register(ctx, duplicateReq)
		assert.Error(t, err)
		assert.Equal(t, user.ErrUserAlreadyExists, err)
		assert.Nil(t, duplicateResponse)
	})
}

func TestUserIntegration_RepositoryOperations(t *testing.T) {
	// Setup
	testDB := NewTestDatabase(t)
	defer testDB.Cleanup()

	userRepo := NewTestRepository(testDB.DB)
	ctx := context.Background()

	// Test create user
	t.Run("create user via repository", func(t *testing.T) {
		testUser := &user.User{
			ID:           uuid.New(),
			Email:        "repo@example.com",
			PasswordHash: "hashedpassword",
			FirstName:    "Repo",
			LastName:     "Test",
			Role:         user.UserRoleUser,
			Status:       user.UserStatusActive,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err := userRepo.Create(ctx, testUser)
		require.NoError(t, err)
	})

	// Test get user by email
	t.Run("get user by email", func(t *testing.T) {
		foundUser, err := userRepo.GetByEmail(ctx, "repo@example.com")
		require.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, "repo@example.com", foundUser.Email)
		assert.Equal(t, "Repo", foundUser.FirstName)
		assert.Equal(t, "Test", foundUser.LastName)
	})

	// Test get user by ID
	t.Run("get user by ID", func(t *testing.T) {
		foundUser, err := userRepo.GetByEmail(ctx, "repo@example.com")
		require.NoError(t, err)

		userByID, err := userRepo.GetByID(ctx, foundUser.ID)
		require.NoError(t, err)
		assert.NotNil(t, userByID)
		assert.Equal(t, foundUser.ID, userByID.ID)
		assert.Equal(t, foundUser.Email, userByID.Email)
	})

	// Test update user
	t.Run("update user", func(t *testing.T) {
		foundUser, err := userRepo.GetByEmail(ctx, "repo@example.com")
		require.NoError(t, err)

		foundUser.FirstName = "Updated"
		foundUser.UpdatedAt = time.Now()

		err = userRepo.Update(ctx, foundUser)
		require.NoError(t, err)

		// Verify update
		updatedUser, err := userRepo.GetByID(ctx, foundUser.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated", updatedUser.FirstName)
	})

	// Test list users
	t.Run("list users", func(t *testing.T) {
		users, err := userRepo.List(ctx, 0, 10)
		require.NoError(t, err)
		assert.Len(t, users, 1) // We created one user in this test
	})
}
