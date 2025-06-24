package budget

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, budget *Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*Budget, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Budget), args.Error(1)
}

func (m *MockRepository) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]Budget, error) {
	args := m.Called(ctx, userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Budget), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, budget *Budget) error {
	args := m.Called(ctx, budget)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) CreateCategory(ctx context.Context, budgetCategory *BudgetCategory) error {
	args := m.Called(ctx, budgetCategory)
	return args.Error(0)
}

func (m *MockRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*BudgetCategory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BudgetCategory), args.Error(1)
}

func (m *MockRepository) GetCategoriesByBudgetID(ctx context.Context, budgetID uuid.UUID) ([]BudgetCategory, error) {
	args := m.Called(ctx, budgetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]BudgetCategory), args.Error(1)
}

func (m *MockRepository) UpdateCategory(ctx context.Context, budgetCategory *BudgetCategory) error {
	args := m.Called(ctx, budgetCategory)
	return args.Error(0)
}

func (m *MockRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetBudgetSummary(ctx context.Context, budgetID uuid.UUID) (*BudgetSummary, error) {
	args := m.Called(ctx, budgetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BudgetSummary), args.Error(1)
}

func (m *MockRepository) UpdateSpentAmount(ctx context.Context, budgetID, categoryID uuid.UUID, amount float64) error {
	args := m.Called(ctx, budgetID, categoryID, amount)
	return args.Error(0)
}

func (m *MockRepository) GetActiveBudgetsByUser(ctx context.Context, userID uuid.UUID) ([]Budget, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Budget), args.Error(1)
}

func TestNewService(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	assert.NotNil(t, service)
}

func TestCreateBudget(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	req := &CreateBudgetRequest{
		Name:        "Monthly Budget",
		Description: "My monthly budget",
		PeriodType:  PeriodTypeMonthly,
		StartDate:   time.Now(),
		TotalAmount: 1000.0,
		Currency:    "USD",
	}

	expectedBudget := &Budget{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		PeriodType:  req.PeriodType,
		StartDate:   req.StartDate,
		TotalAmount: req.TotalAmount,
		Currency:    req.Currency,
		IsActive:    true,
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*budget.Budget")).Return(nil).Run(func(args mock.Arguments) {
		budget := args.Get(1).(*Budget)
		budget.ID = expectedBudget.ID
		budget.CreatedAt = time.Now()
		budget.UpdatedAt = time.Now()
	})

	result, err := service.CreateBudget(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.TotalAmount, result.TotalAmount)
	assert.Equal(t, userID, result.UserID)
	mockRepo.AssertExpectations(t)
}

func TestCreateBudget_ValidationError(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	// Test empty name
	req := &CreateBudgetRequest{
		Name:        "",
		Description: "My monthly budget",
		PeriodType:  PeriodTypeMonthly,
		StartDate:   time.Now(),
		TotalAmount: 1000.0,
	}

	result, err := service.CreateBudget(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "budget name is required")

	// Test negative amount
	req.Name = "Valid Budget"
	req.TotalAmount = -100

	result, err = service.CreateBudget(ctx, userID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "total amount must be positive")
}

