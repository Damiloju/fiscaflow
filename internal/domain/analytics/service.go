package analytics

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Service defines the interface for analytics business logic
type Service interface {
	// Categorization operations
	CategorizeTransaction(ctx context.Context, req *CategorizationRequest) (*CategorizationResponse, error)
	CreateCategorizationRule(ctx context.Context, req *CreateCategorizationRuleRequest) (*CategorizationRuleResponse, error)
	GetCategorizationRule(ctx context.Context, id uuid.UUID) (*CategorizationRuleResponse, error)
	ListCategorizationRules(ctx context.Context, offset, limit int) ([]CategorizationRuleResponse, error)
	UpdateCategorizationRule(ctx context.Context, id uuid.UUID, req *UpdateCategorizationRuleRequest) (*CategorizationRuleResponse, error)
	DeleteCategorizationRule(ctx context.Context, id uuid.UUID) error

	// Spending analysis operations
	AnalyzeSpending(ctx context.Context, userID uuid.UUID, req *SpendingAnalysisRequest) (*SpendingAnalysisResponse, error)
	GetSpendingInsights(ctx context.Context, userID uuid.UUID, periodStart, periodEnd time.Time) ([]SpendingInsight, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new analytics service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CategorizeTransaction categorizes a transaction using rule-based and ML approaches
func (s *service) CategorizeTransaction(ctx context.Context, req *CategorizationRequest) (*CategorizationResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "analytics.CategorizeTransaction",
		trace.WithAttributes(
			attribute.String("description", req.Description),
			attribute.String("merchant", req.Merchant),
			attribute.Float64("amount", req.Amount),
		),
	)
	defer span.End()

	// First, try rule-based categorization
	ruleMatch, err := s.categorizeByRules(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to categorize by rules: %w", err)
	}

	if ruleMatch != nil && ruleMatch.Confidence > 0.8 {
		span.SetAttributes(
			attribute.String("category_id", ruleMatch.CategoryID.String()),
			attribute.Float64("confidence", ruleMatch.Confidence),
			attribute.String("source", "rule"),
		)
		return ruleMatch, nil
	}

	// If no high-confidence rule match, try ML-based categorization
	mlMatch, err := s.categorizeByML(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("failed to categorize by ML: %w", err)
	}

	if mlMatch != nil {
		span.SetAttributes(
			attribute.String("category_id", mlMatch.CategoryID.String()),
			attribute.Float64("confidence", mlMatch.Confidence),
			attribute.String("source", "ml"),
		)
		return mlMatch, nil
	}

	// If no categorization found, return a default response
	defaultResponse := &CategorizationResponse{
		CategoryID:           uuid.Nil,
		CategoryName:         "Uncategorized",
		Confidence:           0.0,
		CategorizationSource: "manual",
	}

	span.SetAttributes(
		attribute.String("category_id", "uncategorized"),
		attribute.Float64("confidence", 0.0),
		attribute.String("source", "manual"),
	)

	return defaultResponse, nil
}

// categorizeByRules categorizes a transaction using rule-based matching
func (s *service) categorizeByRules(ctx context.Context, req *CategorizationRequest) (*CategorizationResponse, error) {
	rules, err := s.repo.GetActiveCategorizationRules(ctx)
	if err != nil {
		return nil, err
	}

	// Sort rules by priority (higher priority first)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority > rules[j].Priority
	})

	text := strings.ToLower(req.Description)
	if req.Merchant != "" {
		text += " " + strings.ToLower(req.Merchant)
	}

	for _, rule := range rules {
		if !rule.IsActive {
			continue
		}

		var matched bool
		switch rule.PatternType {
		case "exact":
			matched = strings.Contains(text, strings.ToLower(rule.Pattern))
		case "keyword":
			keywords := strings.Split(strings.ToLower(rule.Pattern), " ")
			matched = true
			for _, keyword := range keywords {
				if !strings.Contains(text, keyword) {
					matched = false
					break
				}
			}
		case "regex":
			re, err := regexp.Compile(strings.ToLower(rule.Pattern))
			if err != nil {
				continue // Skip invalid regex
			}
			matched = re.MatchString(text)
		}

		if matched {
			// Get category name
			category, err := s.repo.GetCategoryByID(ctx, rule.CategoryID)
			if err != nil {
				continue
			}

			confidence := s.calculateRuleConfidence(&rule, req.Amount)

			return &CategorizationResponse{
				CategoryID:           rule.CategoryID,
				CategoryName:         category.Name,
				Confidence:           confidence,
				CategorizationSource: "rule",
				MatchedPattern:       rule.Pattern,
			}, nil
		}
	}

	return nil, nil
}

