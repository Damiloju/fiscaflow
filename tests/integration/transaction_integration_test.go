package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"fiscaflow/internal/api/handlers"
	"fiscaflow/internal/api/middleware"
	"fiscaflow/internal/domain/transaction"
	"fiscaflow/internal/domain/user"
)

// TestTransactionRepository is a test implementation of the transaction.Repository interface
type TestTransactionRepository struct {
	db *gorm.DB
}

func NewTestTransactionRepository(db *gorm.DB) transaction.Repository {
	return &TestTransactionRepository{db: db}
}

func (r *TestTransactionRepository) CreateTransaction(ctx context.Context, t *transaction.Transaction) error {
	testTransaction := &TestTransaction{
		ID:                       t.ID.String(),
		UserID:                   t.UserID.String(),
		AccountID:                t.AccountID.String(),
		Amount:                   t.Amount,
		Currency:                 t.Currency,
		Description:              t.Description,
		Merchant:                 t.Merchant,
		Location:                 t.Location,
		TransactionDate:          t.TransactionDate,
		PostedDate:               t.PostedDate,
		Status:                   string(t.Status),
		CategorizationSource:     string(t.CategorizationSource),
		CategorizationConfidence: t.CategorizationConfidence,
		Tags:                     r.tagsToString(t.Tags),
		Notes:                    t.Notes,
		ReceiptURL:               t.ReceiptURL,
		CreatedAt:                t.CreatedAt,
		UpdatedAt:                t.UpdatedAt,
	}

	if t.FamilyID != nil {
		familyID := t.FamilyID.String()
		testTransaction.FamilyID = &familyID
	}

	if t.CategoryID != nil {
		categoryID := t.CategoryID.String()
		testTransaction.CategoryID = &categoryID
	}

	return r.db.WithContext(ctx).Create(testTransaction).Error
}

func (r *TestTransactionRepository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error) {
	var testTransaction TestTransaction
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&testTransaction).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, transaction.ErrTransactionNotFound
		}
		return nil, err
	}

	return r.testTransactionToTransaction(&testTransaction), nil
}

func (r *TestTransactionRepository) GetTransactionsByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]transaction.Transaction, error) {
	var testTransactions []TestTransaction
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID.String()).
		Order("transaction_date DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&testTransactions).Error
	if err != nil {
		return nil, err
	}

	transactions := make([]transaction.Transaction, len(testTransactions))
	for i, tt := range testTransactions {
		transactions[i] = *r.testTransactionToTransaction(&tt)
	}

	return transactions, nil
}

func (r *TestTransactionRepository) GetTransactionsByAccount(ctx context.Context, accountID uuid.UUID, offset, limit int) ([]transaction.Transaction, error) {
	var testTransactions []TestTransaction
	err := r.db.WithContext(ctx).
		Where("account_id = ?", accountID.String()).
		Order("transaction_date DESC, created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&testTransactions).Error
	if err != nil {
		return nil, err
	}

	transactions := make([]transaction.Transaction, len(testTransactions))
	for i, tt := range testTransactions {
		transactions[i] = *r.testTransactionToTransaction(&tt)
	}

	return transactions, nil
}

func (r *TestTransactionRepository) UpdateTransaction(ctx context.Context, t *transaction.Transaction) error {
	testTransaction := &TestTransaction{
		ID:                       t.ID.String(),
		UserID:                   t.UserID.String(),
		AccountID:                t.AccountID.String(),
		Amount:                   t.Amount,
		Currency:                 t.Currency,
		Description:              t.Description,
		Merchant:                 t.Merchant,
		Location:                 t.Location,
		TransactionDate:          t.TransactionDate,
		PostedDate:               t.PostedDate,
		Status:                   string(t.Status),
		CategorizationSource:     string(t.CategorizationSource),
		CategorizationConfidence: t.CategorizationConfidence,
		Tags:                     r.tagsToString(t.Tags),
		Notes:                    t.Notes,
		ReceiptURL:               t.ReceiptURL,
		CreatedAt:                t.CreatedAt,
		UpdatedAt:                t.UpdatedAt,
	}

	if t.FamilyID != nil {
		familyID := t.FamilyID.String()
		testTransaction.FamilyID = &familyID
	}

	if t.CategoryID != nil {
		categoryID := t.CategoryID.String()
		testTransaction.CategoryID = &categoryID
	}

	return r.db.WithContext(ctx).Save(testTransaction).Error
}

func (r *TestTransactionRepository) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&TestTransaction{}, "id = ?", id.String()).Error
}

