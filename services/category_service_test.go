package services

import (
	"expense-tracker-api/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository is a mock implementation of CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByID(id int) (*domain.Category, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) GetAll() ([]*domain.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category *domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func (m *MockCategoryRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCategoryService_CreateCategory(t *testing.T) {
	t.Run("Successful creation", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		categoryService := NewCategoryService(mockRepo)

		mockRepo.On("Create", mock.AnythingOfType("*domain.Category")).Return(nil).Run(func(args mock.Arguments) {
			category := args.Get(0).(*domain.Category)
			category.ID = 1
		})

		category, err := categoryService.CreateCategory("Food")
		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, "Food", category.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty name", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		categoryService := NewCategoryService(mockRepo)

		category, err := categoryService.CreateCategory("")
		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Equal(t, domain.ErrInvalidInput, err)
	})
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	t.Run("Successful update", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		categoryService := NewCategoryService(mockRepo)

		existingCategory := &domain.Category{ID: 1, Name: "Old Name"}
		mockRepo.On("GetByID", 1).Return(existingCategory, nil)
		mockRepo.On("Update", mock.AnythingOfType("*domain.Category")).Return(nil)

		category, err := categoryService.UpdateCategory(1, "New Name")
		assert.NoError(t, err)
		assert.NotNil(t, category)
		assert.Equal(t, "New Name", category.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Category not found", func(t *testing.T) {
		mockRepo := new(MockCategoryRepository)
		categoryService := NewCategoryService(mockRepo)

		mockRepo.On("GetByID", 1).Return(nil, domain.ErrNotFound)

		category, err := categoryService.UpdateCategory(1, "New Name")
		assert.Error(t, err)
		assert.Nil(t, category)
		assert.Equal(t, domain.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