// categorizeByML categorizes a transaction using simple ML (keyword frequency analysis)
func (s *service) categorizeByML(ctx context.Context, req *CategorizationRequest) (*CategorizationResponse, error) {
	// Simple ML approach: keyword frequency analysis based on historical data
	// In a real implementation, this would use a trained model

	text := strings.ToLower(req.Description)
	if req.Merchant != "" {
		text += " " + strings.ToLower(req.Merchant)
	}

	// Get historical categorization data for similar transactions
	similarTransactions, err := s.repo.GetSimilarTransactions(ctx, text, 10)
	if err != nil {
		return nil, err
	}

	if len(similarTransactions) == 0 {
		return nil, nil
	}

	// Count category frequencies
	categoryCounts := make(map[uuid.UUID]int)
	totalCount := 0

	for _, tx := range similarTransactions {
		if tx.CategoryID != nil {
			categoryCounts[*tx.CategoryID]++
			totalCount++
		}
	}

	if totalCount == 0 {
		return nil, nil
	}

	// Find the most common category
	var bestCategoryID uuid.UUID
	var bestCount int

	for categoryID, count := range categoryCounts {
		if count > bestCount {
			bestCategoryID = categoryID
			bestCount = count
		}
	}

	if bestCount == 0 {
		return nil, nil
	}

	// Calculate confidence based on frequency
	confidence := float64(bestCount) / float64(totalCount)

	// Adjust confidence based on amount similarity
	amountSimilarity := s.calculateAmountSimilarity(req.Amount, similarTransactions)
	confidence = (confidence + amountSimilarity) / 2

	// Get category name
	category, err := s.repo.GetCategoryByID(ctx, bestCategoryID)
	if err != nil {
		return nil, err
	}

	return &CategorizationResponse{
		CategoryID:           bestCategoryID,
		CategoryName:         category.Name,
		Confidence:           confidence,
		CategorizationSource: "ml",
	}, nil
}

// calculateRuleConfidence calculates confidence for rule-based categorization
func (s *service) calculateRuleConfidence(rule *CategorizationRule, amount float64) float64 {
	baseConfidence := 0.8

	// Adjust confidence based on pattern type
	switch rule.PatternType {
	case "exact":
		baseConfidence = 0.9
	case "keyword":
		baseConfidence = 0.85
	case "regex":
		baseConfidence = 0.8
	}

	// Adjust confidence based on amount (if amount is typical for the category)
	// This is a simplified approach - in reality, you'd have historical data
	if amount > 0 && amount < 1000 {
		baseConfidence += 0.05
	}

	return math.Min(baseConfidence, 1.0)
}

// calculateAmountSimilarity calculates similarity based on transaction amounts
func (s *service) calculateAmountSimilarity(amount float64, transactions []Transaction) float64 {
	if len(transactions) == 0 {
		return 0.5
	}

	// Calculate average amount
	var totalAmount float64
	var count int
	for _, tx := range transactions {
		totalAmount += tx.Amount
		count++
	}

	if count == 0 {
		return 0.5
	}

	avgAmount := totalAmount / float64(count)

	// Calculate similarity based on how close the amount is to the average
	diff := math.Abs(amount - avgAmount)
	similarity := 1.0 - (diff / avgAmount)

	return math.Max(0.0, math.Min(1.0, similarity))
}

