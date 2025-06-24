package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Service defines the interface for transaction business logic
type Service interface {
	// Transaction operations
	CreateTransaction(ctx context.Context, userID uuid.UUID, req *CreateTransactionRequest) (*TransactionResponse, error)
	GetTransaction(ctx context.Context, userID, transactionID uuid.UUID) (*TransactionResponse, error)
	GetTransactions(ctx context.Context, userID uuid.UUID, offset, limit int) ([]TransactionResponse, error)
	UpdateTransaction(ctx context.Context, userID, transactionID uuid.UUID, req *UpdateTransactionRequest) (*TransactionResponse, error)
	DeleteTransaction(ctx context.Context, userID, transactionID uuid.UUID) error

	// Category operations
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*Category, error)
	GetCategory(ctx context.Context, categoryID uuid.UUID) (*Category, error)
	GetCategories(ctx context.Context, offset, limit int) ([]Category, error)
	GetDefaultCategories(ctx context.Context) ([]Category, error)
	UpdateCategory(ctx context.Context, categoryID uuid.UUID, req *CreateCategoryRequest) (*Category, error)
	DeleteCategory(ctx context.Context, categoryID uuid.UUID) error

	// Account operations
	CreateAccount(ctx context.Context, userID uuid.UUID, req *CreateAccountRequest) (*Account, error)
	GetAccount(ctx context.Context, userID, accountID uuid.UUID) (*Account, error)
	GetAccounts(ctx context.Context, userID uuid.UUID) ([]Account, error)
	UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, req *CreateAccountRequest) (*Account, error)
	DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new transaction service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Transaction operations

// CreateTransaction creates a new transaction
func (s *service) CreateTransaction(ctx context.Context, userID uuid.UUID, req *CreateTransactionRequest) (*TransactionResponse, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "CreateTransaction",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("account_id", req.AccountID.String()),
			attribute.Float64("amount", req.Amount),
		),
	)
	defer span.End()

	// Validate amount
	if req.Amount == 0 {
		span.RecordError(errors.New("amount cannot be zero"))
		span.SetStatus(codes.Error, "amount cannot be zero")
		return nil, errors.New("amount cannot be zero")
	}

	// Validate account exists and belongs to user
	account, err := s.repo.GetAccountByID(ctx, req.AccountID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get account")
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if account.UserID != userID {
		span.RecordError(errors.New("account does not belong to user"))
		span.SetStatus(codes.Error, "account does not belong to user")
		return nil, errors.New("account does not belong to user")
	}

	// Validate category if provided
	if req.CategoryID != nil {
		_, err := s.repo.GetCategoryByID(ctx, *req.CategoryID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to get category")
			return nil, fmt.Errorf("failed to get category: %w", err)
		}
	}

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "USD"
	}

	// Create transaction
	transaction := &Transaction{
		UserID:          userID,
		AccountID:       req.AccountID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Description:     req.Description,
		Merchant:        req.Merchant,
		Location:        req.Location,
		TransactionDate: req.TransactionDate,
		PostedDate:      req.PostedDate,
		Status:          TransactionStatusPending,
		Tags:            req.Tags,
		Notes:           req.Notes,
	}

	if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create transaction")
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	span.SetStatus(codes.Ok, "transaction created successfully")
	return s.toTransactionResponse(transaction), nil
}

// GetTransaction retrieves a transaction by ID
func (s *service) GetTransaction(ctx context.Context, userID, transactionID uuid.UUID) (*TransactionResponse, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "GetTransaction",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("transaction_id", transactionID.String()),
		),
	)
	defer span.End()

	transaction, err := s.repo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get transaction")
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Check if transaction belongs to user
	if transaction.UserID != userID {
		span.RecordError(errors.New("transaction does not belong to user"))
		span.SetStatus(codes.Error, "transaction does not belong to user")
		return nil, errors.New("transaction does not belong to user")
	}

	span.SetStatus(codes.Ok, "transaction retrieved successfully")
	return s.toTransactionResponse(transaction), nil
}

