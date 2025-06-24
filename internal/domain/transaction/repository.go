package transaction

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for transaction data operations
type Repository interface {
	// Transaction operations
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetTransactionsByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]Transaction, error)
	GetTransactionsByAccount(ctx context.Context, accountID uuid.UUID, offset, limit int) ([]Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
	DeleteTransaction(ctx context.Context, id uuid.UUID) error

	// Category operations
	CreateCategory(ctx context.Context, category *Category) error
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)
	GetCategories(ctx context.Context, offset, limit int) ([]Category, error)
	GetDefaultCategories(ctx context.Context) ([]Category, error)
	UpdateCategory(ctx context.Context, category *Category) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error

	// Account operations
	CreateAccount(ctx context.Context, account *Account) error
	GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error)
	GetAccountsByUser(ctx context.Context, userID uuid.UUID) ([]Account, error)
	UpdateAccount(ctx context.Context, account *Account) error
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new transaction repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Transaction operations

// CreateTransaction creates a new transaction
func (r *repository) CreateTransaction(ctx context.Context, transaction *Transaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// GetTransactionByID retrieves a transaction by ID
func (r *repository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	var transaction Transaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}
	return &transaction, nil
}

// GetTransactionsByUser retrieves transactions for a user with pagination
func (r *repository) GetTransactionsByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("transaction_date DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

// GetTransactionsByAccount retrieves transactions for an account with pagination
func (r *repository) GetTransactionsByAccount(ctx context.Context, accountID uuid.UUID, offset, limit int) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Order("transaction_date DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

// UpdateTransaction updates a transaction
func (r *repository) UpdateTransaction(ctx context.Context, transaction *Transaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

// DeleteTransaction deletes a transaction
func (r *repository) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Transaction{}, id).Error
}

// Category operations

// CreateCategory creates a new category
func (r *repository) CreateCategory(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// GetCategoryByID retrieves a category by ID
func (r *repository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &category, nil
}

// GetCategories retrieves all categories with pagination
func (r *repository) GetCategories(ctx context.Context, offset, limit int) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).
		Order("sort_order ASC, name ASC").
		Offset(offset).
		Limit(limit).
		Find(&categories).Error
	return categories, err
}

// GetDefaultCategories retrieves default categories
func (r *repository) GetDefaultCategories(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := r.db.WithContext(ctx).
		Where("is_default = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}

// UpdateCategory updates a category
func (r *repository) UpdateCategory(ctx context.Context, category *Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// DeleteCategory deletes a category
func (r *repository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Category{}, id).Error
}

// Account operations

// CreateAccount creates a new account
func (r *repository) CreateAccount(ctx context.Context, account *Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

// GetAccountByID retrieves an account by ID
func (r *repository) GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
	var account Account
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

// GetAccountsByUser retrieves all accounts for a user
func (r *repository) GetAccountsByUser(ctx context.Context, userID uuid.UUID) ([]Account, error) {
	var accounts []Account
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&accounts).Error
	return accounts, err
}

// UpdateAccount updates an account
func (r *repository) UpdateAccount(ctx context.Context, account *Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

// DeleteAccount deletes an account
func (r *repository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Account{}, id).Error
}

// Custom errors
var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrAccountNotFound     = errors.New("account not found")
)