// CreateCategorizationRule creates a new categorization rule
func (s *service) CreateCategorizationRule(ctx context.Context, req *CreateCategorizationRuleRequest) (*CategorizationRuleResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "analytics.CreateCategorizationRule",
		trace.WithAttributes(
			attribute.String("category_id", req.CategoryID.String()),
			attribute.String("pattern", req.Pattern),
			attribute.String("pattern_type", req.PatternType),
		),
	)
	defer span.End()

	// Validate pattern
	if err := s.validatePattern(req.Pattern, req.PatternType); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	rule := &CategorizationRule{
		CategoryID:  req.CategoryID,
		Pattern:     req.Pattern,
		PatternType: req.PatternType,
		Priority:    req.Priority,
		IsActive:    true,
	}

	if err := s.repo.CreateCategorizationRule(ctx, rule); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.String("rule_id", rule.ID.String()))
	return s.toCategorizationRuleResponse(rule), nil
}

// GetCategorizationRule retrieves a categorization rule
func (s *service) GetCategorizationRule(ctx context.Context, id uuid.UUID) (*CategorizationRuleResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "analytics.GetCategorizationRule",
		trace.WithAttributes(attribute.String("rule_id", id.String())),
	)
	defer span.End()

	rule, err := s.repo.GetCategorizationRuleByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return s.toCategorizationRuleResponse(rule), nil
}

// ListCategorizationRules retrieves categorization rules
func (s *service) ListCategorizationRules(ctx context.Context, offset, limit int) ([]CategorizationRuleResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "analytics.ListCategorizationRules",
		trace.WithAttributes(
			attribute.Int("offset", offset),
			attribute.Int("limit", limit),
		),
	)
	defer span.End()

	rules, err := s.repo.GetCategorizationRules(ctx, offset, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	responses := make([]CategorizationRuleResponse, len(rules))
	for i, rule := range rules {
		responses[i] = *s.toCategorizationRuleResponse(&rule)
	}

	span.SetAttributes(attribute.Int("rules_count", len(responses)))
	return responses, nil
}

// UpdateCategorizationRule updates a categorization rule
func (s *service) UpdateCategorizationRule(ctx context.Context, id uuid.UUID, req *UpdateCategorizationRuleRequest) (*CategorizationRuleResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "analytics.UpdateCategorizationRule",
		trace.WithAttributes(attribute.String("rule_id", id.String())),
	)
	defer span.End()

	rule, err := s.repo.GetCategorizationRuleByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Update fields
	if req.Pattern != nil {
		rule.Pattern = *req.Pattern
	}
	if req.PatternType != nil {
		rule.PatternType = *req.PatternType
	}
	if req.Priority != nil {
		rule.Priority = *req.Priority
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	// Validate pattern if updated
	if req.Pattern != nil || req.PatternType != nil {
		if err := s.validatePattern(rule.Pattern, rule.PatternType); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
	}

	if err := s.repo.UpdateCategorizationRule(ctx, rule); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return s.toCategorizationRuleResponse(rule), nil
}