// GetTransactions retrieves transactions for a user with pagination
func (s *service) GetTransactions(ctx context.Context, userID uuid.UUID, offset, limit int) ([]TransactionResponse, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "GetTransactions",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.Int("offset", offset),
			attribute.Int("limit", limit),
		),
	)
	defer span.End()

	// Set default limit if not provided
	if limit <= 0 {
		limit = 50
	}

	transactions, err := s.repo.GetTransactionsByUser(ctx, userID, offset, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get transactions")
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	responses := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = *s.toTransactionResponse(&transaction)
	}

	span.SetStatus(codes.Ok, "transactions retrieved successfully")
	return responses, nil
}

// UpdateTransaction updates a transaction
func (s *service) UpdateTransaction(ctx context.Context, userID, transactionID uuid.UUID, req *UpdateTransactionRequest) (*TransactionResponse, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "UpdateTransaction",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("transaction_id", transactionID.String()),
		),
	)
	defer span.End()

	transaction, err := s.repo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get transaction")
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Check if transaction belongs to user
	if transaction.UserID != userID {
		span.RecordError(errors.New("transaction does not belong to user"))
		span.SetStatus(codes.Error, "transaction does not belong to user")
		return nil, errors.New("transaction does not belong to user")
	}

	// Update fields if provided
	if req.CategoryID != nil {
		// Validate category exists
		_, err := s.repo.GetCategoryByID(ctx, *req.CategoryID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to get category")
			return nil, fmt.Errorf("failed to get category: %w", err)
		}
		transaction.CategoryID = req.CategoryID
	}

	if req.Amount != nil {
		if *req.Amount == 0 {
			span.RecordError(errors.New("amount cannot be zero"))
			span.SetStatus(codes.Error, "amount cannot be zero")
			return nil, errors.New("amount cannot be zero")
		}
		transaction.Amount = *req.Amount
	}

	if req.Currency != "" {
		transaction.Currency = req.Currency
	}

	if req.Description != "" {
		transaction.Description = req.Description
	}

	if req.Merchant != "" {
		transaction.Merchant = req.Merchant
	}

	if req.Location != "" {
		transaction.Location = req.Location
	}

	if req.TransactionDate != nil {
		transaction.TransactionDate = *req.TransactionDate
	}

	if req.PostedDate != nil {
		transaction.PostedDate = req.PostedDate
	}

	if req.Status != nil {
		transaction.Status = *req.Status
	}

	if req.Tags != nil {
		transaction.Tags = req.Tags
	}

	if req.Notes != "" {
		transaction.Notes = req.Notes
	}

	transaction.UpdatedAt = time.Now()

	if err := s.repo.UpdateTransaction(ctx, transaction); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update transaction")
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	span.SetStatus(codes.Ok, "transaction updated successfully")
	return s.toTransactionResponse(transaction), nil
}

// DeleteTransaction deletes a transaction
func (s *service) DeleteTransaction(ctx context.Context, userID, transactionID uuid.UUID) error {
	ctx, span := otel.Tracer("transaction").Start(ctx, "DeleteTransaction",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("transaction_id", transactionID.String()),
		),
	)
	defer span.End()

	transaction, err := s.repo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get transaction")
		return fmt.Errorf("failed to get transaction: %w", err)
	}

	// Check if transaction belongs to user
	if transaction.UserID != userID {
		span.RecordError(errors.New("transaction does not belong to user"))
		span.SetStatus(codes.Error, "transaction does not belong to user")
		return errors.New("transaction does not belong to user")
	}

	if err := s.repo.DeleteTransaction(ctx, transactionID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete transaction")
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	span.SetStatus(codes.Ok, "transaction deleted successfully")
	return nil
}

// Category operations