func (r *TestTransactionRepository) CreateCategory(ctx context.Context, c *transaction.Category) error {
	testCategory := &TestCategory{
		ID:          c.ID.String(),
		Name:        c.Name,
		Description: c.Description,
		Icon:        c.Icon,
		Color:       c.Color,
		IsDefault:   c.IsDefault,
		IsActive:    c.IsActive,
		SortOrder:   c.SortOrder,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}

	if c.ParentID != nil {
		parentID := c.ParentID.String()
		testCategory.ParentID = &parentID
	}

	return r.db.WithContext(ctx).Create(testCategory).Error
}

func (r *TestTransactionRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*transaction.Category, error) {
	var testCategory TestCategory
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&testCategory).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, transaction.ErrCategoryNotFound
		}
		return nil, err
	}

	return r.testCategoryToCategory(&testCategory), nil
}

func (r *TestTransactionRepository) GetCategories(ctx context.Context, offset, limit int) ([]transaction.Category, error) {
	var testCategories []TestCategory
	err := r.db.WithContext(ctx).
		Order("sort_order ASC, name ASC").
		Offset(offset).
		Limit(limit).
		Find(&testCategories).Error
	if err != nil {
		return nil, err
	}

	categories := make([]transaction.Category, len(testCategories))
	for i, tc := range testCategories {
		categories[i] = *r.testCategoryToCategory(&tc)
	}

	return categories, nil
}

func (r *TestTransactionRepository) GetDefaultCategories(ctx context.Context) ([]transaction.Category, error) {
	var testCategories []TestCategory
	err := r.db.WithContext(ctx).
		Where("is_default = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&testCategories).Error
	if err != nil {
		return nil, err
	}

	categories := make([]transaction.Category, len(testCategories))
	for i, tc := range testCategories {
		categories[i] = *r.testCategoryToCategory(&tc)
	}

	return categories, nil
}

func (r *TestTransactionRepository) UpdateCategory(ctx context.Context, c *transaction.Category) error {
	testCategory := &TestCategory{
		ID:          c.ID.String(),
		Name:        c.Name,
		Description: c.Description,
		Icon:        c.Icon,
		Color:       c.Color,
		IsDefault:   c.IsDefault,
		IsActive:    c.IsActive,
		SortOrder:   c.SortOrder,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}

	if c.ParentID != nil {
		parentID := c.ParentID.String()
		testCategory.ParentID = &parentID
	}

	return r.db.WithContext(ctx).Save(testCategory).Error
}

func (r *TestTransactionRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&TestCategory{}, "id = ?", id.String()).Error
}

func (r *TestTransactionRepository) CreateAccount(ctx context.Context, a *transaction.Account) error {
	testAccount := &TestAccount{
		ID:                a.ID.String(),
		UserID:            a.UserID.String(),
		Name:              a.Name,
		Type:              string(a.Type),
		Institution:       a.Institution,
		AccountNumberHash: a.AccountNumberHash,
		Balance:           a.Balance,
		Currency:          a.Currency,
		IsActive:          a.IsActive,
		PlaidAccountID:    a.PlaidAccountID,
		LastSyncAt:        a.LastSyncAt,
		Settings:          a.Settings,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}

	if a.FamilyID != nil {
		familyID := a.FamilyID.String()
		testAccount.FamilyID = &familyID
	}

	return r.db.WithContext(ctx).Create(testAccount).Error
}

