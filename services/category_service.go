package services

import (
	"expense-tracker-api/domain"
)

type CategoryService struct {
	categoryRepo domain.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo domain.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(name string) (*domain.Category, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	category := &domain.Category{
		Name: name,
	}

	err := s.categoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategories retrieves all categories
func (s *CategoryService) GetCategories() ([]*domain.Category, error) {
	return s.categoryRepo.GetAll()
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(id int) (*domain.Category, error) {
	return s.categoryRepo.GetByID(id)
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(categoryID int, name string) (*domain.Category, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	// Verify category exists
	category, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	category.Name = name
	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(categoryID int) error {
	// Verify category exists
	_, err := s.categoryRepo.GetByID(categoryID)
	if err != nil {
		return domain.ErrNotFound
	}

	return s.categoryRepo.Delete(categoryID)
}