func TestGetBudget(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()

	expectedBudget := &Budget{
		ID:          budgetID,
		UserID:      userID,
		Name:        "Monthly Budget",
		Description: "My monthly budget",
		PeriodType:  PeriodTypeMonthly,
		StartDate:   time.Now(),
		TotalAmount: 1000.0,
		Currency:    "USD",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, budgetID).Return(expectedBudget, nil)

	result, err := service.GetBudget(ctx, userID, budgetID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBudget.ID, result.ID)
	assert.Equal(t, expectedBudget.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetBudget_Unauthorized(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()
	otherUserID := uuid.New()

	expectedBudget := &Budget{
		ID:     budgetID,
		UserID: otherUserID, // Different user
		Name:   "Monthly Budget",
	}

	mockRepo.On("GetByID", mock.Anything, budgetID).Return(expectedBudget, nil)

	result, err := service.GetBudget(ctx, userID, budgetID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized access to budget")
	mockRepo.AssertExpectations(t)
}

func TestListBudgets(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	expectedBudgets := []Budget{
		{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        "Monthly Budget",
			TotalAmount: 1000.0,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        "Yearly Budget",
			TotalAmount: 12000.0,
			IsActive:    true,
		},
	}

	mockRepo.On("GetByUserID", mock.Anything, userID, 0, 20).Return(expectedBudgets, nil)

	result, err := service.ListBudgets(ctx, userID, 0, 20)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedBudgets[0].Name, result[0].Name)
	assert.Equal(t, expectedBudgets[1].Name, result[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBudget(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()

	existingBudget := &Budget{
		ID:          budgetID,
		UserID:      userID,
		Name:        "Old Budget Name",
		Description: "Old description",
		StartDate:   time.Now(),
		TotalAmount: 1000.0,
		IsActive:    true,
	}

	newName := "Updated Budget Name"
	newDescription := "Updated description"
	req := &UpdateBudgetRequest{
		Name:        &newName,
		Description: &newDescription,
	}

	mockRepo.On("GetByID", mock.Anything, budgetID).Return(existingBudget, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*budget.Budget")).Return(nil)

	result, err := service.UpdateBudget(ctx, userID, budgetID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newName, result.Name)
	assert.Equal(t, newDescription, result.Description)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBudget(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()

	existingBudget := &Budget{
		ID:     budgetID,
		UserID: userID,
		Name:   "Budget to Delete",
	}

	mockRepo.On("GetByID", mock.Anything, budgetID).Return(existingBudget, nil)
	mockRepo.On("Delete", mock.Anything, budgetID).Return(nil)

	err := service.DeleteBudget(ctx, userID, budgetID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddBudgetCategory(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()
	categoryID := uuid.New()

	existingBudget := &Budget{
		ID:     budgetID,
		UserID: userID,
		Name:   "Monthly Budget",
	}

	req := &CreateBudgetCategoryRequest{
		CategoryID:      categoryID,
		AllocatedAmount: 500.0,
		AlertThreshold:  0.8,
	}

	expectedBudgetCategory := &BudgetCategory{
		ID:              uuid.New(),
		BudgetID:        budgetID,
		CategoryID:      categoryID,
		AllocatedAmount: 500.0,
		AlertThreshold:  0.8,
		IsActive:        true,
	}

	mockRepo.On("GetByID", mock.Anything, budgetID).Return(existingBudget, nil)
	mockRepo.On("CreateCategory", mock.Anything, mock.AnythingOfType("*budget.BudgetCategory")).Return(nil).Run(func(args mock.Arguments) {
		category := args.Get(1).(*BudgetCategory)
		category.ID = expectedBudgetCategory.ID
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()
	})

	result, err := service.AddBudgetCategory(ctx, userID, budgetID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, categoryID, result.CategoryID)
	assert.Equal(t, 500.0, result.AllocatedAmount)
	mockRepo.AssertExpectations(t)
}

func TestGetBudgetSummary(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()

	existingBudget := &Budget{
		ID:     budgetID,
		UserID: userID,
		Name:   "Monthly Budget",
	}

	expectedSummary := &BudgetSummary{
		TotalAllocated:   1000.0,
		TotalSpent:       750.0,
		SpendingProgress: 0.75,
		Alerts:           []BudgetAlert{},
	}

	mockRepo.On("GetByID", mock.Anything, budgetID).Return(existingBudget, nil)
	mockRepo.On("GetBudgetSummary", mock.Anything, budgetID).Return(expectedSummary, nil)

	result, err := service.GetBudgetSummary(ctx, userID, budgetID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSummary.TotalAllocated, result.TotalAllocated)
	assert.Equal(t, expectedSummary.TotalSpent, result.TotalSpent)
	assert.Equal(t, expectedSummary.SpendingProgress, result.SpendingProgress)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBudgetFromTransaction(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	budgetID := uuid.New()
	categoryID := uuid.New()
	amount := 100.0

	existingBudget := &Budget{
		ID:     budgetID,
		UserID: userID,
		Name:   "Monthly Budget",
	}

	mockRepo.On("GetByID", mock.Anything, budgetID).Return(existingBudget, nil)
	mockRepo.On("UpdateSpentAmount", mock.Anything, budgetID, categoryID, amount).Return(nil)

	err := service.UpdateBudgetFromTransaction(ctx, userID, budgetID, categoryID, amount)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
