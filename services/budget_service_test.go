package services

import (
	"expense-tracker-api/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBudgetRepository is a mock implementation of BudgetRepository
type MockBudgetRepository struct {
	mock.Mock
}

func (m *MockBudgetRepository) Create(budget *domain.Budget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockBudgetRepository) GetByID(id int) (*domain.Budget, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) GetAll() ([]*domain.Budget, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) GetByMonth(month, year int) (*domain.Budget, error) {
	args := m.Called(month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Budget), args.Error(1)
}

func (m *MockBudgetRepository) Update(budget *domain.Budget) error {
	args := m.Called(budget)
	return args.Error(0)
}

func (m *MockBudgetRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockExpenseRepositoryForBudget is a mock for expense repository used in budget service
type MockExpenseRepositoryForBudget struct {
	mock.Mock
}

func (m *MockExpenseRepositoryForBudget) Create(expense *domain.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepositoryForBudget) GetByID(id int) (*domain.Expense, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepositoryForBudget) GetAll(filter *domain.ExpenseFilter) ([]*domain.Expense, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepositoryForBudget) Update(expense *domain.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepositoryForBudget) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockExpenseRepositoryForBudget) GetTotalByMonth(month, year int) (float64, error) {
	args := m.Called(month, year)
	return args.Get(0).(float64), args.Error(1)
}

func TestBudgetService_CreateOrUpdateBudget(t *testing.T) {
	t.Run("Successful creation", func(t *testing.T) {
		mockBudgetRepo := new(MockBudgetRepository)
		mockExpenseRepo := new(MockExpenseRepositoryForBudget)
		budgetService := NewBudgetService(mockBudgetRepo, mockExpenseRepo)

		mockBudgetRepo.On("GetByMonth", 1, 2024).Return(nil, domain.ErrNotFound)
		mockBudgetRepo.On("Create", mock.AnythingOfType("*domain.Budget")).Return(nil).Run(func(args mock.Arguments) {
			budget := args.Get(0).(*domain.Budget)
			budget.ID = 1
		})

		budget, err := budgetService.CreateOrUpdateBudget(1, 2024, 5000.0)
		assert.NoError(t, err)
		assert.NotNil(t, budget)
		assert.Equal(t, 5000.0, budget.BudgetAmount)
		mockBudgetRepo.AssertExpectations(t)
	})

	t.Run("Invalid month", func(t *testing.T) {
		mockBudgetRepo := new(MockBudgetRepository)
		mockExpenseRepo := new(MockExpenseRepositoryForBudget)
		budgetService := NewBudgetService(mockBudgetRepo, mockExpenseRepo)

		budget, err := budgetService.CreateOrUpdateBudget(13, 2024, 5000.0)
		assert.Error(t, err)
		assert.Nil(t, budget)
		assert.Equal(t, domain.ErrInvalidInput, err)
	})

	t.Run("Negative budget amount", func(t *testing.T) {
		mockBudgetRepo := new(MockBudgetRepository)
		mockExpenseRepo := new(MockExpenseRepositoryForBudget)
		budgetService := NewBudgetService(mockBudgetRepo, mockExpenseRepo)

		budget, err := budgetService.CreateOrUpdateBudget(1, 2024, -100.0)
		assert.Error(t, err)
		assert.Nil(t, budget)
		assert.Equal(t, domain.ErrInvalidInput, err)
	})
}

func TestBudgetService_GetBudgetByMonth(t *testing.T) {
	t.Run("Budget within limit", func(t *testing.T) {
		mockBudgetRepo := new(MockBudgetRepository)
		mockExpenseRepo := new(MockExpenseRepositoryForBudget)
		budgetService := NewBudgetService(mockBudgetRepo, mockExpenseRepo)

		budget := &domain.Budget{
			ID:           1,
			Month:        1,
			Year:         2024,
			BudgetAmount: 5000.0,
		}
		mockBudgetRepo.On("GetByMonth", 1, 2024).Return(budget, nil)
		mockExpenseRepo.On("GetTotalByMonth", 1, 2024).Return(3000.0, nil)

		budgetStatus, err := budgetService.GetBudgetByMonth(1, 2024)
		assert.NoError(t, err)
		assert.NotNil(t, budgetStatus)
		assert.Equal(t, "within_budget", budgetStatus.Status)
		assert.Equal(t, 2000.0, budgetStatus.Remaining)
		mockBudgetRepo.AssertExpectations(t)
		mockExpenseRepo.AssertExpectations(t)
	})

	t.Run("Budget exceeded", func(t *testing.T) {
		mockBudgetRepo := new(MockBudgetRepository)
		mockExpenseRepo := new(MockExpenseRepositoryForBudget)
		budgetService := NewBudgetService(mockBudgetRepo, mockExpenseRepo)

		budget := &domain.Budget{
			ID:           1,
			Month:        1,
			Year:         2024,
			BudgetAmount: 5000.0,
		}
		mockBudgetRepo.On("GetByMonth", 1, 2024).Return(budget, nil)
		mockExpenseRepo.On("GetTotalByMonth", 1, 2024).Return(6000.0, nil)

		budgetStatus, err := budgetService.GetBudgetByMonth(1, 2024)
		assert.NoError(t, err)
		assert.NotNil(t, budgetStatus)
		assert.Equal(t, "exceeded", budgetStatus.Status)
		assert.Equal(t, -1000.0, budgetStatus.Remaining)
		mockBudgetRepo.AssertExpectations(t)
		mockExpenseRepo.AssertExpectations(t)
	})
}