// DeleteCategorizationRule deletes a categorization rule
func (s *service) DeleteCategorizationRule(ctx context.Context, id uuid.UUID) error {
	ctx, span := otel.Tracer("").Start(ctx, "analytics.DeleteCategorizationRule",
		trace.WithAttributes(attribute.String("rule_id", id.String())),
	)
	defer span.End()

	if err := s.repo.DeleteCategorizationRule(ctx, id); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// AnalyzeSpending analyzes spending for a user
func (s *service) AnalyzeSpending(ctx context.Context, userID uuid.UUID, req *SpendingAnalysisRequest) (*SpendingAnalysisResponse, error) {
	ctx, span := otel.Tracer("").Start(ctx, "analytics.AnalyzeSpending",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("start_date", req.StartDate.Format("2006-01-02")),
			attribute.String("end_date", req.EndDate.Format("2006-01-02")),
		),
	)
	defer span.End()

	// Get transactions for the period
	transactions, err := s.repo.GetTransactionsByPeriod(ctx, userID, req.StartDate, req.EndDate)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate basic metrics
	var totalSpent, totalIncome float64
	categorySpending := make(map[uuid.UUID]*CategorySpending)

	for _, tx := range transactions {
		if tx.Amount < 0 {
			totalSpent += math.Abs(tx.Amount)
		} else {
			totalIncome += tx.Amount
		}

		if tx.CategoryID != nil {
			if spending, exists := categorySpending[*tx.CategoryID]; exists {
				spending.Amount += math.Abs(tx.Amount)
				spending.TransactionCount++
			} else {
				category, _ := s.repo.GetCategoryByID(ctx, *tx.CategoryID)
				categoryName := "Uncategorized"
				if category != nil {
					categoryName = category.Name
				}

				categorySpending[*tx.CategoryID] = &CategorySpending{
					CategoryID:       *tx.CategoryID,
					CategoryName:     categoryName,
					Amount:           math.Abs(tx.Amount),
					TransactionCount: 1,
				}
			}
		}
	}

	// Calculate percentages
	for _, spending := range categorySpending {
		if totalSpent > 0 {
			spending.Percentage = (spending.Amount / totalSpent) * 100
		}
	}

	// Get top categories
	topCategories := s.getTopCategories(categorySpending, 5)

	// Generate spending trends
	spendingTrends := s.generateSpendingTrends(transactions, req.GroupBy)

	// Generate insights
	insights := s.generateSpendingInsights(transactions, categorySpending, totalSpent, totalIncome)

	response := &SpendingAnalysisResponse{
		PeriodStart:       req.StartDate,
		PeriodEnd:         req.EndDate,
		TotalSpent:        totalSpent,
		TotalIncome:       totalIncome,
		NetAmount:         totalIncome - totalSpent,
		CategoryBreakdown: s.mapToSlice(categorySpending),
		TopCategories:     topCategories,
		SpendingTrends:    spendingTrends,
		Insights:          insights,
	}

	span.SetAttributes(
		attribute.Float64("total_spent", totalSpent),
		attribute.Float64("total_income", totalIncome),
		attribute.Int("insights_count", len(insights)),
	)

	return response, nil
}

// GetSpendingInsights generates spending insights for a user
func (s *service) GetSpendingInsights(ctx context.Context, userID uuid.UUID, periodStart, periodEnd time.Time) ([]SpendingInsight, error) {
	ctx, span := otel.Tracer("").Start(ctx, "analytics.GetSpendingInsights",
		trace.WithAttributes(
			attribute.String("user_id", userID.String()),
			attribute.String("period_start", periodStart.Format("2006-01-02")),
			attribute.String("period_end", periodEnd.Format("2006-01-02")),
		),
	)
	defer span.End()

	// Get transactions for the period
	transactions, err := s.repo.GetTransactionsByPeriod(ctx, userID, periodStart, periodEnd)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Calculate category spending
	categorySpending := make(map[uuid.UUID]*CategorySpending)
	var totalSpent, totalIncome float64

	for _, tx := range transactions {
		if tx.Amount < 0 {
			totalSpent += math.Abs(tx.Amount)
		} else {
			totalIncome += tx.Amount
		}

		if tx.CategoryID != nil {
			if spending, exists := categorySpending[*tx.CategoryID]; exists {
				spending.Amount += math.Abs(tx.Amount)
				spending.TransactionCount++
			} else {
				category, _ := s.repo.GetCategoryByID(ctx, *tx.CategoryID)
				categoryName := "Uncategorized"
				if category != nil {
					categoryName = category.Name
				}

				categorySpending[*tx.CategoryID] = &CategorySpending{
					CategoryID:       *tx.CategoryID,
					CategoryName:     categoryName,
					Amount:           math.Abs(tx.Amount),
					TransactionCount: 1,
				}
			}
		}
	}

	insights := s.generateSpendingInsights(transactions, categorySpending, totalSpent, totalIncome)

	span.SetAttributes(attribute.Int("insights_count", len(insights)))
	return insights, nil
}

// Helper methods

func (s *service) validatePattern(pattern, patternType string) error {
	switch patternType {
	case "exact", "keyword":
		if pattern == "" {
			return fmt.Errorf("pattern cannot be empty")
		}
	case "regex":
		if _, err := regexp.Compile(pattern); err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
	default:
		return fmt.Errorf("invalid pattern type: %s", patternType)
	}
	return nil
}

