package transaction

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
	userID uuid.UUID
}

// Implement Repository interface methods for mockRepository
func (m *mockRepository) CreateTransaction(ctx context.Context, t *Transaction) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *mockRepository) GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	args := m.Called(ctx, id)
	if tr, ok := args.Get(0).(*Transaction); ok {
		return tr, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockRepository) GetTransactionsByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]Transaction, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).([]Transaction), args.Error(1)
}
func (m *mockRepository) GetTransactionsByAccount(ctx context.Context, accountID uuid.UUID, offset, limit int) ([]Transaction, error) {
	return nil, nil
}
func (m *mockRepository) UpdateTransaction(ctx context.Context, t *Transaction) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *mockRepository) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *mockRepository) CreateCategory(ctx context.Context, c *Category) error { return nil }
func (m *mockRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	return &Category{}, nil
}
func (m *mockRepository) GetCategories(ctx context.Context, offset, limit int) ([]Category, error) {
	return nil, nil
}
func (m *mockRepository) GetDefaultCategories(ctx context.Context) ([]Category, error) {
	return nil, nil
}
func (m *mockRepository) UpdateCategory(ctx context.Context, c *Category) error  { return nil }
func (m *mockRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockRepository) CreateAccount(ctx context.Context, a *Account) error    { return nil }
func (m *mockRepository) GetAccountByID(ctx context.Context, id uuid.UUID) (*Account, error) {
	return &Account{ID: id, UserID: m.userID}, nil
}
func (m *mockRepository) GetAccountsByUser(ctx context.Context, userID uuid.UUID) ([]Account, error) {
	return nil, nil
}
func (m *mockRepository) UpdateAccount(ctx context.Context, a *Account) error   { return nil }
func (m *mockRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error { return nil }

func TestTransactionService_CreateGetUpdateDelete(t *testing.T) {
	repo := new(mockRepository)
	userID := uuid.New()
	repo.userID = userID
	svc := NewService(repo)
	ctx := context.Background()
	accountID := uuid.New()
	transactionID := uuid.New()

	// Mock GetAccountByID to return an account belonging to user
	repo.On("GetAccountByID", mock.Anything, accountID).Return(&Account{ID: accountID, UserID: userID}, nil)

	// Create
	createReq := &CreateTransactionRequest{
		AccountID:       accountID,
		Amount:          100.0,
		Currency:        "USD",
		Description:     "Test transaction",
		TransactionDate: time.Now(),
	}
	repo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*transaction.Transaction")).Return(nil)
	resp, err := svc.CreateTransaction(ctx, userID, createReq)
	assert.NoError(t, err)
	assert.Equal(t, createReq.Amount, resp.Amount)
	assert.Equal(t, createReq.Description, resp.Description)

	// Get
	tr := &Transaction{ID: transactionID, UserID: userID, AccountID: accountID, Amount: 100.0, Description: "Test transaction"}
	repo.On("GetTransactionByID", mock.Anything, transactionID).Return(tr, nil)
	getResp, err := svc.GetTransaction(ctx, userID, transactionID)
	assert.NoError(t, err)
	assert.Equal(t, transactionID, getResp.ID)

	// List
	transactions := []Transaction{*tr}
	repo.On("GetTransactionsByUser", mock.Anything, userID, 0, 10).Return(transactions, nil)
	listResp, err := svc.GetTransactions(ctx, userID, 0, 10)
	assert.NoError(t, err)
	assert.Len(t, listResp, 1)

	// Update
	updateReq := &UpdateTransactionRequest{Description: "Updated"}
	repo.On("UpdateTransaction", mock.Anything, mock.AnythingOfType("*transaction.Transaction")).Return(nil)
	updated, err := svc.UpdateTransaction(ctx, userID, transactionID, updateReq)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.Description)

	// Delete
	repo.On("DeleteTransaction", mock.Anything, transactionID).Return(nil)
	err = svc.DeleteTransaction(ctx, userID, transactionID)
	assert.NoError(t, err)
}

func TestTransactionService_CreateTransaction_AccountNotFound(t *testing.T) {
	repo := new(mockRepository)
	svc := NewService(repo)
	ctx := context.Background()
	userID := uuid.New()
	repo.userID = userID
	accountID := uuid.New()

	repo.On("GetAccountByID", mock.Anything, accountID).Return(nil, errors.New("not found"))
	repo.On("CreateTransaction", mock.Anything, mock.Anything).Return(errors.New("should not be called"))
	createReq := &CreateTransactionRequest{
		AccountID:       accountID,
		Amount:          100.0,
		Currency:        "USD",
		Description:     "Test transaction",
		TransactionDate: time.Now(),
	}
	resp, err := svc.CreateTransaction(ctx, userID, createReq)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
 