func (r *TestTransactionRepository) GetAccountByID(ctx context.Context, id uuid.UUID) (*transaction.Account, error) {
	var testAccount TestAccount
	err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&testAccount).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, transaction.ErrAccountNotFound
		}
		return nil, err
	}

	return r.testAccountToAccount(&testAccount), nil
}

func (r *TestTransactionRepository) GetAccountsByUser(ctx context.Context, userID uuid.UUID) ([]transaction.Account, error) {
	var testAccounts []TestAccount
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID.String()).
		Order("name ASC").
		Find(&testAccounts).Error
	if err != nil {
		return nil, err
	}

	accounts := make([]transaction.Account, len(testAccounts))
	for i, ta := range testAccounts {
		accounts[i] = *r.testAccountToAccount(&ta)
	}

	return accounts, nil
}

func (r *TestTransactionRepository) UpdateAccount(ctx context.Context, a *transaction.Account) error {
	testAccount := &TestAccount{
		ID:                a.ID.String(),
		UserID:            a.UserID.String(),
		Name:              a.Name,
		Type:              string(a.Type),
		Institution:       a.Institution,
		AccountNumberHash: a.AccountNumberHash,
		Balance:           a.Balance,
		Currency:          a.Currency,
		IsActive:          a.IsActive,
		PlaidAccountID:    a.PlaidAccountID,
		LastSyncAt:        a.LastSyncAt,
		Settings:          a.Settings,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}

	if a.FamilyID != nil {
		familyID := a.FamilyID.String()
		testAccount.FamilyID = &familyID
	}

	return r.db.WithContext(ctx).Save(testAccount).Error
}

func (r *TestTransactionRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&TestAccount{}, "id = ?", id.String()).Error
}

// Helper methods for converting between test models and domain models

func (r *TestTransactionRepository) testTransactionToTransaction(tt *TestTransaction) *transaction.Transaction {
	transactionID, _ := uuid.Parse(tt.ID)
	userID, _ := uuid.Parse(tt.UserID)
	accountID, _ := uuid.Parse(tt.AccountID)

	t := &transaction.Transaction{
		ID:                       transactionID,
		UserID:                   userID,
		AccountID:                accountID,
		Amount:                   tt.Amount,
		Currency:                 tt.Currency,
		Description:              tt.Description,
		Merchant:                 tt.Merchant,
		Location:                 tt.Location,
		TransactionDate:          tt.TransactionDate,
		PostedDate:               tt.PostedDate,
		Status:                   transaction.TransactionStatus(tt.Status),
		CategorizationSource:     transaction.CategorizationSource(tt.CategorizationSource),
		CategorizationConfidence: tt.CategorizationConfidence,
		Tags:                     r.stringToTags(tt.Tags),
		Notes:                    tt.Notes,
		ReceiptURL:               tt.ReceiptURL,
		CreatedAt:                tt.CreatedAt,
		UpdatedAt:                tt.UpdatedAt,
	}

	if tt.FamilyID != nil {
		familyID, _ := uuid.Parse(*tt.FamilyID)
		t.FamilyID = &familyID
	}

	if tt.CategoryID != nil {
		categoryID, _ := uuid.Parse(*tt.CategoryID)
		t.CategoryID = &categoryID
	}

	return t
}

func (r *TestTransactionRepository) testCategoryToCategory(tc *TestCategory) *transaction.Category {
	categoryID, _ := uuid.Parse(tc.ID)

	c := &transaction.Category{
		ID:          categoryID,
		Name:        tc.Name,
		Description: tc.Description,
		Icon:        tc.Icon,
		Color:       tc.Color,
		IsDefault:   tc.IsDefault,
		IsActive:    tc.IsActive,
		SortOrder:   tc.SortOrder,
		CreatedAt:   tc.CreatedAt,
		UpdatedAt:   tc.UpdatedAt,
	}

	if tc.ParentID != nil {
		parentID, _ := uuid.Parse(*tc.ParentID)
		c.ParentID = &parentID
	}

	return c
}

