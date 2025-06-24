package transaction

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID         uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	FamilyID   *uuid.UUID `json:"family_id" gorm:"type:uuid"`
	AccountID  uuid.UUID  `json:"account_id" gorm:"type:uuid;not null"`
	CategoryID *uuid.UUID `json:"category_id" gorm:"type:uuid"`

	Amount      float64 `json:"amount" gorm:"type:decimal(15,2);not null"`
	Currency    string  `json:"currency" gorm:"default:'USD'"`
	Description string  `json:"description" gorm:"not null"`
	Merchant    string  `json:"merchant"`
	Location    string  `json:"location" gorm:"type:jsonb"`

	TransactionDate time.Time         `json:"transaction_date" gorm:"not null"`
	PostedDate      *time.Time        `json:"posted_date"`
	Status          TransactionStatus `json:"status" gorm:"default:'pending'"`

	CategorizationSource     CategorizationSource `json:"categorization_source" gorm:"default:'manual'"`
	CategorizationConfidence *float64             `json:"categorization_confidence"`

	Tags       []string `json:"tags" gorm:"type:text[]"`
	Notes      string   `json:"notes"`
	ReceiptURL string   `json:"receipt_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusPosted    TransactionStatus = "posted"
	TransactionStatusCancelled TransactionStatus = "cancelled"
	TransactionStatusDisputed  TransactionStatus = "disputed"
)

// CategorizationSource represents how the transaction was categorized
type CategorizationSource string

const (
	CategorizationSourceManual         CategorizationSource = "manual"
	CategorizationSourceML             CategorizationSource = "ml"
	CategorizationSourcePlaid          CategorizationSource = "plaid"
	CategorizationSourceUserCorrection CategorizationSource = "user_correction"
)

// Category represents a transaction category
type Category struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Color       string     `json:"color"`
	ParentID    *uuid.UUID `json:"parent_id" gorm:"type:uuid"`
	IsDefault   bool       `json:"is_default" gorm:"default:false"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Account represents a financial account
type Account struct {
	ID                uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID            uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	FamilyID          *uuid.UUID  `json:"family_id" gorm:"type:uuid"`
	Name              string      `json:"name" gorm:"not null"`
	Type              AccountType `json:"type" gorm:"not null"`
	Institution       string      `json:"institution"`
	AccountNumberHash string      `json:"account_number_hash"`
	Balance           float64     `json:"balance" gorm:"type:decimal(15,2);default:0.00"`
	Currency          string      `json:"currency" gorm:"default:'USD'"`
	IsActive          bool        `json:"is_active" gorm:"default:true"`
	PlaidAccountID    string      `json:"plaid_account_id"`
	LastSyncAt        *time.Time  `json:"last_sync_at"`
	Settings          string      `json:"settings" gorm:"type:jsonb;default:'{}'"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeCreditCard AccountType = "credit_card"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeLoan       AccountType = "loan"
	AccountTypeOther      AccountType = "other"
)

// CreateTransactionRequest represents a request to create a new transaction
type CreateTransactionRequest struct {
	AccountID       uuid.UUID  `json:"account_id" binding:"required"`
	CategoryID      *uuid.UUID `json:"category_id"`
	Amount          float64    `json:"amount" binding:"required"`
	Currency        string     `json:"currency"`
	Description     string     `json:"description" binding:"required"`
	Merchant        string     `json:"merchant"`
	Location        string     `json:"location"`
	TransactionDate time.Time  `json:"transaction_date" binding:"required"`
	PostedDate      *time.Time `json:"posted_date"`
	Tags            []string   `json:"tags"`
	Notes           string     `json:"notes"`
}

// UpdateTransactionRequest represents a request to update a transaction
type UpdateTransactionRequest struct {
	CategoryID      *uuid.UUID         `json:"category_id"`
	Amount          *float64           `json:"amount"`
	Currency        string             `json:"currency"`
	Description     string             `json:"description"`
	Merchant        string             `json:"merchant"`
	Location        string             `json:"location"`
	TransactionDate *time.Time         `json:"transaction_date"`
	PostedDate      *time.Time         `json:"posted_date"`
	Status          *TransactionStatus `json:"status"`
	Tags            []string           `json:"tags"`
	Notes           string             `json:"notes"`
}

// TransactionResponse represents a transaction response
type TransactionResponse struct {
	ID                       uuid.UUID            `json:"id"`
	UserID                   uuid.UUID            `json:"user_id"`
	FamilyID                 *uuid.UUID           `json:"family_id"`
	AccountID                uuid.UUID            `json:"account_id"`
	CategoryID               *uuid.UUID           `json:"category_id"`
	Amount                   float64              `json:"amount"`
	Currency                 string               `json:"currency"`
	Description              string               `json:"description"`
	Merchant                 string               `json:"merchant"`
	Location                 string               `json:"location"`
	TransactionDate          time.Time            `json:"transaction_date"`
	PostedDate               *time.Time           `json:"posted_date"`
	Status                   TransactionStatus    `json:"status"`
	CategorizationSource     CategorizationSource `json:"categorization_source"`
	CategorizationConfidence *float64             `json:"categorization_confidence"`
	Tags                     []string             `json:"tags"`
	Notes                    string               `json:"notes"`
	ReceiptURL               string               `json:"receipt_url"`
	CreatedAt                time.Time            `json:"created_at"`
	UpdatedAt                time.Time            `json:"updated_at"`
}

// CreateCategoryRequest represents a request to create a new category
type CreateCategoryRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	Color       string     `json:"color"`
	ParentID    *uuid.UUID `json:"parent_id"`
	IsDefault   bool       `json:"is_default"`
	SortOrder   int        `json:"sort_order"`
}

// CreateAccountRequest represents a request to create a new account
type CreateAccountRequest struct {
	Name              string      `json:"name" binding:"required"`
	Type              AccountType `json:"type" binding:"required"`
	Institution       string      `json:"institution"`
	AccountNumberHash string      `json:"account_number_hash"`
	Balance           float64     `json:"balance"`
	Currency          string      `json:"currency"`
	PlaidAccountID    string      `json:"plaid_account_id"`
}

// TableName specifies the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}

// TableName specifies the table name for Category
func (Category) TableName() string {
	return "categories"
}

// TableName specifies the table name for Account
func (Account) TableName() string {
	return "accounts"
}
