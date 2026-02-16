package services

import (
	"expense-tracker-api/domain"
	"time"
)

type ExpenseService struct {
	expenseRepo  domain.ExpenseRepository
	categoryRepo domain.CategoryRepository
	budgetRepo   domain.BudgetRepository
}

// NewExpenseService creates a new expense service
func NewExpenseService(expenseRepo domain.ExpenseRepository, categoryRepo domain.CategoryRepository, budgetRepo domain.BudgetRepository) *ExpenseService {
	return &ExpenseService{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
		budgetRepo:   budgetRepo,
	}
}

// CreateExpense creates a new expense
func (s *ExpenseService) CreateExpense(expense *domain.Expense) (*domain.Expense, error) {
	// Validate payment mode
	if !expense.PaymentMode.IsValid() {
		return nil, domain.ErrInvalidPaymentMode
	}

	// Verify category exists
	_, err := s.categoryRepo.GetByID(expense.CategoryID)
	if err != nil {
		return nil, domain.ErrInvalidCategory
	}

	if expense.ExpenseDate.IsZero() {
		expense.ExpenseDate = time.Now()
	}

	err = s.expenseRepo.Create(expense)
	if err != nil {
		return nil, err
	}

	// Check budget status
	s.checkBudget(expense)

	return expense, nil
}

// checkBudget checks if the monthly budget is exceeded
func (s *ExpenseService) checkBudget(expense *domain.Expense) {
	month := int(expense.ExpenseDate.Month())
	year := expense.ExpenseDate.Year()

	budget, err := s.budgetRepo.GetByMonth(month, year)
	if err == nil && budget != nil {
		spentAmount, err := s.expenseRepo.GetTotalByMonth(month, year)
		if err == nil && spentAmount > budget.BudgetAmount {
			expense.Warning = "Warning: Monthly budget exceeded!"
		}
	}
}

// GetExpenses retrieves expenses with optional filters
func (s *ExpenseService) GetExpenses(filter *domain.ExpenseFilter) ([]*domain.Expense, error) {
	if filter == nil {
		filter = &domain.ExpenseFilter{}
	}
	return s.expenseRepo.GetAll(filter)
}

// GetExpenseByID retrieves an expense by ID
func (s *ExpenseService) GetExpenseByID(expenseID int) (*domain.Expense, error) {
	expense, err := s.expenseRepo.GetByID(expenseID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return expense, nil
}

// UpdateExpense updates an expense
func (s *ExpenseService) UpdateExpense(expense *domain.Expense) (*domain.Expense, error) {
	// Verify expense exists
	existingExpense, err := s.expenseRepo.GetByID(expense.ID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Validate payment mode if provided
	if expense.PaymentMode != "" && !expense.PaymentMode.IsValid() {
		return nil, domain.ErrInvalidPaymentMode
	}

	// Verify category exists if category is being updated
	if expense.CategoryID != 0 {
		_, err := s.categoryRepo.GetByID(expense.CategoryID)
		if err != nil {
			return nil, domain.ErrInvalidCategory
		}
	}

	// Update fields
	if expense.CategoryID != 0 {
		existingExpense.CategoryID = expense.CategoryID
	}
	if expense.Amount != 0 {
		existingExpense.Amount = expense.Amount
	}
	if expense.Description != "" {
		existingExpense.Description = expense.Description
	}
	if expense.PaymentMode != "" {
		existingExpense.PaymentMode = expense.PaymentMode
	}
	if !expense.ExpenseDate.IsZero() {
		existingExpense.ExpenseDate = expense.ExpenseDate
	}

	err = s.expenseRepo.Update(existingExpense)
	if err != nil {
		return nil, err
	}

	// Check budget status for updated expense
	s.checkBudget(existingExpense)

	return existingExpense, nil
}

// DeleteExpense deletes an expense
func (s *ExpenseService) DeleteExpense(expenseID int) error {
	// Verify expense exists
	_, err := s.expenseRepo.GetByID(expenseID)
	if err != nil {
		return domain.ErrNotFound
	}

	return s.expenseRepo.Delete(expenseID)
}