func (r *TestTransactionRepository) testAccountToAccount(ta *TestAccount) *transaction.Account {
	accountID, _ := uuid.Parse(ta.ID)
	userID, _ := uuid.Parse(ta.UserID)

	a := &transaction.Account{
		ID:                accountID,
		UserID:            userID,
		Name:              ta.Name,
		Type:              transaction.AccountType(ta.Type),
		Institution:       ta.Institution,
		AccountNumberHash: ta.AccountNumberHash,
		Balance:           ta.Balance,
		Currency:          ta.Currency,
		IsActive:          ta.IsActive,
		PlaidAccountID:    ta.PlaidAccountID,
		LastSyncAt:        ta.LastSyncAt,
		Settings:          ta.Settings,
		CreatedAt:         ta.CreatedAt,
		UpdatedAt:         ta.UpdatedAt,
	}

	if ta.FamilyID != nil {
		familyID, _ := uuid.Parse(*ta.FamilyID)
		a.FamilyID = &familyID
	}

	return a
}

func (r *TestTransactionRepository) tagsToString(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	data, _ := json.Marshal(tags)
	return string(data)
}

func (r *TestTransactionRepository) stringToTags(tagsStr string) []string {
	if tagsStr == "" || tagsStr == "[]" {
		return []string{}
	}
	var tags []string
	json.Unmarshal([]byte(tagsStr), &tags)
	return tags
}

// Integration test functions

