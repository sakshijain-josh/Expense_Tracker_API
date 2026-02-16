package services

import (
	"expense-tracker-api/domain"
)

type BudgetService struct {
	budgetRepo  domain.BudgetRepository
	expenseRepo domain.ExpenseRepository
}

// NewBudgetService creates a new budget service
func NewBudgetService(budgetRepo domain.BudgetRepository, expenseRepo domain.ExpenseRepository) *BudgetService {
	return &BudgetService{
		budgetRepo:  budgetRepo,
		expenseRepo: expenseRepo,
	}
}

// CreateOrUpdateBudget creates or updates a monthly budget
func (s *BudgetService) CreateOrUpdateBudget(month, year int, budgetAmount float64) (*domain.Budget, error) {
	if month < 1 || month > 12 {
		return nil, domain.ErrInvalidInput
	}

	if budgetAmount < 0 {
		return nil, domain.ErrInvalidInput
	}

	// Check if budget already exists
	existingBudget, err := s.budgetRepo.GetByMonth(month, year)
	if err == nil && existingBudget != nil {
		// Update existing budget
		existingBudget.BudgetAmount = budgetAmount
		err = s.budgetRepo.Update(existingBudget)
		if err != nil {
			return nil, err
		}
		return existingBudget, nil
	}

	// Create new budget
	budget := &domain.Budget{
		Month:        month,
		Year:         year,
		BudgetAmount: budgetAmount,
	}

	err = s.budgetRepo.Create(budget)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

// GetBudgets retrieves all budgets
func (s *BudgetService) GetBudgets() ([]*domain.Budget, error) {
	return s.budgetRepo.GetAll()
}

// GetBudgetByMonth retrieves a budget for a specific month with status
func (s *BudgetService) GetBudgetByMonth(month, year int) (*domain.BudgetStatus, error) {
	budget, err := s.budgetRepo.GetByMonth(month, year)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Calculate spent amount for the month
	spentAmount, err := s.expenseRepo.GetTotalByMonth(month, year)
	if err != nil {
		return nil, err
	}

	remaining := budget.BudgetAmount - spentAmount
	status := "within_budget"
	if spentAmount > budget.BudgetAmount {
		status = "exceeded"
	}

	return &domain.BudgetStatus{
		Budget:      budget,
		SpentAmount: spentAmount,
		Remaining:   remaining,
		Status:      status,
	}, nil
}

// DeleteBudget deletes a budget
func (s *BudgetService) DeleteBudget(budgetID int) error {
	// Verify budget exists
	_, err := s.budgetRepo.GetByID(budgetID)
	if err != nil {
		return domain.ErrNotFound
	}

	return s.budgetRepo.Delete(budgetID)
}