func (s *service) getTopCategories(categorySpending map[uuid.UUID]*CategorySpending, limit int) []CategorySpending {
	// Convert map to slice
	categories := make([]CategorySpending, 0, len(categorySpending))
	for _, spending := range categorySpending {
		categories = append(categories, *spending)
	}

	// Sort by amount (descending)
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Amount > categories[j].Amount
	})

	// Return top N categories
	if len(categories) > limit {
		return categories[:limit]
	}
	return categories
}

func (s *service) generateSpendingTrends(transactions []Transaction, groupBy string) []SpendingTrend {
	// Simplified trend generation
	// In a real implementation, this would analyze historical data

	trends := []SpendingTrend{
		{
			Period: "Week 1",
			Amount: 500.0,
			Change: 0.0,
			Trend:  "stable",
		},
		{
			Period: "Week 2",
			Amount: 550.0,
			Change: 10.0,
			Trend:  "increasing",
		},
		{
			Period: "Week 3",
			Amount: 480.0,
			Change: -12.7,
			Trend:  "decreasing",
		},
	}

	return trends
}

func (s *service) generateSpendingInsights(transactions []Transaction, categorySpending map[uuid.UUID]*CategorySpending, totalSpent, totalIncome float64) []SpendingInsight {
	var insights []SpendingInsight

	// High spending category insight
	var highestSpending *CategorySpending
	var highestAmount float64

	for _, spending := range categorySpending {
		if spending.Amount > highestAmount {
			highestAmount = spending.Amount
			highestSpending = spending
		}
	}

	if highestSpending != nil && highestSpending.Percentage > 30 {
		insights = append(insights, SpendingInsight{
			Type:        "pattern",
			Title:       "High Spending Category",
			Description: fmt.Sprintf("You're spending %.1f%% of your budget on %s", highestSpending.Percentage, highestSpending.CategoryName),
			Severity:    "medium",
			Data: map[string]interface{}{
				"category_id":   highestSpending.CategoryID.String(),
				"category_name": highestSpending.CategoryName,
				"percentage":    highestSpending.Percentage,
				"amount":        highestSpending.Amount,
			},
			CreatedAt: time.Now(),
		})
	}

	// Spending vs Income insight
	if totalIncome > 0 {
		spendingRatio := totalSpent / totalIncome
		if spendingRatio > 0.9 {
			insights = append(insights, SpendingInsight{
				Type:        "trend",
				Title:       "High Spending Ratio",
				Description: fmt.Sprintf("You're spending %.1f%% of your income", spendingRatio*100),
				Severity:    "high",
				Data: map[string]interface{}{
					"spending_ratio": spendingRatio,
					"total_spent":    totalSpent,
					"total_income":   totalIncome,
				},
				CreatedAt: time.Now(),
			})
		}
	}

	// Transaction frequency insight
	if len(transactions) > 50 {
		insights = append(insights, SpendingInsight{
			Type:        "pattern",
			Title:       "High Transaction Frequency",
			Description: fmt.Sprintf("You have %d transactions in this period", len(transactions)),
			Severity:    "low",
			Data: map[string]interface{}{
				"transaction_count": len(transactions),
			},
			CreatedAt: time.Now(),
		})
	}

	return insights
}

func (s *service) mapToSlice(categorySpending map[uuid.UUID]*CategorySpending) []CategorySpending {
	result := make([]CategorySpending, 0, len(categorySpending))
	for _, spending := range categorySpending {
		result = append(result, *spending)
	}
	return result
}

func (s *service) toCategorizationRuleResponse(rule *CategorizationRule) *CategorizationRuleResponse {
	// Get category name
	category, _ := s.repo.GetCategoryByID(context.Background(), rule.CategoryID)
	categoryName := "Unknown"
	if category != nil {
		categoryName = category.Name
	}

	return &CategorizationRuleResponse{
		ID:           rule.ID,
		CategoryID:   rule.CategoryID,
		CategoryName: categoryName,
		Pattern:      rule.Pattern,
		PatternType:  rule.PatternType,
		Priority:     rule.Priority,
		IsActive:     rule.IsActive,
		CreatedAt:    rule.CreatedAt,
		UpdatedAt:    rule.UpdatedAt,
	}
}
