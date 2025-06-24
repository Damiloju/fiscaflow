package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDatabase represents a test database setup
// (reused by both user and transaction integration tests)
type TestDatabase struct {
	DB *gorm.DB
}

// NewTestDatabase creates a new test database using SQLite in memory
func NewTestDatabase(t *testing.T) *TestDatabase {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the schema with all models needed for integration tests
	err = db.AutoMigrate(
		&TestUser{}, &TestUserSession{},
		&TestTransaction{}, &TestCategory{}, &TestAccount{},
	)
	require.NoError(t, err)

	return &TestDatabase{DB: db}
}

// Cleanup cleans up the test database
func (td *TestDatabase) Cleanup() {
	td.DB.Exec("DELETE FROM transactions")
	td.DB.Exec("DELETE FROM categories")
	td.DB.Exec("DELETE FROM accounts")
	td.DB.Exec("DELETE FROM user_sessions")
	td.DB.Exec("DELETE FROM users")
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
	Timezone         string     `json:"timezone" gorm:"default:'UTC'"`
	Locale           string     `json:"locale" gorm:"default:'en-US'"`
	Role             string     `json:"role" gorm:"default:'user'"`
	Status           string     `json:"status" gorm:"default:'active'"`
	EmailVerified    bool       `json:"email_verified" gorm:"default:false"`
	PhoneVerified    bool       `json:"phone_verified" gorm:"default:false"`
	TwoFactorEnabled bool       `json:"two_factor_enabled" gorm:"default:false"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// TableName specifies the table name for TestUser
func (TestUser) TableName() string {
	return "users"
}

// TestUserSession is a SQLite-compatible version of the UserSession model
// DeviceInfo is stored as a JSON string
type TestUserSession struct {
	ID              string     `json:"id" gorm:"type:text;primary_key"`
	UserID          string     `json:"user_id" gorm:"type:text;not null"`
	RefreshToken    string     `json:"refresh_token" gorm:"unique;not null"`
	AccessTokenHash string     `json:"access_token_hash"`
	DeviceInfo      string     `json:"device_info" gorm:"type:text"` // Store as JSON string for SQLite
	IPAddress       string     `json:"ip_address"`
	UserAgent       string     `json:"user_agent"`
	ExpiresAt       time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt       *time.Time `json:"revoked_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// TableName specifies the table name for TestUserSession
func (TestUserSession) TableName() string {
	return "user_sessions"
}

// TestTransaction is a SQLite-compatible version of the Transaction model for integration tests
type TestTransaction struct {
	ID         string  `json:"id" gorm:"type:text;primary_key"`
	UserID     string  `json:"user_id" gorm:"type:text;not null"`
	FamilyID   *string `json:"family_id" gorm:"type:text"`
	AccountID  string  `json:"account_id" gorm:"type:text;not null"`
	CategoryID *string `json:"category_id" gorm:"type:text"`

	Amount      float64 `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency    string  `json:"currency" gorm:"default:'USD'"`
	Description string  `json:"description" gorm:"not null"`
	Merchant    string  `json:"merchant"`
	Location    string  `json:"location" gorm:"type:text"`

	TransactionDate time.Time  `json:"transaction_date" gorm:"not null"`
	PostedDate      *time.Time `json:"posted_date"`
	Status          string     `json:"status" gorm:"default:'pending'"`

	CategorizationSource     string   `json:"categorization_source" gorm:"default:'manual'"`
	CategorizationConfidence *float64 `json:"categorization_confidence"`

	Tags       string `json:"tags" gorm:"type:text"` // Store as JSON string for SQLite
	Notes      string `json:"notes"`
	ReceiptURL string `json:"receipt_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for TestTransaction
func (TestTransaction) TableName() string {
	return "transactions"
}

// TestCategory is a SQLite-compatible version of the Category model for integration tests
type TestCategory struct {
	ID          string    `json:"id" gorm:"type:text;primary_key"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Color       string    `json:"color"`
	ParentID    *string   `json:"parent_id" gorm:"type:text"`
	IsDefault   bool      `json:"is_default" gorm:"default:false"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	SortOrder   int       `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for TestCategory
func (TestCategory) TableName() string {
	return "categories"
}

// TestAccount is a SQLite-compatible version of the Account model for integration tests
type TestAccount struct {
	ID                string     `json:"id" gorm:"type:text;primary_key"`
	UserID            string     `json:"user_id" gorm:"type:text;not null"`
	FamilyID          *string    `json:"family_id" gorm:"type:text"`
	Name              string     `json:"name" gorm:"not null"`
	Type              string     `json:"type" gorm:"not null"`
	Institution       string     `json:"institution"`
	AccountNumberHash string     `json:"account_number_hash"`
	Balance           float64    `json:"balance" gorm:"type:decimal(15,2);default:0.00"`
	Currency          string     `json:"currency" gorm:"default:'USD'"`
	IsActive          bool       `json:"is_active" gorm:"default:true"`
	PlaidAccountID    string     `json:"plaid_account_id"`
	LastSyncAt        *time.Time `json:"last_sync_at"`
	Settings          string     `json:"settings" gorm:"type:text;default:'{}'"` // Store as JSON string for SQLite
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// TableName specifies the table name for TestAccount
func (TestAccount) TableName() string {
	return "accounts"
}
