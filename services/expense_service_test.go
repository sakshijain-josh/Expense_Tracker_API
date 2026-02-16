package services

import (
	"expense-tracker-api/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExpenseRepository is a mock implementation of ExpenseRepository
type MockExpenseRepository struct {
	mock.Mock
}

func (m *MockExpenseRepository) Create(expense *domain.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) GetByID(id int) (*domain.Expense, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepository) GetAll(filter *domain.ExpenseFilter) ([]*domain.Expense, error) {
	args := m.Called(filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Expense), args.Error(1)
}

func (m *MockExpenseRepository) Update(expense *domain.Expense) error {
	args := m.Called(expense)
	return args.Error(0)
}

func (m *MockExpenseRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockExpenseRepository) GetTotalByMonth(month, year int) (float64, error) {
	args := m.Called(month, year)
	return args.Get(0).(float64), args.Error(1)
}

// MockCategoryRepository for expense service tests
type MockCategoryRepositoryForExpense struct {
	mock.Mock
}

func (m *MockCategoryRepositoryForExpense) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepositoryForExpense) GetByID(id int) (*domain.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForExpense) GetAll() ([]*domain.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepositoryForExpense) Update(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepositoryForExpense) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestExpenseService_CreateExpense(t *testing.T) {
	t.Run("Successful creation", func(t *testing.T) {
		mockExpenseRepo := new(MockExpenseRepository)
		mockCategoryRepo := new(MockCategoryRepositoryForExpense)
		expenseService := NewExpenseService(mockExpenseRepo, mockCategoryRepo)

		category := &domain.Category{ID: 1}
		mockCategoryRepo.On("GetByID", 1).Return(category, nil)
		mockExpenseRepo.On("Create", mock.AnythingOfType("*domain.Expense")).Return(nil).Run(func(args mock.Arguments) {
			expense := args.Get(0).(*domain.Expense)
			expense.ID = 1
		})

		expense := &domain.Expense{
			CategoryID:  1,
			Amount:      100.0,
			Description: "Test expense",
			PaymentMode: domain.PaymentModeUPI,
		}

		createdExpense, err := expenseService.CreateExpense(expense)
		assert.NoError(t, err)
		assert.NotNil(t, createdExpense)
		mockExpenseRepo.AssertExpectations(t)
		mockCategoryRepo.AssertExpectations(t)
	})

	t.Run("Invalid payment mode", func(t *testing.T) {
		mockExpenseRepo := new(MockExpenseRepository)
		mockCategoryRepo := new(MockCategoryRepositoryForExpense)
		expenseService := NewExpenseService(mockExpenseRepo, mockCategoryRepo)

		expense := &domain.Expense{
			CategoryID:  1,
			Amount:      100.0,
			PaymentMode: domain.PaymentMode("Invalid"),
		}

		createdExpense, err := expenseService.CreateExpense(expense)
		assert.Error(t, err)
		assert.Nil(t, createdExpense)
		assert.Equal(t, domain.ErrInvalidPaymentMode, err)
	})

	t.Run("Invalid category", func(t *testing.T) {
		mockExpenseRepo := new(MockExpenseRepository)
		mockCategoryRepo := new(MockCategoryRepositoryForExpense)
		expenseService := NewExpenseService(mockExpenseRepo, mockCategoryRepo)

		mockCategoryRepo.On("GetByID", 1).Return(nil, domain.ErrNotFound)

		expense := &domain.Expense{
			CategoryID:  1,
			Amount:      100.0,
			PaymentMode: domain.PaymentModeUPI,
		}

		createdExpense, err := expenseService.CreateExpense(expense)
		assert.Error(t, err)
		assert.Nil(t, createdExpense)
		assert.Equal(t, domain.ErrInvalidCategory, err)
		mockCategoryRepo.AssertExpectations(t)
	})
}