// CreateCategory creates a new category
func (s *service) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*Category, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "CreateCategory",
		trace.WithAttributes(
			attribute.String("name", req.Name),
		),
	)
	defer span.End()

	category := &Category{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		ParentID:    req.ParentID,
		IsDefault:   req.IsDefault,
		SortOrder:   req.SortOrder,
	}

	if err := s.repo.CreateCategory(ctx, category); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create category")
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	span.SetStatus(codes.Ok, "category created successfully")
	return category, nil
}

// GetCategory retrieves a category by ID
func (s *service) GetCategory(ctx context.Context, categoryID uuid.UUID) (*Category, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "GetCategory",
		trace.WithAttributes(
			attribute.String("category_id", categoryID.String()),
		),
	)
	defer span.End()

	category, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get category")
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	span.SetStatus(codes.Ok, "category retrieved successfully")
	return category, nil
}

// GetCategories retrieves all categories with pagination
func (s *service) GetCategories(ctx context.Context, offset, limit int) ([]Category, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "GetCategories",
		trace.WithAttributes(
			attribute.Int("offset", offset),
			attribute.Int("limit", limit),
		),
	)
	defer span.End()

	if limit <= 0 {
		limit = 100
	}

	categories, err := s.repo.GetCategories(ctx, offset, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get categories")
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	span.SetStatus(codes.Ok, "categories retrieved successfully")
	return categories, nil
}

// GetDefaultCategories retrieves default categories
func (s *service) GetDefaultCategories(ctx context.Context) ([]Category, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "GetDefaultCategories")
	defer span.End()

	categories, err := s.repo.GetDefaultCategories(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get default categories")
		return nil, fmt.Errorf("failed to get default categories: %w", err)
	}

	span.SetStatus(codes.Ok, "default categories retrieved successfully")
	return categories, nil
}

// UpdateCategory updates a category
func (s *service) UpdateCategory(ctx context.Context, categoryID uuid.UUID, req *CreateCategoryRequest) (*Category, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "UpdateCategory",
		trace.WithAttributes(
			attribute.String("category_id", categoryID.String()),
		),
	)
	defer span.End()

	category, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get category")
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Update fields
	category.Name = req.Name
	category.Description = req.Description
	category.Icon = req.Icon
	category.Color = req.Color
	category.ParentID = req.ParentID
	category.IsDefault = req.IsDefault
	category.SortOrder = req.SortOrder
	category.UpdatedAt = time.Now()

	if err := s.repo.UpdateCategory(ctx, category); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update category")
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	span.SetStatus(codes.Ok, "category updated successfully")
	return category, nil
}

// DeleteCategory deletes a category
func (s *service) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	ctx, span := otel.Tracer("transaction").Start(ctx, "DeleteCategory",
		trace.WithAttributes(
			attribute.String("category_id", categoryID.String()),
		),
	)
	defer span.End()

	if err := s.repo.DeleteCategory(ctx, categoryID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete category")
		return fmt.Errorf("failed to delete category: %w", err)
	}

	span.SetStatus(codes.Ok, "category deleted successfully")
	return nil
}

// Account operations

// CreateAccount creates a new account
func (s *service) CreateAccount(ctx context.Context, userID uuid.UUID, req *CreateAccountRequest) (*Account, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "CreateAccount",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("name", req.Name),
			attribute.String("type", string(req.Type)),
		),
	)
	defer span.End()

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "USD"
	}

	account := &Account{
		UserID:            userID,
		Name:              req.Name,
		Type:              req.Type,
		Institution:       req.Institution,
		AccountNumberHash: req.AccountNumberHash,
		Balance:           req.Balance,
		Currency:          req.Currency,
		PlaidAccountID:    req.PlaidAccountID,
	}

	if err := s.repo.CreateAccount(ctx, account); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create account")
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	span.SetStatus(codes.Ok, "account created successfully")
	return account, nil
}

