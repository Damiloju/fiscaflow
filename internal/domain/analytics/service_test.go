package analytics_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"fiscaflow/internal/domain/analytics"
	"fiscaflow/internal/domain/analytics/mocks"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)
	service := analytics.NewService(mockRepo)

	assert.NotNil(t, service)
}

func TestCategorizeTransaction_RuleBased(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mocks.NewMockRepository(ctrl)
	service := analytics.NewService(mockRepo)

	rule := analytics.CategorizationRule{
		ID:          uuid.New(),
		CategoryID:  uuid.New(),
		Pattern:     "walmart",
		PatternType: "keyword",
		Priority:    1,
		IsActive:    true,
	}
	mockRepo.EXPECT().GetActiveCategorizationRules(gomock.Any()).Return([]analytics.CategorizationRule{rule}, nil)
	mockRepo.EXPECT().GetCategoryByID(gomock.Any(), rule.CategoryID).Return(&analytics.Category{ID: rule.CategoryID, Name: "Food & Groceries"}, nil)

	request := &analytics.CategorizationRequest{
		Description: "Walmart grocery purchase",
		Amount:      150.50,
		Merchant:    "Walmart",
	}
	resp, err := service.CategorizeTransaction(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, rule.CategoryID, resp.CategoryID)
	assert.Equal(t, "Food & Groceries", resp.CategoryName)
	assert.Equal(t, "rule", resp.CategorizationSource)
}