func TestTransactionIntegration_CreateAndRetrieve(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repositories
	userRepo := NewTestRepository(db.DB)
	transactionRepo := NewTestTransactionRepository(db.DB)

	// Create a test user
	testUser := &user.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Test",
		LastName:     "User",
		Role:         user.UserRoleUser,
		Status:       user.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create a test account
	testAccount := &transaction.Account{
		ID:        uuid.New(),
		UserID:    testUser.ID,
		Name:      "Test Checking Account",
		Type:      transaction.AccountTypeChecking,
		Balance:   1000.00,
		Currency:  "USD",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = transactionRepo.CreateAccount(context.Background(), testAccount)
	require.NoError(t, err)

	// Create a test transaction
	testTransaction := &transaction.Transaction{
		ID:                   uuid.New(),
		UserID:               testUser.ID,
		AccountID:            testAccount.ID,
		Amount:               50.00,
		Currency:             "USD",
		Description:          "Grocery shopping",
		Merchant:             "Walmart",
		TransactionDate:      time.Now(),
		Status:               transaction.TransactionStatusPosted,
		CategorizationSource: transaction.CategorizationSourceManual,
		Tags:                 []string{"groceries", "food"},
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Test creating transaction
	err = transactionRepo.CreateTransaction(context.Background(), testTransaction)
	require.NoError(t, err)

	// Test retrieving transaction by ID
	retrievedTransaction, err := transactionRepo.GetTransactionByID(context.Background(), testTransaction.ID)
	require.NoError(t, err)
	assert.Equal(t, testTransaction.ID, retrievedTransaction.ID)
	assert.Equal(t, testTransaction.UserID, retrievedTransaction.UserID)
	assert.Equal(t, testTransaction.AccountID, retrievedTransaction.AccountID)
	assert.Equal(t, testTransaction.Amount, retrievedTransaction.Amount)
	assert.Equal(t, testTransaction.Description, retrievedTransaction.Description)
	assert.Equal(t, testTransaction.Merchant, retrievedTransaction.Merchant)
	assert.Equal(t, testTransaction.Status, retrievedTransaction.Status)
	assert.Equal(t, testTransaction.Tags, retrievedTransaction.Tags)

	// Test retrieving transactions by user
	userTransactions, err := transactionRepo.GetTransactionsByUser(context.Background(), testUser.ID, 0, 10)
	require.NoError(t, err)
	assert.Len(t, userTransactions, 1)
	assert.Equal(t, testTransaction.ID, userTransactions[0].ID)

	// Test retrieving transactions by account
	accountTransactions, err := transactionRepo.GetTransactionsByAccount(context.Background(), testAccount.ID, 0, 10)
	require.NoError(t, err)
	assert.Len(t, accountTransactions, 1)
	assert.Equal(t, testTransaction.ID, accountTransactions[0].ID)
}

func TestTransactionIntegration_UpdateAndDelete(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repositories
	userRepo := NewTestRepository(db.DB)
	transactionRepo := NewTestTransactionRepository(db.DB)

	// Create a test user
	testUser := &user.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Test",
		LastName:     "User",
		Role:         user.UserRoleUser,
		Status:       user.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create a test account
	testAccount := &transaction.Account{
		ID:        uuid.New(),
		UserID:    testUser.ID,
		Name:      "Test Checking Account",
		Type:      transaction.AccountTypeChecking,
		Balance:   1000.00,
		Currency:  "USD",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = transactionRepo.CreateAccount(context.Background(), testAccount)
	require.NoError(t, err)

	// Create a test transaction
	testTransaction := &transaction.Transaction{
		ID:                   uuid.New(),
		UserID:               testUser.ID,
		AccountID:            testAccount.ID,
		Amount:               50.00,
		Currency:             "USD",
		Description:          "Grocery shopping",
		Merchant:             "Walmart",
		TransactionDate:      time.Now(),
		Status:               transaction.TransactionStatusPosted,
		CategorizationSource: transaction.CategorizationSourceManual,
		Tags:                 []string{"groceries"},
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	err = transactionRepo.CreateTransaction(context.Background(), testTransaction)
	require.NoError(t, err)

	// Test updating transaction
	testTransaction.Description = "Updated grocery shopping"
	testTransaction.Amount = 75.00
	testTransaction.Tags = []string{"groceries", "updated"}
	testTransaction.UpdatedAt = time.Now()

	err = transactionRepo.UpdateTransaction(context.Background(), testTransaction)
	require.NoError(t, err)

	// Verify update
	retrievedTransaction, err := transactionRepo.GetTransactionByID(context.Background(), testTransaction.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated grocery shopping", retrievedTransaction.Description)
	assert.Equal(t, 75.00, retrievedTransaction.Amount)
	assert.Equal(t, []string{"groceries", "updated"}, retrievedTransaction.Tags)

	// Test deleting transaction
	err = transactionRepo.DeleteTransaction(context.Background(), testTransaction.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = transactionRepo.GetTransactionByID(context.Background(), testTransaction.ID)
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrTransactionNotFound, err)
}

func TestTransactionIntegration_CategoryOperations(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repository
	transactionRepo := NewTestTransactionRepository(db.DB)

	// Create a test category
	testCategory := &transaction.Category{
		ID:          uuid.New(),
		Name:        "Groceries",
		Description: "Food and household items",
		Icon:        "shopping-cart",
		Color:       "#4CAF50",
		IsDefault:   true,
		IsActive:    true,
		SortOrder:   1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test creating category
	err := transactionRepo.CreateCategory(context.Background(), testCategory)
	require.NoError(t, err)

	// Test retrieving category by ID
	retrievedCategory, err := transactionRepo.GetCategoryByID(context.Background(), testCategory.ID)
	require.NoError(t, err)
	assert.Equal(t, testCategory.ID, retrievedCategory.ID)
	assert.Equal(t, testCategory.Name, retrievedCategory.Name)
	assert.Equal(t, testCategory.Description, retrievedCategory.Description)
	assert.Equal(t, testCategory.IsDefault, retrievedCategory.IsDefault)

	// Test updating category
	testCategory.Description = "Updated description"
	testCategory.Color = "#FF5722"
	testCategory.UpdatedAt = time.Now()

	err = transactionRepo.UpdateCategory(context.Background(), testCategory)
	require.NoError(t, err)

	// Verify update
	retrievedCategory, err = transactionRepo.GetCategoryByID(context.Background(), testCategory.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", retrievedCategory.Description)
	assert.Equal(t, "#FF5722", retrievedCategory.Color)

	// Test retrieving all categories
	categories, err := transactionRepo.GetCategories(context.Background(), 0, 10)
	require.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, testCategory.ID, categories[0].ID)

	// Test retrieving default categories
	defaultCategories, err := transactionRepo.GetDefaultCategories(context.Background())
	require.NoError(t, err)
	assert.Len(t, defaultCategories, 1)
	assert.Equal(t, testCategory.ID, defaultCategories[0].ID)

	// Test deleting category
	err = transactionRepo.DeleteCategory(context.Background(), testCategory.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = transactionRepo.GetCategoryByID(context.Background(), testCategory.ID)
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrCategoryNotFound, err)
}

func TestTransactionIntegration_AccountOperations(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repositories
	userRepo := NewTestRepository(db.DB)
	transactionRepo := NewTestTransactionRepository(db.DB)

	// Create a test user
	testUser := &user.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		FirstName:    "Test",
		LastName:     "User",
		Role:         user.UserRoleUser,
		Status:       user.UserStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := userRepo.Create(context.Background(), testUser)
	require.NoError(t, err)

	// Create a test account
	testAccount := &transaction.Account{
		ID:          uuid.New(),
		UserID:      testUser.ID,
		Name:        "Test Checking Account",
		Type:        transaction.AccountTypeChecking,
		Institution: "Test Bank",
		Balance:     1000.00,
		Currency:    "USD",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test creating account
	err = transactionRepo.CreateAccount(context.Background(), testAccount)
	require.NoError(t, err)

	// Test retrieving account by ID
	retrievedAccount, err := transactionRepo.GetAccountByID(context.Background(), testAccount.ID)
	require.NoError(t, err)
	assert.Equal(t, testAccount.ID, retrievedAccount.ID)
	assert.Equal(t, testAccount.UserID, retrievedAccount.UserID)
	assert.Equal(t, testAccount.Name, retrievedAccount.Name)
	assert.Equal(t, testAccount.Type, retrievedAccount.Type)
	assert.Equal(t, testAccount.Balance, retrievedAccount.Balance)

	// Test updating account
	testAccount.Balance = 1500.00
	testAccount.Name = "Updated Checking Account"
	testAccount.UpdatedAt = time.Now()

	err = transactionRepo.UpdateAccount(context.Background(), testAccount)
	require.NoError(t, err)

	// Verify update
	retrievedAccount, err = transactionRepo.GetAccountByID(context.Background(), testAccount.ID)
	require.NoError(t, err)
	assert.Equal(t, 1500.00, retrievedAccount.Balance)
	assert.Equal(t, "Updated Checking Account", retrievedAccount.Name)

	// Test retrieving accounts by user
	userAccounts, err := transactionRepo.GetAccountsByUser(context.Background(), testUser.ID)
	require.NoError(t, err)
	assert.Len(t, userAccounts, 1)
	assert.Equal(t, testAccount.ID, userAccounts[0].ID)

	// Test deleting account
	err = transactionRepo.DeleteAccount(context.Background(), testAccount.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = transactionRepo.GetAccountByID(context.Background(), testAccount.ID)
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrAccountNotFound, err)
}

func TestTransactionIntegration_APIHandlers(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repositories
	userRepo := NewTestRepository(db.DB)
	transactionRepo := NewTestTransactionRepository(db.DB)

	// Create transaction service
	transactionService := transaction.NewService(transactionRepo)

	// Setup user service
	userService := user.NewService(userRepo, "test-secret")

	// Register the test user using the service
	password := "password123"
	registerResp, err := userService.Register(context.Background(), &user.CreateUserRequest{
		Email:     "test@example.com",
		Password:  password,
		FirstName: "Test",
		LastName:  "User",
	})
	require.NoError(t, err)
	testUserID := registerResp.ID

	// Login to get a real JWT
	loginResp, err := userService.Login(context.Background(), &user.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	})
	require.NoError(t, err)
	accessToken := loginResp.AccessToken

	// Create a test account
	testAccount := &transaction.Account{
		ID:        uuid.New(),
		UserID:    testUserID,
		Name:      "Test Checking Account",
		Type:      transaction.AccountTypeChecking,
		Balance:   1000.00,
		Currency:  "USD",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = transactionRepo.CreateAccount(context.Background(), testAccount)
	require.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create transaction handler
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Register routes with authentication middleware
	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(userService))
	transactionHandler.RegisterRoutes(api)

	// Test creating transaction via API
	createRequest := transaction.CreateTransactionRequest{
		AccountID:       testAccount.ID,
		Amount:          50.00,
		Currency:        "USD",
		Description:     "Grocery shopping",
		Merchant:        "Walmart",
		TransactionDate: time.Now(),
		Tags:            []string{"groceries", "food"},
	}

	requestBody, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/api/v1/transactions", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createResponse transaction.TransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &createResponse)
	require.NoError(t, err)
	assert.Equal(t, testAccount.ID, createResponse.AccountID)
	assert.Equal(t, 50.00, createResponse.Amount)
	assert.Equal(t, "Grocery shopping", createResponse.Description)
	assert.Equal(t, "Walmart", createResponse.Merchant)
	assert.Equal(t, []string{"groceries", "food"}, createResponse.Tags)

	// Test retrieving transaction via API
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/transactions/%s", createResponse.ID), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var getResponse transaction.TransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &getResponse)
	require.NoError(t, err)
	assert.Equal(t, createResponse.ID, getResponse.ID)
	assert.Equal(t, createResponse.Amount, getResponse.Amount)
	assert.Equal(t, createResponse.Description, getResponse.Description)

	// Test listing transactions via API
	req, _ = http.NewRequest("GET", "/api/v1/transactions?limit=10&offset=0", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var listResponse []transaction.TransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	require.NoError(t, err)
	assert.Len(t, listResponse, 1)
	assert.Equal(t, createResponse.ID, listResponse[0].ID)

	// Test updating transaction via API
	updateRequest := transaction.UpdateTransactionRequest{
		Description: "Updated grocery shopping",
		Amount:      &[]float64{75.00}[0],
		Tags:        []string{"groceries", "updated"},
	}

	requestBody, _ = json.Marshal(updateRequest)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/transactions/%s", createResponse.ID), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse transaction.TransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
	require.NoError(t, err)
	assert.Equal(t, "Updated grocery shopping", updateResponse.Description)
	assert.Equal(t, 75.00, updateResponse.Amount)
	assert.Equal(t, []string{"groceries", "updated"}, updateResponse.Tags)

	// Test deleting transaction via API
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/transactions/%s", createResponse.ID), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify transaction is deleted
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/transactions/%s", createResponse.ID), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestTransactionIntegration_ErrorHandling(t *testing.T) {
	// Setup test database
	db := NewTestDatabase(t)
	defer db.Cleanup()

	// Setup repository
	transactionRepo := NewTestTransactionRepository(db.DB)

	// Test getting non-existent transaction
	_, err := transactionRepo.GetTransactionByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrTransactionNotFound, err)

	// Test getting non-existent category
	_, err = transactionRepo.GetCategoryByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrCategoryNotFound, err)

	// Test getting non-existent account
	_, err = transactionRepo.GetAccountByID(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, transaction.ErrAccountNotFound, err)

	// Test deleting non-existent transaction
	err = transactionRepo.DeleteTransaction(context.Background(), uuid.New())
	assert.NoError(t, err) // GORM doesn't return error for deleting non-existent records

	// Test deleting non-existent category
	err = transactionRepo.DeleteCategory(context.Background(), uuid.New())
	assert.NoError(t, err)

	// Test deleting non-existent account
	err = transactionRepo.DeleteAccount(context.Background(), uuid.New())
	assert.NoError(t, err)
}
