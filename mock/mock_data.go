package mock

import (
	"expense-tracker-api/domain"
	"time"
)

// GenerateMockCategory creates a mock category
func GenerateMockCategory(id int, name string) *domain.Category {
	return &domain.Category{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
	}
}

// GenerateMockExpense creates a mock expense
func GenerateMockExpense(id, categoryID int, amount float64, paymentMode domain.PaymentMode) *domain.Expense {
	return &domain.Expense{
		ID:          id,
		CategoryID:  categoryID,
		Amount:      amount,
		Description: "Mock expense description",
		PaymentMode: paymentMode,
		ExpenseDate: time.Now(),
		CreatedAt:   time.Now(),
	}
}

// GenerateMockBudget creates a mock budget
func GenerateMockBudget(id, month, year int, budgetAmount float64) *domain.Budget {
	return &domain.Budget{
		ID:           id,
		Month:        month,
		Year:         year,
		BudgetAmount: budgetAmount,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// GenerateMockCategories generates multiple mock categories
func GenerateMockCategories(names []string) []*domain.Category {
	categories := make([]*domain.Category, len(names))
	for i, name := range names {
		categories[i] = GenerateMockCategory(i+1, name)
	}
	return categories
}

// GenerateMockExpenses generates multiple mock expenses
func GenerateMockExpenses(categoryID int, count int) []*domain.Expense {
	expenses := make([]*domain.Expense, count)
	paymentModes := []domain.PaymentMode{domain.PaymentModeUPI, domain.PaymentModeCash}
	for i := 0; i < count; i++ {
		expenses[i] = GenerateMockExpense(i+1, categoryID, float64((i+1)*100), paymentModes[i%2])
	}
	return expenses
}