// GetAccount retrieves an account by ID
func (s *service) GetAccount(ctx context.Context, userID, accountID uuid.UUID) (*Account, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "GetAccount",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("account_id", accountID.String()),
		),
	)
	defer span.End()

	account, err := s.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get account")
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Check if account belongs to user
	if account.UserID != userID {
		span.RecordError(errors.New("account does not belong to user"))
		span.SetStatus(codes.Error, "account does not belong to user")
		return nil, errors.New("account does not belong to user")
	}

	span.SetStatus(codes.Ok, "account retrieved successfully")
	return account, nil
}

// GetAccounts retrieves all accounts for a user
func (s *service) GetAccounts(ctx context.Context, userID uuid.UUID) ([]Account, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "GetAccounts",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
		),
	)
	defer span.End()

	accounts, err := s.repo.GetAccountsByUser(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get accounts")
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	span.SetStatus(codes.Ok, "accounts retrieved successfully")
	return accounts, nil
}

// UpdateAccount updates an account
func (s *service) UpdateAccount(ctx context.Context, userID, accountID uuid.UUID, req *CreateAccountRequest) (*Account, error) {
	ctx, span := otel.Tracer("transaction").Start(ctx, "UpdateAccount",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("account_id", accountID.String()),
		),
	)
	defer span.End()

	account, err := s.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get account")
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Check if account belongs to user
	if account.UserID != userID {
		span.RecordError(errors.New("account does not belong to user"))
		span.SetStatus(codes.Error, "account does not belong to user")
		return nil, errors.New("account does not belong to user")
	}

	// Update fields
	account.Name = req.Name
	account.Type = req.Type
	account.Institution = req.Institution
	account.AccountNumberHash = req.AccountNumberHash
	account.Balance = req.Balance
	account.Currency = req.Currency
	account.PlaidAccountID = req.PlaidAccountID
	account.UpdatedAt = time.Now()

	if err := s.repo.UpdateAccount(ctx, account); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update account")
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	span.SetStatus(codes.Ok, "account updated successfully")
	return account, nil
}

// DeleteAccount deletes an account
func (s *service) DeleteAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	ctx, span := otel.Tracer("transaction").Start(ctx, "DeleteAccount",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("account_id", accountID.String()),
		),
	)
	defer span.End()

	account, err := s.repo.GetAccountByID(ctx, accountID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get account")
		return fmt.Errorf("failed to get account: %w", err)
	}

	// Check if account belongs to user
	if account.UserID != userID {
		span.RecordError(errors.New("account does not belong to user"))
		span.SetStatus(codes.Error, "account does not belong to user")
		return errors.New("account does not belong to user")
	}

	if err := s.repo.DeleteAccount(ctx, accountID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete account")
		return fmt.Errorf("failed to delete account: %w", err)
	}

	span.SetStatus(codes.Ok, "account deleted successfully")
	return nil
}

// Helper methods

// toTransactionResponse converts a Transaction to TransactionResponse
func (s *service) toTransactionResponse(transaction *Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:                       transaction.ID,
		UserID:                   transaction.UserID,
		FamilyID:                 transaction.FamilyID,
		AccountID:                transaction.AccountID,
		CategoryID:               transaction.CategoryID,
		Amount:                   transaction.Amount,
		Currency:                 transaction.Currency,
		Description:              transaction.Description,
		Merchant:                 transaction.Merchant,
		Location:                 transaction.Location,
		TransactionDate:          transaction.TransactionDate,
		PostedDate:               transaction.PostedDate,
		Status:                   transaction.Status,
		CategorizationSource:     transaction.CategorizationSource,
		CategorizationConfidence: transaction.CategorizationConfidence,
		Tags:                     transaction.Tags,
		Notes:                    transaction.Notes,
		ReceiptURL:               transaction.ReceiptURL,
		CreatedAt:                transaction.CreatedAt,
		UpdatedAt:                transaction.UpdatedAt,
	}
}
